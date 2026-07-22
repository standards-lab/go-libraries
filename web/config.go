package web

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/standards-lab/go-libraries/config"
)

const (
	defaultHost              = "0.0.0.0"
	defaultPort              = 8080
	defaultReadTimeout       = time.Minute
	defaultReadHeaderTimeout = 5 * time.Second
	defaultWriteTimeout      = 15 * time.Minute
	defaultIdleTimeout       = 2 * time.Minute
)

type Config struct {
	Host              string          `json:"host"`
	Port              int             `json:"port"`
	ReadTimeout       config.Duration `json:"read_timeout"`
	ReadHeaderTimeout config.Duration `json:"read_header_timeout"`
	WriteTimeout      config.Duration `json:"write_timeout"`
	IdleTimeout       config.Duration `json:"idle_timeout"`
	Env               Env             `json:"-"`
}

func (c *Config) Addr() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c *Config) Merge(src *Config) {
	if src.Host != "" {
		c.Host = src.Host
	}
	if src.Port != 0 {
		c.Port = src.Port
	}
	if src.ReadTimeout != 0 {
		c.ReadTimeout = src.ReadTimeout
	}
	if src.ReadHeaderTimeout != 0 {
		c.ReadHeaderTimeout = src.ReadHeaderTimeout
	}
	if src.WriteTimeout != 0 {
		c.WriteTimeout = src.WriteTimeout
	}
	if src.IdleTimeout != 0 {
		c.IdleTimeout = src.IdleTimeout
	}
}

func (c *Config) Finalize() error {
	c.applyDefaults()
	if err := c.applyEnv(); err != nil {
		return err
	}
	return c.validate()
}

func (c *Config) applyDefaults() {
	if c.Host == "" {
		c.Host = defaultHost
	}
	if c.Port == 0 {
		c.Port = defaultPort
	}
	if c.ReadTimeout == 0 {
		c.ReadTimeout = config.Duration(defaultReadTimeout)
	}
	if c.ReadHeaderTimeout == 0 {
		c.ReadHeaderTimeout = config.Duration(defaultReadHeaderTimeout)
	}
	if c.WriteTimeout == 0 {
		c.WriteTimeout = config.Duration(defaultWriteTimeout)
	}
	if c.IdleTimeout == 0 {
		c.IdleTimeout = config.Duration(defaultIdleTimeout)
	}
}

func (c *Config) applyEnv() error {
	if v := os.Getenv(c.Env.Host); v != "" {
		c.Host = v
	}
	if v := os.Getenv(c.Env.Port); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("%s: %w", c.Env.Port, err)
		}
		c.Port = port
	}
	if err := setDurationFromEnv(&c.ReadTimeout, c.Env.ReadTimeout); err != nil {
		return err
	}
	if err := setDurationFromEnv(&c.ReadHeaderTimeout, c.Env.ReadHeaderTimeout); err != nil {
		return err
	}
	if err := setDurationFromEnv(&c.WriteTimeout, c.Env.WriteTimeout); err != nil {
		return err
	}
	return setDurationFromEnv(&c.IdleTimeout, c.Env.IdleTimeout)
}

func (c *Config) validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	if c.ReadTimeout < 0 {
		return fmt.Errorf("invalid read_timeout: %s", c.ReadTimeout)
	}
	if c.ReadHeaderTimeout < 0 {
		return fmt.Errorf("invalid read_header_timeout: %s", c.ReadHeaderTimeout)
	}
	if c.WriteTimeout < 0 {
		return fmt.Errorf("invalid write_timeout: %s", c.WriteTimeout)
	}
	if c.IdleTimeout < 0 {
		return fmt.Errorf("invalid idle_timeout: %s", c.IdleTimeout)
	}
	return nil
}

func setDurationFromEnv(dest *config.Duration, name string) error {
	if err := dest.Set(os.Getenv(name)); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}
