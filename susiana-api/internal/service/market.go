package service

import (
	"context"
	"fmt"

	"github.com/ArminDashti/susiana-api/internal/domain"
	"github.com/ArminDashti/susiana-api/internal/repository"
)

type StockStore interface {
	UpsertMany(ctx context.Context, stocks []domain.Stock) error
	GetByISIN(ctx context.Context, isin string) (*domain.Stock, error)
	List(ctx context.Context, params domain.ListParams) ([]domain.Stock, int, error)
}

type ETFStore interface {
	UpsertMany(ctx context.Context, etfs []domain.ETF) error
	GetByISIN(ctx context.Context, isin string) (*domain.ETF, error)
	List(ctx context.Context, params domain.ListParams) ([]domain.ETF, int, error)
}

type SyncRunStore interface {
	Start(ctx context.Context, instrumentType domain.InstrumentType) (*domain.SyncRun, error)
	Finish(ctx context.Context, id int64, status domain.SyncStatus, fetchedCount int, errMsg string) error
	ListRecent(ctx context.Context, limit int) ([]domain.SyncRun, error)
}

type MarketService struct {
	provider domain.MarketDataProvider
	stocks   StockStore
	etfs     ETFStore
	syncRuns SyncRunStore
}

func NewMarketService(
	provider domain.MarketDataProvider,
	stocks StockStore,
	etfs ETFStore,
	syncRuns SyncRunStore,
) *MarketService {
	return &MarketService{
		provider: provider,
		stocks:   stocks,
		etfs:     etfs,
		syncRuns: syncRuns,
	}
}

type SyncResult struct {
	Run          *domain.SyncRun `json:"run"`
	FetchedCount int             `json:"fetched_count"`
}

func (s *MarketService) Sync(ctx context.Context, instrumentType domain.InstrumentType) (*SyncResult, error) {
	switch instrumentType {
	case domain.InstrumentTypeStock, domain.InstrumentTypeETF, domain.InstrumentTypeAll:
	default:
		return nil, fmt.Errorf("invalid instrument type %q", instrumentType)
	}

	run, err := s.syncRuns.Start(ctx, instrumentType)
	if err != nil {
		return nil, err
	}

	fetched := 0
	var syncErr error

	switch instrumentType {
	case domain.InstrumentTypeStock:
		fetched, syncErr = s.syncStocks(ctx)
	case domain.InstrumentTypeETF:
		fetched, syncErr = s.syncETFs(ctx)
	case domain.InstrumentTypeAll:
		var nStocks, nETFs int
		nStocks, syncErr = s.syncStocks(ctx)
		if syncErr == nil {
			nETFs, syncErr = s.syncETFs(ctx)
		}
		fetched = nStocks + nETFs
	}

	status := domain.SyncStatusSuccess
	errMsg := ""
	if syncErr != nil {
		status = domain.SyncStatusFailed
		errMsg = syncErr.Error()
	}

	if finishErr := s.syncRuns.Finish(ctx, run.ID, status, fetched, errMsg); finishErr != nil && syncErr == nil {
		return nil, finishErr
	}

	run.Status = status
	run.FetchedCount = fetched
	run.Error = errMsg

	result := &SyncResult{Run: run, FetchedCount: fetched}
	if syncErr != nil {
		return result, syncErr
	}
	return result, nil
}

func (s *MarketService) syncStocks(ctx context.Context) (int, error) {
	stocks, err := s.provider.FetchStocks(ctx)
	if err != nil {
		return 0, err
	}
	if err := s.stocks.UpsertMany(ctx, stocks); err != nil {
		return 0, err
	}
	return len(stocks), nil
}

func (s *MarketService) syncETFs(ctx context.Context) (int, error) {
	etfs, err := s.provider.FetchETFs(ctx)
	if err != nil {
		return 0, err
	}
	if err := s.etfs.UpsertMany(ctx, etfs); err != nil {
		return 0, err
	}
	return len(etfs), nil
}

func (s *MarketService) ListStocks(ctx context.Context, params domain.ListParams) ([]domain.Stock, int, error) {
	return s.stocks.List(ctx, params)
}

func (s *MarketService) GetStock(ctx context.Context, isin string) (*domain.Stock, error) {
	return s.stocks.GetByISIN(ctx, isin)
}

func (s *MarketService) ListETFs(ctx context.Context, params domain.ListParams) ([]domain.ETF, int, error) {
	return s.etfs.List(ctx, params)
}

func (s *MarketService) GetETF(ctx context.Context, isin string) (*domain.ETF, error) {
	return s.etfs.GetByISIN(ctx, isin)
}

func (s *MarketService) ListSyncRuns(ctx context.Context, limit int) ([]domain.SyncRun, error) {
	return s.syncRuns.ListRecent(ctx, limit)
}

// Ensure repository types satisfy interfaces at compile time.
var (
	_ StockStore   = (*repository.StockRepository)(nil)
	_ ETFStore     = (*repository.ETFRepository)(nil)
	_ SyncRunStore = (*repository.SyncRunRepository)(nil)
)
