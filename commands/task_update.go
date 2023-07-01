package commands

import (
	"errors"
	"fmt"
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
	rootCmd.AddCommand(pauseTaskCmd)

	defineFlagsForTaskCUD(newTaskCmd.Flags(), false)
	rootCmd.AddCommand(newTaskCmd)

	defineFlagsForTaskCUD(updateTaskCmd.Flags(), true)
	rootCmd.AddCommand(updateTaskCmd)
}

var updateTaskCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u"},
	Short:   "(u)update [taskid] -flags",
	Long: `Syntax Supported:
update [taskid]: to update the task with given ID
update: to update the current running task
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var taskID api.ID

		if len(args) == 0 {
			runningEvent, err := client.QueryRunningEvent()
			if err != nil {
				return fmt.Errorf("error querying running event: %w", err)
			}
			if runningEvent == nil {
				return errors.New("no running event")
			}
			taskID = runningEvent.TaskID
		} else {
			taskIDInt, err := strconv.ParseInt(args[0], 10, 64)

			if err != nil {
				return fmt.Errorf("error parsing taskID: %w", err)
			}
			taskID = api.ID(taskIDInt)
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

		t, err := client.UpdateTask(api.ID(taskID), input)
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

		token := config.MustGetToken()
		client := api.NewClient(token)

		var desc *string

		if len(args) == 0 {
			return errors.New("task title must be provided as first argument")
		}
		if len(args) > 1 {
			desc = &args[1]
		}

		task, err := client.CreateTask(api.TaskCreateInput{
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
