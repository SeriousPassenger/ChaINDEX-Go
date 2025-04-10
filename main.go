package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
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

	// ----------------------------------------
	// Flags for all commands
	// ----------------------------------------
	testConnectionCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to the config file")

	// ----------------------------------------
	// Add commands to root command
	// ----------------------------------------
	rootCmd.AddCommand(createConfigCmd)
	rootCmd.AddCommand(testConnectionCmd)

	// -----------------------------------------
	// Execute the root command
	// -----------------------------------------
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing command:", err)
	}
}
