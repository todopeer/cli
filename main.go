package main

import (
	"fmt"

	"github.com/flyfy1/diarier_cli/cmd"
)

func main() {
	err := cmd.Run()
	if err != nil && err.Error() == "access denied" {
		fmt.Println("access defined. Probably token expired, try login again.")
	}
}
