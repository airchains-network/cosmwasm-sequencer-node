package air_leveldb

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log"
)

var txDbInstance *leveldb.DB
var blockDbInstance *leveldb.DB
var staticDbInstance *leveldb.DB
var batchesDbInstance *leveldb.DB
var proofDbInstance *leveldb.DB
var publicWitnessDbInstance *leveldb.DB
var daDbInstance *leveldb.DB

// InitTxDb The function initializes a LevelDB database for transactions and returns a boolean indicating
// whether the initialization was successful.
func InitTxDb() bool {
	txDB, err := leveldb.OpenFile("data/leveldb/tx", nil)
	if err != nil {
		log.Fatal("Failed to open transaction LevelDB:", err)
		return false
	}
	txDbInstance = txDB
	return true
}

// InitBlockDb The function initializes a LevelDB database for storing blocks and returns a boolean indicating
// whether the initialization was successful.
func InitBlockDb() bool {
	blockDB, err := leveldb.OpenFile("data/leveldb/blocks", nil)
	if err != nil {
		log.Fatal("Failed to open block LevelDB:", err)
		return false
	}
	blockDbInstance = blockDB
	return true
}

// InitStaticDb The function initializes a static LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitStaticDb() bool {
	staticDB, err := leveldb.OpenFile("data/leveldb/static", nil)
	if err != nil {
		log.Fatal("Failed to open static LevelDB:", err)
		return false
	}
	staticDbInstance = staticDB
	return true
}

// InitBatchesDb The function initializes a batches LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitBatchesDb() bool {
	batchesDB, err := leveldb.OpenFile("data/leveldb/batches", nil)
	if err != nil {
		log.Fatal("Failed to open batches LevelDB:", err)
		return false
	}
	batchesDbInstance = batchesDB
	return true
}

// InitProofDb The function initializes a proof LevelDB database and returns a boolean indicating whether the
// initialization was successful or not.
func InitProofDb() bool {
	proofDB, err := leveldb.OpenFile("data/leveldb/proof", nil)
	if err != nil {
		log.Fatal("Failed to open proof LevelDB:", err)
		return false
	}
	proofDbInstance = proofDB
	return true
}

func InitPublicWitnessDb() bool {
	publicWitnessDB, err := leveldb.OpenFile("data/leveldb/publicWitness", nil)
	if err != nil {
		log.Fatal("Failed to open public witness LevelDB:", err)
		return false
	}
	publicWitnessDbInstance = publicWitnessDB
	return true
}

func InitDaDb() bool {
	daDB, err := leveldb.OpenFile("data/leveldb/da", nil)
	if err != nil {
		log.Fatal("Failed to open da LevelDB:", err)
		return false
	}
	daDbInstance = daDB
	return true
}

// InitDb The function `InitDb` initializes 4 different databases and returns true if all of them are
// successfully initialized, otherwise it returns false.
func InitDb() bool {
	txStatus := InitTxDb()
	blockStatus := InitBlockDb()
	staticStatus := InitStaticDb()
	batchesStatus := InitBatchesDb()
	proofStatus := InitProofDb()
	publicWitnessStatus := InitPublicWitnessDb()
	daDbInstanceStatus := InitDaDb()

	if txStatus && blockStatus && staticStatus && batchesStatus && proofStatus && publicWitnessStatus && daDbInstanceStatus {
		return true
	} else {
		return false
	}
}

// GetTxDbInstance The function GetTxDbInstance returns the instance of the leveldb database.
func GetTxDbInstance() *leveldb.DB {
	return txDbInstance
}

// GetBlockDbInstance The function returns the instance of the block database.
func GetBlockDbInstance() *leveldb.DB {
	return blockDbInstance
}

// GetStaticDbInstance The function `GetStaticDbInstance()` is returning the instance of the LevelDB database that was
// initialized in the `InitStaticDb()` function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetStaticDbInstance() *leveldb.DB {
	return staticDbInstance
}

// GetBatchesDbInstance The function `GetBatchesDbInstance()` is returning the instance of the LevelDB database that was
// initialized in the `InitBatchesDb()` function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetBatchesDbInstance() *leveldb.DB {
	return batchesDbInstance
}

// GetProofDbInstance The function `GetProofDbInstance()` is returning the instance of the LevelDB database that was
// initialized in the `InitProofDb()` function. This allows other parts of the code to access and use
// the LevelDB database instance for performing operations such as reading or writing data.
func GetProofDbInstance() *leveldb.DB {
	return proofDbInstance
}

func GetPublicWitnessDbInstance() *leveldb.DB {
	return publicWitnessDbInstance
}

func GetDaDbInstance() *leveldb.DB {
	return daDbInstance
}
