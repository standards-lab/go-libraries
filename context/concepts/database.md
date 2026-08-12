# Database capability: direction

A planning session with the reference service (2026-08-12) set the direction for the database
capability and a ladder of build sessions; the service-side roadmap is the reference service's
`concepts/data-cqrs-roadmap.md` (volatile cross-reference). This note records the library-side
direction. It deepens the database section reserved in `concepts/module-set.md`, which stands.
Everything here is candidate direction: each build session settles its slice's API in plan mode.

## Settled direction

- The base `database` package is zero-dependency and builds on `database/sql`. A provider
  constructs the `*sql.DB`; the package wraps it with lifecycle integration (startup pings,
  shutdown closes, the wrapper satisfies `lifecycle.ReadinessChecker`) and a dialect seam.
- The dialect seam is the whole provider surface the base needs: a name, a placeholder renderer
  (postgres `$1`, a future mssql `@p1`), and an error mapper. Query and exec helpers route every
  error through the mapper, so a consumer imports its provider only at the composition root.
- `database/postgres` is the first provider sub-module, over pgx through its `database/sql`
  adapter, carrying the driver dependency in its own `go.mod`. Selection stays a typed constant
  and a switch at the composition root.
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
- Errors: sentinels for the constraint-violation classes, and a violation type carrying the
  constraint name while wrapping both the sentinel and the driver error, so `errors.Is` matches
  the class and the constraint name survives for mapping. `sql.ErrNoRows` flows through reads
  unchanged.
- `web` gains its deferred pieces alongside: page-parameter parsing and the page envelope, the
  command result envelope, error-to-status matchers, and ETag/If-Match helpers. The pagination
  split holds: database owns directives, web owns the HTTP side, the consumer translates.

## Open questions (settled per build session, in plan mode)

- The exact members of the DB wrapper, and whether the query and exec helpers hang off it or
  stand free.
- The dialect seam's precise shape; a struct of the three members is the candidate.
- The database config block: fields, env names, and how the password rides the secrets layering.
- The builder API: projection construction, where-directives, sort validation.
- Where the typed write-path errors (conflict, concurrency, validation) live, and how their field
  maps project into problem bodies; a structural interface consumed by `web` is the candidate.
- Dev-tag cadence between rungs and the first sub-module versions (v0.1.0 candidates).

## Build ladder (library slices)

1. **Connectivity**: the base package skeleton (wrapper, dialect seam, config block, lifecycle and
   readiness) and `database/postgres` open and ping.
2. **Migrations**: `database/migrate`; postgres supplies its driver.
3. **Reads**: the query vocabulary, projection map and builder; web page parsing and envelope.
4. **Writes**: the result envelope, the optimistic-concurrency contract, the error model. v0.4.0
   and the first sub-module tags follow once the consuming service proves the surface.
