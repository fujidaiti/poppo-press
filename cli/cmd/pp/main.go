package main

import (
	"fmt"
	"os"

	"github.com/fujidaiti/poppo-press/cli/internal/commands"
)

func main() {
	root := commands.NewRootCmd()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
