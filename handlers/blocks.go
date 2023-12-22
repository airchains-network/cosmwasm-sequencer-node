package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/types"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

func SaveLastProcessedBlock(blockNum int, db *leveldb.DB) {
	lastBlockKey := []byte("lastProcessedBlock")
	lastBlockValue := []byte(strconv.Itoa(blockNum))
	if err := db.Put(lastBlockKey, lastBlockValue, nil); err != nil {
		log.Println("Error saving last processed block to LevelDB:", err)
	}
}

func OldBlocks(startBlock int, numLatestBlock int, db *leveldb.DB, txnDB *leveldb.DB) {
	for i := startBlock; i <= numLatestBlock; i++ {
		fmt.Println("Saving block number:", i)
		rpcUrl := fmt.Sprintf("%s/block?height=%d", common.ExecutionClientJsonRPC, i)
		//resp, err := http.Get("http://192.168.1.33:26657/block?height=" + strconv.Itoa(i))
		resp, err := http.Get(rpcUrl)
		if err != nil {
			log.Println("Error fetching block:", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // Close the response body here
		if err != nil {
			log.Println("Error reading block body:", err)
			continue
		}

		var blockData types.Response

		jsonErr := json.Unmarshal(body, &blockData)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		if len(blockData.Result.Block.Data.Txs) > 0 {

			// fmt.Printf("txsJSON: %s\n", txsJSON)
			SaveTransaction(blockData.Result.Block.Data.Txs, txnDB)
		}

		var responseMap map[string]interface{}
		if err := json.Unmarshal(body, &responseMap); err != nil {
			log.Fatal("Error unmarshalling JSON:", err) // Consider if fatal is appropriate here
		}
		if result, ok := responseMap["result"]; ok {
			resultJSON, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Fatal("Error marshalling JSON:", err)
			}

			blockKey := []byte("Block" + strconv.Itoa(i))

			if err = db.Put(blockKey, resultJSON, nil); err != nil {
				log.Println("Error saving block to LevelDB:", err)
				continue
			}

			SaveLastProcessedBlock(i, db)

		}

	}
}

func GetCurrentBlock() (types.BlockObject, error) {
	rpcUrl := fmt.Sprintf("%s/cosmos/base/tendermint/v1beta1/blocks/latest", common.ExecutionClientTRPC)
	res, err := http.Get(rpcUrl)
	if err != nil {
		return types.BlockObject{}, err
	}
	defer res.Body.Close()

	var data types.BlockObject
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return types.BlockObject{}, err
	}

	return data, nil
}

func watchBlocks(currentBlockHeight int, db *leveldb.DB, txnDB *leveldb.DB) {
	var currentBlock types.BlockObject

	for {
		latestBlock, err := GetCurrentBlock()
		if err != nil {
			fmt.Println("Error fetching current block:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if currentBlock.Block.Header.Height == latestBlock.Block.Header.Height {
			fmt.Println("No new blocks. Waiting for 3 second.")
			time.Sleep(3 * time.Second)
			continue
		}

		fmt.Println("New block:", latestBlock.Block.Header.Height)
		rpcUrl := fmt.Sprintf("%s/block?height=%s", common.ExecutionClientJsonRPC, latestBlock.Block.Header.Height)
		//resp, err := http.Get("http://192.168.1.33:26657/block?height=" + latestBlock.Block.Header.Height)
		resp, err := http.Get(rpcUrl)
		if err != nil {
			fmt.Println("Error fetching block details:", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close() // Explicitly close the body
		if err != nil {
			fmt.Println("Error reading response body:", err)
			continue
		}

		var blockData types.Response

		error := json.Unmarshal(body, &blockData)
		if error != nil {
			log.Fatal(error)
		}

		if len(blockData.Result.Block.Data.Txs) > 0 {

			SaveTransaction(blockData.Result.Block.Data.Txs, txnDB)
		}

		var responseMap map[string]interface{}
		if err := json.Unmarshal(body, &responseMap); err != nil {
			log.Fatal("Error unmarshalling JSON:", err)
		}

		if result, ok := responseMap["result"]; ok {
			resultJSON, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				log.Fatal("Error marshalling JSON:", err)
			}

			blockKey := []byte("Block" + latestBlock.Block.Header.Height)
			if err = db.Put(blockKey, resultJSON, nil); err != nil {
				log.Println("Error saving block to LevelDB:", err)
			}

			height, err := strconv.Atoi(latestBlock.Block.Header.Height)
			if err != nil {
				log.Fatal("Error converting block height to integer:", err)
			}
			SaveLastProcessedBlock(height, db)
		} else {
			fmt.Println("Result key not found in response")
		}

		currentBlock = latestBlock
	}
}

func NewBlocks(curentBlock int, db *leveldb.DB, txnDB *leveldb.DB) {
	// Start watching for new blocks
	watchBlocks(curentBlock, db, txnDB)
}
