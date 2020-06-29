package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/container"
)

func init() {
	initCmdFlags(RunCmd)
}

func initCmdFlags(cmd *cobra.Command) {
	//cmd.PersistentFlags().BoolVarP(&usetty, "ti", "", false, "enable tty")
}

var (
	// InitCmd the root command
	InitCmd = &cobra.Command{
		Use:   "init",
		Short: "Init Command, Don't call it outside.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("InitCmd: args :", args)
			err := container.RunContainerInitProcess()
			// return err
			panic(err)
		},
	}
)
