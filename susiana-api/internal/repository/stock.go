package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ArminDashti/susiana-api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type StockRepository struct {
	pool *pgxpool.Pool
}

func NewStockRepository(pool *pgxpool.Pool) *StockRepository {
	return &StockRepository{pool: pool}
}

func (r *StockRepository) UpsertMany(ctx context.Context, stocks []domain.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin stock upsert: %w", err)
	}
	defer tx.Rollback(ctx)

	const q = `
INSERT INTO stocks (
    isin, tsetmc_code, last, close, open, lower_limit, upper_limit, base_volume,
    volume, turnover, trades, retail_buy_volume, institutional_buy_volume,
    retail_buy_count, institutional_buy_count, ownership_structure, sector_code,
    board_code, security_name, local_ticker, ticker, exchange_code, company_code,
    local_symbol_long, issuer_isin, market_segment, industry_code, as_of,
    trading_status, shares_outstanding, free_float, pe_ratio, eps, address,
    website, fiscal_year_end, synced_at, raw_payload
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
    $21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34,$35,$36,$37,$38
)
ON CONFLICT (isin) DO UPDATE SET
    tsetmc_code = EXCLUDED.tsetmc_code,
    last = EXCLUDED.last,
    close = EXCLUDED.close,
    open = EXCLUDED.open,
    lower_limit = EXCLUDED.lower_limit,
    upper_limit = EXCLUDED.upper_limit,
    base_volume = EXCLUDED.base_volume,
    volume = EXCLUDED.volume,
    turnover = EXCLUDED.turnover,
    trades = EXCLUDED.trades,
    retail_buy_volume = EXCLUDED.retail_buy_volume,
    institutional_buy_volume = EXCLUDED.institutional_buy_volume,
    retail_buy_count = EXCLUDED.retail_buy_count,
    institutional_buy_count = EXCLUDED.institutional_buy_count,
    ownership_structure = EXCLUDED.ownership_structure,
    sector_code = EXCLUDED.sector_code,
    board_code = EXCLUDED.board_code,
    security_name = EXCLUDED.security_name,
    local_ticker = EXCLUDED.local_ticker,
    ticker = EXCLUDED.ticker,
    exchange_code = EXCLUDED.exchange_code,
    company_code = EXCLUDED.company_code,
    local_symbol_long = EXCLUDED.local_symbol_long,
    issuer_isin = EXCLUDED.issuer_isin,
    market_segment = EXCLUDED.market_segment,
    industry_code = EXCLUDED.industry_code,
    as_of = EXCLUDED.as_of,
    trading_status = EXCLUDED.trading_status,
    shares_outstanding = EXCLUDED.shares_outstanding,
    free_float = EXCLUDED.free_float,
    pe_ratio = EXCLUDED.pe_ratio,
    eps = EXCLUDED.eps,
    address = EXCLUDED.address,
    website = EXCLUDED.website,
    fiscal_year_end = EXCLUDED.fiscal_year_end,
    synced_at = EXCLUDED.synced_at,
    raw_payload = EXCLUDED.raw_payload
`

	now := time.Now().UTC()
	for _, s := range stocks {
		if s.ISIN == "" {
			return fmt.Errorf("stock missing isin")
		}
		ownership := nullJSON(s.OwnershipStructure)
		raw := nullJSON(s.RawPayload)
		_, err := tx.Exec(ctx, q,
			s.ISIN, nullString(s.TSETMCCode), s.Last, s.Close, s.Open, s.LowerLimit, s.UpperLimit, s.BaseVolume,
			s.Volume, s.Turnover, s.Trades, s.RetailBuyVolume, s.InstitutionalBuyVolume,
			s.RetailBuyCount, s.InstitutionalBuyCount, ownership, nullString(s.SectorCode),
			nullString(s.BoardCode), nullString(s.SecurityName), nullString(s.LocalTicker), nullString(s.Ticker),
			nullString(s.ExchangeCode), nullString(s.CompanyCode), nullString(s.LocalSymbolLong),
			nullString(s.IssuerISIN), nullString(s.MarketSegment), nullString(s.IndustryCode), s.AsOf,
			nullString(s.TradingStatus), s.SharesOutstanding, s.FreeFloat, s.PERatio, s.EPS,
			nullString(s.Address), nullString(s.Website), nullString(s.FiscalYearEnd), now, raw,
		)
		if err != nil {
			return fmt.Errorf("upsert stock %s: %w", s.ISIN, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit stock upsert: %w", err)
	}
	return nil
}

func (r *StockRepository) GetByISIN(ctx context.Context, isin string) (*domain.Stock, error) {
	const q = `
SELECT isin, tsetmc_code, last, close, open, lower_limit, upper_limit, base_volume,
       volume, turnover, trades, retail_buy_volume, institutional_buy_volume,
       retail_buy_count, institutional_buy_count, ownership_structure, sector_code,
       board_code, security_name, local_ticker, ticker, exchange_code, company_code,
       local_symbol_long, issuer_isin, market_segment, industry_code, as_of,
       trading_status, shares_outstanding, free_float, pe_ratio, eps, address,
       website, fiscal_year_end, synced_at, raw_payload
FROM stocks WHERE isin = $1`

	row := r.pool.QueryRow(ctx, q, isin)
	s, err := scanStock(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s, nil
}

func (r *StockRepository) List(ctx context.Context, params domain.ListParams) ([]domain.Stock, int, error) {
	limit := params.NormalizedLimit()
	offset := params.Offset()

	countQ := `SELECT COUNT(*) FROM stocks`
	listQ := `
SELECT isin, tsetmc_code, last, close, open, lower_limit, upper_limit, base_volume,
       volume, turnover, trades, retail_buy_volume, institutional_buy_volume,
       retail_buy_count, institutional_buy_count, ownership_structure, sector_code,
       board_code, security_name, local_ticker, ticker, exchange_code, company_code,
       local_symbol_long, issuer_isin, market_segment, industry_code, as_of,
       trading_status, shares_outstanding, free_float, pe_ratio, eps, address,
       website, fiscal_year_end, synced_at, raw_payload
FROM stocks`

	args := []any{}
	if params.Ticker != "" {
		countQ += ` WHERE ticker ILIKE $1`
		listQ += ` WHERE ticker ILIKE $1`
		args = append(args, params.Ticker)
	}
	listQ += fmt.Sprintf(` ORDER BY ticker NULLS LAST, isin LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)

	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count stocks: %w", err)
	}

	listArgs := append(append([]any{}, args...), limit, offset)
	rows, err := r.pool.Query(ctx, listQ, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list stocks: %w", err)
	}
	defer rows.Close()

	var out []domain.Stock
	for rows.Next() {
		s, err := scanStock(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *s)
	}
	return out, total, rows.Err()
}

type scannable interface {
	Scan(dest ...any) error
}

func scanStock(row scannable) (*domain.Stock, error) {
	var s domain.Stock
	var ownership, raw []byte
	var tsetmc, sector, board, name, localTicker, ticker, exchange, company, localLong, issuer, market, industry, status, address, website, fiscal *string

	err := row.Scan(
		&s.ISIN, &tsetmc, &s.Last, &s.Close, &s.Open, &s.LowerLimit, &s.UpperLimit, &s.BaseVolume,
		&s.Volume, &s.Turnover, &s.Trades, &s.RetailBuyVolume, &s.InstitutionalBuyVolume,
		&s.RetailBuyCount, &s.InstitutionalBuyCount, &ownership, &sector,
		&board, &name, &localTicker, &ticker, &exchange, &company,
		&localLong, &issuer, &market, &industry, &s.AsOf,
		&status, &s.SharesOutstanding, &s.FreeFloat, &s.PERatio, &s.EPS, &address,
		&website, &fiscal, &s.SyncedAt, &raw,
	)
	if err != nil {
		return nil, err
	}

	s.TSETMCCode = deref(tsetmc)
	s.SectorCode = deref(sector)
	s.BoardCode = deref(board)
	s.SecurityName = deref(name)
	s.LocalTicker = deref(localTicker)
	s.Ticker = deref(ticker)
	s.ExchangeCode = deref(exchange)
	s.CompanyCode = deref(company)
	s.LocalSymbolLong = deref(localLong)
	s.IssuerISIN = deref(issuer)
	s.MarketSegment = deref(market)
	s.IndustryCode = deref(industry)
	s.TradingStatus = deref(status)
	s.Address = deref(address)
	s.Website = deref(website)
	s.FiscalYearEnd = deref(fiscal)
	if len(ownership) > 0 {
		s.OwnershipStructure = json.RawMessage(ownership)
	}
	if len(raw) > 0 {
		s.RawPayload = json.RawMessage(raw)
	}
	return &s, nil
}
