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
		{"space becomes separator", "app", []string{"db host"}, "APP_DB_HOST"},
		{"hyphen becomes separator", "app", []string{"read-timeout"}, "APP_READ_TIMEOUT"},
		{"prefix separators trimmed", "-app-", []string{"x"}, "APP_X"},
		{"separator run collapses", "app", []string{"a--b"}, "APP_A_B"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := config.EnvName(tc.prefix, tc.parts...); got != tc.want {
				t.Errorf("EnvName(%q, %v) = %q, want %q", tc.prefix, tc.parts, got, tc.want)
			}
		})
	}
}
