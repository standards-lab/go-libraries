package web

import (
	"net/http"

	"github.com/standards-lab/go-libraries/lifecycle"
)

const (
	HealthPath = "/healthz"
	ReadyPath  = "/readyz"
)

type Check struct {
	Name    string
	Checker lifecycle.ReadinessChecker
}

type Mounter interface {
	Handle(pattern string, handler http.Handler)
}

type checkResult struct {
	Name  string `json:"name"`
	Ready bool   `json:"ready"`
}

type readyBody struct {
	Status string        `json:"status"`
	Checks []checkResult `json:"checks,omitempty"`
}

func Liveness() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
}

func Readiness(checks ...Check) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		results := make([]checkResult, 0, len(checks))
		ready := true
		for _, check := range checks {
			ok := check.Checker != nil && check.Checker.Ready()
			if !ok {
				ready = false
			}
			results = append(
				results,
				checkResult{Name: check.Name, Ready: ok},
			)
		}

		if !ready {
			_ = WriteProblemWith(
				w, r,
				http.StatusServiceUnavailable,
				"",
				"one or more readiness checks failed",
				map[string]any{"checks": results},
			)
			return
		}

		_ = WriteJSON(w, http.StatusOK, readyBody{
			Status: "ready",
			Checks: results,
		})
	})
}

func RegisterHealth(m Mounter, checks ...Check) {
	m.Handle("GET "+HealthPath, Liveness())
	m.Handle("GET "+ReadyPath, Readiness(checks...))
}
