package main

import "log"

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

	_, err = GetRpcClient(config)

	if err != nil {
		log.Fatalf("Error creating RPC client: %v", err)
	}
}
