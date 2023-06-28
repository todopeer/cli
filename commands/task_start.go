package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

func init() {
	startTaskCmd.Flags().StringVarP(&varDurationOffset, "offset", "o", "", "if provided, start task with offset")
	startTaskCmd.Flags().BoolVarP(&varPomodoro, "pomodoro", "p", false, "if provided, start task with pomodoro")
	rootCmd.AddCommand(startTaskCmd)
}

var startTaskCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s"},
	Short:   "(s) start task",
	Long: `Syntax Supported:
start: to start the previously running task
start [task name]: to start a new task with given name
start [taskID]: to start a new task with given ID
start [taskID] [Description]: to start a task with given ID, add description to the event

Examples:
start "math homework" -p: to start "math homework" task in pomodoro mode
`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		ctx := context.Background()

		var taskID api.ID
		var startTaskOptions []api.StartTaskOptionFunc
		if len(args) == 0 {
			// try getting the previously running event
			e, err := api.QueryLatestEvents(ctx, token)
			if err != nil {
				return fmt.Errorf("query event error: %w", err)
			}
			if e == nil {
				return errors.New("taskID not provided, no event run in recent 2 days")
			}
			taskID = e.TaskID
		} else {
			// if cannot parse as int, then it's task content
			taskIDInt, err := strconv.ParseInt(args[0], 10, 64)
			taskID = api.ID(taskIDInt)
			if err != nil {
				taskContent := args[0]

				createdTask, err := api.CreateTask(ctx, token, api.TaskCreateInput{
					Name: graphql.String(taskContent),
				})
				if err != nil {
					return err
				}
				taskID = createdTask.ID
				log.Printf("created task with ID: %d", taskID)
			} else {
				// if got more string, use it as input to the event desc
				if len(args) > 1 {
					startTaskOptions = append(startTaskOptions, api.StartTaskWithDescription(args[1]))
				}
			}
		}
		if len(varDurationOffset) > 0 {
			offset, err := time.ParseDuration(varDurationOffset)
			if err != nil {
				return fmt.Errorf("error parsing offset: %w", err)
			}
			startTaskOptions = append(startTaskOptions, api.StartTaskWithOffset(offset))
		}

		t, evt, err := api.StartTask(ctx, token, taskID, startTaskOptions...)
		if err != nil {
			return fmt.Errorf("start task error: %w", err)
		}
		fmt.Printf("task(id=%d) started successfully: %s\n", t.ID, t.Name)
		if evt != nil {
			fmt.Printf("\tevent(id=%d) started successfully at: %s\n", evt.ID, evt.StartAt.EventTimeOnly())
		}

		if varPomodoro {
			err = pomodoro(25*time.Minute, 0, makeTaskPauseCallback(token, t.ID))

			if err != nil {
				return err
			}
			fmt.Println("task paused: ", t.Name)
		}
		return nil
	},
}
