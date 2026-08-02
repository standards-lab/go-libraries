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

// Config holds the server's address and timeouts. The port and timeouts are
// tri-state pointers: nil is unset and takes the default, while an explicit
// zero survives the load and means what it says — a disabled timeout, or an
// ephemeral port. Env names the environment variables Finalize reads; it is
// excluded from JSON, so only the composition root can set it.
type Config struct {
	Host              string           `json:"host"`
	Port              *int             `json:"port"`
	ReadTimeout       *config.Duration `json:"read_timeout"`
	ReadHeaderTimeout *config.Duration `json:"read_header_timeout"`
	WriteTimeout      *config.Duration `json:"write_timeout"`
	IdleTimeout       *config.Duration `json:"idle_timeout"`
	Env               Env              `json:"-"`
}

// Addr joins the host and port; an unset port reads as 0.
func (c *Config) Addr() string {
	port := 0
	if c.Port != nil {
		port = *c.Port
	}
	return net.JoinHostPort(c.Host, strconv.Itoa(port))
}

// Merge overlays src's set fields onto the receiver.
func (c *Config) Merge(src *Config) {
	if src.Host != "" {
		c.Host = src.Host
	}
	if src.Port != nil {
		c.Port = src.Port
	}
	if src.ReadTimeout != nil {
		c.ReadTimeout = src.ReadTimeout
	}
	if src.ReadHeaderTimeout != nil {
		c.ReadHeaderTimeout = src.ReadHeaderTimeout
	}
	if src.WriteTimeout != nil {
		c.WriteTimeout = src.WriteTimeout
	}
	if src.IdleTimeout != nil {
		c.IdleTimeout = src.IdleTimeout
	}
}

// Finalize applies defaults, applies the environment overrides named by Env,
// and validates.
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
	if c.Port == nil {
		c.Port = intPtr(defaultPort)
	}
	if c.ReadTimeout == nil {
		c.ReadTimeout = durationPtr(defaultReadTimeout)
	}
	if c.ReadHeaderTimeout == nil {
		c.ReadHeaderTimeout = durationPtr(defaultReadHeaderTimeout)
	}
	if c.WriteTimeout == nil {
		c.WriteTimeout = durationPtr(defaultWriteTimeout)
	}
	if c.IdleTimeout == nil {
		c.IdleTimeout = durationPtr(defaultIdleTimeout)
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
		c.Port = &port
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
	if *c.Port < 0 || *c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", *c.Port)
	}
	if *c.ReadTimeout < 0 {
		return fmt.Errorf("invalid read_timeout: %s", c.ReadTimeout)
	}
	if *c.ReadHeaderTimeout < 0 {
		return fmt.Errorf("invalid read_header_timeout: %s", c.ReadHeaderTimeout)
	}
	if *c.WriteTimeout < 0 {
		return fmt.Errorf("invalid write_timeout: %s", c.WriteTimeout)
	}
	if *c.IdleTimeout < 0 {
		return fmt.Errorf("invalid idle_timeout: %s", c.IdleTimeout)
	}
	return nil
}

func (c *Config) finalized() bool {
	return c.Port != nil &&
		c.ReadTimeout != nil &&
		c.ReadHeaderTimeout != nil &&
		c.WriteTimeout != nil &&
		c.IdleTimeout != nil
}

func setDurationFromEnv(dest **config.Duration, name string) error {
	v := os.Getenv(name)
	if v == "" {
		return nil
	}
	var d config.Duration
	if err := d.Set(v); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	*dest = &d
	return nil
}

func intPtr(v int) *int {
	return &v
}

func durationPtr(v time.Duration) *config.Duration {
	d := config.Duration(v)
	return &d
}
