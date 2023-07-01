package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var showTaskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t"},
	Short:   "(t) [id] show task. If not provided, show current running task",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var taskID api.ID
		if len(args) == 0 {
			e, err := client.QueryRunningEvent()
			if err != nil {
				return fmt.Errorf("error query running event: %w", err)
			}
			taskID = e.TaskID
		} else {
			taskIDInt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			taskID = api.ID(taskIDInt)
		}

		task, events, err := client.GetTaskEvents(taskID)
		if err != nil {
			return fmt.Errorf("error getting task: %w", err)
		}
		task.Output()
		ef := api.EventFormatter{Prefix: "\t", WithDate: true}

		for _, e := range events {
			fmt.Printf("\t")
			ef.Output(&e)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(showTaskCmd)
}
