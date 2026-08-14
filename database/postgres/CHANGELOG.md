# Changelog

All notable changes to the PostgreSQL provider (`github.com/standards-lab/go-libraries/database/postgres`)
are documented here. Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html). This
changelog covers this sub-module only; the base library keeps its own.

## [v0.1.0] - 2026-08-14

The first stable release: the open-and-ping provider from `v0.1.0-dev.1` against base library
`v0.4.0`. This release includes every `v0.1.0-dev.1` change; that prerelease tag is removed.

### Added

- **Open and ping**: `postgres.New` constructs the connection pool from a finalized
  `database.Config` over pgx's `database/sql` adapter — eager `pgx.ParseConfig`, so a malformed
  config is a construction error; the URL composed with `net/url` and the password set post-parse
  as a field, never entering the string; `Options` keys that name connection fields rejected; the
  connect timeout bound to the config's `conn_timeout`. The postgres dialect names itself, renders
  `$n` placeholders, and maps errors as the identity until the write surface lands.
  `postgres.Provider` types the selection constant.

### Changed

- **The base requirement is `github.com/standards-lab/go-libraries v0.4.0`**, the stable base
  this release pairs with; the capability `Finalize` signatures it carries do not affect this
  provider's API.
