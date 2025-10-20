package commands

import "github.com/spf13/cobra"

func newLaterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "later",
		Short: "Read-later queue",
		Long:  "Read-later queue",
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "add <article-id>",
		Short:   "Add an article to read later",
		Args:    cobra.ExactArgs(1),
		Example: "pp later add 202",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	})

	list := &cobra.Command{
		Use:     "list",
		Short:   "List read-later articles",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
		Example: "pp later list --limit 10",
	}
	list.Flags().Int("limit", 0, "max number of items")
	list.Flags().Int("offset", 0, "number of items to skip")
	cmd.AddCommand(list)

	cmd.AddCommand(&cobra.Command{
		Use:     "rm <article-id>",
		Short:   "Remove an article from read later",
		Args:    cobra.ExactArgs(1),
		Example: "pp later rm 202",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	})

	return cmd
}
