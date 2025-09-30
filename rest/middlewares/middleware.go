package middlewares

import (
	"context"
	"errors"
)

var (
	JwtAuthMiddleware    = JwtAuthenticationMiddleware
	CorsMiddleware       = CorsHandlerMiddleware
	GetUserIDFromContext = getUserIDFromContext
)

func getUserIDFromContext(ctx context.Context) (int, error) {
	userID, ok := ctx.Value(userIDKey).(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
