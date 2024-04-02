//go:build go1.22

package restapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"savetabs/sqlc"
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

func sanitizeResourcesWithGroups(data resourcesWithGroups) (err error) {
	for i, rg := range data {
		if rg.Url == nil || *rg.Url == "" {
			err = errors.Join(ErrUrlNotSpecified, fmt.Errorf("error found in resource index %d", i))
			goto end
		}
		if rg.Id == nil {
			data[i].Id = ptr[int64](0)
		}
		if rg.Group == nil || *rg.Group == "" {
			data[i].Group = ptr("<none>")
		}
		if rg.GroupType == nil || *rg.GroupType == "" {
			data[i].GroupType = ptr("invalid")
		}
	}
end:
	return err
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
	err = sanitizeResourcesWithGroups(data)
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

// Middleware to add CORS headers to every response
func (a *API) addCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//w.Header().Set("Access-Control-Allow-Origin", "http://localhost:"+a.Port)

		// Allow specific methods (e.g., GET, POST, OPTIONS)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		// Allow specific headers (e.g., Content-Type, Authorization)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		handler.ServeHTTP(w, r)
	})
}

func (a *API) ListenAndServe() (err error) {
	log.Printf("Server listening on port %s...", a.Port)
	return a.Server.ListenAndServe()
}
