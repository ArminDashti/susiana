package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ArminDashti/susiana-api/internal/config"
	"github.com/ArminDashti/susiana-api/internal/database"
	"github.com/ArminDashti/susiana-api/internal/handler"
	"github.com/ArminDashti/susiana-api/internal/provider"
	"github.com/ArminDashti/susiana-api/internal/repository"
	"github.com/ArminDashti/susiana-api/internal/scheduler"
	"github.com/ArminDashti/susiana-api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	gin.SetMode(cfg.GinMode)

	ctx := context.Background()
	pool, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	if err := database.Migrate(ctx, pool); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	marketProvider, err := provider.New(cfg.MarketProvider, cfg.MarketAPIBaseURL, cfg.MarketAPIKey)
	if err != nil {
		log.Fatalf("provider: %v", err)
	}

	stockRepo := repository.NewStockRepository(pool)
	etfRepo := repository.NewETFRepository(pool)
	syncRepo := repository.NewSyncRunRepository(pool)
	marketSvc := service.NewMarketService(marketProvider, stockRepo, etfRepo, syncRepo)

	sched, err := scheduler.Start(cfg.SyncInterval, marketSvc)
	if err != nil {
		log.Fatalf("scheduler: %v", err)
	}
	defer sched.Stop()

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	handler.NewMarketHandler(marketSvc).Register(r)

	srv := &http.Server{
		Addr:              ":" + cfg.HTTPPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("listening on :%s (provider=%s)", cfg.HTTPPort, cfg.MarketProvider)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
