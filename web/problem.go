package web

import (
	"encoding/json"
	"maps"
	"net/http"
)

const (
	ProblemMediaType = "application/problem+json"
	ProblemTypeBlank = "about:blank"
)

type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

func (p Problem) Write(w http.ResponseWriter) error {
	p.applyDefaults()

	w.Header().Set("Content-Type", ProblemMediaType)
	w.WriteHeader(p.Status)
	return json.NewEncoder(w).Encode(p)
}

func (p *Problem) applyDefaults() {
	if p.Type == "" {
		p.Type = ProblemTypeBlank
	}
	if p.Title == "" {
		p.Title = http.StatusText(p.Status)
	}
}

func WriteProblem(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	title, detail string,
) error {
	return Problem{
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
	}.Write(w)
}

func WriteProblemWith(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	title, detail string,
	extras map[string]any,
) error {
	p := Problem{
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: r.URL.Path,
	}
	p.applyDefaults()

	body := map[string]any{
		"type":   p.Type,
		"title":  p.Title,
		"status": p.Status,
	}
	if p.Detail != "" {
		body["detail"] = p.Detail
	}
	if p.Instance != "" {
		body["instance"] = p.Instance
	}
	maps.Copy(body, extras)

	w.Header().Set("Content-Type", ProblemMediaType)
	w.WriteHeader(p.Status)
	return json.NewEncoder(w).Encode(body)
}
