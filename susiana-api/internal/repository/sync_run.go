package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ArminDashti/susiana-api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncRunRepository struct {
	pool *pgxpool.Pool
}

func NewSyncRunRepository(pool *pgxpool.Pool) *SyncRunRepository {
	return &SyncRunRepository{pool: pool}
}

func (r *SyncRunRepository) Start(ctx context.Context, instrumentType domain.InstrumentType) (*domain.SyncRun, error) {
	const q = `
INSERT INTO sync_runs (instrument_type, started_at, status, fetched_count)
VALUES ($1, NOW(), $2, 0)
RETURNING id, instrument_type, started_at, finished_at, status, COALESCE(error, ''), fetched_count`

	var run domain.SyncRun
	err := r.pool.QueryRow(ctx, q, instrumentType, domain.SyncStatusRunning).Scan(
		&run.ID, &run.InstrumentType, &run.StartedAt, &run.FinishedAt, &run.Status, &run.Error, &run.FetchedCount,
	)
	if err != nil {
		return nil, fmt.Errorf("start sync run: %w", err)
	}
	return &run, nil
}

func (r *SyncRunRepository) Finish(ctx context.Context, id int64, status domain.SyncStatus, fetchedCount int, errMsg string) error {
	const q = `
UPDATE sync_runs
SET finished_at = $2, status = $3, fetched_count = $4, error = NULLIF($5, '')
WHERE id = $1`

	_, err := r.pool.Exec(ctx, q, id, time.Now().UTC(), status, fetchedCount, errMsg)
	if err != nil {
		return fmt.Errorf("finish sync run: %w", err)
	}
	return nil
}

func (r *SyncRunRepository) ListRecent(ctx context.Context, limit int) ([]domain.SyncRun, error) {
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	const q = `
SELECT id, instrument_type, started_at, finished_at, status, COALESCE(error, ''), fetched_count
FROM sync_runs
ORDER BY started_at DESC
LIMIT $1`

	rows, err := r.pool.Query(ctx, q, limit)
	if err != nil {
		return nil, fmt.Errorf("list sync runs: %w", err)
	}
	defer rows.Close()

	var out []domain.SyncRun
	for rows.Next() {
		var run domain.SyncRun
		if err := rows.Scan(&run.ID, &run.InstrumentType, &run.StartedAt, &run.FinishedAt, &run.Status, &run.Error, &run.FetchedCount); err != nil {
			return nil, err
		}
		out = append(out, run)
	}
	return out, rows.Err()
}
