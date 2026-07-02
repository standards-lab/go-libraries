# Library topology and naming

How the Go reference libraries are organized and named.

## One base library, providers as sub-modules

The repository releases as a single base library plus a set of provider sub-modules. The base library is
one Go module rooted at the repository, and every capability lives inside it as its own package —
`lifecycle`, `config`, `auth`, `database`, `storage`, `web`, and so on. These packages version and
release together as one artifact, because they co-evolve and depend on each other; keeping them in one
module makes those inter-package edges ordinary imports with no version coordination.

A concrete provider whose weight comes from a third-party SDK is a nested sub-module with its own
`go.mod`, versioned and released on its own schedule. Providers are where independent versioning earns its
cost: each pins a distinct, heavy dependency surface, and a consumer selects one without pulling the
others. The base library carries none of that weight (see `conventions.md`), so depending on it to reach
one capability is effectively free even though it contains the others.

Keeping everything in one repository makes cross-cutting changes atomic — a change spanning the base and a
provider is made and exercised together through `go.work` before any tag is cut.

## Naming

- **Repository / base module:** `github.com/standards-lab/go-libraries` under the `standards-lab`
  organization. The base module is rooted here; its capabilities are packages —
  `github.com/standards-lab/go-libraries/database`, `.../lifecycle`, and so on.
- **Provider sub-modules** are nested directories that are themselves modules, named for the **target API
  or system**, not the SDK: `database/postgres`, `database/mssql`, `storage/s3`, `storage/azureblob`,
  `auth/keycloak`, `auth/entra`. Naming by target system lets the underlying SDK change without renaming
  the module. Each sub-module's path is its directory appended to the base — e.g.
  `github.com/standards-lab/go-libraries/database/postgres`.
- **Capability interfaces** are named per package for what reads best there — `database.DB`,
  `storage.Store`, `auth.Authenticator`/`auth.TokenSource` — not a forced uniform noun (see
  `conventions.md`).
- **Release tags:** the base library is tagged `v<semver>` at the repository root (`v0.1.0`); each
  provider sub-module is tagged with its directory path as prefix, `<path>/v<semver>`
  (`database/postgres/v0.1.0`, `auth/entra/v0.2.1`).

## The cost and the benefit

The topology's cost is a changelog and tag namespace per released artifact — the base plus each provider —
and the discipline of keeping the base's inter-package dependencies acyclic (which Go enforces at compile
time). Its benefit is that the tightly co-evolving base moves as one coherent version with zero internal
release ripple, while each provider stays independently consumable: a consumer takes the base and one
provider without pulling the rest.
