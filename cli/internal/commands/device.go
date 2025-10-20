package commands

import (
	"io"
	"net/http"

	"github.com/fujidaiti/poppo-press/cli/internal/config"
	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
	"github.com/spf13/cobra"
)

func newDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "device",
		Short: "Manage devices",
		Long:  "Manage devices",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			req, err := hc.NewRequest(cmd.Context(), http.MethodGet, "/v1/devices", nil)
			if err != nil {
				return err
			}
			resp, err := hc.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			b, _ := io.ReadAll(resp.Body)
			if len(b) > 0 && b[len(b)-1] != '\n' {
				b = append(b, '\n')
			}
			_, _ = cmd.OutOrStdout().Write(b)
			return nil
		},
		Example: "pp device list",
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "revoke <id>",
		Short:   "Revoke a device",
		Args:    cobra.ExactArgs(1),
		Example: "pp device revoke 5",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			req, err := hc.NewRequest(cmd.Context(), http.MethodDelete, "/v1/devices/"+args[0], nil)
			if err != nil {
				return err
			}
			resp, err := hc.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			return nil
		},
	})

	return cmd
}
