package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
	//"github.com/OpenKikCoc/mydocker/container"
)

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
		Run: func(cmd *cobra.Command, args []string) {
			v := args[0]
			log.Println("start: ", v)
			switch v {
			case "1":
				test1()
			case "2":
				test2()
			case "3":
				test3()
			default:
				log.Printf("not support %v\n", v)
			}
		},
	}
)

func main() {
	if RootCmd.Execute() != nil {
		os.Exit(1)
	}
}

func test1() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(-1)
}

func test2() {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err := syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		log.Printf("syscall.Mount proc error: %v\n", err)
	}
	argv := []string{"/bin/sh"}
	if err := syscall.Exec("/bin/sh", argv, os.Environ()); err != nil {
		log.Println(err.Error())
	}
}

func NewParentProcess(tty bool, command string) *exec.Cmd {
	log.Println("NewParentProcess: command", command)
	args := []string{"2"}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      0,
				Size:        1,
			},
		},
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}

func test3() {
	//Run
	parent := NewParentProcess(true, "/bin/sh")
	log.Printf("Before Run parent: %+v\n", parent)
	if err := parent.Start(); err != nil {
		log.Println("Run parent.Start()", err)
	}
	log.Printf("After Run parent: %+v\n", parent)
	log.Println("Run parent.Start() has finished")
	parent.Wait()
	os.Exit(-1)
}
