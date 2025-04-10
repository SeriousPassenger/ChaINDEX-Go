package main

import (
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "chaindex",
		Short: "ChaINDEX CLI",
	}

	// ------------------------------------------------------
	// create-config command
	// ------------------------------------------------------
	createConfigCmd := &cobra.Command{
		Use:   "create-config",
		Short: "Create a sample config file",
		Run: func(cmd *cobra.Command, args []string) {
			if err := CreateSampleConfig(); err != nil {
				cmd.Println("Error creating sample config:", err)
				return
			}
			fmt.Println("Sample config file created successfully.")
		},
	}

	// ------------------------------------------------------
	// test-connection command
	// ------------------------------------------------------
	var configFile string
	testConnectionCmd := &cobra.Command{
		Use:   "test-connection",
		Short: "Test the connection to the RPC server",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}

			// Now load and test the connection
			config, err := GetConfig(configFile)
			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}

			if err := TestConnection(config); err != nil {
				cmd.Println("Connection test failed:", err)
			} else {
				cmd.Println("Connection test succeeded!")
			}
		},
	}
	// Define the --config (or -c) flag
	testConnectionCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")

	testGetBlocksBatchCmd := &cobra.Command{
		Use:   "test-get-blocks-batch",
		Short: "Test the get blocks batch functionality",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}

			// Now load and test the connection
			config, err := GetConfig(configFile)
			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}

			client, err := GetClient(config)

			if err != nil {
				cmd.Println("Error getting client:", err)
				return
			}

			// Example block numbers to fetch
			blockNumbers := []*big.Int{
				big.NewInt(1),
			}

			blocks, err := GetBlocksBatch(client, blockNumbers)

			if err != nil {
				cmd.Println("Error fetching blocks:", err)
				return
			}

			err = SaveStructToJSONFile(blocks, "blocks_debug.json")

			if err != nil {
				cmd.Println("Error saving blocks to file:", err)
				return
			}

			fmt.Println("Blocks fetched and saved successfully.")
		},
	}

	testGetBlocksBatchCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")

	testGetReceiptsBatchCmd := &cobra.Command{
		Use:   "test-get-receipts-batch",
		Short: "Test the get receipts batch functionality",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}

			// Now load and test the connection
			config, err := GetConfig(configFile)
			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}

			client, err := GetClient(config)

			if err != nil {
				cmd.Println("Error getting client:", err)
				return
			}

			// Example transaction hashes to fetch
			transactions := []string{
				"0x9c223f125a698edadd81665082f4de89a20d44ad267e83cf2210a28225a5c89a",
			}

			receipts, err := GetReceiptsBatch(client, transactions)

			if err != nil {
				cmd.Println("Error fetching receipts:", err)
				return
			}

			err = SaveStructToJSONFile(receipts, "receipts_debug.json")
			if err != nil {
				cmd.Println("Error saving receipts to file:", err)
				return
			}

			fmt.Println("Receipts fetched and saved successfully.")
		},
	}
	testGetReceiptsBatchCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")

	// Add subcommands to root
	rootCmd.AddCommand(createConfigCmd)
	rootCmd.AddCommand(testConnectionCmd)
	rootCmd.AddCommand(testGetBlocksBatchCmd)
	rootCmd.AddCommand(testGetReceiptsBatchCmd)

	// Execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}
