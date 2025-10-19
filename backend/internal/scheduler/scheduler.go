package scheduler

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/fujidaiti/poppo-press/backend/internal/fetcher"
)

type Scheduler struct {
	c *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{c: cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)))}
}

func (s *Scheduler) Start() { s.c.Start() }
func (s *Scheduler) Stop()  { ctx := s.c.Stop(); <-ctx.Done() }

func (s *Scheduler) HourlyFetch(database *sql.DB) error {
	_, err := s.c.AddFunc("0 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := fetcher.FetchAllSources(ctx, database, nil); err != nil {
			log.Printf("fetch job error: %v", err)
		}
	})
	return err
}
