package middleware

import (
	"context"
	"net/http"

	"github.com/Nezent/go-queue/common"
)

const (
	UserIDKey ctxKey = "userID"
	RoleKey   ctxKey = "userRole"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			common.RespondJSON(w, http.StatusUnauthorized, common.ErrorResponse("Unauthorized - token missing"))
			return
		}

		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			common.RespondJSON(w, http.StatusUnauthorized, common.ErrorResponse("Unauthorized - invalid token prefix"))
			return
		}

		token := authHeader[len(prefix):]

		userID, role, err := common.ParseJWT(token)
		if err != nil {
			common.RespondJSON(w, http.StatusUnauthorized, common.ErrorResponse("Unauthorized - invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, RoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extracts the user ID from context
func GetUserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(UserIDKey).(string)
	return id, ok
}

// GetUserRole extracts the user role from context
func GetUserRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}
