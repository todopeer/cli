package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var pauseTaskCmd = &cobra.Command{
	Use:     "pause",
	Aliases: []string{"p"},
	Short:   "pause(p) current running task/event",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		client := api.NewClient(token)

		// try getting the current running task
		user, err := client.Me()
		if err != nil {
			return err
		}

		if user.RunningTaskID == nil {
			return ErrNoRunningTaskNeedID
		}
		taskID := int64(*user.RunningTaskID)

		t, err := client.UpdateTask(api.ID(taskID), api.TaskUpdateInput{
			Status: &api.TaskStatusPaused,
		})
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) paused: %s\n", t.ID, t.Name)
		return err
	},
}
