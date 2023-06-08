package cmd

import (
	"context"
	"fmt"

	"github.com/flyfy1/diarier_cli/api"
	"github.com/flyfy1/diarier_cli/services/config"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		err := api.Logout(ctx, token)
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
