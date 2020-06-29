package container

import (
	"log"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {
	log.Printf("RunContainerInitProcess %s, %s\n", command, args)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		log.Println(err.Error())
	}
	return nil
}
