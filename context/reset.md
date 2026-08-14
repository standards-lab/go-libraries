# reset · infrastructure-composition

- **Status:** closeout
- **Session:** start
- **Branch:** infrastructure-composition

## Disposition

- **Built:** the `web` routing layer (`Group`, `Module`, `Router`), re-derived from the prior
  reference material with its catalogued flaws corrected: a group tree compiles once at wiring
  time and seals, so nothing recomposes per request and mutating a compiled group panics; a
  group becomes servable and mounts through one path each, with multi-level prefixes
  first-class; probes register on the router's native mux, outside every module's middleware;
  registration mistakes panic at wiring time. Coordinated with composition-root sessions in
  go-service-template and go-service on companion branches.
- **Built:** the configuration revision settled during consumer validation: capability configs
  implement `Finalize(envPrefix string) error` and compose their own override names, `Env` is
  recorded output rather than caller-seeded input, `NewEnv("")` returns the zero `Env`, and
  `Options.EnvPrefix` threads the prefix through `Load`, deriving `EnvVar` when unset. The
  earlier seed-then-finalize pairing silently disabled every override when the seed was
  forgotten.
- **Released:** `v0.4.0` and `database/postgres/v0.1.0` at this closeout — the changelogs
  aggregate their `-dev.1` sections, the postgres base require bumps to `v0.4.0`, and the
  `v0.4.0-dev.1` and `database/postgres/v0.1.0-dev.1` tags are deleted. The template and the
  reference service pin the stable versions at their own closeouts.
- **Integrated:** `design/config.md` restates the contract, the finalize order, and the
  environment-name composition with the reasoning for the prefix parameter; `design/web.md`
  gains the "Route groups and modules" section.
- **Integrated:** the release cadence in `concepts/database.md` — this release ships mid-effort
  so the template releases on a stable base; the remaining rungs cut `v0.5.0-dev.N` and
  postgres `v0.2.0-dev.N`, and the effort closes at `v0.5.0`.
- **Retained:** `concepts/middleware-split.md`, `concepts/module-set.md`, and
  `concepts/second-provider.md`, untouched by this session.

## Next-focus

Rung 2 of the data ladder: migrations, coordinated with the reference service. Plan the
`database/migrate` sub-module in plan mode — the golang-migrate wrapper behind a small migrator
API (up, down, steps, version, force; embedded `NNNN_<slug>.{up,down}.sql`), postgres supplying
its driver — settling the sub-module naming question first. The rung's merge cuts
`v0.5.0-dev.1` and `database/postgres/v0.2.0-dev.1`.
