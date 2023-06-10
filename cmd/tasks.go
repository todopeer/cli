package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var taskCmd = &cobra.Command{
	Use:     "task",
	Aliases: []string{"t"},
	Short:   "(t) show tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		tasks, err := api.QueryTasks(ctx, token, nil)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			outputTask(t)
		}
		return nil
	},
}

var dueDate string
var newTaskCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "(n) create new task",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var dueTime *graphql.String

		if dueDate != "" {
			// dueTimeData, err := time.Parse(time.DateOnly, dueDate)
			// if err != nil {
			// 	return err
			// }
			dueTime = (*graphql.String)(&dueDate)
		}
		ctx := context.Background()
		token := config.MustGetToken()
		var desc *string

		if len(args) == 0 {
			return errors.New("task title must be provided as first argument")
		}
		if len(args) > 1 {
			desc = &args[1]
		}

		task, err := api.CreateTask(ctx, token, api.TaskCreateInput{
			Name:        graphql.String(args[0]),
			Description: (*graphql.String)(desc),
			DueDate:     dueTime,
		})
		if err != nil {
			return err
		}

		outputTask(task)
		return nil
	},
}

var removeTaskCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "remove (rm) a task by its ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		if len(args) == 0 {
			log.Fatal("taskID must be provided")
		}

		taskID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		t, err := api.RemoveTask(ctx, token, api.ID(taskID))
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) removed successfully: %s\n", t.ID, t.Name)
		return err
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
			log.Fatal("taskID must be provided")
		}

		taskID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		t, err := api.StartTask(ctx, token, api.ID(taskID))
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

func init() {
	rootCmd.AddCommand(taskCmd)
	rootCmd.AddCommand(removeTaskCmd)
	rootCmd.AddCommand(startTaskCmd)

	newTaskCmd.Flags().StringVarP(&dueDate, "due", "d", "", "Due date for the task (format: 2006-01-02)")
	rootCmd.AddCommand(newTaskCmd)
}
