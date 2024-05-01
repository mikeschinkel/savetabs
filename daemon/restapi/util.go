package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/safehtml"
)

type Context = context.Context
type Buffer = bytes.Buffer

func urlForRequest(r *http.Request) string {
	r.URL.Host = r.Host
	r.URL.Scheme = "http"
	return r.URL.String()
}

// sendError wraps sending of an error in the Error format, and
// handling the failure to marshal that.
func (a *API) sendError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	msg = fmt.Sprintf("ERROR[%d]: '%s' from %s.", code, msg, urlForRequest(r))
	html, _, _ := a.Views.GetErrorHTML(nil, errors.New(msg))
	_, _ = fmt.Fprint(w, html.String())
}

// sendJSON sends a success code and json encoded content
func sendJSON(w http.ResponseWriter, code int, content any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(content)
}

// sendHTML sends a success code of 200 and the HTML content provided
func (a *API) sendHTML(w http.ResponseWriter, html safehtml.HTML) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, html.String())
}

//goland:noinspection GoUnusedFunction
func deleteElement[T any](slice []T, index int) []T {
	// Copy the elements following the index one position to the left.
	copy(slice[index:], slice[index+1:])
	// Return the slice without the last element.
	return slice[:len(slice)-1]
}

func (a *API) sendView(ctx Context, w http.ResponseWriter, r *http.Request, fn func(ctx Context) (safehtml.HTML, int, error)) {
	out, status, err := fn(ctx)
	if err != nil {
		a.sendError(w, r, status, err.Error())
		return
	}
	a.sendHTML(w, out)
}

func toUpperSlice(s []string) []string {
	for i := range s {
		s[i] = strings.ToUpper(s[i])
	}
	return s
}
