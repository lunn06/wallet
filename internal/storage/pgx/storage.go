package pgx

import (
	"context"
	"log/slog"
	"runtime"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lunn06/wallet/pkg/semaphore"

	_ "github.com/lib/pq"
)

var (
	maxCons = runtime.NumCPU() - runtime.NumCPU()/4
)

// Storage is pgxpool.Pool wrapper (that need to embed in other storages)
// access to which is regulated by semaphore.Semaphore
type Storage struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
	sem    *semaphore.Semaphore
}

func NewStorage(dns string, logger *slog.Logger) (*Storage, error) {
	cfg, err := pgxpool.ParseConfig(dns)
	if err != nil {
		return nil, err
	}

	// Register pgxdecimal to handle postgres NUMERIC via decimal.Decimal
	cfg.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}

	cfg.MaxConns = int32(maxCons)

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	sem := semaphore.New(maxCons)
	logger.Info("PgxStorage created")

	return &Storage{pool: db, logger: logger, sem: sem}, nil
}

// Do method make access to inner pgxpool.Pool via semaphore
func (s *Storage) Do(f func(db *pgxpool.Pool) error) error {
	resultCh := make(chan error)

	go func() {
		s.sem.Acquire()
		defer s.sem.Release()

		result := f(s.pool)
		resultCh <- result
	}()

	return <-resultCh
}

func (s *Storage) Close(ctx context.Context) error {
	done := make(chan struct{}, 1)
	go func() {
		s.pool.Close()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
