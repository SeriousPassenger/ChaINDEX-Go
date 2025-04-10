package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func GetRpcClient(config *Config) (*rpc.Client, error) {
	client, err := rpc.Dial(config.Rpc.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC server: %w", err)
	}

	return client, nil
}

func GetEthClient(config *Config) (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.Rpc.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	return client, nil
}

func TestConnection(config *Config) error {
	client, err := GetEthClient(config)
	if err != nil {
		return err
	}

	blockNumber, err := GetLastBlockNumber(client)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully connected to RPC server. Last block number: %s\n", blockNumber)
	return nil
}
