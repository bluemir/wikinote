package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bluemir/go-utils/auth"
	"github.com/bluemir/wikinote/backend"
)

func NewUserAssignCommand() *cobra.Command {
	var UserAssignCmd = &cobra.Command{
		Use: "assign USER to ROLE",
		Args: func(cmd *cobra.Command, args []string) error {
			// argument validataion
			if len(args) != 3 {
				return fmt.Errorf("invaild arguments")
			}
			if !strings.EqualFold(args[1], "to") {
				return fmt.Errorf("invaild arguements")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			role := args[2]

			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}

			u, ok, err := b.Auth().GetUser(username)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("user not found")
			}

			u.Role = auth.Role(role)

			return b.Auth().UpdateUser(u)

		},
	}
	// vaildate
	return UserAssignCmd
}
