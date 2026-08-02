# Changelog

All notable changes to the base library (`github.com/standards-lab/go-libraries`) are documented here.
Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html). This changelog covers the
base module only; each provider sub-module keeps its own.

## [v0.3.1] - 2026-08-02

### Changed

- **`web.Server.Shutdown` before a successful `Start` is a no-op** that leaves the server
  startable. Under `Run`'s drain-what-started semantics, the drain after a failed startup
  invokes every shutdown hook — including the server's, whose bind never succeeded — so a
  never-started drain is a normal path, not misuse. A held-port startup failure now reports the
  single `startup: listen …` error instead of joining a `shutdown: server not started` red
  herring.

## [v0.3.0] - 2026-08-02

Composition ergonomics, driven by the template baseline — the library's first end-to-end
consumer.

### Changed

- **`lifecycle.Coordinator` is a `Run`-based host** (breaking): consumers declare hooks and
  monitors while the coordinator waits, then make one blocking `Run(ctx, drainTimeout)` call that
  owns the whole sequence. `New` takes no context (the caller's signal context goes to `Run`);
  `OnStartup` and `OnShutdown` hooks are `func(context.Context) error`, launched concurrently by
  `Run` rather than at registration; new `OnReady` hooks fire once startup succeeds, and new
  `Monitor` channels end the run on their first non-nil error. A startup failure drains what did
  start and never flips readiness; `Run` returns nil exactly when a signal-driven exit drained
  cleanly, and otherwise the joined errors wrapped by phase (`startup:`, `run:`, `shutdown:`).
  `WaitForStartup`, `Shutdown`, and `Context` are gone — `Run` owns them.
- **`web.Server.Start` takes a context** (breaking): the bind goes through `net.ListenConfig`,
  so it is cancellable, and `Start`/`Shutdown` now carry the lifecycle hook signature —
  `lc.OnStartup(srv.Start)` and `lc.OnShutdown(srv.Shutdown)` register as bare method values.
- **`web.RequestLogger` demotes successful probe requests to debug**: a 2xx on `/healthz` or
  `/readyz` is orchestrator heartbeat and stays out of info-level production logs; a failing
  probe still logs at info.
- **`web.NewServer` panics diagnostically on an unfinalized `Config`** — naming the fix — instead
  of nil-dereferencing a tri-state pointer.
- Error strings name their operation and operand, never their package (`server not started`
  drops its `web:` prefix); package prefixes remain on panic messages only.

### Added

- `config.Duration.Duration()` returns the value as a `time.Duration`, removing the conversion
  at call sites such as `lc.Run(ctx, cfg.ShutdownTimeout.Duration())`.

## [v0.2.0] - 2026-07-31

### Changed

- **`web.Config` is tri-state** (breaking): `Port` is `*int` and the four timeouts are
  `*config.Duration`. `nil` is unset and takes the default; an explicit zero survives the load and
  means what it says — a disabled timeout, or an ephemeral port read back through `Server.Addr`. A
  file and the environment express both states identically.
- **`lifecycle.Coordinator` is a state machine** (breaking): `OnStartup` panics once `WaitForStartup`
  has been entered and `OnShutdown` panics once `Shutdown` has begun; `Shutdown` runs its hooks
  exactly once, with a repeated call blocking until the first completes and returning the same error;
  `Ready` is false the moment shutdown begins, including when it overtakes a `WaitForStartup` still
  in flight; the timeout error wraps `context.DeadlineExceeded`; zero registered hooks return nil
  deterministically, and completion within the timeout always wins.
- `logging.Config.Finalize` normalizes before applying defaults, so a whitespace-only value takes the
  default instead of failing validation.

### Fixed

- `config.Duration` treats JSON null as a no-op and decodes bare numbers exactly as int64
  nanoseconds; a fractional number is rejected toward the string form.
- `config.Load` validates `Options.OverlayPattern` and fails a malformed pattern instead of silently
  skipping the overlays it names.
- `config.EnvName` sanitizes segments: runs of characters outside A-Z and 0-9 collapse to single
  underscores.
- `web.Server.Shutdown` before a successful `Start` returns an error and leaves the server startable,
  instead of poisoning a later `Start` into serving nothing.
- The `web.RequestLogger` wrapper delegates `io.ReaderFrom`, restoring the zero-copy path for file
  responses, and logs a panicking request at error level with the panic value before re-panicking.
- `web.Problem` defaults a zero `Status` to 500 instead of panicking in `WriteHeader`; an unknown
  status code omits the `title` member; and `WriteProblemWith` re-seeds `status` after the extras
  copy, so the body always matches the status line.

## [v0.1.0] - 2026-07-28

### Added

- Base module `github.com/standards-lab/go-libraries`.
- `lifecycle` package: a process-lifecycle coordinator with a caller-provided root context, concurrent
  startup hooks, a readiness contract (`ReadinessChecker`), and two-phase, timeout-bounded graceful
  shutdown.
- `config` package: a generic layered configuration loader (`Load`) over a `Config`/`Merge`/`Finalize`
  contract, layering a base file, an environment overlay, `secrets.json`, and a secrets overlay, plus an
  `EnvName` helper for composing environment-variable override names and a `Duration` type that carries
  a `time.Duration` through JSON as a string (`"30s"`).
- `logging` package: a `Config`/`Env` pair implementing the configuration contract and a `New`
  constructor returning a `*slog.Logger` over a caller-supplied `io.Writer`, with a string-backed `Level`
  that delegates parsing to `slog.Level` and a `Format` selecting the text or JSON handler.
- `web` package: an HTTP bootstrap (`Server`) that binds before serving so a bind failure reaches the
  caller, reports the bound address, and surfaces a later serve failure on a channel; a `Config`
  implementing the configuration contract; RFC 9457 problem responses (`Problem`, `WriteProblem`,
  `WriteProblemWith`) and a JSON writer (`WriteJSON`); `/healthz` and `/readyz` handlers
  (`Liveness`, `Readiness`, `RegisterHealth`) that aggregate `lifecycle.ReadinessChecker` participants;
  and middleware primitives (`Middleware`, `Chain`) with a `RequestLogger` that emits one record per
  request through a caller-supplied `*slog.Logger`.
