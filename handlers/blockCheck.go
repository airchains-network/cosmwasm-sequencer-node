package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	"github.com/airchains-network/cosmwasm-sequencer-node/types"
	"github.com/syndtr/goleveldb/leveldb"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func getLastProcessedBlock(db *leveldb.DB) int {
	lastBlockKey := []byte("lastProcessedBlock")
	data, err := db.Get(lastBlockKey, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			// If not found, return 0 indicating start from the beginning
			return 0
		}
		logs.LogMessage("ERROR:", fmt.Sprintf("Error reading last processed block from LevelDB : %s", err.Error()))
		os.Exit(0)
	}
	lastBlockNum, _ := strconv.Atoi(string(data))
	return lastBlockNum
}

func BlockCheck(wg *sync.WaitGroup, ldb *leveldb.DB, ldt *leveldb.DB) {
	defer wg.Done()

	rpcUrl := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/latest", common.ExecutionClientTRPC)
	res, resErr := http.Get(rpcUrl)
	if resErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in comswasm RPC : %s"+resErr.Error()))
		os.Exit(0)
	}
	defer res.Body.Close()

	bodyBlockHeight, bodyBlockHeightErr := ioutil.ReadAll(res.Body)
	if bodyBlockHeightErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Reading in comswasm RPC : %s"+bodyBlockHeightErr.Error()))
		os.Exit(0)
	}

	var blockHeight types.BlockObject
	error := json.Unmarshal(bodyBlockHeight, &blockHeight)
	if error != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Unmarshal in comswasm RPC : %s"+error.Error()))
		os.Exit(0)
	}
	latestBlock := blockHeight.Block.Header.Height

	numLatestBlock, err := strconv.Atoi(latestBlock)
	if err != nil {
		panic(err)
	}
	lastBlock := getLastProcessedBlock(ldb)
	startBlock := lastBlock + 1

	OldBlocks(startBlock, numLatestBlock, ldb, ldt)
	NewBlocks(numLatestBlock, ldb, ldt)

}
