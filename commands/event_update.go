package commands

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
	"github.com/todopeer/cli/util/dt"
)

var (
	varStartAtStr string
	varEndAtStr   string

	varDayoffsetStr string
	varNewTaskIDStr string
)

func defineFlagsForEvent(s *pflag.FlagSet, isUpdate bool) {
	s.StringVarP(&varStartAtStr, "startAt", "s", "", "startAt for event")
	s.StringVarP(&varEndAtStr, "endAt", "e", "", "endAt")
	s.StringVarP(&varDayoffsetStr, "offset", "D", "", "if provided, the start / end date would be based on offset day")
	s.StringVarP(&varDescription, "desc", "d", "", "description")
	s.StringVarP(&varNewTaskIDStr, "taskID", "t", "", "new taskID to assign")
}

func init() {
	rootCmd.AddCommand(removeEventCmd)

	defineFlagsForEvent(updateEventCmd.Flags(), true)
	rootCmd.AddCommand(updateEventCmd)
}

var removeEventCmd = &cobra.Command{
	Use:     "remove-event",
	Aliases: []string{"re"},
	Short:   "remove-event (re) a event by its ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		if len(args) == 0 {
			log.Fatal("eventID must be provided")
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

var updateEventCmd = &cobra.Command{
	Use:     "update-event",
	Aliases: []string{"ue"},
	Short:   "(ue)update-event [event-id] -flags",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		if len(args) == 0 {
			log.Fatal("eventID must be provided")
		}

		eventID, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		event, err := api.GetEvent(ctx, token, api.ID(eventID))
		if err != nil {
			return err
		}
		startTime, err := dt.FromTime(string(event.StartAt))
		if err != nil {
			return fmt.Errorf("err parse e.startAt: %w", err)
		}
		endTimeP, err := dt.FromTimePtr((*string)(event.EndAt))
		if err != nil {
			return fmt.Errorf("err parse e.endAt: %w", err)
		}

		dayOffset := 0
		if len(varDayoffsetStr) > 0 {
			dayOffset, err = getDayOffset(varDayoffsetStr)
			if err != nil {
				return err
			}
		}

		input := api.EventUpdateInput{}
		if varDescription != "" {
			input.Description = graphql.NewString(graphql.String(varDescription))
		}
		input.StartAt, err = parsePointOfTime(&startTime, dayOffset, varStartAtStr)
		if err != nil {
			return fmt.Errorf("err parse startInput: %w", err)
		}
		input.EndAt, err = parsePointOfTime(endTimeP, dayOffset, varEndAtStr)
		if err != nil {
			return fmt.Errorf("err parse endInput: %w", err)
		}

		if len(varNewTaskIDStr) > 0 {
			taskID, err := strconv.ParseInt(varNewTaskIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("err parse taskID: %w", err)
			}

			input.TaskID = (*api.ID)(&taskID)
		}

		e, err := api.UpdateEvent(ctx, token, api.ID(eventID), input)
		if err != nil {
			return err
		}
		fmt.Println("event successfully updated")
		return outputEvent(e, true)
	},
}

func outputEvent(e *api.Event, withDate bool) error {
	format := time.TimeOnly
	if withDate {
		format = time.DateTime
	}

	from, err := dt.FromTime(string(e.StartAt))
	if err != nil {
		return err
	}
	fromStr := from.Local().Format(format)

	toStr := "doing"
	if e.EndAt != nil {
		endT, err := dt.FromTime(string(*e.EndAt))
		if err != nil {
			return err
		}
		toStr = endT.Local().Format(format)
	}
	fmt.Printf("%d: %s - %s", e.ID, fromStr, toStr)

	if len(e.Description) > 0 {
		fmt.Println(":", e.Description)
	} else {
		fmt.Println()
	}

	return nil
}

func parsePointOfTime(dateReference *time.Time, dayOffset int, s string) (*graphql.String, error) {
	if len(s) == 0 {
		return nil, nil
	}

	relDuration, err := tryParseDuration(s)
	if err != nil {
		return nil, err
	}

	if dateReference == nil {
		now := time.Now()
		dateReference = &now
	} else {
		localDate := dateReference.Local()
		dateReference = &localDate
	}

	nd := dateReference.AddDate(0, 0, dayOffset)
	dateReference = &nd

	if relDuration != nil {
		// this is relative time
		ts := dt.ToTime(dateReference.Add(*relDuration))
		return (*graphql.String)(&ts), nil
	}

	// the format would be "HH:MM" or "HH:MM:SS"
	parts := strings.Split(s, ":")
	if len(parts) < 2 || len(parts) > 3 {
		return nil, fmt.Errorf("expect time in format of HH:MM[:SS], got: %s", s)
	}

	// build duration
	hms := [3]int{}

	for i, part := range parts {
		v, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("time isn't number: %s. Parse error: %w", s, err)
		}
		hms[i] = v
	}
	t := time.Date(dateReference.Year(), dateReference.Month(), dateReference.Day(),
		hms[0], hms[1], hms[2], 0, dateReference.Location())

	ts := dt.ToTime(t)
	return (*graphql.String)(&ts), err
}

func tryParseDuration(s string) (*time.Duration, error) {
	if s[0] == 'p' || s[0] == 'n' {
		duration, err := time.ParseDuration(s)
		if s[0] == 'p' {
			duration = -duration
		}

		return &duration, err
	}
	return nil, nil
}
