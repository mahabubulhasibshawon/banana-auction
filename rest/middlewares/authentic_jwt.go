package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"banana-auction/config"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const (
	userIDKey ctxKey = "userID"
)

func JwtAuthenticationMiddleware(next http.Handler) http.Handler {
	cfg := config.GetConfig()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header: Expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JwtSecretKey), nil
		})

		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JWT: %v", err), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid JWT: Token is not valid", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid JWT: Unable to parse claims", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid JWT: user_id claim missing or invalid", http.StatusUnauthorized)
			return
		}

		userID := int(userIDFloat)
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
