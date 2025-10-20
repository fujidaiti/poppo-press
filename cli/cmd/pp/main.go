package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:           "pp",
		Short:         "Poppo Press CLI",
		Long:          "Poppo Press CLI - script-friendly interface to the Poppo Press API",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.SetHelpTemplate(`{{.Long}}

Usage:
  {{.UseLine}}

Available Commands:
{{range .Commands}}{{if (and .IsAvailableCommand (not .IsAdditionalHelpTopicCommand))}}  {{rpad .Name .NamePadding}} {{.Short}}
{{end}}{{end}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}
`)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
