package main

import (
	"fmt"
	"os"

	"github.com/fujidaiti/poppo-press/cli/internal/commands"
	"github.com/fujidaiti/poppo-press/cli/internal/diag"
	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
)

func main() {
	root := commands.NewRootCmd()
	if err := root.Execute(); err != nil {
		if he, ok := err.(*httpc.Error); ok {
			fmt.Fprintln(os.Stderr, diag.FormatError(he))
			os.Exit(he.Code)
		}
		fmt.Fprintln(os.Stderr, diag.FormatError(err))
		os.Exit(1)
	}
}
