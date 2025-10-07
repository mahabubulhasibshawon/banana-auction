package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"banana-auction/config"

	"github.com/golang-jwt/jwt/v5"
)

var userIDKey = "userID"

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().JwtSecretKey), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid claims", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user_id", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), userIDKey, int(userIDFloat)))
		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in request")
	}
	return userID, nil
}
