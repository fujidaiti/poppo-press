package commands

import (
	"github.com/spf13/cobra"
)

// NewRootCmd constructs the root Cobra command for the Poppo Press CLI.
// The command disables default completion commands and sets a custom help template
// aligned with the raw, script-friendly UX constraints (no colors/pager).
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "pp",
		Short:         "Poppo Press CLI",
		Long:          "Poppo Press CLI - script-friendly interface to the Poppo Press API",
		SilenceUsage:  true,
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
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

	return root
}
