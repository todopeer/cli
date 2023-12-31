package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var meCmd = &cobra.Command{
	Use:   "my",
	Short: "show current user & task info",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		user, task, event, err := client.MeWithTaskEvent()
		if err != nil {
			return fmt.Errorf("error loading running task: %w", err)
		}

		if flagSimpleOutput {
			fmt.Println(task.Name)
		} else {
			fmt.Printf("%s - %s\n", user.Name, user.Email)
			if task != nil {
				fmt.Println("\tCurrent task: ")
				task.Output()

				if event != nil {
					fmt.Println("\tCurrent event: ")
					api.EventFormatter{}.Output(event)
				}
			} else {
				fmt.Println("no running task")
			}
		}
		return nil
	},
}

func init() {
	meCmd.Flags().BoolVarP(&flagSimpleOutput, "name-only", "N", false, "(N) if set, output task name only. Useful when pipeline with others")
	rootCmd.AddCommand(meCmd)
}
