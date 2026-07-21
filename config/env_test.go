package config_test

import (
	"testing"

	"github.com/standards-lab/go-libraries/config"
)

func TestEnvName(t *testing.T) {
	cases := []struct {
		name   string
		prefix string
		parts  []string
		want   string
	}{
		{"prefix and parts", "app", []string{"db", "host"}, "APP_DB_HOST"},
		{"prefix only", "app", nil, "APP"},
		{"empty part dropped", "app", []string{"", "host"}, "APP_HOST"},
		{"whitespace trimmed", "app", []string{"  db  ", "host"}, "APP_DB_HOST"},
		{"empty prefix dropped", "", []string{"db"}, "DB"},
		{"all empty", "", nil, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := config.EnvName(tc.prefix, tc.parts...); got != tc.want {
				t.Errorf("EnvName(%q, %v) = %q, want %q", tc.prefix, tc.parts, got, tc.want)
			}
		})
	}
}
