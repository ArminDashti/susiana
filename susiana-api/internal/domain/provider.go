package domain

import "context"

// MarketDataProvider fetches stock and ETF quotes from an external source.
type MarketDataProvider interface {
	FetchStocks(ctx context.Context) ([]Stock, error)
	FetchETFs(ctx context.Context) ([]ETF, error)
}
