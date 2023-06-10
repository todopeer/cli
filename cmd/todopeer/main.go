package main

import (
	"fmt"

	"github.com/todopeer/cli/commands"
)

func main() {
	err := commands.Run()
	if err != nil && err.Error() == "access denied" {
		fmt.Println("access defined. Probably token expired, try login again.")
	}
}
