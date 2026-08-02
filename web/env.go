package web

import "github.com/standards-lab/go-libraries/config"

// Env names the environment variables [Config.Finalize] reads. An empty name
// disables that one override; the zero value disables all of them.
type Env struct {
	Host              string
	Port              string
	ReadTimeout       string
	ReadHeaderTimeout string
	WriteTimeout      string
	IdleTimeout       string
}

// NewEnv composes the standard override names from a prefix: SERVER_HOST,
// SERVER_PORT, and the four timeout names under whatever prefix
// [config.EnvName] produces.
func NewEnv(prefix string) Env {
	return Env{
		Host: config.EnvName(
			prefix, "server", "host",
		),
		Port: config.EnvName(
			prefix, "server", "port",
		),
		ReadTimeout: config.EnvName(
			prefix, "server", "read", "timeout",
		),
		ReadHeaderTimeout: config.EnvName(
			prefix, "server", "read", "header", "timeout",
		),
		WriteTimeout: config.EnvName(
			prefix, "server", "write", "timeout",
		),
		IdleTimeout: config.EnvName(
			prefix, "server", "idle", "timeout",
		),
	}
}
