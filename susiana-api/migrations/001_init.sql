-- +migrate Up
CREATE TABLE IF NOT EXISTS stocks (
    isin                    VARCHAR(12) PRIMARY KEY,
    tsetmc_code             VARCHAR(32),
    last                    DOUBLE PRECISION,
    close                   DOUBLE PRECISION,
    open                    DOUBLE PRECISION,
    lower_limit             DOUBLE PRECISION,
    upper_limit             DOUBLE PRECISION,
    base_volume             BIGINT,
    volume                  BIGINT,
    turnover                DOUBLE PRECISION,
    trades                  BIGINT,
    retail_buy_volume       BIGINT,
    institutional_buy_volume BIGINT,
    retail_buy_count        BIGINT,
    institutional_buy_count BIGINT,
    ownership_structure     JSONB,
    sector_code             VARCHAR(32),
    board_code              VARCHAR(32),
    security_name           TEXT,
    local_ticker            VARCHAR(64),
    ticker                  VARCHAR(64),
    exchange_code           VARCHAR(8),
    company_code            VARCHAR(8),
    local_symbol_long       VARCHAR(64),
    issuer_isin             VARCHAR(12),
    market_segment          VARCHAR(64),
    industry_code           VARCHAR(32),
    as_of                   TIMESTAMPTZ,
    trading_status          VARCHAR(64),
    shares_outstanding      BIGINT,
    free_float              DOUBLE PRECISION,
    pe_ratio                DOUBLE PRECISION,
    eps                     DOUBLE PRECISION,
    address                 TEXT,
    website                 TEXT,
    fiscal_year_end         VARCHAR(32),
    synced_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw_payload             JSONB
);

COMMENT ON COLUMN stocks.last IS 'آخرین معامله — Last traded price (LTP)';
COMMENT ON COLUMN stocks.close IS 'قیمت پایانی — Official closing price';
COMMENT ON COLUMN stocks.open IS 'اولین قیمت — Opening / first trade price';
COMMENT ON COLUMN stocks.lower_limit IS 'قیمت مجاز (حد پایین) — Daily price band lower limit';
COMMENT ON COLUMN stocks.upper_limit IS 'قیمت مجاز (حد بالا) — Daily price band upper limit';
COMMENT ON COLUMN stocks.base_volume IS 'حجم مبنا — Base volume';
COMMENT ON COLUMN stocks.volume IS 'حجم معاملات — Traded volume';
COMMENT ON COLUMN stocks.turnover IS 'ارزش معاملات — Traded value';
COMMENT ON COLUMN stocks.trades IS 'تعداد معاملات — Number of trades';
COMMENT ON COLUMN stocks.retail_buy_volume IS 'حجم خرید حقیقی — Retail buy volume';
COMMENT ON COLUMN stocks.institutional_buy_volume IS 'حجم خرید حقوقی — Institutional buy volume';
COMMENT ON COLUMN stocks.retail_buy_count IS 'تعداد خرید حقیقی — Retail buy trade count';
COMMENT ON COLUMN stocks.institutional_buy_count IS 'تعداد خرید حقوقی — Institutional buy trade count';
COMMENT ON COLUMN stocks.ownership_structure IS 'ترکیب سهامداران — Shareholder composition';
COMMENT ON COLUMN stocks.sector_code IS 'کد گروه صنعت — Sector / industry group code';
COMMENT ON COLUMN stocks.board_code IS 'کد تابلو — Exchange board code';
COMMENT ON COLUMN stocks.security_name IS 'نام سهم — Instrument / security name';
COMMENT ON COLUMN stocks.local_ticker IS 'نماد سهم به فارسی — Local-language ticker';
COMMENT ON COLUMN stocks.ticker IS 'نماد سهم به انگلیسی — Latin ticker symbol';
COMMENT ON COLUMN stocks.isin IS 'کد 12 رقمی نماد — Instrument ISIN';
COMMENT ON COLUMN stocks.tsetmc_code IS 'کد TSETMC — TSETMC instrument ID';
COMMENT ON COLUMN stocks.exchange_code IS 'کد 5 رقمی نماد — Short exchange instrument code';
COMMENT ON COLUMN stocks.company_code IS 'کد 4 رقمی شرکت — Issuer short code';
COMMENT ON COLUMN stocks.local_symbol_long IS 'نماد 30 رقمی فارسی — Long local symbol';
COMMENT ON COLUMN stocks.issuer_isin IS 'کد 12 رقمی شرکت — Issuer / company ISIN';
COMMENT ON COLUMN stocks.market_segment IS 'بازار — Market / segment';
COMMENT ON COLUMN stocks.industry_code IS 'کد زیر گروه صنعت — Sub-industry code';
COMMENT ON COLUMN stocks.as_of IS 'آخرین اطلاعات قیمت — Quote timestamp';
COMMENT ON COLUMN stocks.trading_status IS 'وضعیت — Trading status';
COMMENT ON COLUMN stocks.shares_outstanding IS 'تعداد سهام — Shares outstanding';
COMMENT ON COLUMN stocks.free_float IS 'سهام شناور — Free float';
COMMENT ON COLUMN stocks.pe_ratio IS 'P/E — Price-to-earnings';
COMMENT ON COLUMN stocks.eps IS 'EPS — Earnings per share';
COMMENT ON COLUMN stocks.address IS 'نشانی — Registered address';
COMMENT ON COLUMN stocks.website IS 'وب سایت — Corporate website';
COMMENT ON COLUMN stocks.fiscal_year_end IS 'سال مالی — Fiscal year end';

CREATE UNIQUE INDEX IF NOT EXISTS stocks_tsetmc_code_uidx ON stocks (tsetmc_code) WHERE tsetmc_code IS NOT NULL AND tsetmc_code <> '';
CREATE INDEX IF NOT EXISTS stocks_ticker_idx ON stocks (ticker);

CREATE TABLE IF NOT EXISTS etfs (
    isin                    VARCHAR(12) PRIMARY KEY,
    tsetmc_code             VARCHAR(32),
    last                    DOUBLE PRECISION,
    close                   DOUBLE PRECISION,
    open                    DOUBLE PRECISION,
    lower_limit             DOUBLE PRECISION,
    upper_limit             DOUBLE PRECISION,
    base_volume             BIGINT,
    volume                  BIGINT,
    turnover                DOUBLE PRECISION,
    trades                  BIGINT,
    retail_buy_volume       BIGINT,
    institutional_buy_volume BIGINT,
    retail_buy_count        BIGINT,
    institutional_buy_count BIGINT,
    ownership_structure     JSONB,
    sector_code             VARCHAR(32),
    board_code              VARCHAR(32),
    security_name           TEXT,
    local_ticker            VARCHAR(64),
    ticker                  VARCHAR(64),
    exchange_code           VARCHAR(8),
    company_code            VARCHAR(8),
    local_symbol_long       VARCHAR(64),
    issuer_isin             VARCHAR(12),
    market_segment          VARCHAR(64),
    industry_code           VARCHAR(32),
    as_of                   TIMESTAMPTZ,
    trading_status          VARCHAR(64),
    authorized_units        BIGINT,
    units_outstanding       BIGINT,
    redemption_nav          DOUBLE PRECISION,
    synced_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    raw_payload             JSONB
);

COMMENT ON COLUMN etfs.last IS 'آخرین معامله — Last traded price (LTP)';
COMMENT ON COLUMN etfs.close IS 'قیمت پایانی — Official closing price';
COMMENT ON COLUMN etfs.open IS 'اولین قیمت — Opening / first trade price';
COMMENT ON COLUMN etfs.lower_limit IS 'قیمت مجاز (حد پایین) — Daily price band lower limit';
COMMENT ON COLUMN etfs.upper_limit IS 'قیمت مجاز (حد بالا) — Daily price band upper limit';
COMMENT ON COLUMN etfs.base_volume IS 'حجم مبنا — Base volume';
COMMENT ON COLUMN etfs.volume IS 'حجم معاملات — Traded volume';
COMMENT ON COLUMN etfs.turnover IS 'ارزش معاملات — Traded value';
COMMENT ON COLUMN etfs.trades IS 'تعداد معاملات — Number of trades';
COMMENT ON COLUMN etfs.retail_buy_volume IS 'حجم خرید حقیقی — Retail buy volume';
COMMENT ON COLUMN etfs.institutional_buy_volume IS 'حجم خرید حقوقی — Institutional buy volume';
COMMENT ON COLUMN etfs.retail_buy_count IS 'تعداد خرید حقیقی — Retail buy trade count';
COMMENT ON COLUMN etfs.institutional_buy_count IS 'تعداد خرید حقوقی — Institutional buy trade count';
COMMENT ON COLUMN etfs.ownership_structure IS 'ترکیب سهامداران — Shareholder composition';
COMMENT ON COLUMN etfs.sector_code IS 'کد گروه صنعت — Sector / industry group code';
COMMENT ON COLUMN etfs.board_code IS 'کد تابلو — Exchange board code';
COMMENT ON COLUMN etfs.security_name IS 'نام سهم — Instrument / security name';
COMMENT ON COLUMN etfs.local_ticker IS 'نماد سهم به فارسی — Local-language ticker';
COMMENT ON COLUMN etfs.ticker IS 'نماد سهم به انگلیسی — Latin ticker symbol';
COMMENT ON COLUMN etfs.isin IS 'کد 12 رقمی نماد — Instrument ISIN';
COMMENT ON COLUMN etfs.tsetmc_code IS 'کد TSETMC — TSETMC instrument ID';
COMMENT ON COLUMN etfs.exchange_code IS 'کد 5 رقمی نماد — Short exchange instrument code';
COMMENT ON COLUMN etfs.company_code IS 'کد 4 رقمی شرکت — Issuer short code';
COMMENT ON COLUMN etfs.local_symbol_long IS 'نماد 30 رقمی فارسی — Long local symbol';
COMMENT ON COLUMN etfs.issuer_isin IS 'کد 12 رقمی شرکت — Issuer / company ISIN';
COMMENT ON COLUMN etfs.market_segment IS 'بازار — Market / segment';
COMMENT ON COLUMN etfs.industry_code IS 'کد زیر گروه صنعت — Sub-industry code';
COMMENT ON COLUMN etfs.as_of IS 'آخرین اطلاعات قیمت — Quote timestamp';
COMMENT ON COLUMN etfs.trading_status IS 'وضعیت — Trading status';
COMMENT ON COLUMN etfs.authorized_units IS 'سقف واحدهای صندوق — Fund unit cap / authorized units';
COMMENT ON COLUMN etfs.units_outstanding IS 'تعداد واحدهای صادر شده — Issued / outstanding units';
COMMENT ON COLUMN etfs.redemption_nav IS 'NAV ابطال — Redemption NAV';

CREATE UNIQUE INDEX IF NOT EXISTS etfs_tsetmc_code_uidx ON etfs (tsetmc_code) WHERE tsetmc_code IS NOT NULL AND tsetmc_code <> '';
CREATE INDEX IF NOT EXISTS etfs_ticker_idx ON etfs (ticker);

CREATE TABLE IF NOT EXISTS sync_runs (
    id               BIGSERIAL PRIMARY KEY,
    instrument_type  VARCHAR(16) NOT NULL,
    started_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at      TIMESTAMPTZ,
    status           VARCHAR(16) NOT NULL,
    error            TEXT,
    fetched_count    INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS sync_runs_started_at_idx ON sync_runs (started_at DESC);
