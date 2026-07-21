package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBaseName       = "config.json"
	defaultSecretsName    = "secrets.json"
	defaultOverlayPattern = "%s.%s.json"
)

type Config[T any] interface {
	*T
	Merge(src *T)
	Finalize() error
}

type Options struct {
	Dir            string
	EnvVar         string
	BaseName       string
	SecretsName    string
	OverlayPattern string
}

func (o *Options) withDefaults() {
	if o.Dir == "" {
		o.Dir = "."
	}
	if o.BaseName == "" {
		o.BaseName = defaultBaseName
	}
	if o.SecretsName == "" {
		o.SecretsName = defaultSecretsName
	}
	if o.OverlayPattern == "" {
		o.OverlayPattern = defaultOverlayPattern
	}
}

func (o Options) overlay(name, env string) string {
	stem := strings.TrimSuffix(name, filepath.Ext(name))
	return fmt.Sprintf(o.OverlayPattern, stem, env)
}

func Load[T any, PT Config[T]](opts Options) (PT, error) {
	opts.withDefaults()

	var env string
	if opts.EnvVar != "" {
		env = os.Getenv(opts.EnvVar)
	}

	names := []string{opts.BaseName}
	if env != "" {
		names = append(names, opts.overlay(opts.BaseName, env))
	}
	names = append(names, opts.SecretsName)
	if env != "" {
		names = append(names, opts.overlay(opts.SecretsName, env))
	}

	cfg := PT(new(T))
	for _, name := range names {
		path := filepath.Join(opts.Dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		layer := new(T)
		if err := json.Unmarshal(data, layer); err != nil {
			return nil, fmt.Errorf("parse %s: %w", path, err)
		}
		cfg.Merge(layer)
	}

	if err := cfg.Finalize(); err != nil {
		return nil, fmt.Errorf("finalize config: %w", err)
	}
	return cfg, nil
}
