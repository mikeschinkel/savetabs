package restapi

import (
	"net/http"
)

func (a *API) GetReadyz(w http.ResponseWriter, r *http.Request) {
	// TODO: Check health of database
	sendJSON(w, http.StatusOK, newJSONResponse(true))
}
