package commands

import "github.com/spf13/cobra"

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and store token",
		Long:  "Login and store token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Example: "pp login --device devbox --username admin --password secret",
	}
	cmd.Flags().String("device", "", "device name")
	cmd.Flags().String("username", "", "username")
	cmd.Flags().String("password", "", "password")
	_ = cmd.MarkFlagRequired("device")
	_ = cmd.MarkFlagRequired("username")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}
