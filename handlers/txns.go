package handlers

import (
	"bufio"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func SaveTransaction(txn []interface{}, db *leveldb.DB) {
	fmt.Println("Saving Transactions ......⏳")
	fmt.Println("Total Transactions to be saved: ", len(txn))

	for i, tx := range txn {
		fmt.Println("Processing Transaction: ", i)

		hash, err := ComputeTransactionHash(tx.(string))
		if err != nil {
			log.Println("Error computing transaction hash:", err)
			continue
		}

		rpcUrl := fmt.Sprintf("%s/cosmos/tx/v1beta1/txs/%s", common.ExecutionClientTRPC, hash)
		respo, err := http.Get(rpcUrl)
		if err != nil {
			log.Println("HTTP request failed for transaction hash:", err)
			continue
		}
		if respo != nil {
			bodyTxnHash, err := ioutil.ReadAll(respo.Body)
			respo.Body.Close() // Close immediately after reading
			if err != nil {
				log.Println("Error reading response body:", err)
				continue
			}

			fileOpen, err := os.Open("data/transactionCount.txt")
			if err != nil {
				logs.LogMessage("ERROR:", fmt.Sprintf("Failed to read file: %s"+err.Error()))
				os.Exit(0)
			}
			defer fileOpen.Close()

			scanner := bufio.NewScanner(fileOpen)

			transactionNumberBytes := ""

			for scanner.Scan() {
				transactionNumberBytes = scanner.Text()
			}

			transactionNumber, err := strconv.Atoi(strings.TrimSpace(string(transactionNumberBytes)))
			if err != nil {
				logs.LogMessage("ERROR:", fmt.Sprintf("Invalid transaction number : %s"+err.Error()))
				os.Exit(0)
			}

			// fmt.Println("Transaction Data:", string(bodyTxnHash))

			//txnKey := []byte(hash)
			txnsKey := fmt.Sprintf("txns-%d", transactionNumber+1)
			if err = db.Put([]byte(txnsKey), bodyTxnHash, nil); err != nil {
				log.Println("Error saving txn to LevelDB:", err)
				continue
			} else {
				err = os.WriteFile("data/transactionCount.txt", []byte(strconv.Itoa(transactionNumber+1)), 0666)
				if err != nil {
					logs.LogMessage("ERROR:", fmt.Sprintf("Failed to update transaction number: %s"+err.Error()))
					os.Exit(0)
				}
				fmt.Println("Transaction saved successfully:", txnsKey)
			}

		} else {
			log.Println("Received nil response for transaction hash:", hash)
		}
	}
}

func ComputeTransactionHash(base64Tx string) (string, error) {
	// Decode the base64 transaction data
	txBytes, err := base64.StdEncoding.DecodeString(base64Tx)
	if err != nil {
		return "", err
	}

	// Compute the SHA-256 hash
	hash := sha256.Sum256(txBytes)

	// Convert the hash to a hexadecimal string
	txHash := hex.EncodeToString(hash[:])

	return txHash, nil
}
