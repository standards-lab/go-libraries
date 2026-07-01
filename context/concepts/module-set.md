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
- **web** — the HTTP layer: RFC 9457 problem responses, middleware, the authorization enforcement point.

Open questions to settle as the work reaches each module: whether `storage` is its own module or folds
into another; how much of the `web` layer belongs in libraries versus the reference architecture; the
exact split between `core` and a separate config module.
