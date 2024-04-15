//go:build go1.22

package restapi

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

type swagger = openapi3.T

// Declare that *API implements ServerInterface
var _ ServerInterface = (*API)(nil)

type API struct {
	Port    string
	NextId  int64
	Mux     *http.ServeMux
	Swagger *swagger
	Handler http.Handler
	Server  *http.Server
	Lock    sync.Mutex
}

func NewAPI(port string, s *swagger) *API {
	api := &API{
		Port:    port,
		NextId:  1000,
		Mux:     http.NewServeMux(),
		Swagger: s,
	}

	// We now register our api above as the handler for the interface
	HandlerFromMux(api, api.Mux)
	// Report any panics
	h := api.catchPanic(api.Mux)
	// Authenticate request
	h = api.requireAuth(h)
	// Validation request against the OpenAPI schema.
	h = middleware.OapiRequestValidatorWithOptions(api.Swagger, api.openApiOptions())(h)
	// Add CORS security headers
	h = api.addCORSHeaders(h)
	// Log requests
	api.Handler = api.addRequestLogging(h)
	api.Server = &http.Server{
		Handler: api.Handler,
		Addr:    net.JoinHostPort("0.0.0.0", api.Port),
	}
	return api
}

func (a *API) openApiOptions() *middleware.Options {
	return &middleware.Options{
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			log.Printf("HTTP ERROR[%d]: %s", statusCode, message)
			http.Error(w, message, statusCode)
		},
	}
}

func (a *API) PostGroups(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *API) ListenAndServe() (err error) {
	log.Printf("Server listening on port %s...", a.Port)
	return a.Server.ListenAndServe()
}
