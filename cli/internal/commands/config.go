package commands

import (
	"fmt"

	"github.com/fujidaiti/poppo-press/cli/internal/config"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
		Long:  "Manage CLI configuration",
	}

	tz := &cobra.Command{Use: "tz", Short: "Timezone settings", Long: "Timezone settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			if c.Timezone == "" {
				fmt.Fprintln(cmd.OutOrStdout(), "")
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), c.Timezone)
			return nil
		},
	}
	tz.AddCommand(&cobra.Command{
		Use:     "set <IANA-TZ>",
		Short:   "Set CLI timezone",
		Args:    cobra.ExactArgs(1),
		Example: "pp config tz set Asia/Tokyo",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			c.Timezone = args[0]
			return config.Save(c)
		},
	})
	cmd.AddCommand(tz)

	return cmd
}
