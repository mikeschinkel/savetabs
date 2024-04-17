package restapi

import (
	"log/slog"
	"net/http"
)

// Middleware to log request and URL
func (a *API) addRequestLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request", "method", r.Method, "url", r.URL)
		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
