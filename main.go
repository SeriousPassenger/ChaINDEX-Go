package main

func main() {
	// -----------------------------------------
	// Entrypoint for the CLI
	// -----------------------------------------

	rootCmd := GetCobraRootCmd()
	rootCmd.Execute()
}
