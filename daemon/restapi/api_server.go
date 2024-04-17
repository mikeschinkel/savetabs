//go:build go1.22

// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.1-0.20240325090356-a14414f04fdd DO NOT EDIT.
package restapi

import (
	"fmt"
	"net/http"

	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Health Check
	// (GET /healthz)
	GetHealthz(w http.ResponseWriter, r *http.Request)
	// Return the HTML for a paginated table of links with optional filtering criteria in query parameters
	// (GET /html/links)
	GetLinks(w http.ResponseWriter, r *http.Request, params GetLinksParams)
	// Return the HTML for the Menu
	// (GET /html/menu)
	GetHtmlMenu(w http.ResponseWriter, r *http.Request)
	// Return the HTML for the Menu
	// (GET /html/menu/{menuItem})
	GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem)
	// Adds multiple links, each with group info
	// (POST /links/with-groups)
	PostLinksWithGroups(w http.ResponseWriter, r *http.Request)
	// Readiness Check
	// (GET /readyz)
	GetReadyz(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetHealthz operation middleware
func (siw *ServerInterfaceWrapper) GetHealthz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHealthz(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetLinks operation middleware
func (siw *ServerInterfaceWrapper) GetLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetLinksParams

	// ------------- Optional query parameter "gt" -------------

	err = runtime.BindQueryParameter("form", true, false, "gt", r.URL.Query(), &params.Gt)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "gt", Err: err})
		return
	}

	// ------------- Optional query parameter "g" -------------

	err = runtime.BindQueryParameter("form", true, false, "g", r.URL.Query(), &params.G)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "g", Err: err})
		return
	}

	// ------------- Optional query parameter "c" -------------

	err = runtime.BindQueryParameter("form", true, false, "c", r.URL.Query(), &params.C)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "c", Err: err})
		return
	}

	// ------------- Optional query parameter "t" -------------

	err = runtime.BindQueryParameter("form", true, false, "t", r.URL.Query(), &params.T)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "t", Err: err})
		return
	}

	// ------------- Optional query parameter "k" -------------

	err = runtime.BindQueryParameter("form", true, false, "k", r.URL.Query(), &params.K)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "k", Err: err})
		return
	}

	// ------------- Optional query parameter "b" -------------

	err = runtime.BindQueryParameter("form", true, false, "b", r.URL.Query(), &params.B)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "b", Err: err})
		return
	}

	// ------------- Optional query parameter "m" -------------

	err = runtime.BindQueryParameter("form", true, false, "m", r.URL.Query(), &params.M)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "m", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetLinks(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetHtmlMenu operation middleware
func (siw *ServerInterfaceWrapper) GetHtmlMenu(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHtmlMenu(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetHtmlMenuMenuItem operation middleware
func (siw *ServerInterfaceWrapper) GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "menuItem" -------------
	var menuItem MenuItem

	err = runtime.BindStyledParameterWithOptions("simple", "menuItem", r.PathValue("menuItem"), &menuItem, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "menuItem", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHtmlMenuMenuItem(w, r, menuItem)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostLinksWithGroups operation middleware
func (siw *ServerInterfaceWrapper) PostLinksWithGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostLinksWithGroups(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetReadyz operation middleware
func (siw *ServerInterfaceWrapper) GetReadyz(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetReadyz(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/healthz", wrapper.GetHealthz)
	m.HandleFunc("GET "+options.BaseURL+"/html/links", wrapper.GetLinks)
	m.HandleFunc("GET "+options.BaseURL+"/html/menu", wrapper.GetHtmlMenu)
	m.HandleFunc("GET "+options.BaseURL+"/html/menu/{menuItem}", wrapper.GetHtmlMenuMenuItem)
	m.HandleFunc("POST "+options.BaseURL+"/links/with-groups", wrapper.PostLinksWithGroups)
	m.HandleFunc("GET "+options.BaseURL+"/readyz", wrapper.GetReadyz)

	return m
}
