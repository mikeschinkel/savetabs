//go:build go1.22

package restapi

import (
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
)

type swagger = openapi3.T

// Declare that *API implements ServerInterface and Persister interfaces
var _ ServerInterface = (*API)(nil)
var _ Persister = (*API)(nil)

type API struct {
	Port    string
	Mux     *http.ServeMux
	Swagger *swagger
	Handler http.Handler
	Server  *http.Server
	Lock    sync.Mutex
	Views   Viewer
}

type APIArgs struct {
	Port    string
	Swagger *swagger
	Views   Viewer
}

func NewAPI(args APIArgs) *API {
	if args.Port == "" {
		args.Port = DefaultPort
	}
	api := &API{
		Port:    args.Port,
		Swagger: args.Swagger,
		Views:   args.Views,
		Mux:     http.NewServeMux(),
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
			slog.Error("HTTP ERROR", "status_code", statusCode, "error_msg", message)
			http.Error(w, message, statusCode)
		},
	}
}

func (a *API) ListenAndServe() (err error) {
	slog.Info("SaveTabs server listening", "port", a.Port)
	return a.Server.ListenAndServe()
}
func (a *API) Shutdown(ctx Context) (err error) {
	slog.Info("SaveTabs server shutting down")
	return a.Server.Shutdown(ctx)
}
