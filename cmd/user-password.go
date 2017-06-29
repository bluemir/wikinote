package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
	usertype "github.com/bluemir/wikinote/backend/user"
)

func NewUserPasswordCommand() *cobra.Command {
	UserPasswordConfig := struct {
		Password string
	}{}
	UserPasswordCmd := &cobra.Command{
		Use: "password NAME",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) < 1 {
				return fmt.Errorf("Must put username")
			}

			name := args[0]
			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}

			u, err := b.User().Get(name)
			if err != nil {
				return err
			}

			// if cmd.Flags().Changed("password")
			u.Password = usertype.NewPassword(UserPasswordConfig.Password)
			// else ask password
			return b.User().Put(u)
		},
	}
	UserPasswordCmd.Flags().StringVarP(&UserPasswordConfig.Password, "password", "p", "", "Password. WANNING!")

	return UserPasswordCmd
}
