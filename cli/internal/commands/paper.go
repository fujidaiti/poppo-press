package commands

import "github.com/spf13/cobra"

func newPaperCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "paper",
		Short: "Daily editions",
		Long:  "Daily editions",
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "read",
		Short:   "Read today's or given edition",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
		Example: "pp paper read --date 2025-10-20",
	})

	list := &cobra.Command{
		Use:     "list",
		Short:   "List recent editions",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
		Example: "pp paper list --limit 10",
	}
	list.Flags().Int("limit", 0, "max number of items")
	list.Flags().Int("offset", 0, "number of items to skip")
	cmd.AddCommand(list)

	return cmd
}
