package lifecycle

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ReadinessChecker interface {
	Ready() bool
}

type Coordinator struct {
	ctx      context.Context
	cancel   context.CancelFunc
	startup  sync.WaitGroup
	mu       sync.Mutex
	shutdown []func(context.Context)
	ready    bool
}

func New(ctx context.Context) *Coordinator {
	cctx, cancel := context.WithCancel(ctx)
	return &Coordinator{
		ctx:    cctx,
		cancel: cancel,
	}
}

func (c *Coordinator) Context() context.Context {
	return c.ctx
}

func (c *Coordinator) OnStartup(fn func()) {
	c.startup.Go(fn)
}

func (c *Coordinator) OnShutdown(fn func(context.Context)) {
	c.mu.Lock()
	c.shutdown = append(c.shutdown, fn)
	c.mu.Unlock()
}

func (c *Coordinator) WaitForStartup() {
	c.startup.Wait()
	c.mu.Lock()
	c.ready = true
	c.mu.Unlock()
}

func (c *Coordinator) Ready() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ready
}

func (c *Coordinator) Shutdown(timeout time.Duration) error {
	c.mu.Lock()
	c.ready = false
	hooks := make([]func(context.Context), len(c.shutdown))
	copy(hooks, c.shutdown)
	c.mu.Unlock()

	c.cancel()

	drainCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, fn := range hooks {
		wg.Go(func() { fn(drainCtx) })
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-drainCtx.Done():
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}
