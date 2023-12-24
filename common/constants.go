package common

import "os"

const BatchSize = 25
const BlockDelay = 5

const ExecutionClientTRPC = "http://localhost:1317"
const ExecutionClientJsonRPC = "http://localhost:26657"
var DaClientRPC = os.Getenv("DA_CLIENT_RPC")
const SettlementClientRPC = "http://localhost:8080"
