package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/cgroups"
	"github.com/OpenKikCoc/mydocker/cgroups/subsystems"
	"github.com/OpenKikCoc/mydocker/container"
)

var (
	usetty   bool
	memory   string
	cpushare string
	cpuset   string
)

func init() {
	runCmdFlags(RunCmd)
}

func runCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&usetty, "ti", "", false, "enable tty")
	cmd.PersistentFlags().StringVarP(&memory, "mem", "m", "", "memory limit")
	cmd.PersistentFlags().StringVarP(&cpushare, "cpushare", "", "", "cpushare limit")
	cmd.PersistentFlags().StringVarP(&memory, "cpuset", "", "", "cpuset limit")
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
			fmt.Println("Hello, World! args:", args)
			resConf := &subsystems.ResourceConfig{
				MemoryLimit: memory,
				CpuSet:      cpuset,
				CpuShare:    cpushare,
			}
			Run(usetty, args, resConf)
		},
	}
)

func Run(usetty bool, args []string, conf *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(usetty)
	if parent == nil {
		log.Println("New parent error")
	}
	// parent.Run will wait
	if err := parent.Start(); err != nil {
		log.Println("Run parent.Start()", err)
	}
	log.Println("Run parent.Start() has finished")

	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(conf)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(args, writePipe)
	parent.Wait()
}

func sendInitCommand(args []string, writePipe *os.File) {
	command := strings.Join(args, " ")
	log.Printf("sendInitCommand is %s\n", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
