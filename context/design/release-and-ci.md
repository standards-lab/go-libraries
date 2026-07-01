# Releases and CI

How modules are versioned, released, built, and tested.

## Per-module releases

Each module is released independently by pushing a tag `<module>/v<semver>`. There is no umbrella version
for the repository. `.github/workflows/release.yml` derives the module from the tag prefix, reads that
module's `CHANGELOG.md`, and cuts a GitHub release with `taiki-e/create-gh-release-action`. The workflow
triggers on both `v*` and `*/v*`, so a future root module and the prefixed submodules are both covered.

Each module keeps its own `CHANGELOG.md` in Keep-a-Changelog form, with dated headings
(`## [vX.Y.Z] - YYYY-MM-DD`) so the release action can slice the right section.

## Coordinating a change across modules

`go.work` (committed, at the repository root) resolves sibling modules during local development, so a
change spanning modules is built and tested together before anything is tagged. Pinned `require` versions
are the committed steady state. A `replace` directive is only a transient bridge while a parent module is
unreleased; it is removed when the parent is tagged. The release ripple is bottom-up: tag the lower
module, bump the consumer's `require`, note it in the consumer's changelog, tag the consumer.

## CI

`.github/workflows/ci.yml` runs a per-module matrix: `go vet`, then `go test -race`, then golangci-lint,
each with `working-directory` set to the module. The repository is public, so modules resolve through the
normal Go proxy and checksum database; CI carries no `GOPRIVATE` or `.netrc` configuration.

The module matrix, the `go.work` use-list, and `mise`'s `GO_MODULES` are three lists that name the same
modules; they are updated together whenever a module is added. The CI matrix is seeded with `core` and its
steps are guarded on the module's `go.mod`, so it stays green until `core` lands.

## Tasks

`mise.toml` defines the developer tasks — `build`, `test`, `vet`, `fmt`, `tidy`, `lint` — each looping
over `GO_MODULES`. `mise run test` builds and tests every module.
