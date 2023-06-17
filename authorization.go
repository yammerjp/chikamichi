package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func signSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func parseToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return signSecret(), nil
	})
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
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := parseToken(tokenStr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		auds, err := token.Claims.GetAudience()
		if err != nil {
			status := http.StatusBadRequest
			m := map[string]string{"chikamichi_error": "failed to get jwt audiences"}
			if err := json.NewEncoder(w).Encode(m); err != nil {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
			return
		}
		if auds == nil || len(auds) != 1 {
			status := http.StatusBadRequest
			m := map[string]string{"chikamichi_error": "failed to get a jwt audience"}
			if err := json.NewEncoder(w).Encode(m); err != nil {
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
			return
		}

		aud := Audience(auds[0])
		ctx := context.WithValue(r.Context(), AudienceKey{}, aud)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
