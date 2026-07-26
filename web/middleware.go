package web

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middleware ...Middleware) http.Handler {
	for _, m := range slices.Backward(middleware) {
		if m == nil {
			continue
		}
		handler = m(handler)
	}
	return handler
}
