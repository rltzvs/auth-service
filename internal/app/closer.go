package app

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type Closer struct {
	mu    sync.Mutex
	funcs []func(ctx context.Context) error
}

func NewCloser() *Closer {
	return &Closer{}
}

func (c *Closer) Add(f func(ctx context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	funcs := c.funcs
	c.mu.Unlock()

	g, ctx := errgroup.WithContext(ctx)
	for i := len(funcs) - 1; i >= 0; i-- {
		f := funcs[i]
		g.Go(func() error {
			return f(ctx)
		})
	}

	if err := g.Wait(); err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.Is(ctx.Err(), context.Canceled) {
			return fmt.Errorf("shutdown timeout: %w", ctx.Err())
		}
		return fmt.Errorf("shutdown completed with errors: %w", err)
	}

	return nil
}
