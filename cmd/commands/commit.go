package commands

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/OpenKikCoc/mydocker/container"
	//"github.com/OpenKikCoc/mydocker/cgroups"
)

func init() {
	commitCmdFlags(RootCmd)
}

func commitCmdFlags(cmd *cobra.Command) {
	//cmd.PersistentFlags().StringVar(&configFile, "configFile", "config.toml", "config file (default is ./config.toml)")
}

var (
	// CommitCmd the root command
	CommitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Commit Command",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("Missing container name and image name")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CommitCmd")
			commitContainer(args[0], args[1])
		},
	}
)

func commitContainer(containerName, imageName string) {
	mntURL := fmt.Sprintf(container.MntUrl, containerName)
	mntURL += "/"

	imageTar := container.RootUrl + "/" + imageName + ".tar"

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Printf("Tar folder %s error %v\n", mntURL, err)
	}
}
