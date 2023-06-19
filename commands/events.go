package commands

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
	"github.com/todopeer/cli/util/dt"
	"github.com/todopeer/cli/util/maps"
)

func getDayOffset(s string) (int, error) {
	offset, err := strconv.Atoi(s[1:])
	if err != nil {
		return 0, err
	}

	switch s[0] {
	case 'n':
		return offset, err
	case 'p':
		return -offset, err
	default:
		return 0, fmt.Errorf("unknown offset: %s", s)
	}
}

var listEventsCommand = &cobra.Command{
	Use:   "day",
	Short: "show events of a day, default to today.",
	Long:  "Can pass in `p[n]` to see n-th day before today, or [YYYY-MM-DD] to see a specific day",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		token := config.MustGetToken()
		ctx := context.Background()

		now := time.Now()
		dayForQuery := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

		if len(args) > 0 {
			param := args[0]
			if param[0] == 'p' {
				dayOffset, err := getDayOffset(param)
				if err != nil {
					return err
				}
				dayForQuery = dayForQuery.Add(time.Duration(dayOffset) * time.Hour * 24)
			} else {
				// expect to be a specific
				dayForQuery, err = dt.FromDate(param)
				if err != nil {
					return err
				}
			}
		}

		result, err := api.QueryEvents(ctx, token, dayForQuery, 1)
		if err != nil {
			return err
		}

		taskIDMap := map[api.ID]api.Task{}
		for _, t := range result.Tasks {
			taskIDMap[t.ID] = t
		}

		taskSummary := map[api.ID]time.Duration{}

		for _, e := range result.Events {
			t := taskIDMap[e.TaskID]
			start, err := dt.FromTime(string(e.StartAt))
			if err != nil {
				return err
			}
			startS := toLocalTimeStr(start)

			if e.EndAt == nil {
				fmt.Printf("[%d]: %s ~ doing -- %s\n", e.ID, startS, t.Name)
				taskSummary[t.ID] += time.Since(start)
			} else {
				end, err := dt.FromTimePtr((*string)(e.EndAt))
				if err != nil {
					return err
				}
				endS := toLocalTimeStr(*end)
				fmt.Printf("[%d]: %s ~ %s -- %s\n", e.ID, startS, endS, taskIDMap[e.TaskID].Name)
				taskSummary[t.ID] += end.Sub(start)
			}
		}

		// then show a summary on time spent
		fmt.Println()
		fmt.Println("\t*** Summary ***")
		sortedK := maps.SortedKByV(taskSummary)

		for i := len(sortedK) - 1; i >= 0; i-- {
			tid := sortedK[i]
			spent := taskSummary[tid]
			fmt.Printf("[%d]%s: %s\n", tid, taskIDMap[tid].Name, formatDuration(spent))
		}

		return nil
	},
}

func toLocalTimeStr(t time.Time) string {
	return t.Local().Format(time.TimeOnly)
}

func formatDuration(d time.Duration) string {
	res := ""
	if d >= time.Hour {
		res = strconv.Itoa(int(d/time.Hour)) + "h"
		d %= time.Hour
	}
	if d >= time.Minute {
		res = res + strconv.Itoa(int(d/time.Minute)) + "m"
		d %= time.Minute
	}
	if d >= time.Second {
		res = res + strconv.Itoa(int(d/time.Second)) + "s"
	}
	return res
}

func init() {
	rootCmd.AddCommand(listEventsCommand)
}
