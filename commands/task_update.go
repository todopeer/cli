package commands

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

func defineFlagsForTaskCUD(s *pflag.FlagSet, isUpdate bool) {
	s.StringVarP(&varDueDate, "due", "D", "", "due date for the task (format: 2006-01-02)")
	s.StringVarP(&varDescription, "desc", "d", "", "task description")

	if isUpdate {
		s.BoolVarP(&varTriggerPause, "pause", "p", false, "if set, mark as pause")
		s.StringVarP(&varName, "name", "n", "", "task name")
	}
}

func init() {
	rootCmd.AddCommand(removeTaskCmd)
	rootCmd.AddCommand(pauseTaskCmd)

	defineFlagsForTaskCUD(newTaskCmd.Flags(), false)
	rootCmd.AddCommand(newTaskCmd)

	defineFlagsForTaskCUD(updateTaskCmd.Flags(), true)
	rootCmd.AddCommand(updateTaskCmd)
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

var pauseTaskCmd = &cobra.Command{
	Use:   "pause",
	Aliases: []string{"p"},
	Short: "pause(p) current running task/event",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		ctx := context.Background()

		// try getting the current running task
		user, err := api.Me(ctx, token)
		if err != nil {
			return err
		}

		if user.RunningTaskID == nil {
			return errors.New("no running task. taskID must be provided")
		}
		taskID := int64(*user.RunningTaskID)

		t, err := api.UpdateTask(ctx, token, api.ID(taskID), api.TaskUpdateInput{
			Status: &api.TaskStatusPaused,
		})
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) paused: %s\n", t.ID, t.Name)
		return err
	},
}

var (
	varName         string
	varDescription  string
	varDueDate      string
	varTriggerPause bool
)

var updateTaskCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "(u)update [taskid] -flags",
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

		input := api.TaskUpdateInput{}
		if varDueDate != "" {
			input.DueDate = graphql.NewString(graphql.String(varDueDate))
		}

		if varName != "" {
			input.Name = graphql.NewString(graphql.String(varName))
		}

		if varDescription != "" {
			input.Description = graphql.NewString(graphql.String(varDescription))
		}

		if varTriggerPause {
			input.Status = &api.TaskStatusNotStarted
		}

		t, err := api.UpdateTask(ctx, token, api.ID(taskID), input)
		if err != nil {
			return err
		}
		fmt.Printf("task(id=%d) successfully updated: %s\n", t.ID, t.Name)
		return err
	},
}

var newTaskCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Short:   "(a) add new task",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var dueTime *graphql.String

		if varDueDate != "" {
			dueTime = (*graphql.String)(&varDueDate)
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

		task.Output()
		return nil
	},
}
