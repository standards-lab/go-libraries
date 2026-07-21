# reset · build-config-package

- **Status:** closeout
- **Session:** start
- **Branch:** build-config-package

## Disposition

- **Built the base library's `config` package** (`config/`): a generic `Load[T]` that layers a base file,
  an environment overlay, `secrets.json`, and a secrets overlay onto a caller's config type, driven by the
  `Config[T]` constraint (`Merge`/`Finalize`); an `EnvName` helper composes override-variable names.
  Black-box tests (`-race`, 9 cases), `doc.go`, and a `config` entry appended to the **unreleased**
  `v0.1.0` `CHANGELOG.md`. It lands as a package of the existing base module — no new module, no
  `go.work`/`mise`/`ci` synced-list change.
- **Promoted** the settled config design into a new `design/config.md`: the loader-in-the-base-library
  posture (vs. the baseline's per-consumer hand-rolled `Load`), the generic `Load`/`Config[T]` contract,
  the `Merge` rule (non-empty-source-wins, reflection-free, correct-from-zero) and the canonical
  `Finalize` order (defaults → env → validate, once, after all layers, failing on a bad env value), the
  `EnvName`/`Env` convention, and config's ephemeral lifecycle. This resolves the baseline inconsistencies
  the re-derivation set out to unify (Finalize ordering, `Env` naming, fail-vs-swallow on bad env).
- **Integrated** the `config` bullet in `concepts/module-set.md` — marked built, pointed at the code and
  `design/config.md`, and dropped the now-settled "config package shape" open question.
- **Retained:** the rest of `concepts/module-set.md` (auth, database, storage, web, and the provider set)
  and its remaining open questions — all still unbuilt, settled when each capability is reached.

## Next-focus

Build a **minimal `web`** package — the third step in the README build order (lifecycle → config → web).
Start narrow: the `net/http` bootstrap plus `/healthz` and `/readyz`, with `/readyz` surfacing the
`lifecycle` readiness signal (`ReadinessChecker`) and the server configured through a `web` config that
implements the `config` contract. Defer problem responses, the success envelope, middleware, and the
authorization enforcement point to later steps. Settle the `web` config shape and the server's lifecycle
wiring as it is built. Start here next session with `marathon start`.
