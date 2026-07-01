# Library topology and naming

How the Go reference libraries are organized and named.

## A single monorepo

All Go capability modules live in one repository, `go-libraries`, each module versioned and released on
its own schedule. Keeping them in one repository makes cross-module changes atomic: they are made and
exercised together (through `go.work`) before any tag is cut, against a single CI and release pipeline. It
also keeps the organization from fragmenting into a repository per library while still demonstrating
independent, per-module versioning.

## Naming

- **Repository:** `go-libraries` under the `standards-lab` organization.
- **Module path:** the directory path appended to the repository root —
  `github.com/standards-lab/go-libraries/<module>`. The directory name is the module-path suffix.
- **Capability modules** are flat top-level directories: `core`, `auth`, `database`, `web`, and so on.
- **Vendor implementations** are nested directories that are themselves modules, named for the **target
  system**, not the SDK: `auth/keycloak`, `auth/entra`, `database/postgres`. Naming by target system lets
  the underlying SDK change without renaming the module.
- **Release tags:** `<module>/v<semver>`, the prefix being the full module directory path — `core/v0.1.0`,
  `auth/v0.2.1`, `auth/keycloak/v0.1.0`.

## The cost and the benefit

The monorepo's cost is the tag-prefix namespacing and a changelog per module; its benefit is atomic
cross-module PRs and a single pipeline. Each module remains independently consumable: a consumer takes one
module without pulling the rest.
