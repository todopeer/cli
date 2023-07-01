package commands

import (
	"errors"
	"log"

	"github.com/Shopify/hoff"
	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var listTaskCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "(l) list tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		input := api.QueryTaskInput{}
		var err error
		input.Status, err = hoff.MapError(statusForQuery, taskStatusShortToInput)
		if err != nil {
			return err
		}
		log.Printf("loading status: %v", input.Status)

		tasks, err := client.QueryTasks(input)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			t.Output()
		}
		return nil
	},
}

func taskStatusShortToInput(statusShort string, _ int) (api.TaskStatus, error) {
	r, found := mapStatusShort2TaskStatus[statusShort]
	if found {
		return r, nil
	}
	return "", errors.New("unknown status short: " + statusShort)
}

func init() {
	listTaskCmd.Flags().StringArrayVar(&statusForQuery, "status", []string{"n", "i", "p"}, "n: not_started; i: doing; d: done; p: paused")
	rootCmd.AddCommand(listTaskCmd)
}
