package scheduler

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/fujidaiti/poppo-press/backend/internal/aggregator"
	"github.com/fujidaiti/poppo-press/backend/internal/config"
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

func (s *Scheduler) DailyAssemble(database *sql.DB, cfg config.Config) error {
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		loc = time.Local
	}
	// Run at 08:00 local time
	_, err = s.c.AddFunc("0 0 8 * *", func() {
		now := time.Now().In(loc)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := aggregator.AssembleDailyEdition(ctx, database, loc, now); err != nil {
			log.Printf("assemble job error: %v", err)
		}
	})
	return err
}
