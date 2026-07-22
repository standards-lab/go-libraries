package web_test

import (
	"testing"

	"github.com/standards-lab/go-libraries/web"
)

func TestNewEnv_ComposesNamesFromPrefix(t *testing.T) {
	env := web.NewEnv("herald")

	for _, tc := range []struct {
		got  string
		want string
	}{
		{env.Host, "HERALD_SERVER_HOST"},
		{env.Port, "HERALD_SERVER_PORT"},
		{env.ReadTimeout, "HERALD_SERVER_READ_TIMEOUT"},
		{env.ReadHeaderTimeout, "HERALD_SERVER_READ_HEADER_TIMEOUT"},
		{env.WriteTimeout, "HERALD_SERVER_WRITE_TIMEOUT"},
		{env.IdleTimeout, "HERALD_SERVER_IDLE_TIMEOUT"},
	} {
		if tc.got != tc.want {
			t.Errorf("got %q, want %q", tc.got, tc.want)
		}
	}
}

func TestNewEnv_EmptyPrefixDropsTheSegment(t *testing.T) {
	env := web.NewEnv("")
	if got, want := env.Port, "SERVER_PORT"; got != want {
		t.Errorf("Port = %q, want %q", got, want)
	}
}
