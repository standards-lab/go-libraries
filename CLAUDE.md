# go-libraries

The standards-lab organization's Go reference libraries: a public, multi-module monorepo of layered,
independently versioned capability libraries. This is the library level of the reference-architecture
effort — the worked example for how to design, layer, and release shared libraries. Managed with the
marathon workflow; start from `context/README.md`.

## Conventions are settled in the repository

The design and conventions for these libraries are recorded in `context/design/` — that is the authority.
Keep them there; do not restate them here.

## Role boundary

go-libraries is a marathon **code** project (`.claude/marathon.toml` declares `kind = "code"`). The
developer owns the production Go source — they apply it and answer for it. The agent writes everything
else: tests, godoc and `doc.go`, prose documentation, the files in `context/`, the implementation guide,
and the reset file.

## Repository specifics

- **Module layout** — the repository is one base library (a single module rooted at
  `github.com/standards-lab/go-libraries`) whose capabilities are packages, plus provider sub-modules.
  Only vendor implementations are nested submodules with their own `go.mod`, named for the target system
  (`auth/keycloak`, `database/postgres`), not the SDK.
- **Local development** uses the committed root `go.work`; pinned `require` versions are the committed
  steady state. A `replace` directive is only a transient bridge while the base carries unreleased changes
  a provider needs, removed when the base is tagged.
- **Releases** — the base library is tagged `v<semver>` at the root from the root `CHANGELOG.md`; each
  provider sub-module is tagged `<path>/v<semver>` from its own `CHANGELOG.md`, cut by
  `.github/workflows/release.yml`.
- **Tests** are co-located `{file}_test.go` files in an external black-box package (`package <pkg>_test`)
  that exercise the public API.
- **Tasks** run through `mise` (`build`, `test`, `vet`, `fmt`, `tidy`, `lint`).
- **Public repo.** Modules resolve through the public Go proxy; CI carries no private-module config.
