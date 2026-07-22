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
// address once started, which lets a caller configure port 0 and read back the
// port that was assigned. A serve failure after startup arrives on
// [Server.Err], a buffered channel that is closed when serving stops;
// http.ErrServerClosed is the expected end of a shutdown and is not reported.
//
// # Lifecycle wiring
//
// The package registers no lifecycle hooks of its own and carries no shutdown
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
// [Config] holds the listen address and the server's four timeouts, and
// implements the config package's Merge and Finalize contract, so it loads as
// part of an application's configuration rather than on its own. Finalize
// applies defaults, then the environment overrides named by the [Env] value the
// configuration carries, then validates. [NewEnv] composes the standard names
// from a prefix; an empty name disables that one override, so a zero Env means
// no environment overrides at all.
//
// # Health
//
// [Liveness] reports that the process is up and serving HTTP and checks nothing
// else, which is what makes an unanswered probe the signal. [Readiness]
// aggregates the [Check] values the caller supplies — typically the lifecycle
// coordinator alongside each subsystem that reports readiness — and answers 503
// unless every one of them is ready. A Check with a nil Checker reports not
// ready, so a subsystem that failed to construct fails the probe.
// [RegisterHealth] mounts both endpoints on a [Mounter], which http.ServeMux
// satisfies incidentally.
//
// # Problem responses
//
// Error responses are RFC 9457 problem documents. The type member identifies
// the problem's semantics and is the member a client branches on, with title
// advisory and status an advisory copy of the status line. This package defines
// no type URIs of its own, because a type URI names an application's
// vocabulary rather than a library's: every problem it emits is
// [ProblemTypeBlank], meaning no semantics beyond the HTTP status code, and a
// consumer supplies its own URI through [Problem.Write] or the extras map of
// [WriteProblemWith]. An empty title defaults to the status phrase, which is
// what RFC 9457 asks of an about:blank problem in any case.
package web
