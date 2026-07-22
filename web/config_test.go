package web_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/standards-lab/go-libraries/config"
	"github.com/standards-lab/go-libraries/web"
)

// Instantiating Load proves *web.Config satisfies the config.Config contract at
// compile time; the constraint cannot be written as an ordinary interface
// assertion because it carries a type element.
var _ = config.Load[web.Config]

const (
	envHost              = "GO_LIBRARIES_WEB_TEST_HOST"
	envPort              = "GO_LIBRARIES_WEB_TEST_PORT"
	envReadTimeout       = "GO_LIBRARIES_WEB_TEST_READ_TIMEOUT"
	envReadHeaderTimeout = "GO_LIBRARIES_WEB_TEST_READ_HEADER_TIMEOUT"
	envWriteTimeout      = "GO_LIBRARIES_WEB_TEST_WRITE_TIMEOUT"
	envIdleTimeout       = "GO_LIBRARIES_WEB_TEST_IDLE_TIMEOUT"
)

func testEnv() web.Env {
	return web.Env{
		Host:              envHost,
		Port:              envPort,
		ReadTimeout:       envReadTimeout,
		ReadHeaderTimeout: envReadHeaderTimeout,
		WriteTimeout:      envWriteTimeout,
		IdleTimeout:       envIdleTimeout,
	}
}

func TestConfig_MergeSourceWinsWithoutClearing(t *testing.T) {
	base := web.Config{
		Host:        "base",
		Port:        3000,
		ReadTimeout: config.Duration(time.Minute),
	}
	base.Merge(&web.Config{Port: 9000})

	if base.Host != "base" {
		t.Errorf("Host = %q, want base (the source omits it and must not clear it)", base.Host)
	}
	if base.Port != 9000 {
		t.Errorf("Port = %d, want 9000 (the source sets it)", base.Port)
	}
	if time.Duration(base.ReadTimeout) != time.Minute {
		t.Errorf("ReadTimeout = %s, want 1m (the source omits it)", base.ReadTimeout)
	}
}

func TestConfig_FinalizeAppliesDefaults(t *testing.T) {
	var cfg web.Config
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("Finalize: %v", err)
	}

	if cfg.Host != "0.0.0.0" || cfg.Port != 8080 {
		t.Errorf("Addr() = %q, want 0.0.0.0:8080", cfg.Addr())
	}
	for _, want := range []struct {
		name  string
		got   config.Duration
		value time.Duration
	}{
		{"ReadTimeout", cfg.ReadTimeout, time.Minute},
		{"ReadHeaderTimeout", cfg.ReadHeaderTimeout, 5 * time.Second},
		{"WriteTimeout", cfg.WriteTimeout, 15 * time.Minute},
		{"IdleTimeout", cfg.IdleTimeout, 2 * time.Minute},
	} {
		if time.Duration(want.got) != want.value {
			t.Errorf("%s = %s, want %s", want.name, want.got, want.value)
		}
	}
}

func TestConfig_FinalizeEnvOverridesFiles(t *testing.T) {
	t.Setenv(envHost, "from-env")
	t.Setenv(envPort, "9443")
	t.Setenv(envReadTimeout, "45s")

	cfg := web.Config{Host: "from-file", Port: 8080, Env: testEnv()}
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("Finalize: %v", err)
	}

	if cfg.Host != "from-env" {
		t.Errorf("Host = %q, want from-env", cfg.Host)
	}
	if cfg.Port != 9443 {
		t.Errorf("Port = %d, want 9443", cfg.Port)
	}
	if time.Duration(cfg.ReadTimeout) != 45*time.Second {
		t.Errorf("ReadTimeout = %s, want 45s", cfg.ReadTimeout)
	}
	if time.Duration(cfg.WriteTimeout) != 15*time.Minute {
		t.Errorf("WriteTimeout = %s, want the default (no override set)", cfg.WriteTimeout)
	}
}

func TestConfig_FinalizeEmptyEnvValueLeavesConfigured(t *testing.T) {
	// An empty variable reads as unset, not as a request to clear the value.
	t.Setenv(envReadTimeout, "")

	cfg := web.Config{ReadTimeout: config.Duration(30 * time.Second), Env: testEnv()}
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("Finalize: %v", err)
	}
	if time.Duration(cfg.ReadTimeout) != 30*time.Second {
		t.Errorf("ReadTimeout = %s, want 30s", cfg.ReadTimeout)
	}
}

func TestConfig_FinalizeZeroEnvSkipsOverrides(t *testing.T) {
	// A zero Env names no variables, so nothing in the environment applies —
	// which is what makes config.Load[web.Config] usable on its own.
	t.Setenv(envHost, "from-env")
	t.Setenv(envPort, "9443")

	var cfg web.Config
	if err := cfg.Finalize(); err != nil {
		t.Fatalf("Finalize: %v", err)
	}
	if cfg.Host != "0.0.0.0" || cfg.Port != 8080 {
		t.Errorf("Addr() = %q, want the defaults with a zero Env", cfg.Addr())
	}
}

func TestConfig_FinalizeMalformedDurationNamesTheVariable(t *testing.T) {
	t.Setenv(envWriteTimeout, "1hour")

	cfg := web.Config{Env: testEnv()}
	err := cfg.Finalize()
	if err == nil {
		t.Fatal("Finalize returned nil for a malformed duration override")
	}
	if !strings.Contains(err.Error(), envWriteTimeout) {
		t.Errorf("error = %v, want it to name %s", err, envWriteTimeout)
	}
}

func TestConfig_FinalizeMalformedPortNamesTheVariable(t *testing.T) {
	t.Setenv(envPort, "http")

	cfg := web.Config{Env: testEnv()}
	err := cfg.Finalize()
	if err == nil {
		t.Fatal("Finalize returned nil for a malformed port override")
	}
	if !strings.Contains(err.Error(), envPort) {
		t.Errorf("error = %v, want it to name %s", err, envPort)
	}
}

func TestConfig_FinalizeRejectsPortOutOfRange(t *testing.T) {
	cfg := web.Config{Port: 70000}
	if err := cfg.Finalize(); err == nil {
		t.Fatal("Finalize returned nil for port 70000")
	}
}

func TestConfig_FinalizeRejectsNegativeTimeout(t *testing.T) {
	cfg := web.Config{ReadTimeout: config.Duration(-time.Second)}
	if err := cfg.Finalize(); err == nil {
		t.Fatal("Finalize returned nil for a negative read_timeout")
	}
}

func TestConfig_AddrBracketsIPv6(t *testing.T) {
	cfg := web.Config{Host: "::1", Port: 8080}
	if got, want := cfg.Addr(), "[::1]:8080"; got != want {
		t.Errorf("Addr() = %q, want %q", got, want)
	}
}

func TestConfig_LoadsThroughConfigLoad(t *testing.T) {
	dir := t.TempDir()
	body := `{"host":"127.0.0.1","port":9000,"read_timeout":"90s"}`
	if err := os.WriteFile(filepath.Join(dir, "config.json"), []byte(body), 0o600); err != nil {
		t.Fatalf("write config.json: %v", err)
	}

	cfg, err := config.Load[web.Config](config.Options{Dir: dir})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if got, want := cfg.Addr(), "127.0.0.1:9000"; got != want {
		t.Errorf("Addr() = %q, want %q", got, want)
	}
	if time.Duration(cfg.ReadTimeout) != 90*time.Second {
		t.Errorf("ReadTimeout = %s, want 1m30s", cfg.ReadTimeout)
	}
	if time.Duration(cfg.IdleTimeout) != 2*time.Minute {
		t.Errorf("IdleTimeout = %s, want the default", cfg.IdleTimeout)
	}
}
