package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var meCmd = &cobra.Command{
	Use:   "my",
	Short: "show current user & task info",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		token, err := config.ReadToken()
		if err != nil {
			return err
		}

		user, task, err := api.MeWithTask(ctx, token)
		if err != nil {
			return err
		}
		fmt.Printf("%s - %s\n", user.Email, user.Email)
		if task != nil {
			fmt.Println("\tCurrent task: ")
			task.Output()
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(meCmd)
}
