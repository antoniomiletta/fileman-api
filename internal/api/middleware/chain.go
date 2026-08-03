package middleware

import "slices"

import "net/http"

func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for _, middleware := range slices.Backward(middlewares) {
			final = middleware(final)
		}

		return final
	}
}
