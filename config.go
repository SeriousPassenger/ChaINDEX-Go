package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	SampleConfigPath = "config.toml.sample"

	// Default values for the RPC configuration
	DefaultRpcUrl    = "http://localhost:8080" // Default URL for the RPC server
	DefaultRpcDelay  = 1000                    // Default delay between requests in milliseconds
	DefaultBatchSize = 1                       // Default batch size for requests

	// Default values for the scanning configuration
	DefaultFromBlock = 1
	DefaultToBlock   = 100
	ScanBalances     = true
	ScanReceipts     = true
	ScanContractCode = true
	DefaultOutputDir = "output"
)

type RpcConfig struct {
	Url       string `toml:"url"`        // URL for the RPC server
	Delay     int    `toml:"delay"`      // Delay between requests in milliseconds
	BatchSize int    `toml:"batch_size"` // Batch size for requests
}

type ScanConfig struct {
	FromBlock    string `toml:"from_block"`         // Starting block for scanning
	ToBlock      string `toml:"to_block"`           // Ending block for scanning
	Balances     bool   `toml:"scan_balances"`      // Flag to indicate if balances should be scanned
	Receipts     bool   `toml:"scan_receipts"`      // Flag to indicate if receipts should be scanned
	ContractCode bool   `toml:"scan_contract_code"` // Flag to indicate if contract code should be scanned
	OutputDir    string `toml:"output_dir"`         // Directory to save the output files
}

type Config struct {
	Rpc  RpcConfig  `toml:"rpc"`  // RPC configuration
	Scan ScanConfig `toml:"scan"` // Scanning configuration
}

func CreateSampleConfig() error {
	sampleConfig := new(Config)

	sampleConfig.Rpc.Url = DefaultRpcUrl
	sampleConfig.Rpc.Delay = DefaultRpcDelay
	sampleConfig.Rpc.BatchSize = DefaultBatchSize

	sampleConfig.Scan.FromBlock = fmt.Sprintf("%d", DefaultFromBlock)
	sampleConfig.Scan.ToBlock = fmt.Sprintf("%d", DefaultToBlock)
	sampleConfig.Scan.Balances = ScanBalances
	sampleConfig.Scan.Receipts = ScanReceipts
	sampleConfig.Scan.ContractCode = ScanContractCode
	sampleConfig.Scan.OutputDir = DefaultOutputDir

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
