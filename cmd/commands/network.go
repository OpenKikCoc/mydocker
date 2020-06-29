package commands

import (
	"github.com/spf13/cobra"
)

var (
	driver string
	subnet string
)

func init() {
	networkCmdFlags(RootCmd)
}

func networkCmdFlags(cmd *cobra.Command) {
	//cmd.PersistentFlags().StringVarP(&driver, "driver", "", "", "network driver")
	//cmd.PersistentFlags().StringVarP(&subnet, "subnet", "", "", "subnet cidr")
}

var (
	// NetworkCmd the root command
	NetworkCmd = &cobra.Command{
		Use:   "network",
		Short: "Network Command",
		/*
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println("network create")
				network.Init()
				err := network.CreateNetwork(context.String("driver"), context.String("subnet"), context.Args()[0])
				if err != nil {
					fmt.Errorf("create network error: %+v\n", err)
					panic(err)
				}
			},
		*/
	}
)
