package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ArminDashti/susiana-api/internal/domain"
)

// HTTPProvider fetches market data from a configurable REST API.
// Endpoint paths and JSON field mapping are placeholders until the real API is provided.
type HTTPProvider struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	stocksPath string
	etfsPath   string
}

type HTTPOption func(*HTTPProvider)

func WithStocksPath(path string) HTTPOption {
	return func(p *HTTPProvider) { p.stocksPath = path }
}

func WithETFsPath(path string) HTTPOption {
	return func(p *HTTPProvider) { p.etfsPath = path }
}

func NewHTTP(baseURL, apiKey string, opts ...HTTPOption) *HTTPProvider {
	p := &HTTPProvider{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		// Placeholder paths — update when the real API is known.
		stocksPath: "/stocks",
		etfsPath:   "/etfs",
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *HTTPProvider) FetchStocks(ctx context.Context) ([]domain.Stock, error) {
	body, err := p.get(ctx, p.stocksPath)
	if err != nil {
		return nil, err
	}

	var stocks []domain.Stock
	if err := json.Unmarshal(body, &stocks); err != nil {
		// Try wrapped envelope: {"data": [...]}
		var envelope struct {
			Data []domain.Stock `json:"data"`
		}
		if err2 := json.Unmarshal(body, &envelope); err2 != nil {
			return nil, fmt.Errorf("decode stocks response: %w", err)
		}
		stocks = envelope.Data
	}

	return stocks, nil
}

func (p *HTTPProvider) FetchETFs(ctx context.Context) ([]domain.ETF, error) {
	body, err := p.get(ctx, p.etfsPath)
	if err != nil {
		return nil, err
	}

	var etfs []domain.ETF
	if err := json.Unmarshal(body, &etfs); err != nil {
		var envelope struct {
			Data []domain.ETF `json:"data"`
		}
		if err2 := json.Unmarshal(body, &envelope); err2 != nil {
			return nil, fmt.Errorf("decode etfs response: %w", err)
		}
		etfs = envelope.Data
	}

	return etfs, nil
}

func (p *HTTPProvider) get(ctx context.Context, path string) ([]byte, error) {
	url := p.baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get %s: %w", path, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("market api %s returned %d: %s", path, resp.StatusCode, truncate(string(body), 256))
	}
	return body, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// New selects the configured provider implementation.
func New(providerName, baseURL, apiKey string) (domain.MarketDataProvider, error) {
	switch providerName {
	case "stub":
		return NewStub(), nil
	case "http":
		return NewHTTP(baseURL, apiKey), nil
	default:
		return nil, fmt.Errorf("unknown market provider %q", providerName)
	}
}
