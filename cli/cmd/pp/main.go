package main

import (
	"fmt"
	"os"

	"github.com/fujidaiti/poppo-press/cli/internal/commands"
	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
)

func main() {
	root := commands.NewRootCmd()
	if err := root.Execute(); err != nil {
		// Map known http client error to exit code; otherwise generic 1
		if he, ok := err.(*httpc.Error); ok {
			os.Exit(he.Code)
		}
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
