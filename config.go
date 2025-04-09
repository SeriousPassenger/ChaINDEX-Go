package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type RpcConfig struct {
	Url     string `toml:"url"`     // URL for the RPC server
	Delay   int    `toml:"delay"`   // Delay between requests in milliseconds
	Timeout int    `toml:"timeout"` // Timeout for requests in seconds
	Retry   int    `toml:"retry"`   // Number of retries for failed request
}

const (
	SampleConfigPath = "config.toml.sample"

	// Default values for the RPC configuration
	DefaultRpcUrl     = "http://localhost:8080" // Default URL for the RPC server
	DefaultRpcDelay   = 1000                    // Default delay between requests in milliseconds
	DefaultRpcTimeout = 5                       // Default timeout for requests in seconds
	DefaultRpcRetry   = 3                       // Default number of retries for failed requests
)

type Config struct {
	Rpc RpcConfig `toml:"rpc"` // RPC configuration
}

func CreateSampleConfig() error {
	sampleConfig := new(Config)

	sampleConfig.Rpc.Url = DefaultRpcUrl
	sampleConfig.Rpc.Delay = DefaultRpcDelay
	sampleConfig.Rpc.Timeout = DefaultRpcTimeout
	sampleConfig.Rpc.Retry = DefaultRpcRetry

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
