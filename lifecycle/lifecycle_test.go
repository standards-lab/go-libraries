package lifecycle_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/standards-lab/go-libraries/lifecycle"
)

var _ lifecycle.ReadinessChecker = (*lifecycle.Coordinator)(nil)

// failsafe bounds every wait for an event that should occur, so a broken
// coordinator fails the test instead of hanging it.
const failsafe = 2 * time.Second

// recvOrFail receives from ch, or fails the test if nothing arrives within the
// failsafe window.
func recvOrFail[T any](t *testing.T, ch <-chan T, what string) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(failsafe):
		t.Fatalf("timed out waiting for %s", what)
		var zero T
		return zero
	}
}

func TestOnStartup_RunConcurrently(t *testing.T) {
	lc := lifecycle.New(context.Background())

	const n = 5
	arrived := make(chan struct{}, n)
	release := make(chan struct{})
	var count atomic.Int64

	for range n {
		lc.OnStartup(func() {
			count.Add(1)
			arrived <- struct{}{}
			<-release
		})
	}

	// Every hook must reach its arrival send before any is released, which only
	// holds if they run simultaneously rather than one after another.
	for range n {
		recvOrFail(t, arrived, "startup hook arrival")
	}
	close(release)

	lc.WaitForStartup()

	if got := count.Load(); got != n {
		t.Fatalf("ran %d startup hooks, want %d", got, n)
	}
}

func TestReady_FlipsAfterStartup(t *testing.T) {
	lc := lifecycle.New(context.Background())

	started := make(chan struct{})
	release := make(chan struct{})
	lc.OnStartup(func() {
		close(started)
		<-release
	})

	recvOrFail(t, started, "startup hook to start")
	if lc.Ready() {
		t.Fatal("Ready() is true before WaitForStartup returned")
	}

	waited := make(chan struct{})
	go func() {
		lc.WaitForStartup()
		close(waited)
	}()

	// WaitForStartup must block while the hook is held.
	select {
	case <-waited:
		t.Fatal("WaitForStartup returned before its startup hook completed")
	case <-time.After(50 * time.Millisecond):
	}
	if lc.Ready() {
		t.Fatal("Ready() is true while a startup hook is still running")
	}

	close(release)
	recvOrFail(t, waited, "WaitForStartup to return")
	if !lc.Ready() {
		t.Fatal("Ready() is false after WaitForStartup returned")
	}
}

func TestShutdown_InvokesHooksAfterCancel(t *testing.T) {
	lc := lifecycle.New(context.Background())

	type observation struct{ rootErr, drainErr error }
	obs := make(chan observation, 1)
	lc.OnShutdown(func(drainCtx context.Context) {
		obs <- observation{rootErr: lc.Context().Err(), drainErr: drainCtx.Err()}
	})

	if err := lc.Shutdown(failsafe); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}

	o := recvOrFail(t, obs, "shutdown hook invocation")
	if o.rootErr == nil {
		t.Error("root context was not cancelled when the shutdown hook ran")
	}
	if o.drainErr != nil {
		t.Errorf("drain context was already cancelled when the hook ran: %v", o.drainErr)
	}
}

func TestShutdown_RunsAllHooks(t *testing.T) {
	lc := lifecycle.New(context.Background())

	const n = 5
	var count atomic.Int64
	for range n {
		lc.OnShutdown(func(context.Context) { count.Add(1) })
	}

	if err := lc.Shutdown(failsafe); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}
	if got := count.Load(); got != n {
		t.Fatalf("ran %d shutdown hooks, want %d", got, n)
	}
}

func TestShutdown_TimesOut(t *testing.T) {
	lc := lifecycle.New(context.Background())

	// The hook outlives the timeout; releasing it only at cleanup keeps the
	// hooks-done path closed, so Shutdown must return via the deadline.
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	lc.OnShutdown(func(context.Context) {
		<-release
	})

	if err := lc.Shutdown(20 * time.Millisecond); err == nil {
		t.Fatal("Shutdown returned nil for a hook that outlived the timeout")
	}
}

func TestReady_FalseAfterShutdown(t *testing.T) {
	lc := lifecycle.New(context.Background())

	lc.WaitForStartup()
	if !lc.Ready() {
		t.Fatal("Ready() is false after WaitForStartup with no hooks registered")
	}

	if err := lc.Shutdown(failsafe); err != nil {
		t.Fatalf("Shutdown returned error: %v", err)
	}
	if lc.Ready() {
		t.Fatal("Ready() is true after Shutdown")
	}
}

func TestNew_DerivesFromParentContext(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	lc := lifecycle.New(parent)

	select {
	case <-lc.Context().Done():
		t.Fatal("coordinator context is done before the parent was cancelled")
	default:
	}

	cancelParent()

	select {
	case <-lc.Context().Done():
	case <-time.After(failsafe):
		t.Fatal("coordinator context was not cancelled after the parent was cancelled")
	}

	// An explicit Shutdown still drains cleanly after a parent-driven cancel.
	if err := lc.Shutdown(failsafe); err != nil {
		t.Fatalf("Shutdown after parent cancellation returned error: %v", err)
	}
}
