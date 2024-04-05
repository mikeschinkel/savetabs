package restapi

import (
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
}

func (a *API) GetHealthz(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, statusResponse{Status: "OK"})
}

func (a *API) GetReadyz(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, http.StatusOK, statusResponse{Status: "OK"})
}
