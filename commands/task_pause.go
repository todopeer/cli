package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var pauseTaskCmd = &cobra.Command{
	Use:     "pause",
	Aliases: []string{"p"},
	Short:   "pause(p) current running task/event. If an ID is provided, pause that task instead",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var taskID api.ID
		if len(args) > 0 {
			taskIDInt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("unable to parse taskID: %s, error: %w", args[0], err)
			}
			taskID = api.ID(taskIDInt)
		} else {
			// try getting the current running task
			user, err := client.Me()
			if err != nil {
				return err
			}

			if user.RunningTaskID == nil {
				return ErrNoRunningTaskNeedID
			}

			taskID = *user.RunningTaskID
		}

		t, err := client.UpdateTask(taskID, api.TaskUpdateInput{
			Status: &api.TaskStatusPaused,
		})
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) paused: %s\n", t.ID, t.Name)
		return err
	},
}
