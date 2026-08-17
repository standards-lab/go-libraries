# Module set

The libraries the repository is expected to grow, kept provisional until each is built and settled. This
session re-derived the layout from the `ref-go-libraries` baseline (core, auth, authz, database, web)
rather than carrying it forward — the names, boundaries, and release topology changed deliberately.

## Shape: one base library, providers as sub-modules

The repository is a single base library (one module at the repository root) whose capabilities are
packages, plus provider sub-modules that carry the heavy SDKs. Rationale over a module-per-capability
layout: every base concern here is near-stdlib, so folding them into one released artifact collapses the
inter-package release ripple to zero while keeping provider isolation where the dependency weight actually
lives. The unit of reuse is still one capability: importing one capability compiles no other and pulls
no vendor SDK.

## Base library packages

- **lifecycle** — built (`lifecycle/`); the code and its `doc.go` are authoritative for the package
  API. The settled conventions are in `design/conventions.md`.
- **config** — built (`config/`); the code and its `doc.go` are authoritative for the package API.
  Settled in `design/config.md`.
- **auth** — `Authenticator`/`TokenSource` behavior interfaces; providers (Keycloak self-hosted, Entra
  and others managed) as nested sub-modules. Authorization (RBAC/ABAC) as an in-package `auth/authz`,
  with the enforcement point in `web`.
- **database** — built (`database/`, `database/seed`, and the `database/postgres` provider); the code and
  its `doc.go` files are authoritative. The standard tier is ISO/IEC 9075; the persistence query
  vocabulary is under way. Further engines (sqlite, mssql) as nested sub-modules when a proof or a
  consumer earns them. Direction in `concepts/database.md`.
- **storage** — the `storage.Store` interface; providers per API family (S3, Azure Blob) as nested
  sub-modules.
- **logging** — built (`logging/`); the code and its `doc.go` are authoritative for the package API.
  Settled in `design/logging.md`.
- **web** — partly built (`web/`). The bootstrap, the health endpoints, and the middleware chain are in: a
  `Server` that binds before it serves, a `Config` implementing the configuration contract, RFC 9457
  problem responses, a JSON writer, `/healthz` and `/readyz` aggregating `lifecycle.ReadinessChecker`
  participants, and `Middleware`/`Chain` with a `RequestLogger`. It is one flat package — a split is
  earned by dependency weight, not by topic — and it defines no problem type URIs, leaving that
  vocabulary to consumers. Settled in `design/web.md`; the code and its `doc.go` are authoritative for the
  package API. Still to come: the rest of the middleware set, error-to-status mapping, the success
  envelope, the HTTP query-param and page-response contracts, and the authorization enforcement point.

## Provider sub-modules (provisional — scaffolded only when built)

Named for the target API/system, one per API rather than per deployment: `database/postgres` ↔
`database/mssql`; `storage/s3` (minio locally ↔ AWS managed) ↔ `storage/azureblob` (azurite locally ↔
Azure managed); `auth/keycloak` (self-hosted) ↔ `auth/entra` (managed). The self-hosted↔managed seam is a
config change within a provider wherever the API is shared. Selection is direct typed construction — each
provider owns a `Provider` constant and the consumer switches over it at the composition root; no registry
and no `Register()`. (The baseline built a registry, then removed it as unused; this repository starts
without one.)

## Decisions carried into the layout

- **No `core` module.** The baseline's `core` was a grab-bag; its concerns become distinct base packages
  (`lifecycle`, `config`). Its `result` success envelope had no cross-capability consumer and moves to
  `web`; its root grab-bag (`bytes`, `parse`, `workers`) is not ported — `parse`'s json-fence extraction
  is LLM-specific and belongs nowhere here.
- **Pagination decomposes by layer.** The baseline's shared `core/pagination` mixed a wire type (json
  tags, string parsing) into the driver-neutral query builder. Instead, `database` owns a pure
  persistence query vocabulary (page + sort as plain directives) and `web` owns the HTTP query-param
  parsing and JSON page-response envelope; the service translates at the seam. The SQL a page directive
  renders to is the standard form, not dialect-dispatched, so nothing an engine does reaches `web`.
- **Capability interfaces named per package.** No forced uniform noun (the baseline's `database.System`
  read awkwardly, and `auth` is not one encapsulated interface). `database.DB` keeps the package name
  `database` (renaming to `sql` would collide with stdlib `database/sql`); `storage.Store`; `auth` keeps
  behavior interfaces.

## Open questions to settle as each capability is built

- The exact members of `storage.Store`: the operations the target APIs share are the standard tier this
  library establishes for object storage; everything else stays native.
- The persistence query vocabulary (`database`) and the HTTP page-response contract (`web`).
- Final storage provider API choices. One family by default; the second only when a proof or a consumer
  earns it.
- The structure of `web`'s success envelope. Problem responses and the middleware chain are settled
  (`design/web.md`); the envelope waits for a domain handler.
- Whether middleware earns a package of its own, and when — see `concepts/middleware-split.md`.
