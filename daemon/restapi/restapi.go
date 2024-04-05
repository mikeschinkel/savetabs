//go:build go1.22

package restapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"savetabs/sqlc"
	"savetabs/ui"
)

type swagger = openapi3.T

type API struct {
	Port    string
	NextId  int64
	Mux     *http.ServeMux
	Swagger *swagger
	Handler http.Handler
	Server  *http.Server
	Lock    sync.Mutex
}

func (a *API) GetHtmlBrowse(w http.ResponseWriter, r *http.Request) {
	out, err := ui.BrowseHTML(r.Host)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendHTML(w, out)
}

func (a *API) GetHtmlGroupTypesTypeNameGroups(w http.ResponseWriter, r *http.Request, typeName GroupTypeName) {
	ctx := context.Background()
	out, err := ui.GroupsByGroupTypeHTML(ctx, r.Host, typeName)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	sendHTML(w, out)
}

// Declare that *API implements ServerInterface
var _ ServerInterface = (*API)(nil)

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
	api.Handler = requireAuth(h)
	api.Server = &http.Server{
		Handler: api.Handler,
		Addr:    net.JoinHostPort("0.0.0.0", api.Port),
	}
	return api
}

func (a *API) PostResourcesWithGroups(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// TODO: Find a better status result than "Bad Gateway"
		sendError(w, http.StatusBadGateway, err.Error())
		return
	}
	var data resourcesWithGroups
	err = json.Unmarshal(body, &data)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	data, err = sanitizeResourcesWithGroups(data)
	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}
	ds := sqlc.GetDatastore()
	db, ok := ds.DB().(*sqlc.NestedDBTX)
	if !ok {
		sendError(w, http.StatusInternalServerError, "DB not a NestedDBTX")
		return
	}
	err = db.Exec(func(tx *sql.Tx) (err error) {
		err = upsertResources(context.TODO(), ds, data)
		switch {
		case err == nil:
			goto end
		case errors.Is(err, ErrFailedToUnmarshal):
			sendError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrFailedUpsertResources):
			// TODO: Break out errors into different status for different reasons
			fallthrough
		default:
			sendError(w, http.StatusInternalServerError, err.Error())
		}
	end:
		return err
	})
}

func (a *API) PostGroups(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (a *API) ListenAndServe() (err error) {
	log.Printf("Server listening on port %s...", a.Port)
	return a.Server.ListenAndServe()
}
