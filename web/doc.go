// Package web provides the HTTP layer: a net/http server bound to a
// caller-supplied handler, RFC 9457 problem responses, a JSON writer, and the
// liveness and readiness endpoints an orchestrator probes.
//
// # Server
//
// [NewServer] wraps an http.Server built from a [Config] and a handler the
// caller composes. [Server.Start] binds the listener on the calling goroutine
// and only then serves in the background, so a bind failure is returned to the
// caller instead of being lost in a goroutine. [Server.Addr] reports the bound
// address once started; a configured port 0 binds an ephemeral port, and Addr
// reads back the assignment.
//
// [Server.Err] belongs to a successful serve session: after Start returns nil,
// a serve failure arrives on it, and the channel closes when serving stops.
// http.ErrServerClosed is the expected end of a shutdown and is not reported.
// [Server.Shutdown] before a successful Start returns an error and leaves the
// server startable; once it has served, a Server is single-use — construct a
// new one to serve again.
//
// # Lifecycle wiring
//
// The package registers no lifecycle hooks of its own and holds no shutdown
// timeout. A composition root wires [Server.Start] as a startup hook and
// [Server.Shutdown] as a shutdown hook, taking the coordinator's
// timeout-bounded drain context directly:
//
//	lc.OnStartup(func() {
//		if err := srv.Start(); err != nil {
//			log.Fatalf("http: %v", err)
//		}
//	})
//	lc.OnShutdown(func(ctx context.Context) {
//		if err := srv.Shutdown(ctx); err != nil {
//			log.Printf("http: shutdown: %v", err)
//		}
//	})
//
// # Configuration
//
// [Config] holds the host, the port, and the server's four timeouts, and
// implements the config package's Merge and Finalize contract, so it loads as
// part of an application's configuration rather than on its own. The port and
// timeouts are pointers: nil is unset and takes the default, while an explicit
// zero survives the load and means what it says — a disabled timeout, or an
// ephemeral port. A file and the environment express both states identically.
// Finalize applies defaults, then the environment overrides named by the
// configuration's [Env] value, then validates. [NewEnv] composes the standard
// names from a prefix; an empty name disables that one override, and a zero Env
// means no environment overrides at all.
//
// # Health
//
// [Liveness] reports that the process is up and serving HTTP and checks nothing
// else; an unanswered probe is the liveness signal. [Readiness] aggregates the
// [Check] values the caller supplies — typically the lifecycle coordinator
// alongside each subsystem that reports readiness — and answers 503 unless
// every one of them is ready. A Check with a nil Checker reports not ready, so
// a subsystem that failed to construct fails the probe. [RegisterHealth] mounts
// both endpoints on a [Mounter]. The patterns use net/http.ServeMux's
// method-scoped syntax ("GET /healthz"), which http.ServeMux handles directly;
// any other Mounter translates them.
//
// # Middleware
//
// A [Middleware] wraps one http.Handler in another, and [Chain] composes a set
// of them around a handler in argument order: in Chain(h, a, b), a sees the
// request first. A nil entry is skipped, so a caller can build a chain with
// conditional entries without filtering it first.
//
// [RequestLogger] emits one record per request through a *slog.Logger the
// caller supplies — method, path, status, duration, and remote address — at
// info level, or at error level with the panic value attached when the handler
// panics (the panic then continues to net/http's recovery). Every request logs
// at info on the normal path; the middleware does not judge whether a 5xx was
// the application's own failure.
//
// The middleware wraps the ResponseWriter to capture the status. The wrapper
// records the first status written, implements Unwrap so flushing and hijacking
// work through http.ResponseController, and delegates io.ReaderFrom so a
// handler serving files keeps the zero-copy path. http.Pusher is not available
// through the wrapper.
//
// # Problem responses
//
// Error responses are RFC 9457 problem documents. The type member identifies
// the problem's semantics and is the member a client branches on, with title
// advisory and status an advisory copy of the status line. This package defines
// no type URIs of its own: every problem it emits is [ProblemTypeBlank], and a
// consumer supplies its own URI through [Problem.Write] or the extras map of
// [WriteProblemWith]. Extras may add or override any member except status,
// which always matches the status line. A zero Status defaults to 500. An empty
// title defaults to the status phrase, and is omitted for a code outside the
// standard table.
package web
