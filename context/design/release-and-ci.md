# Releases and CI

How the libraries are versioned, released, built, and tested.

## Releases: the base library and each provider

There are two kinds of releasable artifact. The **base library** is released by pushing a root tag
`v<semver>`; that single version covers every capability package in it. Each **provider sub-module** is
released by pushing a tag `<path>/v<semver>` (`database/postgres/v0.1.0`). There is no umbrella version
spanning the base and the providers — a consumer pulling one provider is not entangled with another's
version.

`.github/workflows/release.yml` derives the artifact from the tag: a bare `v*` tag reads the root
`CHANGELOG.md`; a `<path>/v*` tag reads that sub-module's `CHANGELOG.md`. It cuts a GitHub release with
`taiki-e/create-gh-release-action`. Each artifact keeps its own `CHANGELOG.md` in Keep-a-Changelog form
with dated headings (`## [vX.Y.Z] - YYYY-MM-DD`) so the release action can slice the right section: the
root `CHANGELOG.md` for the base library, one per provider sub-module.

## Coordinating a change across the base and a provider

`go.work` (committed, at the repository root) resolves the base module and the provider sub-modules during
local development, so a change spanning them is built and tested together before anything is tagged.
Because the whole base is one module, a change touching several capability packages needs no internal
version coordination — it releases as one base version.

A provider depends on the base through a pinned `require`. Pinned versions are the committed steady state;
a `replace` directive is only a transient bridge while the base carries unreleased changes a provider
needs, removed once the base is tagged. The release ripple is bottom-up across the base→provider edge: tag
the base, bump the provider's `require`, note it in the provider's changelog, tag the provider.

## CI

`.github/workflows/ci.yml` runs a per-module matrix: `go vet`, then `go test -race`, then golangci-lint,
each with `working-directory` set to the module. The repository is public, so modules resolve through the
normal Go proxy and checksum database; CI carries no `GOPRIVATE` or `.netrc` configuration.

The CI matrix, the `go.work` use-list, and `mise`'s `GO_MODULES` are three lists that name the same
modules — the base module (`.`) and each provider sub-module — and are updated together whenever a module
is added. Each matrix entry's steps are guarded on the module's `go.mod`, so an entry stays green until
that module lands.

## Tasks

`mise.toml` defines the developer tasks — `build`, `test`, `vet`, `fmt`, `tidy`, `lint` — each looping
over `GO_MODULES`. `mise run test` builds and tests the base library and every provider sub-module.
