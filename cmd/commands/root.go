package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

const usage = `mydocker is a simple container runtime implementation.
		The purpose of this project is to learn how docker works and how to write a docker by ourselves
		Enjoy it, just for fun.`

func init() {
	rootCmdFlags(RootCmd)
}

func rootCmdFlags(cmd *cobra.Command) {
	//cmd.PersistentFlags().StringVar(&configFile, "configFile", "config.toml", "config file (default is ./config.toml)")
}

var (
	// RootCmd the root command
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "Root Command",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			fmt.Println("RootCmd")
			return nil
		},
	}
)
