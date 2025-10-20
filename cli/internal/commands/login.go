package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/fujidaiti/poppo-press/cli/internal/auth"
	"github.com/fujidaiti/poppo-press/cli/internal/config"
	"github.com/spf13/cobra"
)

func newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login and store token",
		Long:  "Login and store token",
		RunE: func(cmd *cobra.Command, args []string) error {
			device, _ := cmd.Flags().GetString("device")
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")

			// env overrides per roadmap
			if v := os.Getenv("PP_USERNAME"); v != "" {
				username = v
			}
			if v := os.Getenv("PP_PASSWORD"); v != "" {
				password = v
			}
			if v := os.Getenv("PP_TOKEN"); v != "" {
				c, err := config.Load()
				if err != nil {
					return err
				}
				c.Token = v
				if err := config.Save(c); err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), "Login successful. Token saved for device \""+device+"\".")
				return nil
			}

			if username == "" || password == "" {
				return fmt.Errorf("username and password are required (or provide PP_TOKEN)")
			}

			c, err := config.Load()
			if err != nil {
				return err
			}
			if c.Server == "" {
				return fmt.Errorf("server not configured; run 'pp init --server <url>'")
			}
			ac := auth.New(c.Server)
			resp, err := ac.Login(context.Background(), auth.LoginRequest{Username: username, Password: password, DeviceName: device})
			if err != nil {
				return err
			}
			c.Token = resp.Token
			if err := config.Save(c); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Login successful. Token saved for device \""+device+"\".")
			return nil
		},
		Example: "pp login --device devbox --username admin --password secret",
	}
	cmd.Flags().String("device", "", "device name")
	cmd.Flags().String("username", "", "username")
	cmd.Flags().String("password", "", "password")
	_ = cmd.MarkFlagRequired("device")
	return cmd
}
