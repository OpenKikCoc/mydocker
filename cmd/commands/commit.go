package commands

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
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
			if len(args) < 1 {
				return errors.New("requires at least one arg")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CommitCmd")
			fmt.Println("imageName = ", args[0])
			commitContainer(args[0])
		},
	}
)

func commitContainer(imageName string) {
	mntURL := "/root/mnt"
	imageTar := "/root/" + imageName + ".tar"
	fmt.Printf("%s", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		log.Printf("Tar folder %s error %v\n", mntURL, err)
	}
}
