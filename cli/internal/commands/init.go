package commands

import "github.com/spf13/cobra"

func newInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize CLI configuration",
		Long:  "Initialize CLI configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Example: "pp init --server http://localhost:8080",
	}
	cmd.Flags().String("server", "", "API base URL")
	_ = cmd.MarkFlagRequired("server")
	return cmd
}
