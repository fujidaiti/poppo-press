package commands

import "github.com/spf13/cobra"

func newDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "Manage devices",
		Long:  "Manage devices",
	}

	cmd.AddCommand(&cobra.Command{
		Use:     "list",
		Short:   "List devices",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
		Example: "pp device list",
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "revoke <id>",
		Short:   "Revoke a device",
		Args:    cobra.ExactArgs(1),
		Example: "pp device revoke 5",
		RunE:    func(cmd *cobra.Command, args []string) error { return nil },
	})

	return cmd
}
