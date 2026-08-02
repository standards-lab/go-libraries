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
// produces both overlay names from the base and secrets stems; Load validates
// the pattern and fails a malformed one rather than silently skipping the
// overlays it names. Every file is optional: a missing file is skipped, so a
// load with no files present yields a configuration carrying only what Finalize
// supplies. Any other read error, or malformed JSON, stops the load.
//
// Later sources win: the secrets overlay overrides the secrets file, which
// overrides the environment overlay, which overrides the base. A set source
// field always wins over an unset receiver. Finalize runs once, after every
// file has been merged.
//
// # Environment overrides
//
// A capability pairs its configuration with an Env struct naming the variables
// its Finalize reads. Env fields are excluded from JSON, so a file cannot set
// them: the composition root seeds a child's Env before the child's Finalize
// runs, typically in the parent configuration's own Finalize. A direct [Load]
// of a capability's configuration therefore applies no environment overrides;
// once a parent has seeded the names, the overrides Finalize reads take
// precedence over every file.
//
// [EnvName] composes an environment-variable name from a prefix and parts. Each
// segment upper-cases, and runs of characters outside A-Z and 0-9 collapse to
// single underscores ("db host" becomes DB_HOST), so a caller can vary the
// prefix without the capability hard-coding one.
//
// # Durations
//
// [Duration] is a time.Duration that takes part in JSON configuration. It
// unmarshals from the string form time.ParseDuration accepts ("1m30s") or from
// a bare integer number of nanoseconds; a fractional number is rejected — the
// string form expresses sub-second durations — and JSON null leaves the value
// unchanged. It marshals as the string form, and [Duration.Duration] returns
// the value as a time.Duration where an API wants the standard type.
//
// # Ephemeral lifecycle
//
// A configuration initializes subsystems and is then discarded: a caller loads
// it, constructs subsystems from its values, and lets it go. Runtime code holds
// the values it needs, not the configuration.
package config
