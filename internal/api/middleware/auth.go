package middleware

import "net/http"

func Auth(next http.Handler) http.Handler
