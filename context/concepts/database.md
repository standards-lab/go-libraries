# Database capability: direction

A planning session with the reference service (2026-08-12) set the direction for the database
capability and a ladder of build sessions; the service-side roadmap is the reference service's
`concepts/data-cqrs-roadmap.md` (volatile cross-reference). This note records the library-side
direction for the rungs still ahead. Rung 1 is built: the `database` package and
`database/postgres` express the wrapper, the dialect seam, the config block, and the lifecycle
and readiness integration — their `doc.go` files carry the settled rationale. Everything below is
candidate direction; each build session settles its slice's API in plan mode.

## Settled direction (remaining rungs)

- Migrations: a `database/migrate` nested sub-module wraps golang-migrate's core and iofs source
  behind a small migrator API (up, down, steps, version, force; embedded `NNNN_<slug>.{up,down}.sql`
  files). Each provider supplies its golang-migrate database driver, so migrate depends on no
  provider. Open question: whether the sub-module naming rule (named for the target system) admits
  `migrate`, which names the capability rather than a vendor.
- The persistence query vocabulary is page and sort directives plus a page result. Filters are
  exact-match directives keyed by projected field names. A projection map and query builder
  (`database/query`, base module) generate SQL through the dialect's placeholder renderer and
  reject unprojected field names.
- Optimistic concurrency: the base fixes a result-shape contract — an update or delete statement
  yields one row `(version, found, matched)`, and a scan helper maps the outcomes to
  `sql.ErrNoRows`, a version-mismatch sentinel, or the new version. The SQL text stays with the
  consumer, per dialect.
- Errors: sentinels for the constraint-violation classes join `ErrNotReady` and
  `ErrConnectionFailed` in the base, and a violation type carries the constraint name while
  wrapping both the sentinel and the driver error, so `errors.Is` matches the class and the
  constraint name survives for mapping. Each dialect's `MapError` wraps the base sentinels —
  declared in the base precisely so consumers stay dialect-agnostic. `sql.ErrNoRows` flows
  through reads unchanged.
- `web` gains its deferred pieces alongside: page-parameter parsing and the page envelope, the
  command result envelope, error-to-status matchers, and ETag/If-Match helpers. The pagination
  split holds: database owns directives, web owns the HTTP side, the consumer translates.

## Releases

The dev-release discipline from `documented-layers` applies across the effort. Each rung's merge
cuts prerelease tags — `v0.4.0-dev.N` for the base, `database/postgres/v0.1.0-dev.N` (and later
`database/migrate/v0.1.0-dev.N`) for the sub-modules — so the reference service's CI resolves the
packages while the surface is still in flux; prerelease tags resolve through the proxy and are
never selected by `@latest`. The provider's transient `replace` stays until the minor release —
it governs only in-repo builds, and local service development rides the gitignored `go.work`
against the repos on disk — but each provider dev tag bumps its `require` to the matching base
dev tag, because consumers read only the require. Changelog entries land under their dev
version's dated heading; the semantic release aggregates the `-dev.N` sections into one and culls
them with the purged tags.
`v0.4.0` and the sub-modules' `v0.1.0` are cut once the full data and CQRS effort closes (all
seven rungs, the surface proven by three domains and the branch-domain composition), paired with
the reference service's own `v0.1.0` and a template patch re-pinning to `v0.4.0`; the
coordinator's release sweep covers the set, and the dev tags purge at the release. CHANGELOG
headings stay pinned to the anticipated versions, dated when actually released.

## Open questions (settled per build session, in plan mode)

- The `migrate` sub-module naming rule (above).
- The builder API: projection construction, where-directives, sort validation.
- Where the typed write-path errors (conflict, concurrency, validation) live, and how their field
  maps project into problem bodies; a structural interface consumed by `web` is the candidate.

## Build ladder (library slices)

1. **Connectivity** — built (this session): the base package skeleton (wrapper, dialect seam,
   config block, lifecycle and readiness) and `database/postgres` open and ping.
2. **Migrations**: `database/migrate`; postgres supplies its driver.
3. **Reads**: the query vocabulary, projection map and builder; web page parsing and envelope.
4. **Writes**: the result envelope, the optimistic-concurrency contract, the error model. The
   consuming service proves the surface and the coordinated snapshot follows.
