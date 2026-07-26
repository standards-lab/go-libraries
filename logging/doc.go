// Package logging constructs the *slog.Logger a process writes through, from a
// configuration that takes part in the layered load.
//
// The standard library already owns the level vocabulary: slog.Level parses
// "debug", "INFO", and offsets such as "warn+2" from text, case-insensitively.
// This package adds only what slog has no opinion about — the configuration
// seam and the handler choice — so it defines no level enum of its own.
//
// # Configuration
//
// [Config] holds a [Level] and a [Format] and implements the config package's
// Merge and Finalize contract, so it loads as part of an application's
// configuration rather than on its own. Finalize applies defaults (info, text),
// then the environment overrides named by the [Env] value the configuration
// carries, then lower-cases both fields, then validates. Normalizing after the
// environment read is what lets a value arrive from a file or a shell in any
// casing and validate identically. [NewEnv] composes the standard names from a
// prefix; an empty name disables that one override, so a zero Env means no
// environment overrides at all.
//
// Level is a string rather than a slog.Level because slog.LevelInfo is the zero
// value of that type, and a configuration that merges layers has to tell "set
// to info" from "unset". An empty Level is unset; [Level.Slog] rejects it, and
// Finalize has already replaced it with a default by the time validation runs.
//
// # Construction
//
// [New] returns a logger writing to the caller's io.Writer, with
// slog.NewJSONHandler selected by [FormatJSON] and slog.NewTextHandler
// otherwise. The writer is a parameter rather than a configuration field
// because a configuration carries values and is discarded; a caller that needs
// to capture output supplies its own writer.
//
// New returns no error, because Finalize is where a configuration is validated.
// A caller that builds a Config literal and skips Finalize gets info rather
// than a failure, so a misconfiguration costs the process a level, not its
// logs.
//
// # Scope
//
// This package constructs a logger; it does not carry one. Subsystems take a
// *slog.Logger as an ordinary dependency, and the HTTP request logger lives in
// the web package, so a consumer that never serves HTTP does not compile
// net/http to get a logger.
package logging
