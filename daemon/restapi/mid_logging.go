package restapi

import (
	"log/slog"
	"net/http"
	"net/url"
)

// Middleware to log request and URL
func (a *API) addRequestLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := url.QueryUnescape(r.URL.String())
		if err != nil {
			u = r.URL.String()
		}
		slog.Info("Request", "method", r.Method, "url", u)
		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
