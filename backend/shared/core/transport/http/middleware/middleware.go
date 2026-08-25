package core_http_middleware

import (
	"net/http"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func Chain(h http.Handler, m ...Middleware) http.Handler {
	for _, v := range slices.Backward(m) {
		h = v(h)
	}

	return h
}
