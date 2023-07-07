package commands

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

var (
	varContinuePomo bool

	defaultPomoSize  = 25 * time.Minute
	defaultBreakSize = 5 * time.Minute
)

const (
	msgPomoStart  = "starting new pomodoro"
	msgPomoDone   = "the pomodoro is done"
	msgBreakStart = "task paused, starting break"
	msgBreakDone  = "the break is done"
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
		duration := defaultPomoSize
		if len(args) > 0 {
			duration, err = time.ParseDuration(args[0])
			if err != nil {
				return fmt.Errorf("err parse duration: %w", err)
			}
		}

		var startVal time.Duration
		var callback func() error

		if varContinuePomo {
			token := config.MustGetToken()
			client := api.NewClient(token)

			_, task, event, err := client.MeWithTaskEvent()
			if err != nil {
				return fmt.Errorf("error getting current task: %w", err)
			}
			if event == nil || task == nil {
				return ErrNoRunningEvent
			}

			fmt.Printf("continue pomo - task: %s; event start at: %s\n", task.Name, event.StartAt.EventTimeOnly())
			startVal = time.Since((time.Time)(event.StartAt))
			if startVal > duration {
				return fmt.Errorf("event is long enough that it should already completed the pomodoro")
			}
			callback = makeTaskPauseCallback(client, task.ID)
		}

		return pomodoro(duration, startVal, msgPomoDone, callback)
	},
}

func makeTaskPauseCallback(client *api.Client, taskID api.ID) func() error {
	return func() error {
		_, err := client.UpdateTask(taskID, api.TaskUpdateInput{
			Status: &api.TaskStatusPaused,
		})
		if err != nil {
			return err
		}

		tryToSayWithLog(msgBreakStart)
		pomodoro(defaultBreakSize, 0, msgBreakDone, nil)
		return err
	}
}

func tryToSayWithLog(msg string) {
	// TODO: this only works for Mac. Need to get others OS work
	if _, err := exec.LookPath("say"); err == nil {
		exec.Command("say", msg).Run()
	}
	log.Print(msg)
}

func pomodoro(duration time.Duration, startVal time.Duration, message string, callback func() error) (err error) {
	tryToSayWithLog(msgPomoStart)
	start := time.Now().Add(-startVal)
	for {
		time.Sleep(time.Second)
		diff := time.Since(start)

		fmt.Printf("\r%s/%s ", toMMSS(diff), toMMSS(duration))
		percentage := float32(diff) / float32(duration)
		if percentage >= 1 {
			break
		}

		showPercentage(percentage, 25)
	}
	tryToSayWithLog(message)

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
