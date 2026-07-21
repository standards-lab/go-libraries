# reset · build-base-lifecycle

- **Status:** closeout
- **Session:** start
- **Branch:** build-base-lifecycle

## Disposition

- **Built the base module and its first package, `lifecycle`** (`go.mod`, `lifecycle/`): the
  process-lifecycle coordinator — caller-provided root context (`New(ctx)`), concurrent startup hooks, a
  `ReadinessChecker` contract, and two-phase, timeout-bounded graceful shutdown. Black-box tests (`-race`,
  7 cases), `doc.go`, and the root `CHANGELOG.md` (seeded `v0.1.0`) accompany it.
- **Wired the base module `.` into the three synced lists**: the `go.work` use-list, `mise`'s
  `GO_MODULES`, and the `ci.yml` matrix (the stale `core` entry became `.` in both the test and lint jobs).
  `mise run test`/`lint` and CI's `go.mod`-guarded steps now activate on the base module.
- **Promoted** the settled lifecycle/context conventions into `design/conventions.md` ("Process lifecycle
  and context ownership"): the composition root owns the root context and traps signals; shutdown is
  coordinator-driven with a fresh drain context; readiness is non-monotonic; leaf subsystems take a plain
  `context.Context`. These refined the previously "locked" lifecycle shape, grounded in the
  ref-go-libraries / herald / personnel-service-demo sources.
- **Integrated** the `lifecycle` bullet in `concepts/module-set.md` — marked built and pointed at the code
  and conventions; the code and its `doc.go` are now authoritative for the package shape.
- **Retained:** the rest of `concepts/module-set.md` (config, auth, database, storage, web, and the
  provider set) and its open questions — all still unbuilt, settled when each capability is reached.

## Next-focus

Build the base library's next package, **`config`**: layered configuration (base + overlay +
`secrets.json`) and the merge/finalize/env contract each capability's config implements — the shared base
primitive the baseline lacked. Settle the `Load` orchestration and the config shape as it is built (open
questions in `concepts/module-set.md`). It lands as a package of the existing base module — no new module,
no synced-list change. Start here next session with `marathon start`.
