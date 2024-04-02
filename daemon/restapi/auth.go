package restapi

import (
	"net/http"
)

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Perform authentication logic here
		if !isAuthenticated(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// If authenticated, call the next handler
		next.ServeHTTP(w, r)
	})
}

func isAuthenticated(r *http.Request) bool {
	// TODO: Implement authentication logic, such as checking for valid tokens
	// Return true if authenticated, false otherwise
	// Example: return r.Header.Get("Authorization") == "Bearer <token>"
	return true
}
