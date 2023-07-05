package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
	"github.com/todopeer/cli/util/gql"
)

func init() {
	rootCmd.AddCommand(doneTaskCmd)
}

func parseTaskIDAndDesc(args []string) (int64, string, error) {
	if len(args) == 0 {
		return 0, "", nil
	}

	taskIDInt, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		if len(args) == 1 {
			return 0, args[0], nil
		} else {
			return 0, "", fmt.Errorf("expect arg1 to be id if single argument. Error: %s", err)
		}
	}

	var desc string
	if len(args) == 2 {
		desc = args[1]
	}
	return taskIDInt, desc, nil
}

var doneTaskCmd = &cobra.Command{
	Use:     "done",
	Aliases: []string{"d"},
	Short:   "done(d) [taskid] [desc] - mark task as done, with optional desc",
	Long:    `If taskid not provided, use current running task.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var taskID api.ID

		taskIDInt, desc, err := parseTaskIDAndDesc(args)
		if err != nil {
			return err
		}

		if taskIDInt == 0 {
			// try getting the current running task
			user, err := client.Me()
			if err != nil {
				return err
			}

			if user.RunningTaskID == nil {
				return ErrNoRunningTaskNeedID
			}
			taskID = *user.RunningTaskID
		} else {
			taskID = api.ID(taskIDInt)
		}

		t, err := client.UpdateTask(taskID, api.TaskUpdateInput{
			Status: &api.TaskStatusDone,
		})
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) successfully done: %s\n", t.ID, t.Name)

		if len(desc) > 0 {
			// use the 2nd arg as input to update the last event attached to this task
			event, err := client.QueryTaskLastEvent(taskID)
			if err != nil {
				return err
			}

			if event == nil {
				return fmt.Errorf("cannot find event for task(id=%d)", taskID)
			}

			event, err = client.UpdateEvent(event.ID, api.EventUpdateInput{Description: gql.ToGqlStringP(desc)})
			if err != nil {
				return fmt.Errorf("update event(id=%d) error: %w", event.ID, err)
			}
			fmt.Printf("event(id=%d) updated desc: %s\n", event.ID, *event.Description)
		}

		return err
	},
}
