package commands

import (
	"errors"

	"github.com/spf13/cobra"
	//"github.com/OpenKikCoc/mydocker/cgroups"
)

func init() {
	removeCmdFlags(RootCmd)
}

func removeCmdFlags(cmd *cobra.Command) {
	//cmd.PersistentFlags().StringVar(&configFile, "configFile", "config.toml", "config file (default is ./config.toml)")
}

var (
	// RemoveCmd the root command
	RemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing container name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			removeContainer(args[0]) // in stop.go
		},
	}
)
