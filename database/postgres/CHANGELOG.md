# Changelog

All notable changes to the PostgreSQL provider (`github.com/standards-lab/go-libraries/database/postgres`)
are documented here. Versions follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html). This
changelog covers this sub-module only; the base library keeps its own.

## [v0.1.0-dev.1] - 2026-08-13

### Added

- **Open and ping**: `postgres.New` constructs the connection pool from a finalized
  `database.Config` over pgx's `database/sql` adapter — eager `pgx.ParseConfig`, so a malformed
  config is a construction error; the URL composed with `net/url` and the password set post-parse
  as a field, never entering the string; `Options` keys that name connection fields rejected; the
  connect timeout bound to the config's `conn_timeout`. The postgres dialect names itself, renders
  `$n` placeholders, and maps errors as the identity until the write surface lands.
  `postgres.Provider` types the selection constant.
