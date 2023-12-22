package settlement_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	air "github.com/airchains-network/cosmwasm-sequencer-node/airdb/air-leveldb"
	"github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	"github.com/airchains-network/cosmwasm-sequencer-node/types"
	"io"
	"net/http"
	"os"
	"time"
)

type PostAddExecutionLayerStruct struct {
	VerificationKey string `json:"verification_key"`
	ChainInfo       string `json:"chain_info"`
}

func AddExecutionLayer() string {

	logs.LogMessage("INFO:", "Adding execution layer")

	if _, err := os.Stat("verificationKey.json"); os.IsNotExist(err) {
		logs.LogMessage("INFO:", "Waiting for verificationKey.json file")
		time.Sleep(5 * time.Second)
		AddExecutionLayer()
	}

	verificationKeyContents, err := os.ReadFile("verificationKey.json")
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error reading verificationKey.json file : %s", err.Error()))
		return "nil"
	}

	verificationKeyContentsAsString := string(verificationKeyContents)

	chainInfoFile, err := os.ReadFile("config/chainInfo.json")
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error reading chainInfo.json file : %s", err.Error()))
		os.Exit(0)
	}

	var chainInfo types.ChainInfoStruct

	err = json.Unmarshal(chainInfoFile, &chainInfo)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error unmarshalling chainInfo.json file : %s", err.Error()))
		os.Exit(0)
	}

	chainInfoAsString, err := json.Marshal(chainInfo.ChainInfo)

	postAddExecutionLayerStruct := PostAddExecutionLayerStruct{
		VerificationKey: verificationKeyContentsAsString,
		ChainInfo:       string(chainInfoAsString),
	}

	jsonData, err := json.Marshal(postAddExecutionLayerStruct)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error marshalling postAddExecutionLayerStruct : %s", err.Error()))
		return "nil"
	}
	rpcUrl := fmt.Sprintf("%s/addexelayer", common.SettlementClientRPC)
	req, err := http.NewRequest("POST", rpcUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error creating request:", err))
		return "nil"
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error sending request:", err))
		return "nil"
	}
	defer resp.Body.Close()

	// Read and unmarshal the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error reading response:", err))
		return "nil"
	}

	var response types.SettlementClientResponseStruct
	err = json.Unmarshal(body, &response)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Error unmarshalling response:", err))
		return "nil"
	}

	if response.Data != "nil" && response.Data != "exist" {
		var settlementChainInfo = types.SettlementLayerChainInfoStruct{
			ChainId:   response.Data,
			ChainName: chainInfo.ChainInfo.Moniker,
		}

		settlementChainInfoBytes, err := json.Marshal(settlementChainInfo)
		if err != nil {
			logs.LogMessage("ERROR:", fmt.Sprintf("Error marshalling settlementChainInfo : %s", err.Error()))
			return "nil"
		}

		err = air.GetStaticDbInstance().Put([]byte("settlementChainInfo"), settlementChainInfoBytes, nil)
		if err != nil {
			logs.LogMessage("ERROR:", fmt.Sprintf("Error in putting settlementChainInfo in static db : %s", err.Error()))
			return "nil"
		}
	}

	return response.Data
}
