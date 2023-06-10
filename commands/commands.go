package commands

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "todopeer",
	Short: "a CLI for interacting with your Todopeer Backend",
}

func Run() error {
	return rootCmd.Execute()
}
