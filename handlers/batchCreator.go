package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	settlement_client "github.com/airchains-network/cosmwasm-sequencer-node/handlers/settlement-client"
	"github.com/airchains-network/cosmwasm-sequencer-node/prover"
	"github.com/airchains-network/cosmwasm-sequencer-node/types"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func BatchGeneration(wg *sync.WaitGroup, lds *leveldb.DB, ldt *leveldb.DB, ldbatch *leveldb.DB, ldDA *leveldb.DB, batchStartIndex []byte) {
	defer wg.Done()

	limit, err := lds.Get([]byte("batchCount"), nil)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in getting batchCount from static db : %s", err.Error()))
		os.Exit(0)
	}
	limitInt, _ := strconv.Atoi(strings.TrimSpace(string(limit)))
	batchStartIndexInt, _ := strconv.Atoi(strings.TrimSpace(string(batchStartIndex)))

	var batch types.BatchStruct

	var From []string
	var To []string
	var Amounts []string
	var TransactionHash []string
	var SenderBalances []string
	var ReceiverBalances []string
	var Messages []string
	var TransactionNonces []string
	var AccountNonces []string

	for i := batchStartIndexInt; i < (common.BatchSize * (limitInt + 1)); i++ {
		findKey := fmt.Sprintf("txns-%d", i+1)
		txData, err := ldt.Get([]byte(findKey), nil)

		if err != nil {
			i--
			time.Sleep(1 * time.Second)
			continue
		}

		var txn types.BatchTransaction

		err = json.Unmarshal(txData, &txn)
		if err != nil {
			logs.LogMessage("ERROR:", fmt.Sprintf("Error in unmarshalling tx data : %s", err.Error()))
			os.Exit(0)
		}

		fromCheck := common.Bech32Decoder(txn.Tx.Body.Messages[0].FromAddress)
		toCheck := common.Bech32Decoder(txn.Tx.Body.Messages[0].ToAddress)
		transactionHashCheck := common.TXHashCheck(txn.TxResponse.TxHash)

		senderBalancesCheck := common.AccountBalanceCheck(txn.Tx.Body.Messages[0].FromAddress, txn.TxResponse.Height)
		receiverBalancesCheck := common.AccountBalanceCheck(txn.Tx.Body.Messages[0].ToAddress, txn.TxResponse.Height)
		accountNoncesCheck := common.AccountNounceCheck(txn.Tx.Body.Messages[0].FromAddress)

		From = append(From, fromCheck)
		To = append(To, toCheck)
		Amounts = append(Amounts, txn.Tx.Body.Messages[0].Amount[0].Amount)
		SenderBalances = append(SenderBalances, senderBalancesCheck)
		ReceiverBalances = append(ReceiverBalances, receiverBalancesCheck)
		TransactionHash = append(TransactionHash, transactionHashCheck)
		Messages = append(Messages, fmt.Sprint(txn.Tx.Body.Messages[0]))
		TransactionNonces = append(TransactionNonces, "0")
		AccountNonces = append(AccountNonces, accountNoncesCheck)
	}

	batch.From = From
	batch.To = To
	batch.Amounts = Amounts
	batch.TransactionHash = TransactionHash
	batch.SenderBalances = SenderBalances
	batch.ReceiverBalances = ReceiverBalances
	batch.Messages = Messages
	batch.TransactionNonces = TransactionNonces
	batch.AccountNonces = AccountNonces

	// add prover here
	witnessVector, currentStatusHash, proofByte, pkErr := prover.GenerateProof(batch, limitInt+1)
	if pkErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in generating proof : %s", pkErr.Error()))
		os.Exit(0)
	}

	// !adding Da client here
	daKeyHash, err := DaCall(batch.TransactionHash, currentStatusHash, limitInt+1, ldDA)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in adding Da client : %s", err.Error()))
		logs.LogMessage("INFO:", "Waiting for 3 seconds")
		time.Sleep(3 * time.Second)
		BatchGeneration(wg, lds, ldt, ldbatch, ldDA, []byte(strconv.Itoa(common.BatchSize*(limitInt+1))))
	}
	logs.LogMessage("SUCCESS:", fmt.Sprintf("Successfully added Da client for Batch %s in the latest phase", daKeyHash))

	addBatchRes := settlement_client.AddBatch(witnessVector, limitInt+1, lds)
	if addBatchRes == "nil" {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in adding batch to settlement client : %s", addBatchRes))
		os.Exit(0)
	}

	status := settlement_client.VerifyBatch(limitInt+1, proofByte, ldDA, lds)
	if !status {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in verifying batch to settlement client : %s", status))
		os.Exit(0)
	}

	logs.LogMessage("SUCCESS:", fmt.Sprintf("Successfully generated proof for Batch %s in the latest phase", strconv.Itoa(limitInt+1)))

	batchJSON, err := json.Marshal(batch)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in marshalling batch data : %s", err.Error()))
		os.Exit(0)
	}

	// !writing batch data to file
	batchKey := fmt.Sprintf("batch-%d", limitInt+1)
	err = ldbatch.Put([]byte(batchKey), batchJSON, nil)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in writing batch data to file : %s", err.Error()))
		os.Exit(0)
	}

	// !updating batchStartIndex in static db
	err = lds.Put([]byte("batchStartIndex"), []byte(strconv.Itoa(common.BatchSize*(limitInt+1))), nil)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in updating batchStartIndex in static db : %s", err.Error()))
		os.Exit(0)
	}

	// !updating batchCount in static db
	err = lds.Put([]byte("batchCount"), []byte(strconv.Itoa(limitInt+1)), nil)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in updating batchCount in static db : %s", err.Error()))
		os.Exit(0)
	}

	err = os.WriteFile("data/batchCount.txt", []byte(strconv.Itoa(limitInt+1)), 0666)
	if err != nil {
		panic("Failed to update batch number: " + err.Error())
	}

	logs.LogMessage("SUCCESS:", fmt.Sprintf("Successfully saved Batch %s in the latest phase", strconv.Itoa(limitInt+1)))

	BatchGeneration(wg, lds, ldt, ldbatch, ldDA, []byte(strconv.Itoa(common.BatchSize*(limitInt+1))))
}
