package config

import (
	"regexp"
	"strings"
)

const envSeparatorPattern = `[^A-Z0-9]+`

var envSeparators = regexp.MustCompile(envSeparatorPattern)

func EnvName(prefix string, parts ...string) string {
	segments := make([]string, 0, len(parts)+1)
	if s := sanitize(prefix); s != "" {
		segments = append(segments, s)
	}
	for _, p := range parts {
		if s := sanitize(p); s != "" {
			segments = append(segments, s)
		}
	}
	return strings.Join(segments, "_")
}

func sanitize(s string) string {
	s = envSeparators.ReplaceAllString(
		strings.ToUpper(strings.TrimSpace(s)),
		"_",
	)
	return strings.Trim(s, "_")
}
