package main

import (
	"encoding/json"
	"fmt"
	"log"
	air "github.com/airchains-network/cosmwasm-sequencer-node/airdb/air-leveldb"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	"github.com/airchains-network/cosmwasm-sequencer-node/handlers"
	settlement_client "github.com/airchains-network/cosmwasm-sequencer-node/handlers/settlement-client"
	"github.com/airchains-network/cosmwasm-sequencer-node/prover"
	"github.com/airchains-network/cosmwasm-sequencer-node/types"
	"github.com/joho/godotenv"
	"os"
	"sync"
	"time"
)

func main() {
	logs.LogMessage("INFO:", "Starting Cosmwasm Seq Indexer")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	daClientRPC := os.Getenv("DA_CLIENT_RPC")

	if daClientRPC == "" {
		logs.LogMessage("ERROR:", "DA_CLIENT_RPC is not set")
		os.Exit(0)
	}

	dbStatus := air.InitDb()
	if !dbStatus {
		logs.LogMessage("ERROR:", "Error in initializing db")
		os.Exit(0)
	}

	prover.CreateVkPk()
	chainId := settlement_client.AddExecutionLayer()
	if chainId == "nil" {
		logs.LogMessage("ERROR:", "Something went wrong while adding execution layer")
		logs.LogMessage("INFO:", "Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
		chainId = settlement_client.AddExecutionLayer()
	} else if chainId == "exist" {
		logs.LogMessage("INFO:", "Chain already exist")
	}

	ldt := air.GetTxDbInstance()
	ldb := air.GetBlockDbInstance()
	lds := air.GetStaticDbInstance()
	ldbatch := air.GetBatchesDbInstance()
	ldda := air.GetDaDbInstance()

	da := types.DAStruct{
		DAKey:             "0",
		DAClientName:      "0",
		BatchNumber:       "0",
		PreviousStateHash: "0",
		CurrentStateHash:  "0",
	}

	daBytes, err := json.Marshal(da)

	_, err = ldda.Get([]byte("batch_0"), nil)
	if err != nil {
		err = ldda.Put([]byte("batch_0"), daBytes, nil)
		if err != nil {
			logs.LogMessage("ERROR:", fmt.Sprintf("Error in saving da in static db : %s", err.Error()))
			os.Exit(0)
		}
	}

	batchStartIndex, err := lds.Get([]byte("batchStartIndex"), nil)
	if err != nil {
		err = lds.Put([]byte("batchStartIndex"), []byte("0"), nil)
		if err != nil {
			logs.LogMessage("ERROR:", fmt.Sprintf("Error in saving batchStartIndex in static db : %s", err.Error()))
			os.Exit(0)
		}
	}

	_, err = lds.Get([]byte("batchCount"), nil)
	if err != nil {
		err = lds.Put([]byte("batchCount"), []byte("0"), nil)
		if err != nil {
			logs.LogMessage("ERROR:", fmt.Sprintf("Error in saving batchCount in static db : %s", err.Error()))
			os.Exit(0)
		}
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		handlers.BlockCheck(&wg, ldb, ldt)
	}()
	go func() {
		defer wg.Done()
		handlers.BatchGeneration(&wg, lds, ldt, ldbatch, ldda, batchStartIndex)
	}()
	wg.Wait()
}
