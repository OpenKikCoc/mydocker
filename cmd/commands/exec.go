package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/container"
	//_ "github.com/OpenKikCoc/mydocker/nsenter" // macos 编译不支持
)

const ENV_EXEC_PID = "mydocker_pid"
const ENV_EXEC_CMD = "mydocker_cmd"

func init() {
	execCmdFlags(RunCmd)
}

func execCmdFlags(cmd *cobra.Command) {
}

var (
	// ExecCmd the root command
	ExecCmd = &cobra.Command{
		Use:   "exec",
		Short: "Exec Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Missing container name or command")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if os.Getenv(ENV_EXEC_PID) != "" {
				log.Printf("pid callback pid %s\n", os.Getgid())
			}
			execContainer(args)
		},
	}
)

func execContainer(args []string) {
	pid, err := getContainerPidByName(args[0])
	if err != nil {
		log.Printf("Exec container getContainerPidByName %s error %v\n", args[0], err)
		return
	}

	cmdStr := strings.Join(args[1:], " ")
	log.Printf("container pid %s\n", pid)
	log.Printf("command %s\n", cmdStr)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	os.Setenv(ENV_EXEC_PID, pid)
	os.Setenv(ENV_EXEC_CMD, cmdStr)
	containerEnvs := getEnvsByPid(pid)
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		log.Printf("Exec container %s error %v", args[0], err)
	}
}

func getContainerPidByName(containerName string) (string, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return "", err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		return "", err
	}
	return containerInfo.Pid, nil
}

func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Read file %s error %v", path, err)
		return nil
	}
	//env split by \u0000
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}
