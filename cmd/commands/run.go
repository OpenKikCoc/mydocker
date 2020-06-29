package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/cgroups"
	"github.com/OpenKikCoc/mydocker/cgroups/subsystems"
	"github.com/OpenKikCoc/mydocker/container"
)

var (
	usetty   bool
	daemon   bool
	memory   string
	cpushare string
	cpuset   string
	volume   string
	name     string
	net      string
	port     []string
	env      []string
)

func init() {
	runCmdFlags(RunCmd)
}

func runCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&usetty, "ti", "", false, "enable tty")
	cmd.PersistentFlags().BoolVarP(&daemon, "d", "d", false, "enable daemon")
	cmd.PersistentFlags().StringVarP(&memory, "mem", "m", "", "memory limit")
	cmd.PersistentFlags().StringVarP(&cpushare, "cpushare", "", "", "cpushare limit")
	cmd.PersistentFlags().StringVarP(&memory, "cpuset", "", "", "cpuset limit")
	cmd.PersistentFlags().StringVarP(&volume, "volume", "v", "", "cpuset limit")
	cmd.PersistentFlags().StringVarP(&name, "name", "", "", "container name")
	cmd.PersistentFlags().StringVarP(&net, "net", "", "", "container net")
	cmd.PersistentFlags().StringSliceVarP(&port, "port", "p", []string{}, "container port")
	cmd.PersistentFlags().StringSliceVarP(&env, "env", "e", []string{}, "container env slice")
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
			// TODO
			Run(usetty, args[1:], resConf, name, volume, args[0], env, net, port)
		},
	}
)

func Run(tty bool, comArray []string, conf *subsystems.ResourceConfig, containerName, volume, imageName string,
	envSlice []string, nw string, portmapping []string) {
	containerID := randStringBytes(10)
	if containerName == "" {
		containerName = containerID
	}
	parent, writePipe := container.NewParentProcess(tty, containerName, volume, imageName, envSlice)
	if parent == nil {
		log.Println("New parent error")
		return
	}
	// parent.Run will wait
	if err := parent.Start(); err != nil {
		log.Println("Run parent.Start()", err)
	}
	log.Println("Run parent.Start() has finished")
	//record container info
	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName, containerID, volume)
	if err != nil {
		log.Printf("Record container info error %v\n", err)
		return
	}
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(conf)
	cgroupManager.Apply(parent.Process.Pid)

	/*
		if nw != "" {
			// config container network
			network.Init()
			containerInfo := &container.ContainerInfo{
				Id:          containerID,
				Pid:         strconv.Itoa(parent.Process.Pid),
				Name:        containerName,
				PortMapping: portmapping,
			}
			if err := network.Connect(nw, containerInfo); err != nil {
				log.Errorf("Error Connect Network %v", err)
				return
			}
		}
	*/
	sendInitCommand(comArray, writePipe)

	if tty {
		parent.Wait()
		deleteContainerInfo(containerName)
		container.DeleteWorkSpace(volume, containerName)
	}
}

func sendInitCommand(args []string, writePipe *os.File) {
	command := strings.Join(args, " ")
	log.Printf("sendInitCommand is %s\n", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, containerName, id, volume string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
		Volume:      volume,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Printf("Record container info error %v\n", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		log.Printf("Mkdir error %s error %v\n", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		log.Printf("Create file %s error %v\n", fileName, err)
		return "", err
	}
	if _, err := file.WriteString(jsonStr); err != nil {
		log.Printf("File write string error %v\n", err)
		return "", err
	}

	return containerName, nil
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Printf("Remove dir %s error %v\n", dirURL, err)
	}
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
