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
			outputTask(t)
		}
		return nil
	},
}

var startTaskCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s"},
	Short:   "(s) start task",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		if len(args) == 0 {
			log.Fatal("taskID, or new task content, must be provided")
		}

		// if cannot parse as int, then it's task content
		taskID, err := strconv.ParseInt(args[0], 10, 64)
		idTypedTaskID := api.ID(taskID)
		if err != nil {
			taskContent := args[0]

			createdTask, err := api.CreateTask(ctx, token, api.TaskCreateInput{
				Name: graphql.String(taskContent),
			})
			if err != nil {
				return err
			}
			idTypedTaskID = createdTask.ID
			log.Printf("created task with ID: %d", idTypedTaskID)
		}

		t, err := api.StartTask(ctx, token, idTypedTaskID)
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) started successfully: %s\n", t.ID, t.Name)
		return err
	},
}

func outputTask(t *api.Task) {
	fmt.Printf("%d\t%s\t%s\t", t.ID, t.Status, t.Name)
	if t.DueDate != nil {
		fmt.Printf("%s\n", *t.DueDate)
	} else {
		fmt.Println()
	}
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
