package middleware

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/sqmmm/finance-app/internal/logger"
)

// TrackerMiddleware adds trackerID to the context
func Tracker(manager logger.LoggerManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := rand.Int63()
			ctx := manager.SetRequestID(r.Context(), requestID)

			w.Header().Set("Request-ID", strconv.Itoa(int(requestID)))

			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
