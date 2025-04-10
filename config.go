package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	SampleConfigPath = "config.toml.sample"

	// Default values for the RPC configuration
	DefaultRpcUrl    = "ws://localhost:8546" // Default URL for the RPC server
	DefaultRpcDelay  = 1000                  // Default delay between requests in milliseconds
	DefaultBatchSize = 1                     // Default batch size for requests

	// Default values for the scanning configuration
	DefaultFromBlock = 1
	DefaultToBlock   = 100
	ScanBalances     = true
	ScanReceipts     = true
	ScanContractCode = true
	DefaultOutputDir = "output"
)

var DefaultFilterAddresses = []string{}

type RpcConfig struct {
	Url       string `toml:"url"`        // URL for the RPC server
	Delay     uint64 `toml:"delay"`      // Delay between requests in milliseconds
	BatchSize uint64 `toml:"batch_size"` // Batch size for requests
}

type ScanConfig struct {
	FromBlock    uint64 `toml:"from_block"`         // Starting block for scanning
	ToBlock      uint64 `toml:"to_block"`           // Ending block for scanning
	Balances     bool   `toml:"scan_balances"`      // Flag to indicate if balances should be scanned
	Receipts     bool   `toml:"scan_receipts"`      // Flag to indicate if receipts should be scanned
	ContractCode bool   `toml:"scan_contract_code"` // Flag to indicate if contract code should be scanned
	OutputDir    string `toml:"output_dir"`         // Directory to save the output files
}

type FilterConfig struct {
	Addresses []string `toml:"addresses"` // List of addresses to filter
}

type Config struct {
	Rpc    RpcConfig    `toml:"rpc"`    // RPC configuration
	Scan   ScanConfig   `toml:"scan"`   // Scanning configuration
	Filter FilterConfig `toml:"filter"` // Filter configuration
}

func CreateSampleConfig() error {
	sampleConfig := new(Config)

	sampleConfig.Rpc.Url = DefaultRpcUrl
	sampleConfig.Rpc.Delay = DefaultRpcDelay
	sampleConfig.Rpc.BatchSize = DefaultBatchSize

	sampleConfig.Scan.FromBlock = DefaultFromBlock
	sampleConfig.Scan.ToBlock = DefaultToBlock
	sampleConfig.Scan.Balances = ScanBalances
	sampleConfig.Scan.Receipts = ScanReceipts
	sampleConfig.Scan.ContractCode = ScanContractCode
	sampleConfig.Scan.OutputDir = DefaultOutputDir

	sampleConfig.Filter.Addresses = DefaultFilterAddresses

	// Marshal the sample configuration to TOML format
	configData, err := toml.Marshal(sampleConfig)

	if err != nil {
		return fmt.Errorf("failed to marshal sample config: %w", err)
	}

	file, err := os.Create(SampleConfigPath)

	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	defer file.Close()

	// Write the sample configuration to the file
	_, err = file.Write(configData)

	if err != nil {
		return fmt.Errorf("failed to write sample config to file: %w", err)
	}

	return nil
}

func GetConfig(path string) (*Config, error) {
	config := new(Config)

	// Read the configuration file
	if _, err := toml.DecodeFile(path, config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return config, nil
}
