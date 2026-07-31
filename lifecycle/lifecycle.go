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

type state string

const (
	stateStartup  state = "STARTUP"
	stateWaiting  state = "WAITING"
	stateRunning  state = "RUNNING"
	stateShutdown state = "SHUTDOWN"
	stateStopped  state = "STOPPED"
)

type Coordinator struct {
	ctx    context.Context
	cancel context.CancelFunc

	mu       sync.Mutex
	state    state
	startup  sync.WaitGroup
	shutdown []func(context.Context)
	stopErr  error
	done     chan struct{}
}

func New(ctx context.Context) *Coordinator {
	cctx, cancel := context.WithCancel(ctx)
	return &Coordinator{
		ctx:    cctx,
		cancel: cancel,
		state:  stateStartup,
		done:   make(chan struct{}),
	}
}

func (c *Coordinator) Context() context.Context {
	return c.ctx
}

func (c *Coordinator) OnStartup(fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state != stateStartup {
		panic("lifecycle: OnStartup after WaitForStartup")
	}
	c.startup.Go(fn)
}

func (c *Coordinator) OnShutdown(fn func(context.Context)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.state == stateShutdown || c.state == stateStopped {
		panic("lifecycle: OnShutdown after Shutdown")
	}
	c.shutdown = append(c.shutdown, fn)
}

func (c *Coordinator) WaitForStartup() {
	c.mu.Lock()
	if c.state == stateStartup {
		c.state = stateWaiting
	}
	c.mu.Unlock()

	c.startup.Wait()

	c.mu.Lock()
	if c.state == stateWaiting {
		c.state = stateRunning
	}
	c.mu.Unlock()
}

func (c *Coordinator) Ready() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state == stateRunning
}

func (c *Coordinator) Shutdown(timeout time.Duration) error {
	c.mu.Lock()
	if c.state == stateShutdown || c.state == stateStopped {
		c.mu.Unlock()
		<-c.done
		return c.stopErr
	}
	c.state = stateShutdown
	hooks := make([]func(context.Context), len(c.shutdown))
	copy(hooks, c.shutdown)
	c.mu.Unlock()

	c.cancel()

	var err error
	if len(hooks) > 0 {
		drainCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var wg sync.WaitGroup
		for _, fn := range hooks {
			wg.Go(func() { fn(drainCtx) })
		}

		drained := make(chan struct{})
		go func() {
			wg.Wait()
			close(drained)
		}()

		select {
		case <-drained:
		case <-drainCtx.Done():
			select {
			case <-drained:
			default:
				err = fmt.Errorf("shutdown timeout after %v: %w", timeout, drainCtx.Err())
			}
		}
	}

	c.mu.Lock()
	c.state = stateStopped
	c.stopErr = err
	c.mu.Unlock()
	close(c.done)

	return err
}
