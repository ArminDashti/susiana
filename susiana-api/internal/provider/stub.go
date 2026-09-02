package provider

import (
	"context"
	"fmt"

	"github.com/ArminDashti/susiana-api/internal/domain"
)

// StubProvider is a placeholder until the real market REST API is configured.
type StubProvider struct{}

func NewStub() *StubProvider {
	return &StubProvider{}
}

func (s *StubProvider) FetchStocks(ctx context.Context) ([]domain.Stock, error) {
	return nil, fmt.Errorf("market data provider is stub: configure MARKET_PROVIDER=http and MARKET_API_BASE_URL when the API is available")
}

func (s *StubProvider) FetchETFs(ctx context.Context) ([]domain.ETF, error) {
	return nil, fmt.Errorf("market data provider is stub: configure MARKET_PROVIDER=http and MARKET_API_BASE_URL when the API is available")
}
