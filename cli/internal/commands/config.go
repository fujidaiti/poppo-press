package commands

import "github.com/spf13/cobra"

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  "Manage CLI configuration",
	}

	tz := &cobra.Command{Use: "tz", Short: "Timezone settings", Long: "Timezone settings"}
	tz.AddCommand(&cobra.Command{
		Use:     "set <IANA-TZ>",
		Short:   "Set CLI timezone",
		Args:    cobra.ExactArgs(1),
		Example: "pp config tz set Asia/Tokyo",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	})
	tz.AddCommand(&cobra.Command{
		Use:     "",
		Short:   "Print current CLI timezone",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
		Example: "pp config tz",
	})
	cmd.AddCommand(tz)

	return cmd
}
