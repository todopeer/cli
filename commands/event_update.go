package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shurcooL/graphql"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
	"github.com/todopeer/cli/util/dt"
	"github.com/todopeer/cli/util/gql"
)

func defineFlagsForEvent(s *pflag.FlagSet, isUpdate bool) {
	s.StringVarP(&varStartAtStr, "startAt", "s", "", "startAt for event")
	s.StringVarP(&varEndAtStr, "endAt", "e", "", "endAt")
	s.StringVarP(&varDayoffsetStr, "offset", "D", "", "if provided, the start / end date would be based on offset day")
	s.StringVarP(&varDescription, "desc", "d", "", "description")
	s.StringVarP(&varNewTaskIDStr, "taskID", "t", "", "new taskID to assign")
}

func init() {
	defineFlagsForEvent(updateEventCmd.Flags(), true)
	rootCmd.AddCommand(updateEventCmd)
}

var updateEventCmd = &cobra.Command{
	Use:     "update-event",
	Aliases: []string{"ue"},
	Short:   "(ue)update-event [event-id] -flags",
	Long: `Syntax Supported:
update-event [event-id]: update event. Details are provided as flags
update-event: update the current running event. Errors if no running event
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		ctx := context.Background()

		var eventID api.ID
		if len(args) == 0 {
			runningEvent, err := api.QueryRunningEvent(ctx, token)
			if err != nil {
				return fmt.Errorf("error querying running event: %w", err)
			}
			if runningEvent == nil {
				return errors.New("no running event")
			}

			eventID = runningEvent.ID
		} else {
			eventIDInt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("error parsing event ID: %w", err)
			}
			eventID = api.ID(eventIDInt)
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
			input.Description = gql.ToGqlStringP(varDescription)
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

	if e.Description != nil {
		fmt.Println(":", *e.Description)
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
		return gql.ToGqlStringP(dt.ToTime(dateReference.Add(*relDuration))), nil
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

	return gql.ToGqlStringP(dt.ToTime(t)), err
}

func tryParseDuration(s string) (*time.Duration, error) {
	if s[0] == 'p' || s[0] == 'n' {
		duration, err := time.ParseDuration(s[1:])
		if s[0] == 'p' {
			duration = -duration
		}

		return &duration, err
	}
	return nil, nil
}
