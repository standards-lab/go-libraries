package web_test

import (
	"context"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/standards-lab/go-libraries/web"
)

// failsafe bounds every wait for an event that should occur, so a broken server
// fails the test instead of hanging it.
const failsafe = 2 * time.Second

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

// startTestServer binds a server on an ephemeral port. Port 0 survives because
// NewServer takes the configuration as written — Finalize, which would replace
// it with the default, belongs to the load path.
func startTestServer(t *testing.T, handler http.Handler) *web.Server {
	t.Helper()

	srv := web.NewServer(web.Config{Host: "127.0.0.1"}, handler)
	if err := srv.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), failsafe)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})
	return srv
}

func TestServer_AddrBeforeStartIsConfigured(t *testing.T) {
	srv := web.NewServer(web.Config{Host: "127.0.0.1", Port: 8080}, http.NewServeMux())
	if got, want := srv.Addr(), "127.0.0.1:8080"; got != want {
		t.Errorf("Addr() = %q, want %q", got, want)
	}
}

func TestServer_StartServesRequests(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, _ *http.Request) {
		_ = web.WriteJSON(w, http.StatusOK, map[string]string{"status": "pong"})
	})

	srv := startTestServer(t, mux)

	if _, port, err := net.SplitHostPort(srv.Addr()); err != nil {
		t.Fatalf("Addr() = %q, want host:port: %v", srv.Addr(), err)
	} else if p, _ := strconv.Atoi(port); p == 0 {
		t.Fatal("Addr() still reports port 0 after Start bound the listener")
	}

	resp, err := http.Get("http://" + srv.Addr() + "/ping")
	if err != nil {
		t.Fatalf("GET /ping: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestServer_StartTwiceFails(t *testing.T) {
	srv := startTestServer(t, http.NewServeMux())

	if err := srv.Start(); err == nil {
		t.Fatal("the second Start returned nil")
	}
}

// TestServer_StartReportsBindFailure is the defect this package exists to fix:
// binding on the calling goroutine means a taken port fails the caller instead
// of being logged by a goroutine nobody is watching.
func TestServer_StartReportsBindFailure(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("occupy a port: %v", err)
	}
	defer func() { _ = occupied.Close() }()

	host, port, err := net.SplitHostPort(occupied.Addr().String())
	if err != nil {
		t.Fatalf("split %q: %v", occupied.Addr(), err)
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		t.Fatalf("parse port %q: %v", port, err)
	}

	srv := web.NewServer(web.Config{Host: host, Port: p}, http.NewServeMux())
	err = srv.Start()
	if err == nil {
		t.Fatal("Start returned nil for an occupied port")
	}
	if !strings.Contains(err.Error(), "listen") {
		t.Errorf("error = %v, want it to mention listen", err)
	}
}

func TestServer_ShutdownWaitsForInFlightRequest(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		<-release
		_ = web.WriteJSON(w, http.StatusOK, map[string]string{"status": "done"})
	})

	srv := startTestServer(t, mux)

	type response struct {
		status int
		err    error
	}
	done := make(chan response, 1)
	go func() {
		resp, err := http.Get("http://" + srv.Addr() + "/slow")
		if err != nil {
			done <- response{err: err}
			return
		}
		defer func() { _ = resp.Body.Close() }()
		_, _ = io.Copy(io.Discard, resp.Body)
		done <- response{status: resp.StatusCode}
	}()

	recvOrFail(t, started, "the handler to start")

	shutdown := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), failsafe)
		defer cancel()
		shutdown <- srv.Shutdown(ctx)
	}()

	// Shutdown must not return while the request is still being served.
	select {
	case err := <-shutdown:
		t.Fatalf("Shutdown returned (%v) while a request was in flight", err)
	case <-time.After(50 * time.Millisecond):
	}

	close(release)

	got := recvOrFail(t, done, "the in-flight request to complete")
	if got.err != nil {
		t.Fatalf("in-flight request failed: %v", got.err)
	}
	if got.status != http.StatusOK {
		t.Errorf("status = %d, want 200 — the request was cut off", got.status)
	}
	if err := recvOrFail(t, shutdown, "Shutdown to return"); err != nil {
		t.Errorf("Shutdown: %v", err)
	}
}

func TestServer_ErrClosesAfterCleanShutdown(t *testing.T) {
	srv := startTestServer(t, http.NewServeMux())

	ctx, cancel := context.WithTimeout(context.Background(), failsafe)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown: %v", err)
	}

	// http.ErrServerClosed is the expected end of Serve, so it is swallowed and
	// the channel simply closes.
	select {
	case err, ok := <-srv.Err():
		if ok {
			t.Errorf("Err() delivered %v after a clean shutdown", err)
		}
	case <-time.After(failsafe):
		t.Fatal("Err() was not closed after shutdown")
	}
}
