package commands

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
)

var debugMode = false

var rootCmd = &cobra.Command{
	Use: "todopeer",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugMode {
			api.SetLogger(log.Printf)
		}
	},
	Short: "a CLI for interacting with your Todopeer Backend",
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug mode")
}

func Run() error {
	return rootCmd.Execute()
}
