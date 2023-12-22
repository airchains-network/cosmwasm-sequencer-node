package types

//type BlockStuct struct {
//	Blockheight       int      `json:"blockHeight"`
//	Blockhash         string   `json:"blockhash"`
//	Parentslot        int      `json:"parentSlot"`
//	Previousblockhash string   `json:"previousBlockhash"`
//	Timestamp         int      `json:"timestamp"`
//	TransactionLength int      `json:"transactionLength"`
//	Transactions      []string `json:"transactions"`
//}
//
//type TransactionStruck struct {
//	Signature   string      `json:"signature"`
//	BlockHeight int         `json:"blockHeight"`
//	BlockHash   string      `json:"blockHash"`
//	Timestamp   int         `json:"timestamp"`
//	Mata        Meta        `json:"mata"`
//	Transaction Transaction `json:"transaction"`
//}

type BatchStruct struct {
	From              []string `json:"from"`
	To                []string `json:"to"`
	Amounts           []string `json:"amounts"`
	TransactionHash   []string `json:"tx_hashes"`
	SenderBalances    []string `json:"sender_balances"`
	ReceiverBalances  []string `json:"receiver_balances"`
	Messages          []string `json:"messages"`
	TransactionNonces []string `json:"tx_nonces"`
	AccountNonces     []string `json:"account_nonces"`
}

type DAStruct struct {
	DAKey             string `json:"da_key"`
	DAClientName      string `json:"da_client_name"`
	BatchNumber       string `json:"batch_number"`
	PreviousStateHash string `json:"previous_state_hash"`
	CurrentStateHash  string `json:"current_state_hash"`
}

type SettlementLayerChainInfoStruct struct {
	ChainId   string `json:"chain_id"`
	ChainName string `json:"chain_name"`
}
