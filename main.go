package main

import (
	"log"
)

func main() {
	// -----------------------------------------
	// Entrypoint for the CLI
	// -----------------------------------------

	// rootCmd := GetCobraRootCmd()
	// rootCmd.Execute()

	// -------------------------------------------
	// Debugging
	// -------------------------------------------

	CreateSampleConfig()

	config, err := GetConfig("config.toml")

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	err = ScanBlocksWithConfig(config)

	if err != nil {
		log.Fatalf("Error scanning blocks: %v", err)
	}

	testTxs := []string{
		"0xfa1b20dff53f445aff3aef41b29d0eb085867fb81495a193fb45faa6a3429952",
	}

	err = ScanReceiptsWithConfig(config, testTxs)

	if err != nil {
		log.Fatalf("Error scanning receipts: %v", err)
	}
}
