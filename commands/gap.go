package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
	"github.com/todopeer/cli/util/gql"
)

func init() {
	rootCmd.AddCommand(gapEventCmd)
}

var gapEventCmd = &cobra.Command{
	Use:     "gap",
	Aliases: []string{"g"},
	Short:   "add gap (g) to the current running event. It would stop the event, update the endtime to minus the hole size, then resume this same task with a new event",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		if len(args) == 0 {
			return errors.New("Please give the gap size as duration")
		}

		// figure out the duration
		duration, err := time.ParseDuration(args[0])
		if err != nil {
			return fmt.Errorf("duration parse error: %w", err)
		}

		event, err := client.QueryRunningEvent()
		if err != nil {
			return err
		}

		if event == nil {
			return ErrNoRunningEvent
		}

		endTime := (api.Time)(time.Now().Add(-duration))

		input := api.EventUpdateInput{EndAt: &endTime}
		if len(args) > 1 {
			input.Description = gql.ToGqlStringP(args[1])
		}

		_, err = client.UpdateEvent(event.ID, input)
		if err != nil {
			return err
		}

		task, _, err := client.StartTask(event.TaskID)
		if err != nil {
			return err
		}
		fmt.Printf("hole added: %s; started new event for task: %s\n", duration, task.Name)

		return nil
	},
}
