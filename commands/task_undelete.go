package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var undeleteTaskCmd = &cobra.Command{
	Use:     "undelete",
	Aliases: []string{"ud"},
	Short:   "undelete (ud) a task by its ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var taskID api.ID
		if len(args) == 0 {
			return errors.New("taskID must be provided")
		}

		taskIDInt, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return fmt.Errorf("taskID parse err: %w", err)
		}
		taskID = api.ID(taskIDInt)

		t, err := client.UndeleteTask(taskID)
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) undeleted successfully: %s\n", t.ID, t.Name)
		return err
	},
}

func init() {
	rootCmd.AddCommand(undeleteTaskCmd)
}
