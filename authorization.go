package main

import (
	"context"
	"net/http"
	"strings"
)

type Jwt struct {
	Issuer     string
	Subject    string
	Audience   string
	Expiration int64
	NotBefore  int64
	IssuedAt   int64
	JwtID      string
	Scopes     []string
}

func validToken(token string) bool {
	return token == "hello"
}

func extractAudience(token string) Audience {
	return Audience("audience:" + token)
}

type AudienceKey struct{}

type Audience string

func WithAuthorization(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if !validToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		audience := extractAudience(token)
		ctx := context.WithValue(r.Context(), AudienceKey{}, audience)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
