# reset · module-topology-redesign

- **Status:** closeout
- **Session:** plan
- **Branch:** module-topology-redesign

## Disposition

- **Re-derived the module topology** (`design/library-topology-and-naming.md`, `concepts/module-set.md`):
  the repository is now **one base library** (a single module rooted at
  `github.com/standards-lab/go-libraries`, capabilities as packages, released as `v<semver>`) plus
  **provider sub-modules** (heavy-SDK implementations, tagged `<path>/v<semver>`). Replaces the earlier
  module-per-capability framing. **There is no `core` module** — its concerns become distinct base
  packages (`lifecycle`, `config`).
- **Settled provider selection and granularity** (`design/conventions.md`): providers are chosen by
  **direct typed construction** (a `Provider` constant + a typed switch at the composition root; no
  registry, no `Register()`, no `init()` side effects), and defined **one per target API**, not per
  deployment — the self-hosted↔managed seam is config within a provider (so storage is `s3` + `azureblob`,
  not four providers).
- **Settled the base dependency policy** (`design/conventions.md`): near-stdlib only (`golang.org/x/…`,
  `google/uuid`); heavy/vendor deps live only in provider sub-modules.
- **Settled capability interface naming** (`design/conventions.md`, `design/library-topology-and-naming.md`):
  per package, not a uniform noun — `database.DB` (package stays `database`, avoiding a stdlib
  `database/sql` clash), `storage.Store`, and `auth.Authenticator`/`auth.TokenSource`.
- **Locked the lifecycle shape** (`concepts/module-set.md`): faithful coordinator — concurrent
  `OnStartup`, `WaitForStartup` flips one readiness signal, `OnShutdown` gated on context cancellation,
  timeout-bounded `Shutdown`; **no startup-hook error handling** (a failing hook fails the process).
- **Recorded relocations** (`concepts/module-set.md`): `result` envelope → `web`; the `bytes`/`parse`/
  `workers` grab-bag not ported; pagination decomposed (`database` owns the persistence query vocabulary,
  `web` owns the HTTP shapes); `logger` deferred to a future `logging` package.
- **Updated cross-cutting notes** (`README.md` capability map, `design/release-and-ci.md`, `CLAUDE.md`
  module-layout/release bullets) to the new topology. Stable context cites no sibling repo; the
  `ref-go-libraries` derivation notes live in `concepts/module-set.md`.
- **Retained:** the open questions in `concepts/module-set.md` — the exact members of `database.DB`/
  `storage.Store`, the `config` package shape, the query and page-response shapes, and the storage
  provider API choices — all unbuilt, settled when each capability is reached.

## Next-focus

Build the base library's first package, **`lifecycle`**. Create the root `go.mod` for
`github.com/standards-lab/go-libraries` (near-stdlib), the `lifecycle` package (the coordinator and its
readiness contract) with its `doc.go` and co-located black-box tests, and the root `CHANGELOG.md` seeded
with `v0.1.0`. Wire the base module in as the first entry of the three synced lists — add `.` to the
`go.work` use-list and to `mise`'s `GO_MODULES`, and change the `ci.yml` matrix's stale `core` entry to
the base module (`.`), confirming its `go.mod`-guarded steps activate. Then `config`, then a minimal
`web`. Start here next session with `marathon start`.
