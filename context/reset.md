# reset · database-connectivity

- **Status:** closeout
- **Session:** start
- **Branch:** database-connectivity

## Disposition

- **Built:** rung 1 of the database build ladder, coordinated with the reference service — the
  `database` package (wrapper over a provider-constructed pool, the `Dialect` seam, the config
  block, lifecycle and readiness integration with a live bounded ping) and `database/postgres`
  (open and ping over pgx's `database/sql` adapter, eager parse, password set post-parse). The
  session also hoisted `config.SetDurationFromEnv` from per-package copies and migrated
  `web.Config` hydration to Go 1.26 `new(expr)`.
- **Integrated:** the rung-1 settled direction from `concepts/database.md` — the code and the two
  `doc.go` files now express the wrapper members, the dialect shape, and the config block the
  note carried as candidates. The note is trimmed to rungs 2–4.
- **Captured:** `concepts/second-provider.md` — an mssql provider as its own session after
  rung 4, when placeholder rendering, error classification, and the optimistic-concurrency
  contract exist to stress the seam; until then agnosticism is design-guarded.
- **Settled:** release cadence — dev tags per rung, the semantic release at the effort's end.
  Each rung's merge cuts `v0.4.0-dev.N` and `database/postgres/v0.1.0-dev.N` so the reference
  service's CI resolves the packages; `v0.4.0`, the sub-modules' `v0.1.0`, the service's own
  `v0.1.0`, and a template patch re-pinning to `v0.4.0` release together when all seven rungs
  close, with the coordinator's sweep following. Changelog entries land under their dev version's
  dated heading (this rung: `v0.4.0-dev.1`, `v0.1.0-dev.1`); the semantic release aggregates the
  dev sections and culls them with the purged tags. The provider's `replace` stays until the
  minor (in-repo builds only); its `require` bumps to the matching base dev tag at each provider
  dev tag.
- **Retained:** `design/` unchanged; `concepts/module-set.md` unchanged — its reserved database
  section now points at built code for the connectivity surface.

## Next-focus

After merge: tag `v0.4.0-dev.1`, bump the provider's base require to it (the `replace` stays),
tag `database/postgres/v0.1.0-dev.1`, and pin the reference service's `go.mod` to both dev tags
(its local development continues on the gitignored `go.work`). The repository then rests until
the reference service returns for rung 2 (migrations) — its next session is infrastructure
composition, so rung 2 is one service session out.
