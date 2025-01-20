package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/safehtml"
	"github.com/mikeschinkel/savetabs/daemon/guard"
	"github.com/mikeschinkel/savetabs/daemon/shared"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Context = context.Context
type Buffer = bytes.Buffer

type jsonResponse struct {
	HTTPStatus int    `json:"-"`
	Success    bool   `json:"success"`
	Message    string `json:"message,omitempty"`
}

func newJSONResponse(ok bool) jsonResponse {
	var resp jsonResponse
	resp.Success = ok
	return resp
}

func urlForRequest(r *http.Request) string {
	r.URL.Host = r.Host
	r.URL.Scheme = "http"
	return r.URL.String()
}

// sendHTMLError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func (a *API) sendHTMLError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	msg = fmt.Sprintf("ERROR[%d]: '%s' from %s.", code, msg, urlForRequest(r))
	hr, err := guard.GetErrorHTML(guard.ErrorParams{
		Host: r.Host,
		Msg:  msg,
	})
	if err != nil {
		shared.Panicf(err.Error())
	}
	w.Header().Set("Content-Type", "text/html")
	if code == 0 {
		code = http.StatusInternalServerError
		slog.Warn("HTTP Status code not set", "error", msg)
	}
	w.WriteHeader(code)
	//goland:noinspection GoDfaErrorMayBeNotNil
	_, _ = fmt.Fprint(w, hr.HTML.String()) // TODO: Change this to use safehtml
}

func (a *API) sendPlainError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	msg = fmt.Sprintf("ERROR[%d]: '%s' from %s.", code, msg, urlForRequest(r))
	w.Header().Set("Content-Type", "text/plain")
	if code == 0 {
		code = http.StatusInternalServerError
		slog.Warn("HTTP Status code not set", "error", msg)
	}
	w.WriteHeader(code)
	_, _ = fmt.Fprint(w, msg) // TODO: Change this to use safehtml
}

// sendJSON sends a success code and json encoded content
func sendJSON(w http.ResponseWriter, code int, content any) {
	w.Header().Set("Content-Type", "application/json")
	if code == 0 {
		code = http.StatusInternalServerError
		slog.Warn("HTTP Status code not set")
	}
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(content)
}

// sendHTML sends a success code of 200 and the HTML content provided
func (a *API) sendHTML(w http.ResponseWriter, html ...safehtml.HTML) {
	a.sendHTMLStatus(w, http.StatusOK, html...)
}
func (a *API) sendHTMLStatus(w http.ResponseWriter, status int, html ...safehtml.HTML) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	for _, h := range html {
		_, _ = fmt.Fprint(w, h.String())
	}
}

func (a *API) sendView(ctx Context, w http.ResponseWriter, r *http.Request, fn func(ctx Context) (guard.HTMLResponse, error)) {
	hr, err := fn(ctx)
	if err != nil {
		//goland:noinspection GoDfaErrorMayBeNotNil
		a.sendHTMLError(w, r, hr.StatusCode, err.Error())
		return
	}
	a.sendHTML(w, hr.HTML)
}
