// Package lifecycle coordinates process startup, readiness, and graceful shutdown.
//
// A [Coordinator] moves through five phases: registration (startup hooks running),
// waiting (inside [Coordinator.WaitForStartup]), running (ready), draining (inside
// [Coordinator.Shutdown]), and stopped. Each call in the API is legal in specific
// phases, and the transitions between them are what the coordinator guarantees.
//
// # Context ownership
//
// The caller owns the root context. A composition root typically builds one from the
// standard library's signal.NotifyContext and passes it to [New]; the coordinator
// derives a cancellable context from it, reachable through [Coordinator.Context].
// Long-lived work observes that context and stops when it is cancelled. The
// coordinator installs no signal handlers of its own.
//
// # Startup
//
// [Coordinator.OnStartup] registers work that runs concurrently from the moment it is
// registered, and it panics once [Coordinator.WaitForStartup] has been entered:
// registration belongs to the composition root's setup phase, and a late registration
// is a programming error. [Coordinator.WaitForStartup] blocks until every startup
// hook has returned and then marks the coordinator ready. Startup hooks return no
// error: a hook that cannot do its job fails the process directly, by panicking or
// exiting.
//
// # Readiness
//
// [Coordinator.Ready] reports whether the coordinator is running: true once
// [Coordinator.WaitForStartup] returns, and false again the moment
// [Coordinator.Shutdown] begins — including when a shutdown overtakes a
// WaitForStartup still in flight — so a readiness probe reports a draining or
// stopped process as not ready. [Coordinator] satisfies [ReadinessChecker], the
// contract a /readyz endpoint consumes.
//
// # Shutdown
//
// [Coordinator.Shutdown] runs in two phases. It first cancels the root context, so
// work watching [Coordinator.Context] stops taking on new work. It then invokes each
// [Coordinator.OnShutdown] hook concurrently, passing a fresh drain context bounded
// by the timeout. The drain context derives from [context.Background], so cleanup
// has the whole timeout regardless of the cancelled root. Registering a hook once
// shutdown has begun panics.
//
// The hooks run exactly once: the first Shutdown call drives the drain, and a
// repeated call blocks until it completes and returns the same error. With no
// registered hooks Shutdown returns nil. Completion within the timeout always
// returns nil; the timeout error wraps [context.DeadlineExceeded]. A timed-out
// Shutdown returns while its unfinished hooks continue on the expired drain
// context — the coordinator cannot stop a goroutine — so hooks that honor their
// context stop promptly.
package lifecycle
