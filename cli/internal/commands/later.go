package commands

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fujidaiti/poppo-press/cli/internal/config"
	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
	"github.com/spf13/cobra"
)

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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			req, err := hc.NewRequest(cmd.Context(), http.MethodPost, "/v1/read-later/"+args[0], nil)
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

	list := &cobra.Command{
		Use:   "list",
		Short: "List read-later articles",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			req, err := hc.NewRequest(cmd.Context(), http.MethodGet, "/v1/read-later", nil)
			if err != nil {
				return err
			}
			resp, err := hc.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			b, _ := io.ReadAll(resp.Body)
			// apply client-side pagination if flags present
			limit, _ := cmd.Flags().GetInt("limit")
			offset, _ := cmd.Flags().GetInt("offset")
			if limit > 0 || offset > 0 {
				var items []map[string]any
				if err := json.Unmarshal(b, &items); err == nil {
					if offset < 0 {
						offset = 0
					}
					if offset > len(items) {
						offset = len(items)
					}
					end := len(items)
					if limit > 0 && offset+limit < end {
						end = offset + limit
					}
					items = items[offset:end]
					if enc, err := json.Marshal(items); err == nil {
						b = enc
					}
				}
			}
			if len(b) > 0 && b[len(b)-1] != '\n' {
				b = append(b, '\n')
			}
			_, _ = cmd.OutOrStdout().Write(b)
			return nil
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			req, err := hc.NewRequest(cmd.Context(), http.MethodDelete, "/v1/read-later/"+args[0], nil)
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
