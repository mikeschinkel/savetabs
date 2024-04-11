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

func (a *API) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem) {
	//TODO implement me
	sendError(w, r, http.StatusInternalServerError, "implement me")
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
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	h := middleware.OapiRequestValidator(api.Swagger)(api.Mux)
	// Use URL logging handler middleware
	h = api.addURLLogging(h)
	// Use middleware to address CORS security
	h = api.addCORSHeaders(h)
	// Add authentication
	h = api.catchPanic(h)
	// Add authentication
	api.Handler = requireAuth(h)
	api.Server = &http.Server{
		Handler: api.Handler,
		Addr:    net.JoinHostPort("0.0.0.0", api.Port),
	}
	return api
}

func (a *API) PostGroups(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *API) ListenAndServe() (err error) {
	log.Printf("Server listening on port %s...", a.Port)
	return a.Server.ListenAndServe()
}
