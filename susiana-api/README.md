# susiana-api

Go (Gin) + PostgreSQL service that syncs **Stock** and **ETF** market data via a pluggable REST client, upserts the latest quote, and exposes read/sync HTTP endpoints.

## Stack

- Go + [Gin](https://github.com/gin-gonic/gin)
- PostgreSQL via [pgx](https://github.com/jackc/pgx)
- Optional interval sync via [robfig/cron](https://github.com/robfig/cron)

## Setup

1. Copy env file and edit values:

```bash
cp .env.example .env
```

2. Ensure PostgreSQL is running and `DATABASE_URL` points at an empty database (schema is applied on startup).

3. Install dependencies (from repo root):

```bash
go mod tidy
```

4. Run the API:

```bash
go run ./cmd/api
```

Health check: `GET http://localhost:8080/health`

## Configuration

| Variable | Description |
|----------|-------------|
| `HTTP_PORT` | Listen port (default `8080`) |
| `DATABASE_URL` | PostgreSQL connection string (**required**) |
| `MARKET_PROVIDER` | `stub` (default) or `http` |
| `MARKET_API_BASE_URL` | Required when `MARKET_PROVIDER=http` |
| `MARKET_API_KEY` | Optional Bearer token for the market API |
| `SYNC_INTERVAL` | e.g. `5m`; empty = manual sync only |
| `GIN_MODE` | `debug` / `release` |

## HTTP API

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/health` | Liveness |
| `GET` | `/api/v1/stocks` | List stocks (`page`, `limit`, `ticker`) |
| `GET` | `/api/v1/stocks/:isin` | Get stock by ISIN |
| `GET` | `/api/v1/etfs` | List ETFs |
| `GET` | `/api/v1/etfs/:isin` | Get ETF by ISIN |
| `POST` | `/api/v1/sync?type=stock\|etf\|all` | Trigger sync |
| `GET` | `/api/v1/sync/runs` | Recent sync runs |

Errors use `{ "error": "..." }`.

## Market data provider (API later)

Until the real REST API is available, keep `MARKET_PROVIDER=stub`. Sync will fail with a clear message.

When you have the API:

1. Set `MARKET_API_BASE_URL` (and `MARKET_API_KEY` if needed).
2. Update request paths and JSON → domain mapping in [`internal/provider/http.go`](internal/provider/http.go) (placeholder paths: `/stocks`, `/etfs`).
3. Set `MARKET_PROVIDER=http`.

Domain field names follow standard market-data terms (`isin`, `ticker`, `last`, `close`, `open`, `volume`, `turnover`, `tsetmc_code`, etc.). See SQL comments in [`migrations/001_init.sql`](migrations/001_init.sql) for Persian labels.

## Project layout

```
cmd/api/                 entrypoint
internal/config/         env config
internal/domain/         models + MarketDataProvider interface
internal/provider/       stub + HTTP client
internal/repository/     PostgreSQL access
internal/service/        sync + query logic
internal/handler/        Gin handlers
internal/scheduler/      optional cron sync
internal/database/       pool + embedded schema migrate
migrations/              SQL reference schema
```
