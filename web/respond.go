package web

import (
	"encoding/json"
	"net/http"
)

const JSONMediaType = "application/json"

func WriteJSON(
	w http.ResponseWriter,
	status int,
	data any,
) error {
	w.Header().Set("Content-Type", JSONMediaType)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}
