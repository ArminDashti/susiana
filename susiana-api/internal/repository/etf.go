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

type ETFRepository struct {
	pool *pgxpool.Pool
}

func NewETFRepository(pool *pgxpool.Pool) *ETFRepository {
	return &ETFRepository{pool: pool}
}

func (r *ETFRepository) UpsertMany(ctx context.Context, etfs []domain.ETF) error {
	if len(etfs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin etf upsert: %w", err)
	}
	defer tx.Rollback(ctx)

	const q = `
INSERT INTO etfs (
    isin, tsetmc_code, last, close, open, lower_limit, upper_limit, base_volume,
    volume, turnover, trades, retail_buy_volume, institutional_buy_volume,
    retail_buy_count, institutional_buy_count, ownership_structure, sector_code,
    board_code, security_name, local_ticker, ticker, exchange_code, company_code,
    local_symbol_long, issuer_isin, market_segment, industry_code, as_of,
    trading_status, authorized_units, units_outstanding, redemption_nav,
    synced_at, raw_payload
) VALUES (
    $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,
    $21,$22,$23,$24,$25,$26,$27,$28,$29,$30,$31,$32,$33,$34
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
    authorized_units = EXCLUDED.authorized_units,
    units_outstanding = EXCLUDED.units_outstanding,
    redemption_nav = EXCLUDED.redemption_nav,
    synced_at = EXCLUDED.synced_at,
    raw_payload = EXCLUDED.raw_payload
`

	now := time.Now().UTC()
	for _, e := range etfs {
		if e.ISIN == "" {
			return fmt.Errorf("etf missing isin")
		}
		ownership := nullJSON(e.OwnershipStructure)
		raw := nullJSON(e.RawPayload)
		_, err := tx.Exec(ctx, q,
			e.ISIN, nullString(e.TSETMCCode), e.Last, e.Close, e.Open, e.LowerLimit, e.UpperLimit, e.BaseVolume,
			e.Volume, e.Turnover, e.Trades, e.RetailBuyVolume, e.InstitutionalBuyVolume,
			e.RetailBuyCount, e.InstitutionalBuyCount, ownership, nullString(e.SectorCode),
			nullString(e.BoardCode), nullString(e.SecurityName), nullString(e.LocalTicker), nullString(e.Ticker),
			nullString(e.ExchangeCode), nullString(e.CompanyCode), nullString(e.LocalSymbolLong),
			nullString(e.IssuerISIN), nullString(e.MarketSegment), nullString(e.IndustryCode), e.AsOf,
			nullString(e.TradingStatus), e.AuthorizedUnits, e.UnitsOutstanding, e.RedemptionNAV, now, raw,
		)
		if err != nil {
			return fmt.Errorf("upsert etf %s: %w", e.ISIN, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit etf upsert: %w", err)
	}
	return nil
}

func (r *ETFRepository) GetByISIN(ctx context.Context, isin string) (*domain.ETF, error) {
	const q = `
SELECT isin, tsetmc_code, last, close, open, lower_limit, upper_limit, base_volume,
       volume, turnover, trades, retail_buy_volume, institutional_buy_volume,
       retail_buy_count, institutional_buy_count, ownership_structure, sector_code,
       board_code, security_name, local_ticker, ticker, exchange_code, company_code,
       local_symbol_long, issuer_isin, market_segment, industry_code, as_of,
       trading_status, authorized_units, units_outstanding, redemption_nav,
       synced_at, raw_payload
FROM etfs WHERE isin = $1`

	row := r.pool.QueryRow(ctx, q, isin)
	e, err := scanETF(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return e, nil
}

func (r *ETFRepository) List(ctx context.Context, params domain.ListParams) ([]domain.ETF, int, error) {
	limit := params.NormalizedLimit()
	offset := params.Offset()

	countQ := `SELECT COUNT(*) FROM etfs`
	listQ := `
SELECT isin, tsetmc_code, last, close, open, lower_limit, upper_limit, base_volume,
       volume, turnover, trades, retail_buy_volume, institutional_buy_volume,
       retail_buy_count, institutional_buy_count, ownership_structure, sector_code,
       board_code, security_name, local_ticker, ticker, exchange_code, company_code,
       local_symbol_long, issuer_isin, market_segment, industry_code, as_of,
       trading_status, authorized_units, units_outstanding, redemption_nav,
       synced_at, raw_payload
FROM etfs`

	args := []any{}
	if params.Ticker != "" {
		countQ += ` WHERE ticker ILIKE $1`
		listQ += ` WHERE ticker ILIKE $1`
		args = append(args, params.Ticker)
	}
	listQ += fmt.Sprintf(` ORDER BY ticker NULLS LAST, isin LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)

	var total int
	if err := r.pool.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count etfs: %w", err)
	}

	listArgs := append(append([]any{}, args...), limit, offset)
	rows, err := r.pool.Query(ctx, listQ, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list etfs: %w", err)
	}
	defer rows.Close()

	var out []domain.ETF
	for rows.Next() {
		e, err := scanETF(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *e)
	}
	return out, total, rows.Err()
}

func scanETF(row scannable) (*domain.ETF, error) {
	var e domain.ETF
	var ownership, raw []byte
	var tsetmc, sector, board, name, localTicker, ticker, exchange, company, localLong, issuer, market, industry, status *string

	err := row.Scan(
		&e.ISIN, &tsetmc, &e.Last, &e.Close, &e.Open, &e.LowerLimit, &e.UpperLimit, &e.BaseVolume,
		&e.Volume, &e.Turnover, &e.Trades, &e.RetailBuyVolume, &e.InstitutionalBuyVolume,
		&e.RetailBuyCount, &e.InstitutionalBuyCount, &ownership, &sector,
		&board, &name, &localTicker, &ticker, &exchange, &company,
		&localLong, &issuer, &market, &industry, &e.AsOf,
		&status, &e.AuthorizedUnits, &e.UnitsOutstanding, &e.RedemptionNAV,
		&e.SyncedAt, &raw,
	)
	if err != nil {
		return nil, err
	}

	e.TSETMCCode = deref(tsetmc)
	e.SectorCode = deref(sector)
	e.BoardCode = deref(board)
	e.SecurityName = deref(name)
	e.LocalTicker = deref(localTicker)
	e.Ticker = deref(ticker)
	e.ExchangeCode = deref(exchange)
	e.CompanyCode = deref(company)
	e.LocalSymbolLong = deref(localLong)
	e.IssuerISIN = deref(issuer)
	e.MarketSegment = deref(market)
	e.IndustryCode = deref(industry)
	e.TradingStatus = deref(status)
	if len(ownership) > 0 {
		e.OwnershipStructure = json.RawMessage(ownership)
	}
	if len(raw) > 0 {
		e.RawPayload = json.RawMessage(raw)
	}
	return &e, nil
}
