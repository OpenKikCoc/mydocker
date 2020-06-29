package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/container"
	//"github.com/OpenKikCoc/mydocker/cgroups"
)

func init() {
	stopCmdFlags(RootCmd)
}

func stopCmdFlags(cmd *cobra.Command) {
	//cmd.PersistentFlags().StringVar(&configFile, "configFile", "config.toml", "config file (default is ./config.toml)")
}

var (
	// StopCmd the root command
	StopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing container name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CommitCmd")
			fmt.Println("imageName = ", args[0])
			stopContainer(args[0])
		},
	}
)

func stopContainer(containerName string) {
	pid, err := getContainerPidByName(containerName)
	if err != nil {
		log.Printf("Get contaienr pid by name %s error %v\n", containerName, err)
		return
	}
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		log.Printf("Conver pid from string to int error %v\n", err)
		return
	}
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		log.Printf("Stop container %s error %v\n", containerName, err)
		return
	}
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Printf("Get container %s info error %v\n", containerName, err)
		return
	}
	containerInfo.Status = container.STOP
	containerInfo.Pid = " "
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Printf("Json marshal %s error %v\n", containerName, err)
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	if err := ioutil.WriteFile(configFilePath, newContentBytes, 0622); err != nil {
		log.Printf("Write file %s error %s\n", configFilePath, err)
	}
}

func getContainerInfoByName(containerName string) (*container.ContainerInfo, error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Printf("Read file %s error %v\n", configFilePath, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes, &containerInfo); err != nil {
		log.Printf("GetContainerInfoByName unmarshal error %v\n", err)
		return nil, err
	}
	return &containerInfo, nil
}

func removeContainer(containerName string) {
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Printf("Get container %s info error %v\n", containerName, err)
		return
	}
	if containerInfo.Status != container.STOP {
		log.Printf("Couldn't remove running container\n")
		return
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Printf("Remove file %s error %v\n", dirURL, err)
		return
	}
	container.DeleteWorkSpace(containerInfo.Volume, containerName)
}
