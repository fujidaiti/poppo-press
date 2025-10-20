package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/fujidaiti/poppo-press/cli/internal/config"
	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
	"github.com/spf13/cobra"
)

func newSourceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "source",
		Short: "Manage sources",
		Long:  "Manage sources",
	}
	cmd.PersistentFlags().Bool("verbose", false, "print HTTP request/response traces")

	cmd.AddCommand(&cobra.Command{
		Use:     "add <url>",
		Short:   "Add a source",
		Args:    cobra.ExactArgs(1),
		Example: "pp source add https://example.com/feed.xml",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			verbose, _ := cmd.Flags().GetBool("verbose")
			var hc *httpc.Client
			if verbose {
				hc, err = httpc.New(c.Server, c.Token, httpc.WithVerbose(cmd.ErrOrStderr()))
			} else {
				hc, err = httpc.New(c.Server, c.Token)
			}
			if err != nil {
				return err
			}
			body, _ := json.Marshal(map[string]string{"url": args[0]})
			req, err := hc.NewRequest(context.Background(), http.MethodPost, "/v1/sources", bytes.NewReader(body))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := hc.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			var out struct {
				ID any `json:"id"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), out.ID)
			return nil
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List sources",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			verbose, _ := cmd.Flags().GetBool("verbose")
			var hc *httpc.Client
			if verbose {
				hc, err = httpc.New(c.Server, c.Token, httpc.WithVerbose(cmd.ErrOrStderr()))
			} else {
				hc, err = httpc.New(c.Server, c.Token)
			}
			if err != nil {
				return err
			}
			req, err := hc.NewRequest(context.Background(), http.MethodGet, "/v1/sources", nil)
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
		Example: "pp source list",
	})

	cmd.AddCommand(&cobra.Command{
		Use:     "rm <id>",
		Short:   "Remove a source",
		Args:    cobra.ExactArgs(1),
		Example: "pp source rm 12",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			verbose, _ := cmd.Flags().GetBool("verbose")
			var hc *httpc.Client
			if verbose {
				hc, err = httpc.New(c.Server, c.Token, httpc.WithVerbose(cmd.ErrOrStderr()))
			} else {
				hc, err = httpc.New(c.Server, c.Token)
			}
			if err != nil {
				return err
			}
			id := args[0]
			if _, err := strconv.ParseInt(id, 10, 64); err != nil {
				return fmt.Errorf("invalid id: %s", id)
			}
			req, err := hc.NewRequest(context.Background(), http.MethodDelete, "/v1/sources/"+id, nil)
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
