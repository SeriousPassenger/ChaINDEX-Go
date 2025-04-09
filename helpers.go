package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/rpc"
)

func HexToBigInt(hexString string) (*big.Int, error) {
	if len(hexString) < 3 || hexString[:2] != "0x" {
		return nil, fmt.Errorf("invalid hex string: %s", hexString)
	}

	// Convert the hex string to a big.Int
	value := new(big.Int)
	value.SetString(hexString[2:], 16)

	return value, nil
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
