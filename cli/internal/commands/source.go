package commands

import "github.com/spf13/cobra"

func newSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage sources",
		Long:  "Manage sources",
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "add <url>",
		Short:   "Add a source",
		Args:    cobra.ExactArgs(1),
		Example: "pp source add https://example.com/feed.xml",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List sources",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
		Example: "pp source list",
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "rm <id>",
		Short:   "Remove a source",
		Args:    cobra.ExactArgs(1),
		Example: "pp source rm 12",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	})

	return cmd
}
