package web

import "github.com/standards-lab/go-libraries/config"

type Env struct {
	Host              string
	Port              string
	ReadTimeout       string
	ReadHeaderTimeout string
	WriteTimeout      string
	IdleTimeout       string
}

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
