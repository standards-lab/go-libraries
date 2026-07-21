// Package config loads layered JSON configuration and defines the contract a
// configuration type implements to take part in that load.
//
// # The contract
//
// A configuration is a type T whose pointer implements [Config]: a Merge that
// overlays another instance's set fields onto the receiver, and a Finalize that
// applies defaults, reads environment-variable overrides, and validates. The
// methods stay concretely typed against T — an implementation writes
// Merge(*T) and Finalize() error against its own type, with no type assertions.
//
// # Layered load
//
// [Load] reads up to four files from a directory, in a fixed precedence, and
// merges each one that exists onto a zero value of T:
//
//   - a base file (config.json),
//   - an environment overlay (config.<env>.json),
//   - a secrets file (secrets.json), and
//   - a secrets overlay (secrets.<env>.json).
//
// The active environment is the value of the [Options.EnvVar] variable; when it
// is empty, both overlays are skipped. A single [Options.OverlayPattern]
// produces both overlay names from the base and secrets stems. Every file is
// optional: a missing file is skipped, so a load with no files present yields a
// configuration carrying only what Finalize supplies. Any other read error, or
// malformed JSON, stops the load.
//
// Later sources win: the secrets file overrides the overlay, which overrides the
// base. Merging onto a zero value is uniform because a set source field always
// wins over an unset receiver, so the base file is not a special case. Finalize
// runs once, after every file has been merged, so the environment overrides it
// reads take precedence over every file.
//
// # Ephemeral lifecycle
//
// A configuration exists to initialize subsystems, not to be retained. A caller
// loads it, constructs subsystems from the values it carries, and discards it;
// runtime code holds the values it needs, not the configuration.
//
// # Environment-variable names
//
// [EnvName] composes an environment-variable name from a prefix and parts. A
// capability pairs its configuration with an Env struct holding the variable
// names its Finalize reads, built with EnvName so a caller can vary the prefix
// without the capability hard-coding one.
package config
