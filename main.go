package main

import (
	"fmt"

	"github.com/todopeer/cli/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil && err.Error() == "access denied" {
		fmt.Println("access defined. Probably token expired, try login again.")
	}
}
