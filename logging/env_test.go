package logging_test

import (
	"testing"

	"github.com/standards-lab/go-libraries/logging"
)

func TestNewEnv_ComposesNamesFromPrefix(t *testing.T) {
	env := logging.NewEnv("herald")

	for _, tc := range []struct {
		got  string
		want string
	}{
		{env.Level, "HERALD_LOG_LEVEL"},
		{env.Format, "HERALD_LOG_FORMAT"},
	} {
		if tc.got != tc.want {
			t.Errorf("got %q, want %q", tc.got, tc.want)
		}
	}
}

func TestNewEnv_EmptyPrefixDropsTheSegment(t *testing.T) {
	env := logging.NewEnv("")
	if got, want := env.Level, "LOG_LEVEL"; got != want {
		t.Errorf("Level = %q, want %q", got, want)
	}
}
