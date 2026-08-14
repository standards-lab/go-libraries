# reset · database-migrations

- **Status:** closeout
- **Session:** start
- **Branch:** database-migrations

## Disposition

- **Built:** `database/seed` in the base module — a seed runner over what only the consumer
  knows: the seed file system, a typed load function per table (`Table` → opaque `Step`), and
  the step order. The `Format` seam decodes by file extension, registered at construction with
  no registry; the strict `JSON` format ships in the package. One transaction per step, outcome
  logged, idempotency in the consumer's SQL. Stdlib-only, so no module bookkeeping changed and
  the base dependency graph is untouched. Proven by the reference service's `db seed` on the
  companion `database-migrations` branch: seven organization rows, idempotent rerun.
- **Culled:** the planned `database/migrate` sub-module and its postgres driver-supply package,
  settled with the developer mid-session — wrapping golang-migrate added constraint but no
  mechanics. Migrations standardize at the service level instead: the reference service consumes
  golang-migrate directly in `cmd/db`, explicit instance construction over the shared pool, no
  blank-import registry. The sub-module naming question the wrapper raised closed with it, and
  the topology note stands unamended. The reference service's toolchain evaluation (its
  roadmap's rung 8) re-asks the promotion question with usage evidence.
- **Integrated:** `concepts/database.md` marks rung 2 built, records the migrations resolution,
  and corrects the release plan — no migrate tags, provider dev tags only when a rung touches
  the provider, and the closing `v0.5.0` set gated on the service's toolchain evaluation. The
  changelog opens `[v0.5.0-dev.1]` with the seed package.
- **Retained:** `concepts/middleware-split.md`, `concepts/module-set.md`, and
  `concepts/second-provider.md`, untouched by this session.

## Next-focus

Rung 3 of the data ladder: reads, coordinated with the reference service. Plan in plan mode the
persistence query vocabulary — page and sort directives, the projection map and query builder
(`database/query`, base module) — and `web`'s page-parameter parsing and page envelope, with the
builder API the recorded open question. The rung's merge cuts `v0.5.0-dev.2`.
