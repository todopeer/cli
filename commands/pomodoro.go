package commands

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pomoCmd)
}

var pomoCmd = &cobra.Command{
	Use:     "pomodoro",
	Aliases: []string{"pomo"},
	Short:   "pomodoro(pomo) [duration]: start a pomodoro",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		duration := 25 * time.Minute
		if len(args) > 0 {
			duration, err = time.ParseDuration(args[0])
			if err != nil {
				return fmt.Errorf("Err parse duration: %w", err)
			}
		}
		return pomodoro(duration)
	},
}

func pomodoro(duration time.Duration) (err error) {
	running := true

	go func() {
		<-time.NewTimer(duration).C
		running = false
	}()

	start := time.Now()
	for running {
		time.Sleep(time.Second)
		diff := time.Since(start)

		fmt.Printf("\r%s/%s ", toMMSS(diff), toMMSS(duration))
		percentage := float32(diff) / float32(duration)

		showPercentage(percentage, 25)
	}
	fmt.Printf("\nDone :)\n")

	return exec.Command("say", "the pomodoro is done").Run()
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
