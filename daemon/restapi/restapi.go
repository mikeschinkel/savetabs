//go:build go1.22

package restapi

import (
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
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

// routeContentType returns the first content type defined for the route
// TODO: Make more robust
func routeContentType(route *routers.Route) (ct string) {
	m := route.Operation.Responses.Map()
	type mt = map[string]*openapi3.MediaType
	for key := range mt(m["200"].Value.Content) {
		ct = key
	}
	return ct
}

func (a *API) openApiOptions() *middleware.Options {
	return &middleware.Options{
		ErrorHandlerWithOpts: func(w http.ResponseWriter, message string, statusCode int, opts middleware.ErrorHandlerOpts) {
			switch routeContentType(opts.Route) {
			case "application/json":
			case "text/html":
				// TODO: Call function in ui package to display error message
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				w.Header().Set("X-Content-Type-Options", "nosniff")
				a.sendError(w, opts.Request, statusCode, message)
			case "text/plain":
				fallthrough
			default:
				slog.Error("HTTP ERROR", "status_code", statusCode, "error_msg", message)
				http.Error(w, message, statusCode)
			}
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
