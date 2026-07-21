// Package lifecycle coordinates process startup, readiness, and graceful shutdown.
//
// A [Coordinator] derives its root context from a context the caller provides, runs
// registered startup hooks concurrently, tracks a single readiness signal, and drives
// a two-phase shutdown bounded by a timeout.
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
// registered. [Coordinator.WaitForStartup] blocks until every startup hook has
// returned and then marks the coordinator ready. Startup hooks return no error: a hook
// that cannot do its job fails the process directly, by panicking or exiting, so the
// coordinator stays a pure orchestrator with no startup error path.
//
// # Readiness
//
// [Coordinator.Ready] reports whether startup has completed and shutdown has not begun.
// It is false until [Coordinator.WaitForStartup] returns, true afterward, and false
// again once [Coordinator.Shutdown] starts, so a readiness probe reports a draining
// process as not ready. [Coordinator] satisfies [ReadinessChecker], the contract a
// /readyz endpoint consumes.
//
// # Shutdown
//
// [Coordinator.Shutdown] runs in two phases. It first cancels the root context, so work
// watching [Coordinator.Context] stops taking on new work. It then invokes each
// [Coordinator.OnShutdown] hook concurrently, passing a fresh drain context bounded by
// the timeout. The drain context is derived from [context.Background], not the cancelled
// root, so cleanup has the whole timeout to finish. Shutdown returns nil once every hook
// returns, or an error if the timeout elapses first. Shutdown hooks need no cancellation
// guard of their own; the coordinator invokes them only after the root context is
// already cancelled.
package lifecycle
