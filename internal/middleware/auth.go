package middleware

import (
	"context"
	"github.com/sqmmm/finance-app/entities"
	"github.com/sqmmm/finance-app/internal/logger"
	"net/http"
)

const UserKey = "user"

// LoggingMiddleware logs all incoming requests
func Auth(manager logger.LoggerManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//todo: check auth
			ctx := context.WithValue(r.Context(), UserKey, &entities.User{ID: 1})
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
