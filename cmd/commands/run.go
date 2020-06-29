package commands

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/container"
)

var (
	usetty bool
)

func init() {
	runCmdFlags(RunCmd)
}

func runCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&usetty, "ti", "", false, "enable tty")
}

var (
	// RunCmd the root command
	RunCmd = &cobra.Command{
		Use:   "run",
		Short: "Run Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one arg")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("RunCmd")
			fmt.Println("Hello, World!")
			Run(usetty, args)
		},
	}
)

func Run(usetty bool, args []string) {
	parent := container.NewParentProcess(usetty, args[0])
	// parent.Run will wait
	if err := parent.Start(); err != nil {
		log.Println("Run parent.Start()", err)
	}
	log.Println("Run parent.Start() has finished")
	parent.Wait()
	os.Exit(-1)
}
