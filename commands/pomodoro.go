package commands

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var (
	varContinuePomo bool
)

func init() {
	pomoCmd.Flags().BoolVarP(&varContinuePomo, "continue", "c", false, "continue pomo")
	rootCmd.AddCommand(pomoCmd)
}

var pomoCmd = &cobra.Command{
	Use:     "pomodoro",
	Aliases: []string{"pomo"},
	Short:   "pomodoro(pomo) [duration]: start a pomodoro with duration. Default to 25m",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		duration := 25 * time.Minute
		if len(args) > 0 {
			duration, err = time.ParseDuration(args[0])
			if err != nil {
				return fmt.Errorf("Err parse duration: %w", err)
			}
		}

		var startVal time.Duration
		var callback func() error

		if varContinuePomo {
			token := config.MustGetToken()
			ctx := context.Background()
			_, task, event, err := api.MeWithTaskEvent(ctx, token)
			if err != nil {
				return fmt.Errorf("error getting current task: %w", err)
			}
			if event == nil || task == nil {
				return fmt.Errorf("no running task")
			}

			fmt.Printf("continue pomo - task: %s; event start at: %s\n", task.Name, event.StartAt.EventTimeOnly())
			startVal = time.Since((time.Time)(event.StartAt))
			if startVal > duration {
				return fmt.Errorf("event is long enough that it should already completed the pomodoro")
			}
			callback = makeTaskPauseCallback(token, task.ID)
		}

		return pomodoro(duration, startVal, callback)
	},
}

func makeTaskPauseCallback(token string, taskID api.ID) func() error {
	return func() error {
		_, err := api.UpdateTask(context.Background(), token, taskID, api.TaskUpdateInput{
			Status: &api.TaskStatusPaused,
		})
		return err
	}
}

func pomodoro(duration time.Duration, startVal time.Duration, callback func() error) (err error) {
	running := true

	go func() {
		<-time.NewTimer(duration - startVal).C
		running = false
	}()

	start := time.Now().Add(-startVal)
	for running {
		time.Sleep(time.Second)
		diff := time.Since(start)

		fmt.Printf("\r%s/%s ", toMMSS(diff), toMMSS(duration))
		percentage := float32(diff) / float32(duration)

		showPercentage(percentage, 25)
	}
	fmt.Printf("\nDone :)\n")

	// TODO: this only works for Mac. Need to get others OS work
	if _, err := exec.LookPath("say"); err == nil {
		exec.Command("say", "the pomodoro is done").Run()
	}
	if callback != nil {
		return callback()
	}

	return nil
}

func showPercentage(percentage float32, size int) {
	cur := int(percentage * float32(size))
	fmt.Print("|")
	for i := 0; i < cur; i++ {
		fmt.Print("*")
	}
	for cur < size {
		fmt.Print(" ")
		cur++
	}
	fmt.Print("|")
}

func toMMSS(d time.Duration) string {
	return fmt.Sprintf("%02d:%02d", d/time.Minute, (d%time.Minute)/time.Second)
}
