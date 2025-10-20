package commands

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/fujidaiti/poppo-press/cli/internal/config"
	"github.com/fujidaiti/poppo-press/cli/internal/httpc"
	"github.com/spf13/cobra"
)

func newPaperCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "paper",
		Short: "Daily editions",
		Long:  "Daily editions",
	}

	read := &cobra.Command{
		Use:     "read",
		Short:   "Read today's or given edition",
		Example: "pp paper read --id 17",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			id, _ := cmd.Flags().GetString("id")
			if id == "" {
				return fmt.Errorf("--id is required for now")
			}
			req, err := hc.NewRequest(cmd.Context(), http.MethodGet, "/v1/editions/"+id, nil)
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
	}
	read.Flags().String("id", "", "edition id")
	cmd.AddCommand(read)

	list := &cobra.Command{
		Use:     "list",
		Short:   "List recent editions",
		Example: "pp paper list --limit 10",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := config.Load()
			if err != nil {
				return err
			}
			hc, err := httpc.New(c.Server, c.Token)
			if err != nil {
				return err
			}
			limit, _ := cmd.Flags().GetInt("limit")
			offset, _ := cmd.Flags().GetInt("offset")
			if limit < 0 {
				limit = 0
			}
			if offset < 0 {
				offset = 0
			}
			pageSize := limit
			if pageSize <= 0 {
				pageSize = 20
			}
			page := 1
			if limit > 0 {
				page = offset/limit + 1
			}
			path := "/v1/editions?page=" + strconv.Itoa(page) + "&pageSize=" + strconv.Itoa(pageSize)
			req, err := hc.NewRequest(cmd.Context(), http.MethodGet, path, nil)
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
	}
	list.Flags().Int("limit", 0, "max number of items")
	list.Flags().Int("offset", 0, "number of items to skip")
	cmd.AddCommand(list)

	return cmd
}
