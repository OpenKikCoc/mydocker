package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/OpenKikCoc/mydocker/container"
)

func main() {
	test3()
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
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{"/bin/sh"}
	if err := syscall.Exec("/bin/sh", argv, os.Environ()); err != nil {
		log.Println(err.Error())
	}
}

func test3() {
	//Run
	parent := container.NewParentProcess(true, "/bin/sh")
	log.Printf("Before Run parent: %+v\n", parent)
	if err := parent.Run(); err != nil {
		log.Println("Run parent.Start()", err)
	}
	log.Printf("After Run parent: %+v\n", parent)
	log.Println("Run parent.Start() has finished")
	parent.Wait()
	os.Exit(-1)
}