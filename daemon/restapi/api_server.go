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
	// Return the HTML for the Browse body, e.g. the list of group types as an HTML list
	// (GET /html/browse)
	GetHtmlBrowse(w http.ResponseWriter, r *http.Request)
	// Return the list of groups for a group type as HTML
	// (GET /html/group-types/{groupTypeName}/groups)
	GetHtmlGroupTypesGroupTypeNameGroups(w http.ResponseWriter, r *http.Request, groupTypeName GroupTypeName)
	// Return the list of resources for a group as HTML
	// (GET /html/groups/{groupType}/{groupSlug})
	GetHtmlGroupsGroupTypeGroupSlug(w http.ResponseWriter, r *http.Request, groupType GroupType, groupSlug GroupSlug)
	// Return the HTML for the Menu
	// (GET /html/menu)
	GetHtmlMenu(w http.ResponseWriter, r *http.Request)
	// Return the HTML for the Menu
	// (GET /html/menu/{menuItem})
	GetHtmlMenuMenuItem(w http.ResponseWriter, r *http.Request, menuItem MenuItem)
	// Readiness Check
	// (GET /readyz)
	GetReadyz(w http.ResponseWriter, r *http.Request)
	// Adds multiple resources, each with group info
	// (POST /resources/with-groups)
	PostResourcesWithGroups(w http.ResponseWriter, r *http.Request)
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

// GetHtmlBrowse operation middleware
func (siw *ServerInterfaceWrapper) GetHtmlBrowse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHtmlBrowse(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetHtmlGroupTypesGroupTypeNameGroups operation middleware
func (siw *ServerInterfaceWrapper) GetHtmlGroupTypesGroupTypeNameGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "groupTypeName" -------------
	var groupTypeName GroupTypeName

	err = runtime.BindStyledParameterWithOptions("simple", "groupTypeName", r.PathValue("groupTypeName"), &groupTypeName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "groupTypeName", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHtmlGroupTypesGroupTypeNameGroups(w, r, groupTypeName)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetHtmlGroupsGroupTypeGroupSlug operation middleware
func (siw *ServerInterfaceWrapper) GetHtmlGroupsGroupTypeGroupSlug(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "groupType" -------------
	var groupType GroupType

	err = runtime.BindStyledParameterWithOptions("simple", "groupType", r.PathValue("groupType"), &groupType, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "groupType", Err: err})
		return
	}

	// ------------- Path parameter "groupSlug" -------------
	var groupSlug GroupSlug

	err = runtime.BindStyledParameterWithOptions("simple", "groupSlug", r.PathValue("groupSlug"), &groupSlug, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "groupSlug", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetHtmlGroupsGroupTypeGroupSlug(w, r, groupType, groupSlug)
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

// PostResourcesWithGroups operation middleware
func (siw *ServerInterfaceWrapper) PostResourcesWithGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostResourcesWithGroups(w, r)
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
	m.HandleFunc("GET "+options.BaseURL+"/html/browse", wrapper.GetHtmlBrowse)
	m.HandleFunc("GET "+options.BaseURL+"/html/group-types/{groupTypeName}/groups", wrapper.GetHtmlGroupTypesGroupTypeNameGroups)
	m.HandleFunc("GET "+options.BaseURL+"/html/groups/{groupType}/{groupSlug}", wrapper.GetHtmlGroupsGroupTypeGroupSlug)
	m.HandleFunc("GET "+options.BaseURL+"/html/menu", wrapper.GetHtmlMenu)
	m.HandleFunc("GET "+options.BaseURL+"/html/menu/{menuItem}", wrapper.GetHtmlMenuMenuItem)
	m.HandleFunc("GET "+options.BaseURL+"/readyz", wrapper.GetReadyz)
	m.HandleFunc("POST "+options.BaseURL+"/resources/with-groups", wrapper.PostResourcesWithGroups)

	return m
}
