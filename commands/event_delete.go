package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/todopeer/cli/api"
	"github.com/todopeer/cli/services/config"
)

func init() {
	rootCmd.AddCommand(deleteEventCmd)
}

var deleteEventCmd = &cobra.Command{
	Use:     "delete-event",
	Aliases: []string{"de"},
	Short:   "delete-event (de) a event by its ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := config.MustGetToken()
		client := api.NewClient(token)

		var eventID api.ID
		var err error
		if len(args) == 0 {
			fmt.Println("eventID not provided, would use current running event")
			evt, err := client.QueryRunningEvent()
			if err != nil {
				return err
			}

			if evt == nil {
				return errors.New("no current running event")
			}
			eventID = evt.ID
		} else {
			eventIDInt, err := strconv.ParseInt(args[0], 10, 64)
			if err != nil {
				return err
			}
			eventID = api.ID(eventIDInt)
		}

		t, err := client.DeleteEvent(eventID)
		if err != nil {
			return err
		}
		fmt.Printf("event(id=%d) deleted successfully\n", t.ID)
		return err
	},
}
