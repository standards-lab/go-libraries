package config

import "strings"

func EnvName(prefix string, parts ...string) string {
	segments := make([]string, 0, len(parts)+1)
	if s := strings.ToUpper(strings.TrimSpace(prefix)); s != "" {
		segments = append(segments, s)
	}
	for _, p := range parts {
		if s := strings.ToUpper(strings.TrimSpace(p)); s != "" {
			segments = append(segments, s)
		}
	}
	return strings.Join(segments, "_")
}
