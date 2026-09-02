package domain

import (
	"encoding/json"
	"time"
)

// QuoteFields holds shared market-data quote and identity columns for stocks and ETFs.
type QuoteFields struct {
	Last                    *float64        `json:"last"`
	Close                   *float64        `json:"close"`
	Open                    *float64        `json:"open"`
	LowerLimit              *float64        `json:"lower_limit"`
	UpperLimit              *float64        `json:"upper_limit"`
	BaseVolume              *int64          `json:"base_volume"`
	Volume                  *int64          `json:"volume"`
	Turnover                *float64        `json:"turnover"`
	Trades                  *int64          `json:"trades"`
	RetailBuyVolume         *int64          `json:"retail_buy_volume"`
	InstitutionalBuyVolume  *int64          `json:"institutional_buy_volume"`
	RetailBuyCount          *int64          `json:"retail_buy_count"`
	InstitutionalBuyCount   *int64          `json:"institutional_buy_count"`
	OwnershipStructure      json.RawMessage `json:"ownership_structure,omitempty"`
	SectorCode              string          `json:"sector_code,omitempty"`
	BoardCode               string          `json:"board_code,omitempty"`
	SecurityName            string          `json:"security_name,omitempty"`
	LocalTicker             string          `json:"local_ticker,omitempty"`
	Ticker                  string          `json:"ticker,omitempty"`
	ISIN                    string          `json:"isin"`
	TSETMCCode              string          `json:"tsetmc_code,omitempty"`
	ExchangeCode            string          `json:"exchange_code,omitempty"`
	CompanyCode             string          `json:"company_code,omitempty"`
	LocalSymbolLong         string          `json:"local_symbol_long,omitempty"`
	IssuerISIN              string          `json:"issuer_isin,omitempty"`
	MarketSegment           string          `json:"market_segment,omitempty"`
	IndustryCode            string          `json:"industry_code,omitempty"`
	AsOf                    *time.Time      `json:"as_of,omitempty"`
	TradingStatus           string          `json:"trading_status,omitempty"`
	SyncedAt                *time.Time      `json:"synced_at,omitempty"`
	RawPayload              json.RawMessage `json:"raw_payload,omitempty"`
}

// Stock is an equity instrument with stock-only fundamentals.
type Stock struct {
	QuoteFields
	SharesOutstanding *int64   `json:"shares_outstanding,omitempty"`
	FreeFloat         *float64 `json:"free_float,omitempty"`
	PERatio           *float64 `json:"pe_ratio,omitempty"`
	EPS               *float64 `json:"eps,omitempty"`
	Address           string   `json:"address,omitempty"`
	Website           string   `json:"website,omitempty"`
	FiscalYearEnd     string   `json:"fiscal_year_end,omitempty"`
}

// ETF is an exchange-traded fund with fund-specific fields.
type ETF struct {
	QuoteFields
	AuthorizedUnits  *int64   `json:"authorized_units,omitempty"`
	UnitsOutstanding *int64   `json:"units_outstanding,omitempty"`
	RedemptionNAV    *float64 `json:"redemption_nav,omitempty"`
}

type InstrumentType string

const (
	InstrumentTypeStock InstrumentType = "stock"
	InstrumentTypeETF   InstrumentType = "etf"
	InstrumentTypeAll   InstrumentType = "all"
)

type SyncStatus string

const (
	SyncStatusRunning   SyncStatus = "running"
	SyncStatusSuccess   SyncStatus = "success"
	SyncStatusFailed    SyncStatus = "failed"
)

// SyncRun records one sync job execution.
type SyncRun struct {
	ID             int64          `json:"id"`
	InstrumentType InstrumentType `json:"instrument_type"`
	StartedAt      time.Time      `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at,omitempty"`
	Status         SyncStatus     `json:"status"`
	Error          string         `json:"error,omitempty"`
	FetchedCount   int            `json:"fetched_count"`
}

// ListParams controls pagination and optional ticker filter.
type ListParams struct {
	Page   int
	Limit  int
	Ticker string
}

func (p ListParams) Offset() int {
	page := p.Page
	if page < 1 {
		page = 1
	}
	return (page - 1) * p.NormalizedLimit()
}

func (p ListParams) NormalizedLimit() int {
	if p.Limit < 1 {
		return 50
	}
	if p.Limit > 200 {
		return 200
	}
	return p.Limit
}

func (p ListParams) NormalizedPage() int {
	if p.Page < 1 {
		return 1
	}
	return p.Page
}
