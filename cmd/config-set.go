package cmd

import (
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
)

func NewConfigSetCommand() *cobra.Command {
	var ConfigSetCmd = &cobra.Command{
		Use: "set key=value",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}

			for _, v := range args {
				logrus.Debugf(v)
				arr := strings.SplitN(v, "=", 2)

				if err := b.Config().Set(arr[0], arr[1]); err != nil {
					return err
				}

				if err := b.SaveConfig(b.Config()); err != nil {
					return err
				}
			}

			return nil
		},
	}

	//UserAddCmd.Flags().StringVar(&UserAddConfig.email, "email", "USERNAME@wikinote", "user email")

	return ConfigSetCmd
}
