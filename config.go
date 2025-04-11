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

	DefaultAccountScanBlockNumber = 48_000_000 // Default block number for account scanning

	DefaultOutputDir = "output"
)

var DefaultFilterAddresses = []string{}

type RpcConfig struct {
	Url   string `toml:"url"`   // URL for the RPC server
	Delay uint64 `toml:"delay"` // Delay between requests in milliseconds
}

type BlockScanConfig struct {
	FromBlock      uint64 `toml:"from_block"`       // Starting block for scanning
	ToBlock        uint64 `toml:"to_block"`         // Ending block for scanning
	OutputFileName string `toml:"output_file_name"` // File name for saving the scanned blocks
	BatchSize      uint64 `toml:"batch_size"`       // Batch size for requests

}

type AccountScanConfig struct {
	BlockNumber    uint64 `toml:"block_number"`     // The state block number to scan
	MaxAccounts    uint64 `toml:"max_accounts"`     // Maximum number of accounts to scan
	StartKey       string `toml:"start_key"`        // Starting key for scanning (used for pagination)
	OutputFileName string `toml:"output_file_name"` // File name for saving the scanned accounts
	BatchSize      uint64 `toml:"batch_size"`       // Batch size for requests
}

type ContractCodeScanConfig struct {
	OutputFileName string `toml:"output_file_name"` // File name for saving the scanned contract codes
	BatchSize      uint64 `toml:"batch_size"`       // Batch size for requests
}

type ReceiptScanConfig struct {
	// TODO: For now, this only accepts the exported BlockScan data,
	// but it should be able to accept any array of transaction hashes
	FullBlocksFile string `toml:"full_blocks_file"` // File containing full blocks for scanning
	// TODO: TransactionHashesFile string `toml:"transaction_hashes_file"` // List of transaction hashes to scan
	OutputFileName string `toml:"output_file_name"` // File name for saving the scanned receipts
	BatchSize      uint64 `toml:"batch_size"`       // Batch size for requests

}

type ScanConfig struct {
	BlockScanConfig        `toml:"block_scan"`         // Configuration for block scanning
	AccountScanConfig      `toml:"account_scan"`       // Configuration for account scanning
	ReceiptScanConfig      `toml:"receipt_scan"`       // Configuration for receipt scanning
	ContractCodeScanConfig `toml:"contract_code_scan"` // Configuration for contract code scanning
	OutputDir              string                      `toml:"output_dir"` // Directory to save the output files
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

	sampleConfig.Scan.BlockScanConfig.FromBlock = DefaultFromBlock
	sampleConfig.Scan.BlockScanConfig.ToBlock = DefaultToBlock
	sampleConfig.Scan.BlockScanConfig.OutputFileName = "full_blocks.json"
	sampleConfig.Scan.BlockScanConfig.BatchSize = DefaultBatchSize

	sampleConfig.Scan.AccountScanConfig.BlockNumber = DefaultAccountScanBlockNumber
	sampleConfig.Scan.AccountScanConfig.MaxAccounts = 100
	sampleConfig.Scan.AccountScanConfig.StartKey = "AAAAAA==" // Base64 encoded string 0x0
	sampleConfig.Scan.AccountScanConfig.OutputFileName = "accounts.json"
	sampleConfig.Scan.AccountScanConfig.BatchSize = DefaultBatchSize

	sampleConfig.Scan.ReceiptScanConfig.FullBlocksFile = "full_blocks.json"
	sampleConfig.Scan.ReceiptScanConfig.OutputFileName = "receipts.json"
	sampleConfig.Scan.ReceiptScanConfig.BatchSize = DefaultBatchSize

	sampleConfig.Scan.ContractCodeScanConfig.OutputFileName = "contract_codes.json"
	sampleConfig.Scan.ContractCodeScanConfig.BatchSize = 1

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
