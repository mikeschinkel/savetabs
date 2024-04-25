package restapi

import (
	"net/http"
	"strings"
)

// Middleware to add CORS headers to every response
func (a *API) addCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+a.Port)

		// Allow specific methods (e.g., GET, POST, OPTIONS)
		w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		}, ", "))

		// Allow specific headers (e.g., Content-Type, Authorization)
		w.Header().Set("Access-Control-Allow-Headers", strings.Join([]string{
			"Content-Type",
			"Authorization",
			"hx-request",
			"hx-target",
			"hx-trigger",
			"hx-sync",
			"hx-current-url",
		}, ", "))

		// TODO: Move to its own middleware to not bypass auth
		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
