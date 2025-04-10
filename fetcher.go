package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
)

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
	Root              string   `json:"root"`
	From              string   `json:"from"`
	EffectiveGasPrice string   `json:"effectiveGasPrice"`
	Type              string   `json:"type"`
	ChainId           string   `json:"chainId"`
	V                 string   `json:"v"`
	R                 string   `json:"r"`
	S                 string   `json:"s"`
}

func GetLastBlockNumber(client *rpc.Client) (string, error) {
	var blockNumber string
	err := client.Call(&blockNumber, "eth_blockNumber")

	if err != nil {
		return "", fmt.Errorf("failed to get last block number: %w", err)
	}

	blockNumberInt, err := HexToBigInt(blockNumber)

	if err != nil {
		return "", fmt.Errorf("failed to convert block number to int: %w", err)
	}

	// Convert the block number to a string
	blockNumber = blockNumberInt.String()

	return blockNumber, nil
}

func GetBlocksBatch(client *rpc.Client, blocks []*big.Int) ([]RpcBlock, error) {
	var batch []rpc.BatchElem

	for _, block := range blocks {
		blockHex := BigIntToHex(block)
		var raw json.RawMessage
		batch = append(batch, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{blockHex, true},
			Result: &raw,
		})
	}

	err := client.BatchCall(batch)

	if err != nil {
		return nil, fmt.Errorf("failed to execute batch call: %w", err)
	}

	responses := make([]RpcBlock, 0)

	for _, elem := range batch {
		if elem.Error != nil {
			log.Printf("error in batch element: %v", elem.Error)
			continue
		}

		raw, ok := elem.Result.(*json.RawMessage)
		if !ok || raw == nil {
			continue
		}

		var response RpcBlock
		if err := json.Unmarshal(*raw, &response); err != nil {
			log.Printf("failed to unmarshal JSON: %v", err)
			continue
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func GetReceiptsBatch(client *rpc.Client, transactions []string) ([]RpcReceipt, error) {
	var batch []rpc.BatchElem

	for _, transaction := range transactions {
		var raw json.RawMessage
		batch = append(batch, rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{transaction},
			Result: &raw,
		})
	}

	err := client.BatchCall(batch)

	if err != nil {
		return nil, fmt.Errorf("failed to execute batch call: %w", err)
	}

	responses := make([]RpcReceipt, 0)

	for _, elem := range batch {
		if elem.Error != nil {
			log.Printf("error in batch element: %v", elem.Error)
			continue
		}

		raw, ok := elem.Result.(*json.RawMessage)
		if !ok || raw == nil {
			continue
		}

		var response RpcReceipt
		if err := json.Unmarshal(*raw, &response); err != nil {
			log.Printf("failed to unmarshal JSON: %v", err)
			continue
		}

		responses = append(responses, response)
	}

	return responses, nil
}
