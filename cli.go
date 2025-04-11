package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func GetCobraRootCmd() *cobra.Command {
	// ------------------------------------------------------
	// Root command
	// ------------------------------------------------------
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

	scanBlocksCmd := &cobra.Command{
		Use:   "scan-blocks",
		Short: "Scan blocks",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}
			config, err := GetConfig(configFile)
			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}
			if err := ScanBlocksWithConfig(config); err != nil {
				cmd.Println("Error scanning full blocks:", err)
				return
			}
			cmd.Println("Blocks scanned successfully.")
		},
	}

	// ------------------------------------------------------
	// scan-receipts command
	// ------------------------------------------------------
	var blockFile string

	scanReceiptsCmd := &cobra.Command{
		Use:   "scan-receipts",
		Short: "Scan receipts",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}
			if blockFile == "" {
				cmd.Println("Error: block file path is required (use --block-file).")
				return
			}
			config, err := GetConfig(configFile)
			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}

			if err := ScanReceiptsWithConfig(config, blockFile); err != nil {
				cmd.Println("Error scanning receipts:", err)
				return
			}
			cmd.Println("Receipts scanned successfully.")
		},
	}

	scanAccountsCmd := &cobra.Command{
		Use:   "scan-accounts",
		Short: "Scan accounts",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}

			config, err := GetConfig(configFile)
			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}

			if err := ScanAllAccounts(config); err != nil {
				cmd.Println("Error scanning accounts:", err)
				return
			}

			cmd.Println("Accounts scanned successfully.")
		},
	}

	var accountsFile string
	scanContractCodeCmd := &cobra.Command{
		Use:   "scan-contract-code",
		Short: "Scan contract code",
		Run: func(cmd *cobra.Command, args []string) {
			if configFile == "" {
				cmd.Println("Error: config file path is required (use --config).")
				return
			}

			if accountsFile == "" {
				cmd.Println("Error: accounts file path is required (use --accounts-file).")
				return
			}

			config, err := GetConfig(configFile)

			if err != nil {
				cmd.Println("Error loading config:", err)
				return
			}

			if err := ScanContractCode(config, accountsFile); err != nil {
				cmd.Println("Error scanning contract code:", err)
				return
			}

			cmd.Println("Contract code scanned successfully.")
		},
	}

	// ----------------------------------------
	// Flags for all commands
	// ----------------------------------------
	testConnectionCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
	scanBlocksCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
	scanReceiptsCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
	scanReceiptsCmd.Flags().StringVarP(&blockFile, "block-file", "b", "", "Path to the block file")
	scanAccountsCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
	scanContractCodeCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")
	scanContractCodeCmd.Flags().StringVarP(&accountsFile, "accounts-file", "a", "", "Path to the accounts file")

	// ----------------------------------------
	// Add commands to root command
	// ----------------------------------------
	rootCmd.AddCommand(createConfigCmd)
	rootCmd.AddCommand(testConnectionCmd)
	rootCmd.AddCommand(scanBlocksCmd)
	rootCmd.AddCommand(scanReceiptsCmd)
	rootCmd.AddCommand(scanAccountsCmd)
	rootCmd.AddCommand(scanContractCodeCmd)

	return rootCmd
}
