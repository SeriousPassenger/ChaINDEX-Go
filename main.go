package main

import (
	"github.com/spf13/cobra"
)

func main() {
	// chaindex create-config: create a sample config file (config.toml.sample)

	var rootCmd = &cobra.Command{}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "create-config",
		Short: "Create a sample config file",
		Run: func(cmd *cobra.Command, args []string) {
			if err := createSampleConfig(); err != nil {
				cmd.Println("Error creating sample config:", err)
			}
		},
	})

	if err := rootCmd.Execute(); err != nil {
		rootCmd.Println("Error executing command:", err)
	} else {
		rootCmd.Println("Command executed successfully.")
	}
}
