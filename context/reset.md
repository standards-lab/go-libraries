# reset · capability-tiers

- **Status:** closeout
- **Session:** plan
- **Branch:** capability-tiers

## Disposition

- **Integrated:** `design/conventions.md` — the two tiers every capability presents (a standard tier
  that is exactly the technology's common standard, and the provider's native API through the handle
  the capability exposes), the rule that a consumer interfaces at the resolution its purpose requires
  and wraps native use beneath the standard tier, the import boundary (a consumer imports a provider
  only from its composition root, `cmd/*`, and packages that declare native use, checked by a lint
  step on the consumer's side), the provider-swap classes each capability design note declares
  (interchangeable, interchangeable with review, schema-bearing), the `Options` map and the
  dual-wrapped error as the native tier at configuration and error level, and the rule that
  `database.Dialect` renders only what engines do differently from the standard.
  `context/README.md` — the aims and the capability map restated in those terms; `database` marked
  built with its Postgres provider; an `observability` entry with OpenTelemetry as its standard tier.
  `design/library-topology-and-naming.md` — `database` in the built list, the sub-module criterion
  stated as the dependency the base must not carry, `database/sqlite` in the naming list, and the
  bound on how providers are added. `design/release-and-ci.md` — the three-list mechanism as what
  makes a second provider cheap; the import boundary runs in consumers. `design/web.md` — the
  package's standard tier named (RFC 9110 and 9457, no providers), the RFC's "extension member"
  disambiguated, and the page contract carrying no engine detail.
- **Integrated:** `concepts/database.md` — rung 1 restated in tier terms; the rung-3 direction
  settled: the builder emits standard SQL:2008 paging, `Dialect` is unchanged, the projection has a
  key field as the ORDER BY tie-breaker and admits expression-backed fields, count and page are two
  statements, an unknown field is a typed error, and `web` parses `page`/`size`/`sort` and writes the
  `items`/`page`/`size`/`total` envelope; `RETURNING` named as native domain SQL; what a provider swap
  changes; rung 5 (second provider) added to the ladder. `concepts/second-provider.md` — rewritten:
  what a second provider is for, paging as the third stress point, SQLite recorded as the cheapest
  candidate beside SQL Server. `concepts/module-set.md` — the upward note to the org context culled,
  `database` marked built, the storage questions reframed as one family by default.
- **Retained:** `design/config.md`, `design/logging.md`, `concepts/middleware-split.md`, unchanged.

## Next-focus

Rung 3 of the data ladder: reads, coordinated with the reference service (this repository first).
The inputs are settled in `concepts/database.md`: `database/query` in the base with page, sort, and
filter directives, a projection with a key field and name-to-expression fields, and a builder that
emits standard SQL — placeholders through the dialect, paging in the SQL:2008 form, count and page as
separate statements, unknown fields rejected with a typed error; `web` gains `page`/`size`/`sort`
parsing into its own directive types and the page envelope. `Dialect` and `database/postgres` are
untouched. The remaining plan-mode detail is the builder API's construction and error types. The
session's deliverables also carry the repository `README.md` (the stale package list, "None exist
yet", RFC 9110 beside 9457) and `CLAUDE.md` (the import boundary; "design, tier, layer, and release").
The rung's merge cuts `v0.5.0-dev.2`.
