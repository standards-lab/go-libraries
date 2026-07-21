package config_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/standards-lab/go-libraries/config"
)

const (
	envSelector     = "GO_LIBRARIES_CONFIG_TEST_ENV"
	envHostOverride = "GO_LIBRARIES_CONFIG_TEST_HOST"
)

type secrets struct {
	Token string `json:"token"`
}

func (s *secrets) merge(src *secrets) {
	if src.Token != "" {
		s.Token = src.Token
	}
}

// testConfig is a synthetic configuration exercising the Config contract: a
// scalar, a scalar with a default, and a nested sub-config. Finalize follows the
// canonical order — defaults, then environment overrides, then validation.
type testConfig struct {
	Host    string  `json:"host"`
	Port    int     `json:"port"`
	Secrets secrets `json:"secrets"`
}

func (c *testConfig) Merge(src *testConfig) {
	if src.Host != "" {
		c.Host = src.Host
	}
	if src.Port != 0 {
		c.Port = src.Port
	}
	c.Secrets.merge(&src.Secrets)
}

func (c *testConfig) Finalize() error {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
	if v := os.Getenv(envHostOverride); v != "" {
		c.Host = v
	}
	if c.Port < 0 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}
	return nil
}

func writeFile(t *testing.T, dir, name, body string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o600); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
}

func TestLoad_LayerPrecedence(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{"host":"base","port":3000,"secrets":{"token":"base-token"}}`)
	writeFile(t, dir, "config.prod.json", `{"host":"prod"}`)
	writeFile(t, dir, "secrets.json", `{"secrets":{"token":"secret-token"}}`)
	writeFile(t, dir, "secrets.prod.json", `{"secrets":{"token":"prod-secret-token"}}`)
	t.Setenv(envSelector, "prod")

	cfg, err := config.Load[testConfig](config.Options{Dir: dir, EnvVar: envSelector})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "prod" {
		t.Errorf("Host = %q, want prod (overlay over base)", cfg.Host)
	}
	if cfg.Port != 3000 {
		t.Errorf("Port = %d, want 3000 (only base set it)", cfg.Port)
	}
	if cfg.Secrets.Token != "prod-secret-token" {
		t.Errorf("Token = %q, want prod-secret-token (secrets overlay wins)", cfg.Secrets.Token)
	}
}

func TestLoad_NoFilesUsesDefaults(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(envSelector, "prod") // selector set, but the directory is empty
	cfg, err := config.Load[testConfig](config.Options{Dir: dir, EnvVar: envSelector})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "localhost" || cfg.Port != 8080 || cfg.Secrets.Token != "" {
		t.Errorf("got %+v, want defaults {localhost 8080 {}}", cfg)
	}
}

func TestLoad_EmptyEnvSkipsOverlays(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{"host":"base"}`)
	writeFile(t, dir, "config.prod.json", `{"host":"prod"}`)
	t.Setenv(envSelector, "") // empty selector value: overlays are skipped

	cfg, err := config.Load[testConfig](config.Options{Dir: dir, EnvVar: envSelector})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "base" {
		t.Errorf("Host = %q, want base (overlay skipped when env is empty)", cfg.Host)
	}
}

func TestLoad_EnvOverrideWinsLast(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{"host":"base"}`)
	writeFile(t, dir, "config.prod.json", `{"host":"prod"}`)
	t.Setenv(envSelector, "prod")
	t.Setenv(envHostOverride, "from-env")

	cfg, err := config.Load[testConfig](config.Options{Dir: dir, EnvVar: envSelector})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "from-env" {
		t.Errorf("Host = %q, want from-env (Finalize env override beats every file)", cfg.Host)
	}
}

func TestLoad_MergeDoesNotClear(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{"host":"base","port":3000}`)
	writeFile(t, dir, "config.prod.json", `{"port":9000}`)
	t.Setenv(envSelector, "prod")

	cfg, err := config.Load[testConfig](config.Options{Dir: dir, EnvVar: envSelector})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "base" {
		t.Errorf("Host = %q, want base (overlay omits host, must not clear it)", cfg.Host)
	}
	if cfg.Port != 9000 {
		t.Errorf("Port = %d, want 9000 (overlay sets it)", cfg.Port)
	}
}

func TestLoad_CustomFilenameAndPattern(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "settings.json", `{"host":"base"}`)
	writeFile(t, dir, "settings-prod.json", `{"host":"prod"}`)
	t.Setenv(envSelector, "prod")

	cfg, err := config.Load[testConfig](config.Options{
		Dir:            dir,
		EnvVar:         envSelector,
		BaseName:       "settings.json",
		OverlayPattern: "%s-%s.json",
	})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Host != "prod" {
		t.Errorf("Host = %q, want prod (custom base name and overlay pattern)", cfg.Host)
	}
}

func TestLoad_MalformedJSON(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{"host": "base"`) // truncated

	_, err := config.Load[testConfig](config.Options{Dir: dir})
	if err == nil {
		t.Fatal("Load returned nil error for malformed JSON")
	}
	if !strings.Contains(err.Error(), "parse") {
		t.Errorf("error = %v, want it to mention parse", err)
	}
}

func TestLoad_FinalizeError(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "config.json", `{"port":-1}`)

	_, err := config.Load[testConfig](config.Options{Dir: dir})
	if err == nil {
		t.Fatal("Load returned nil error when Finalize failed")
	}
	if !strings.Contains(err.Error(), "finalize config") {
		t.Errorf("error = %v, want it wrapped with \"finalize config\"", err)
	}
}
