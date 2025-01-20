package app

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Closer interface {
	Close(ctx context.Context) error
}

// Graceful implements graceful shutdown for app.
type Graceful struct {
	closers []Closer
}

func NewGraceful(closers ...Closer) *Graceful {
	return &Graceful{closers: closers}
}

// Shutdown close closers concurrently
func (g *Graceful) Shutdown(ctx context.Context) error {
	eg, egCtx := errgroup.WithContext(ctx)
	for _, closer := range g.closers {
		eg.Go(func() error {
			return closer.Close(egCtx)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
