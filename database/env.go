package database

import "github.com/standards-lab/go-libraries/config"

// Env names the environment variables [Config.Finalize] reads. An empty name
// disables that one override; the zero value disables all of them.
type Env struct {
	Host            string
	Name            string
	User            string
	Password        string
	Port            string
	MaxOpenConns    string
	MaxIdleConns    string
	ConnMaxLifetime string
	ConnMaxIdleTime string
	ConnTimeout     string
}

// NewEnv composes the standard override names from a prefix under the
// "database" segment: DATABASE_HOST, DATABASE_PORT, and the rest, prefixed by
// whatever [config.EnvName] produces.
func NewEnv(prefix string) Env {
	return Env{
		Host: config.EnvName(
			prefix, "database", "host",
		),
		Name: config.EnvName(
			prefix, "database", "name",
		),
		User: config.EnvName(
			prefix, "database", "user",
		),
		Password: config.EnvName(
			prefix, "database", "password",
		),
		Port: config.EnvName(
			prefix, "database", "port",
		),
		MaxOpenConns: config.EnvName(
			prefix, "database", "max", "open", "conns",
		),
		MaxIdleConns: config.EnvName(
			prefix, "database", "max", "idle", "conns",
		),
		ConnMaxLifetime: config.EnvName(
			prefix, "database", "conn", "max", "lifetime",
		),
		ConnMaxIdleTime: config.EnvName(
			prefix, "database", "conn", "max", "idle", "time",
		),
		ConnTimeout: config.EnvName(
			prefix, "database", "conn", "timeout",
		),
	}
}
