package middleware

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/storage"
	"diploma-1/internal/types"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

var unprotectedEndpoints = map[string]struct{}{
	"/api/user/login":    {},
	"/api/user/register": {},
}

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := unprotectedEndpoints[r.URL.Path]; ok {
			logger.Debugf(r.Context(), "authentication skipped for unprotected endpoint: %v", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if !strings.Contains(tokenHeader, "Bearer") {
			http.Error(w, `invalid token format: must be "Bearer <token>"`, http.StatusUnauthorized)
			return
		}

		tokenString := strings.Split(tokenHeader, " ")[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.Applied.GetJWTSecretKey(), nil
		})
		if err != nil {
			http.Error(w, "unable to parse token", http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if claims["exp"].(float64) <= float64(time.Now().Unix()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}
		user, err := storage.GetUserByLogin(r.Context(), claims["username"].(string))
		if err != nil {
			http.Error(w, "Unable to find user passed in token", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), types.CtxKeyUser, user))

		next.ServeHTTP(w, r)
	})
}
