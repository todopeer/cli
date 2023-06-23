package commands

import(
	"fmt"
	"strconv"
	"context"
	"errors"

	"github.com/todopeer/cli/api"
	"github.com/spf13/cobra"
	"github.com/todopeer/cli/services/config"
)

var removeTaskCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "remove (rm) a task by its ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		var taskID api.ID
		if len(args) == 0 {
			fmt.Println("taskID not provided, would remove the current running task")

			evt, err := api.QueryRunningEvent(ctx, token)
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


		t, err := api.RemoveTask(ctx, token, taskID)
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) removed successfully: %s\n", t.ID, t.Name)
		return err
	},
}
