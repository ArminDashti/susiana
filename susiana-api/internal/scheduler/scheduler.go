package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/ArminDashti/susiana-api/internal/domain"
	"github.com/ArminDashti/susiana-api/internal/service"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func Start(interval time.Duration, svc *service.MarketService) (*Scheduler, error) {
	if interval <= 0 {
		return nil, nil
	}

	c := cron.New()
	spec := "@every " + interval.String()

	_, err := c.AddFunc(spec, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		result, err := svc.Sync(ctx, domain.InstrumentTypeAll)
		if err != nil {
			log.Printf("scheduled sync failed: %v", err)
			return
		}
		log.Printf("scheduled sync ok: fetched=%d run_id=%d", result.FetchedCount, result.Run.ID)
	})
	if err != nil {
		return nil, err
	}

	c.Start()
	log.Printf("scheduler started: interval=%s", interval)
	return &Scheduler{cron: c}, nil
}

func (s *Scheduler) Stop() {
	if s == nil || s.cron == nil {
		return
	}
	ctx := s.cron.Stop()
	<-ctx.Done()
}
