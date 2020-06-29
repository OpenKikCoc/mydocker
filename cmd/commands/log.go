package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/container"
	//"github.com/OpenKikCoc/mydocker/cgroups"
)

func init() {
	logCmdFlags(RunCmd)
}

func logCmdFlags(cmd *cobra.Command) {
}

var (
	// LogCmd the root command
	LogCmd = &cobra.Command{
		Use:   "log",
		Short: "Log Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one arg")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			logContainer(args[0])
		},
	}
)

func logContainer(containerName string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	logFileLocation := dirURL + container.ContainerLogFile
	file, err := os.Open(logFileLocation)
	defer file.Close()
	if err != nil {
		log.Printf("Log container open file %s error %v\n", logFileLocation, err)
		return
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Log container read file %s error %v\n", logFileLocation, err)
		return
	}
	fmt.Fprint(os.Stdout, string(content))
}
