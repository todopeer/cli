package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Shopify/hoff"
	"github.com/shurcooL/graphql"
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
		ctx := context.Background()

		input := api.QueryTaskInput{}
		var err error
		input.Status, err = hoff.MapError(statusForQuery, taskStatusShortToInput)
		if err != nil {
			return err
		}
		log.Printf("loading status: %v", input.Status)

		tasks, err := api.QueryTasks(ctx, token, input)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			t.Output()
		}
		return nil
	},
}

var startTaskCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s"},
	Short:   "(s) start task",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		ctx := context.Background()

		var taskID api.ID
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
			}
		}

		t, err := api.StartTask(ctx, token, taskID)
		if err != nil {
			return fmt.Errorf("start task error: %w", err)
		}
		fmt.Printf("task(id=%d) started successfully: %s\n", t.ID, t.Name)
		return err
	},
}

var (
	statusForQuery            []string
	mapStatusShort2TaskStatus = map[string]api.TaskStatus{
		"n": api.TaskStatusNotStarted,
		"i": api.TaskStatusDoing,
		"d": api.TaskStatusDone,
		"p": api.TaskStatusPaused,
	}
)

func taskStatusShortToInput(statusShort string, _ int) (api.TaskStatus, error) {
	r, found := mapStatusShort2TaskStatus[statusShort]
	if found {
		return r, nil
	}
	return "", errors.New("unknown status short: " + statusShort)
}

func init() {
	listTaskCmd.Flags().StringArrayVar(&statusForQuery, "status", []string{"n", "i", "p"}, "n: not_started; i: in-progress; d: done; p: paused")

	rootCmd.AddCommand(listTaskCmd)
	rootCmd.AddCommand(startTaskCmd)
}
