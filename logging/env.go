package logging

import "github.com/standards-lab/go-libraries/config"

type Env struct {
	Level  string
	Format string
}

func NewEnv(prefix string) Env {
	return Env{
		Level: config.EnvName(
			prefix, "log", "level",
		),
		Format: config.EnvName(
			prefix, "log", "format",
		),
	}
}
