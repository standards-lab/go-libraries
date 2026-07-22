# Changelog

All notable changes to the base library (`github.com/standards-lab/go-libraries`) are documented here. The
format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the library adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html). This changelog covers the base module only;
each provider sub-module keeps its own.

## [v0.1.0] - 2026-07-21

### Added

- Base module `github.com/standards-lab/go-libraries`.
- `lifecycle` package: a process-lifecycle coordinator with a caller-provided root context, concurrent
  startup hooks, a readiness contract (`ReadinessChecker`), and two-phase, timeout-bounded graceful
  shutdown.
- `config` package: a generic layered configuration loader (`Load`) over a `Config`/`Merge`/`Finalize`
  contract, layering a base file, an environment overlay, `secrets.json`, and a secrets overlay, plus an
  `EnvName` helper for composing environment-variable override names and a `Duration` type that carries
  a `time.Duration` through JSON as a string (`"30s"`).
- `web` package: an HTTP bootstrap (`Server`) that binds before serving so a bind failure reaches the
  caller, reports the bound address, and surfaces a later serve failure on a channel; a `Config`
  implementing the configuration contract; RFC 9457 problem responses (`Problem`, `WriteProblem`,
  `WriteProblemWith`) and a JSON writer (`WriteJSON`); and `/healthz` and `/readyz` handlers
  (`Liveness`, `Readiness`, `RegisterHealth`) that aggregate `lifecycle.ReadinessChecker` participants.
