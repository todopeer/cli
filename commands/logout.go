package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		err := client.Logout()
		if err != nil {
			return err
		}

		config.UpdateToken("")
		fmt.Println("Logged out successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
