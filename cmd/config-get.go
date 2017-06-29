package cmd

import (
	"fmt"

	//"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
)

func NewConfigGetCommand() *cobra.Command {
	var ConfigGetCmd = &cobra.Command{
		Use: "get key",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}

			if len(args) != 1 {
				return fmt.Errorf("get key fail")
			}

			str, err := b.Config().Get(args[0])
			fmt.Printf("%s", str)
			return err
		},
	}

	//UserAddCmd.Flags().StringVar(&UserAddConfig.email, "email", "USERNAME@wikinote", "user email")

	return ConfigGetCmd
}
