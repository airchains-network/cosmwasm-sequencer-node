package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	air "github.com/airchains-network/cosmwasm-sequencer-node/airdb/air-leveldb"
	"github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	"github.com/airchains-network/cosmwasm-sequencer-node/types"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
	"strconv"
	"time"
)

func DaCall(transactions []string, currentStateHash string, batchNumber int, ldda *leveldb.DB) (string, error) {

	proofGet, proofGetErr := air.GetProofDbInstance().Get([]byte(fmt.Sprintf("proof_%d", batchNumber)), nil)
	if proofGetErr != nil {
		return "", proofGetErr
	}

	var proofDecode types.ProofStruct

	proofDecodeErr := json.Unmarshal(proofGet, &proofDecode)
	if proofDecodeErr != nil {
		return "", proofDecodeErr
	}

	daGet, daGetErr := ldda.Get([]byte(fmt.Sprintf("batch_%d", batchNumber-1)), nil)
	if daGetErr != nil {
		return "", daGetErr
	}

	var daDecode types.DAStruct
	daDecodeErr := json.Unmarshal(daGet, &daDecode)
	if daDecodeErr != nil {
		return "", daDecodeErr
	}

	DaStruct := types.DAUploadStruct{
		Proof:             proofDecode,
		TxnHashes:         transactions,
		CurrentStateHash:  currentStateHash,
		PreviousStateHash: daDecode.PreviousStateHash,
		MetaData: struct {
			ChainID     string `json:"chainID"`
			BatchNumber int    `json:"batchNumber"`
		}{
			ChainID:     "0",
			BatchNumber: batchNumber,
		},
	}

	payloadJSON, payloadJSONErr := json.Marshal(DaStruct)

	if payloadJSONErr != nil {
		return "", payloadJSONErr
	}

	client := &http.Client{}

	req, reqErr := http.NewRequest("POST", common.DaClientRPC, bytes.NewBuffer(payloadJSON))
	if reqErr != nil {
		return "", reqErr
	}

	req.Header.Set("Content-Type", "application/json")
	res, resErr := client.Do(req)
	if resErr != nil {
		return "", resErr
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error in da call : %s", res.Status))
		time.Sleep(1 * time.Second)
		DaCall(transactions, currentStateHash, batchNumber, ldda)
	}

	var response types.DAResponceStruct
	decodeErr := json.NewDecoder(res.Body).Decode(&response)
	if decodeErr != nil {
		return "", decodeErr
	}

	if response.DaKeyHash == "nil" {
		logs.LogMessage("ERROR:", fmt.Sprintf("DA RPC is not responding, retrying in 1 second"))
		time.Sleep(3 * time.Second)
		DaCall(transactions, currentStateHash, batchNumber, ldda)
	}

	da := types.DAStruct{
		DAKey:             response.DaKeyHash,
		DAClientName:      "celestia",
		BatchNumber:       strconv.Itoa(batchNumber),
		PreviousStateHash: daDecode.CurrentStateHash,
		CurrentStateHash:  currentStateHash,
	}

	daBytes, err := json.Marshal(da)

	batchKey := fmt.Sprintf("batch_%d", batchNumber)
	err = ldda.Put([]byte(batchKey), daBytes, nil)
	if err != nil {
		return "", err
	}

	return response.DaKeyHash, nil
}
