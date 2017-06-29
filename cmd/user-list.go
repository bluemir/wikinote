package cmd

import (
	"fmt"

	"github.com/apcera/termtables"
	"github.com/spf13/cobra"

	"github.com/bluemir/wikinote/backend"
)

func NewUserListCommand() *cobra.Command {
	var UserListCmd = &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := backend.New(backendOpts)
			if err != nil {
				return err
			}
			users, err := b.User().List()
			if err != nil {
				return err
			}

			table := termtables.CreateTable()
			table.AddHeaders("Name", "Email", "Role")
			for _, u := range users {
				table.AddRow(u.Id, u.Email, u.Role)
			}

			fmt.Println(table.Render())
			return nil
		},
	}
	return UserListCmd
}
