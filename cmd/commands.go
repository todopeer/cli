package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "diarier",
	Short: "Diarier is a CLI for interacting with your Diarier Backend",
}

func Run() error {
	return rootCmd.Execute()
}