# reset · build-web-package

- **Status:** closeout
- **Session:** start
- **Branch:** build-web-package

## Disposition

- **Built the base library's `web` package** (`web/`): a `Server` that binds the listener on the calling
  goroutine and only then serves in the background — so a bind failure reaches the caller instead of a
  goroutine nobody watches — reporting the bound address and delivering any later serve failure on a
  buffered `Err()` channel; a `Config`/`Env` pair implementing the `config` contract; RFC 9457 problem
  responses (`Problem`, `WriteProblem`, `WriteProblemWith`) and a `WriteJSON` writer; and `/healthz` and
  `/readyz` handlers aggregating named `lifecycle.ReadinessChecker` participants. Also added
  `config.Duration`, a duration that travels through JSON as a string, replacing the stringly-typed
  timeout fields and per-field accessors both baselines hand-rolled. Black-box tests (`-race`, 33 in
  `web` and 10 in `config`), `web/doc.go`, and `web` plus `Duration` entries in the **unreleased**
  `v0.1.0` `CHANGELOG.md`. Packages of the existing base module — no new module, no `go.work`/`mise`/CI
  synced-list change.
- **Promoted** the settled HTTP posture into a new `design/web.md`: the bootstrap rationale (bind before
  serve, the composition root owns the wiring, no shutdown timeout of its own since the coordinator
  supplies a drain context); the **one flat package** rule — a split is earned by dependency weight, not
  by topic, the same test that makes `database/postgres` a sub-module; health as a readiness *surface*
  that contributes no checker of its own; and the decision that **the library defines no problem type
  URIs**, because a `type` names an application's vocabulary rather than a library's, leaving it to
  consumers and accepting that `/readyz` carries its `checks` extension on an `about:blank` problem.
- **Integrated** the `web` bullet in `concepts/module-set.md` — marked partly built, pointed at the code
  and `design/web.md`, listed what remains — narrowed the open question about `web`'s internal split to
  the envelope and middleware only, and updated the `README.md` build-order line now that `lifecycle`,
  `config`, and a minimal `web` are in.
- **Retained:** the rest of `concepts/module-set.md` (auth, database, storage, and the provider set) and
  its remaining open questions — all still unbuilt, settled when each capability is reached.

## Next-focus

Complete the base library to a releasable **v0.1.0**, then tag it. Three pieces, in order:

1. A **`logging` package** — a `Config`/`Env` implementing the configuration contract (level, format) and
   a constructor returning a `*slog.Logger`. This is the trigger condition `module-set.md` set for a
   logging package: both baselines hand-roll the same `LogLevel` parsing and handler construction, so the
   need is demonstrated rather than predicted.
2. **Middleware primitives in `web`** — a `Middleware` type, a chain helper, and the request logger built
   on the new `logging` package. Still deferred: CORS (no browser client), `Auth`/`Authorize` (blocked on
   `auth`), and error-to-status mapping (blocked on a domain handler).
3. An **ergonomics pass** on the composition root. This session's scratch program measured it: of roughly
   thirty statements needed to stand a service up, only about six encode a decision — where config lives,
   which routes and checks are registered, the drain timeout — while the rest is identical in every
   consumer (trap signals, `lifecycle.New`, register `Start` with a fatal-on-error closure, register
   `Shutdown` with a log-on-error closure, `WaitForStartup`, `select` on `Context().Done()` versus
   `Err()`, then `Shutdown`). Candidate shapes are `srv.Bind(lc)` to register both hooks in one explicit
   call, and `lc.Run(timeout)` to absorb the wait-block-drain sequence — which first needs a way for a
   subsystem to report a fatal runtime error to the coordinator, since that is what the `select` over
   `Err()` exists for. Keep both **additive**: `Start`, `Shutdown`, `OnStartup`, and `OnShutdown` stay
   public and usable on their own, so a helper can be reshaped later without disturbing the foundation.
   Unlike the first two pieces, these shapes are inferred from one scratch program rather than proven by
   a real composition root — build them lightly.

Close by tagging `v0.1.0` from the root `CHANGELOG.md`, which is also the first exercise of
`.github/workflows/release.yml`. That release is what the next level builds against; the ceremony a real
composition root turns out not to need comes back as a `v0.2.0` refinement.
