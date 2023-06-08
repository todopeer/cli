package cmd

import (
	"fmt"
	"log"

	"github.com/flyfy1/diarier_cli/api"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from your account",
	Run: func(cmd *cobra.Command, args []string) {
		err := api.Deauthenticate()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Logged out successfully!")
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}
