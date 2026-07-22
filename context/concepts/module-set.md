# Module set

The libraries the repository is expected to grow, kept provisional until each is built and settled. This
session re-derived the layout from the `ref-go-libraries` baseline (core, auth, authz, database, web)
rather than carrying it forward — the names, boundaries, and release topology changed deliberately.

## Shape: one base library, providers as sub-modules

The repository is a single base library (one module at the repository root) whose capabilities are
packages, plus provider sub-modules that carry the heavy SDKs. Rationale over a module-per-capability
layout: every base concern here is near-stdlib, so folding them into one released artifact collapses the
inter-package release ripple to zero while keeping provider isolation where the dependency weight actually
lives. This softens the org's "library is the unit of reuse" line in letter (one module holds all
capability packages) but holds it in spirit (importing one capability compiles no other and pulls no
vendor SDK) — worth a possible follow-up to the org-level context.

## Base library packages

- **lifecycle** — built (`lifecycle/`). The process-lifecycle coordinator: the caller provides the root
  context (`New(ctx)`), startup hooks run concurrently, and shutdown is two-phase — cancel the root, then
  invoke each hook with a fresh timeout-bounded drain context. Startup hooks carry no error handling (a
  hook that cannot do its job fails the process). Readiness is non-monotonic (ready after startup,
  not-ready once shutdown begins) and satisfies `ReadinessChecker`, the contract `web`'s `/readyz`
  consumes. The settled conventions are in `design/conventions.md`; the code and its `doc.go` are now
  authoritative for the package shape.
- **config** — built (`config/`). The layered configuration loader and the contract each capability's
  config implements: a generic `Load` layers base + environment overlay + `secrets.json` + secrets
  overlay onto a caller's type, which supplies `Merge`/`Finalize`; `EnvName` composes override-variable
  names. In the baseline this was a per-module convention with no shared machinery and no secrets file;
  here it became the real base primitive the consumers had hand-rolled. Settled in `design/config.md`;
  the code and its `doc.go` are now authoritative for the package shape.
- **auth** — `Authenticator`/`TokenSource` behavior interfaces; providers (Keycloak self-hosted, Entra
  and others managed) as nested sub-modules. Authorization (RBAC/ABAC) as an in-package `auth/authz`,
  with the enforcement point in `web`.
- **database** — the `database.DB` interface plus a persistence query vocabulary; SQL drivers (postgres,
  mssql) as nested sub-modules.
- **storage** — the `storage.Store` interface; providers per API family (S3, Azure Blob) as nested
  sub-modules.
- **web** — partly built (`web/`). The bootstrap and the health surface are in: a `Server` that binds
  before it serves, a `Config` implementing the configuration contract, RFC 9457 problem responses, a
  JSON writer, and `/healthz` and `/readyz` aggregating `lifecycle.ReadinessChecker` participants. It is
  one flat package — a split is earned by dependency weight, not by topic — and it defines no problem
  type URIs, leaving that vocabulary to consumers. Settled in `design/web.md`; the code and its `doc.go`
  are authoritative for the package shape. Still to come: middleware, error-to-status mapping, the
  success envelope, the HTTP query-param and page-response shapes, and the authorization enforcement
  point.

Future base packages (e.g. **logging**) are added when a consumer needs them — subsystems already take a
stdlib `*slog.Logger`, so a logging construction helper waits until `web`/observability calls for it.

## Provider sub-modules (provisional — scaffolded only when built)

Named for the target API/system, one per API rather than per deployment: `database/postgres` ↔
`database/mssql`; `storage/s3` (minio locally ↔ AWS managed) ↔ `storage/azureblob` (azurite locally ↔
Azure managed); `auth/keycloak` (self-hosted) ↔ `auth/entra` (managed). The self-hosted↔managed seam is a
config change within a provider wherever the API is shared. Selection is direct typed construction — each
provider owns a `Provider` constant and the consumer switches over it at the composition root; no registry
and no `Register()`. (The baseline built a registry, then removed it as unused; we start without it.)

## Decisions carried into the layout

- **No `core` module.** The baseline's `core` was a grab-bag; its concerns become distinct base packages
  (`lifecycle`, `config`). Its `result` success envelope had no cross-capability consumer and moves to
  `web`; its root grab-bag (`bytes`, `parse`, `workers`) is not ported — `parse`'s json-fence extraction
  is LLM-specific and belongs nowhere here.
- **Pagination decomposes by layer.** The baseline's shared `core/pagination` mixed a wire type (json
  tags, string parsing) into the driver-neutral query builder. Instead, `database` owns a pure
  persistence query vocabulary (page + sort as plain directives) and `web` owns the HTTP query-param
  parsing and JSON page-response envelope; the service translates at the seam.
- **Capability interfaces named per package.** No forced uniform noun (the baseline's `database.System`
  read awkwardly, and `auth` is not one encapsulated interface). `database.DB` keeps the package name
  `database` (renaming to `sql` would collide with stdlib `database/sql`); `storage.Store`; `auth` keeps
  behavior interfaces.

## Open questions to settle as each capability is built

- The exact members of `database.DB` and `storage.Store` (lifecycle + access surface).
- The persistence query vocabulary shape (`database`) and the HTTP page-response shape (`web`).
- Final storage provider API choices, and whether both the S3 and Azure Blob families are demonstrated.
- The shape of `web`'s success envelope and its middleware chain. Problem responses are settled
  (`design/web.md`); the envelope waits for a domain handler, and middleware for the logging story and
  for `auth`.
