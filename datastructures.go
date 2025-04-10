package main

import "math/big"

/* =========================================================================== */
/* ====================== RPC Response Data Structures ======================= */
/* =========================================================================== */

type RpcBlock struct {
	Difficulty       string           `json:"difficulty"`
	ExtraData        string           `json:"extraData"`
	GasLimit         string           `json:"gasLimit"`
	GasUsed          string           `json:"gasUsed"`
	Hash             string           `json:"hash"`
	LogsBloom        string           `json:"logsBloom"`
	Miner            string           `json:"miner"`
	MixHash          string           `json:"mixHash"`
	Nonce            string           `json:"nonce"`
	Number           string           `json:"number"`
	ParentHash       string           `json:"parentHash"`
	ReceiptsRoot     string           `json:"receiptsRoot"`
	Sha3Uncles       string           `json:"sha3Uncles"`
	Size             string           `json:"size"`
	StateRoot        string           `json:"stateRoot"`
	Timestamp        string           `json:"timestamp"`
	TotalDifficulty  string           `json:"totalDifficulty"`
	Transactions     []RpcTransaction `json:"transactions"`
	TransactionsRoot string           `json:"transactionsRoot"`
	Uncles           []string         `json:"uncles"`
}

type RpcTransaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value            string `json:"value"`
	Type             string `json:"type"`
	ChainId          string `json:"chainId"`
	V                string `json:"v"`
	R                string `json:"r"`
	S                string `json:"s"`
}

type RpcLog struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex string   `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         string   `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

type RpcReceipt struct {
	BlockHash         string   `json:"blockHash"`
	BlockNumber       string   `json:"blockNumber"`
	ContractAddress   string   `json:"contractAddress"`
	CumulativeGasUsed string   `json:"cumulativeGasUsed"`
	GasUsed           string   `json:"gasUsed"`
	Status            string   `json:"status"`
	To                string   `json:"to"`
	TransactionHash   string   `json:"transactionHash"`
	TransactionIndex  string   `json:"transactionIndex"`
	Logs              []RpcLog `json:"logs"`
	LogsBloom         string   `json:"logsBloom"`
	From              string   `json:"from"`
	EffectiveGasPrice string   `json:"effectiveGasPrice"`
	Type              string   `json:"type"`
}

/* ========================================================================== */
/* ====================== Converted Data Structures ======================= */
/* ========================================================================== */

type BaseBlock struct {
	Difficulty       *big.Int `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         *big.Int `json:"gasLimit"`
	GasUsed          *big.Int `json:"gasUsed"`
	Hash             string   `json:"hash"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            *big.Int `json:"nonce"`
	Number           *big.Int `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             *big.Int `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        *big.Int `json:"timestamp"`
	TotalDifficulty  *big.Int `json:"totalDifficulty"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}

type Transaction struct {
	BlockHash        string   `json:"blockHash"`
	BlockNumber      *big.Int `json:"blockNumber"`
	From             string   `json:"from"`
	Gas              *big.Int `json:"gas"`
	GasPrice         *big.Int `json:"gasPrice"`
	Hash             string   `json:"hash"`
	Input            string   `json:"input"`
	Nonce            *big.Int `json:"nonce"`
	To               string   `json:"to"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	Value            *big.Int `json:"value"`
	Type             *big.Int `json:"type"`
	ChainId          *big.Int `json:"chainId"`
	V                string   `json:"v"`
	R                string   `json:"r"`
	S                string   `json:"s"`
}

type Log struct {
	Address          string   `json:"address"`
	Topics           []string `json:"topics"`
	Data             string   `json:"data"`
	BlockNumber      *big.Int `json:"blockNumber"`
	TransactionHash  string   `json:"transactionHash"`
	TransactionIndex *big.Int `json:"transactionIndex"`
	BlockHash        string   `json:"blockHash"`
	LogIndex         *big.Int `json:"logIndex"`
	Removed          bool     `json:"removed"`
}

type Receipt struct {
	BlockHash         string   `json:"blockHash"`
	BlockNumber       *big.Int `json:"blockNumber"`
	ContractAddress   string   `json:"contractAddress"`
	CumulativeGasUsed *big.Int `json:"cumulativeGasUsed"`
	GasUsed           *big.Int `json:"gasUsed"`
	Status            string   `json:"status"`
	To                string   `json:"to"`
	TransactionHash   string   `json:"transactionHash"`
	TransactionIndex  *big.Int `json:"transactionIndex"`
	Logs              []Log    `json:"logs"`
	LogsBloom         string   `json:"logsBloom"`
	From              string   `json:"from"`
	EffectiveGasPrice *big.Int `json:"effectiveGasPrice"`
	Type              string   `json:"type"`
}

type FullTransaction struct {
	BaseTransaction Transaction `json:"baseTransaction"`
	Receipt         Receipt     `json:"receipt"`
}

type FullBlock struct {
	BaseBlock        BaseBlock         `json:"baseBlock"`
	FullTransactions []FullTransaction `json:"fullTransactions"`
}

type ContractCode struct {
	Address string `json:"address"`
	Code    string `json:"code"`
}

type BalanceSheet struct {
	Address   string   `json:"address"`
	Balance   *big.Int `json:"balance"`
	UpdatedAt int64    `json:"updatedAt"`
}
