package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var deleteTaskCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"dt"},
	Short:   "delete (dt) a task by its ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var taskID api.ID
		if len(args) == 0 {
			fmt.Println("taskID not provided, would delete the current running task")

			evt, err := client.QueryRunningEvent()
			if err != nil {
				return err
			}

			if evt == nil {
				return errors.New("no current running event")
			}
			taskID = evt.TaskID
		} else {
			taskIDInt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			taskID = api.ID(taskIDInt)
		}

		t, err := client.DeleteTask(taskID)
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) deleted successfully: %s\n", t.ID, t.Name)
		return err
	},
}

func init() {
	rootCmd.AddCommand(deleteTaskCmd)
}
