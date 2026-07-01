# Candidate module set

The capability modules the repository is expected to grow, kept as a candidate until each is built and
settled. The set re-derives the `ref-go-libraries` baseline (core, auth, authz, database, web) rather than
copying it forward, so the names and boundaries may change.

- **core** — promoted to the first build; see `reset.md`.
- **auth** — authentication behind one interface; providers (Keycloak self-hosted, Entra and others
  managed) as nested submodules. Authorization (RBAC/ABAC) as an in-module package (`auth/authz`), with
  the policy-enforcement point in `web`.
- **database** — data access behind one interface; SQL providers (postgres, sql server) as nested
  submodules.
- **storage** — object storage behind one interface (minio / azurite self-hosted ↔ s3 / azure blob
  managed).
- **web** — the HTTP layer: a stdlib `net/http` server, RFC 9457 problem responses, middleware,
  liveness/readiness (`/healthz`, `/readyz`, where `/readyz` surfaces `core`'s hot-start lifecycle), and
  the authorization enforcement point.

Build order starts with `core`, then a minimal `web` — the smallest runnable HTTP surface. From there the
set grows one module at a time, each capability landing as the pattern and integration it encapsulates are
established.

The `web` split is settled in principle: the `net/http` bootstrap, problem-response scaffolding,
middleware, and the `/healthz`/`/readyz` probes are library-level; request routing and domain handlers
belong to the consuming service. Open questions to settle as the work reaches each module: whether
`storage` is its own module or folds into another; the exact split between `core` and a separate config
module.
