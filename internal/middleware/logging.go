package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/sqmmm/finance-app/internal/logger"
)

// LoggingMiddleware logs all incoming requests
func Logging(tr logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return handlers.CombinedLoggingHandler(tr.Writer(), next)
	}
}
