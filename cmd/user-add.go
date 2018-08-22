package cmd

import (
	"fmt"

	"github.com/bluemir/go-utils/auth"
	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
)

func NewUserAddCommand() *cobra.Command {
	var UserAddConfig = struct {
		email    string
		password string
	}{}
	var UserAddCmd = &cobra.Command{
		Use: "add USERNAME",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}
			if len(args) != 1 {
				return fmt.Errorf("invaild name")
			}
			/*
				email := args[0] + "@wikinote"
				//if not set
				if cmd.Flags().Changed("port") {
					email = UserAddConfig.email
				}
			*/
			return b.Auth().CreateUser(&auth.User{
				Name: args[0],
			})
			/*return b.User().New(&backend.User{
				Name:  args[0],
				Email: email,
			}, "")*/
		},
	}

	UserAddCmd.Flags().StringVar(&UserAddConfig.email, "email", "USERNAME@wikinote", "user email")

	return UserAddCmd
}
