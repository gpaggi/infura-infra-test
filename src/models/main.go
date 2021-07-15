package models

// Represents a block
type Block struct {
	BlockNum   int64         `json:"blockNumber"`
	Hash       string        `json:"hash"`
	Difficulty uint64        `json:"difficulty"`
	Timestamp  uint64        `json:"timestamp"`
	TxCount    int           `json:"transactionsCount"`
	Txs        []Transaction `json:"transactions"`
}

// Represents a transaction
type Transaction struct {
	Hash     string `json:"hash"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Nonce    uint64 `json:"nonce"`
	Pending  bool   `json:"pending"`
	Gas      uint64 `json:"gas"`
	GasPrice uint64 `json:"gasPrice"`
}

// Error data structure
type Error struct {
	Code uint64 `json:"code"`
	Msg  string `json:"message"`
}
