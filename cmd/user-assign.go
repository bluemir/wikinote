package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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

			u, err := b.User().Get(username)
			if err != nil {
				return err
			}

			u.Role = role
			return b.User().Put(u)
		},
	}
	// vaildate
	return UserAssignCmd
}
