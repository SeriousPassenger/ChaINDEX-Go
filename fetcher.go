package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

// GetLastBlockNumber retrieves the latest block number from the Ethereum client.
func GetLastBlockNumber(client *ethclient.Client) (string, error) {
	blockNumber, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to get block number: %w", err)
	}

	return blockNumber.Number().String(), nil
}

// GetBlocksBatch retrieves a batch of blocks from the Ethereum client.
func GetBlocksBatch(client *rpc.Client, blocks []*big.Int) ([]RpcBlock, error) {
	var batch []rpc.BatchElem

	for _, block := range blocks {
		blockHex := BigIntToHex(block)
		var raw json.RawMessage
		batch = append(batch, rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []any{blockHex, true},
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

// GetReceiptsBatch retrieves a batch of transaction receipts from the Ethereum client.
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

// GetContractCodeBatch retrieves the contract code for a batch of addresses from the Ethereum client.
func GetContractCodeBatch(client *rpc.Client, addresses []string) ([]RpcContractCode, error) {
	var batch []rpc.BatchElem

	for _, address := range addresses {
		var raw json.RawMessage
		batch = append(batch, rpc.BatchElem{
			Method: "eth_getCode",
			Args:   []interface{}{address, "latest"},
			Result: &raw,
		})
	}

	err := client.BatchCall(batch)

	if err != nil {
		return nil, fmt.Errorf("failed to execute batch call: %w", err)
	}

	responses := make([]RpcContractCode, 0)

	for _, elem := range batch {
		if elem.Error != nil {
			log.Printf("error in batch element: %v", elem.Error)
			continue
		}

		raw, ok := elem.Result.(*json.RawMessage)
		if !ok || raw == nil {
			continue
		}

		var response RpcContractCode
		if err := json.Unmarshal(*raw, &response); err != nil {
			log.Printf("failed to unmarshal JSON: %v", err)
			continue
		}

		responses = append(responses, response)
	}

	return responses, nil
}
