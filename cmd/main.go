package main

import (
	"os"

	cmd "github.com/OpenKikCoc/mydocker/cmd/commands"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.RunCmd,
		cmd.InitCmd,
		// cmd.VersionCmd,
	)
	if rootCmd.Execute() != nil {
		os.Exit(1)
	}
}
