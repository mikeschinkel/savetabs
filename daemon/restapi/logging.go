package restapi

import (
	"log"
	"net/http"
)

// Middleware to log request and URL
func (a *API) addRequestLogging(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}
