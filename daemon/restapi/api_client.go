// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.1-0.20240325090356-a14414f04fdd DO NOT EDIT.
package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetHealthz request
	GetHealthz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHtmlBrowse request
	GetHtmlBrowse(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHtmlGroupTypesGroupTypeNameGroups request
	GetHtmlGroupTypesGroupTypeNameGroups(ctx context.Context, groupTypeName GroupTypeName, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHtmlGroupsGroupTypeGroupSlug request
	GetHtmlGroupsGroupTypeGroupSlug(ctx context.Context, groupType GroupType, groupSlug GroupSlug, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetLinks request
	GetLinks(ctx context.Context, params *GetLinksParams, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHtmlMenu request
	GetHtmlMenu(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetHtmlMenuMenuItem request
	GetHtmlMenuMenuItem(ctx context.Context, menuItem MenuItem, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetReadyz request
	GetReadyz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostResourcesWithGroupsWithBody request with any body
	PostResourcesWithGroupsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostResourcesWithGroups(ctx context.Context, body PostResourcesWithGroupsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetHealthz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthzRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHtmlBrowse(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHtmlBrowseRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHtmlGroupTypesGroupTypeNameGroups(ctx context.Context, groupTypeName GroupTypeName, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHtmlGroupTypesGroupTypeNameGroupsRequest(c.Server, groupTypeName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHtmlGroupsGroupTypeGroupSlug(ctx context.Context, groupType GroupType, groupSlug GroupSlug, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHtmlGroupsGroupTypeGroupSlugRequest(c.Server, groupType, groupSlug)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetLinks(ctx context.Context, params *GetLinksParams, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetLinksRequest(c.Server, params)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHtmlMenu(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHtmlMenuRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetHtmlMenuMenuItem(ctx context.Context, menuItem MenuItem, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHtmlMenuMenuItemRequest(c.Server, menuItem)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetReadyz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetReadyzRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostResourcesWithGroupsWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostResourcesWithGroupsRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostResourcesWithGroups(ctx context.Context, body PostResourcesWithGroupsJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostResourcesWithGroupsRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetHealthzRequest generates requests for GetHealthz
func NewGetHealthzRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/healthz")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHtmlBrowseRequest generates requests for GetHtmlBrowse
func NewGetHtmlBrowseRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/html/browse")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHtmlGroupTypesGroupTypeNameGroupsRequest generates requests for GetHtmlGroupTypesGroupTypeNameGroups
func NewGetHtmlGroupTypesGroupTypeNameGroupsRequest(server string, groupTypeName GroupTypeName) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "groupTypeName", runtime.ParamLocationPath, groupTypeName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/html/group-types/%s/groups", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHtmlGroupsGroupTypeGroupSlugRequest generates requests for GetHtmlGroupsGroupTypeGroupSlug
func NewGetHtmlGroupsGroupTypeGroupSlugRequest(server string, groupType GroupType, groupSlug GroupSlug) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "groupType", runtime.ParamLocationPath, groupType)
	if err != nil {
		return nil, err
	}

	var pathParam1 string

	pathParam1, err = runtime.StyleParamWithLocation("simple", false, "groupSlug", runtime.ParamLocationPath, groupSlug)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/html/groups/%s/%s", pathParam0, pathParam1)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetLinksRequest generates requests for GetLinks
func NewGetLinksRequest(server string, params *GetLinksParams) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/html/links")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	if params != nil {
		queryValues := queryURL.Query()

		if params.G != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "g", runtime.ParamLocationQuery, *params.G); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.C != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "c", runtime.ParamLocationQuery, *params.C); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.T != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "t", runtime.ParamLocationQuery, *params.T); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.K != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "k", runtime.ParamLocationQuery, *params.K); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.B != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "b", runtime.ParamLocationQuery, *params.B); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		if params.M != nil {

			if queryFrag, err := runtime.StyleParamWithLocation("form", true, "m", runtime.ParamLocationQuery, *params.M); err != nil {
				return nil, err
			} else if parsed, err := url.ParseQuery(queryFrag); err != nil {
				return nil, err
			} else {
				for k, v := range parsed {
					for _, v2 := range v {
						queryValues.Add(k, v2)
					}
				}
			}

		}

		queryURL.RawQuery = queryValues.Encode()
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHtmlMenuRequest generates requests for GetHtmlMenu
func NewGetHtmlMenuRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/html/menu")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetHtmlMenuMenuItemRequest generates requests for GetHtmlMenuMenuItem
func NewGetHtmlMenuMenuItemRequest(server string, menuItem MenuItem) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "menuItem", runtime.ParamLocationPath, menuItem)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/html/menu/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetReadyzRequest generates requests for GetReadyz
func NewGetReadyzRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/readyz")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostResourcesWithGroupsRequest calls the generic PostResourcesWithGroups builder with application/json body
func NewPostResourcesWithGroupsRequest(server string, body PostResourcesWithGroupsJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostResourcesWithGroupsRequestWithBody(server, "application/json", bodyReader)
}

// NewPostResourcesWithGroupsRequestWithBody generates requests for PostResourcesWithGroups with any type of body
func NewPostResourcesWithGroupsRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/resources/with-groups")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetHealthzWithResponse request
	GetHealthzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthzResponse, error)

	// GetHtmlBrowseWithResponse request
	GetHtmlBrowseWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHtmlBrowseResponse, error)

	// GetHtmlGroupTypesGroupTypeNameGroupsWithResponse request
	GetHtmlGroupTypesGroupTypeNameGroupsWithResponse(ctx context.Context, groupTypeName GroupTypeName, reqEditors ...RequestEditorFn) (*GetHtmlGroupTypesGroupTypeNameGroupsResponse, error)

	// GetHtmlGroupsGroupTypeGroupSlugWithResponse request
	GetHtmlGroupsGroupTypeGroupSlugWithResponse(ctx context.Context, groupType GroupType, groupSlug GroupSlug, reqEditors ...RequestEditorFn) (*GetHtmlGroupsGroupTypeGroupSlugResponse, error)

	// GetLinksWithResponse request
	GetLinksWithResponse(ctx context.Context, params *GetLinksParams, reqEditors ...RequestEditorFn) (*GetLinksResponse, error)

	// GetHtmlMenuWithResponse request
	GetHtmlMenuWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHtmlMenuResponse, error)

	// GetHtmlMenuMenuItemWithResponse request
	GetHtmlMenuMenuItemWithResponse(ctx context.Context, menuItem MenuItem, reqEditors ...RequestEditorFn) (*GetHtmlMenuMenuItemResponse, error)

	// GetReadyzWithResponse request
	GetReadyzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetReadyzResponse, error)

	// PostResourcesWithGroupsWithBodyWithResponse request with any body
	PostResourcesWithGroupsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostResourcesWithGroupsResponse, error)

	PostResourcesWithGroupsWithResponse(ctx context.Context, body PostResourcesWithGroupsJSONRequestBody, reqEditors ...RequestEditorFn) (*PostResourcesWithGroupsResponse, error)
}

type GetHealthzResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetHealthzResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHealthzResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHtmlBrowseResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetHtmlBrowseResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHtmlBrowseResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHtmlGroupTypesGroupTypeNameGroupsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetHtmlGroupTypesGroupTypeNameGroupsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHtmlGroupTypesGroupTypeNameGroupsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHtmlGroupsGroupTypeGroupSlugResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetHtmlGroupsGroupTypeGroupSlugResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHtmlGroupsGroupTypeGroupSlugResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetLinksResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetLinksResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetLinksResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHtmlMenuResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetHtmlMenuResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHtmlMenuResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetHtmlMenuMenuItemResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetHtmlMenuMenuItemResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHtmlMenuMenuItemResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetReadyzResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r GetReadyzResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetReadyzResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostResourcesWithGroupsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *[]IdObjects
	JSONDefault  *UnexpectedError
}

// Status returns HTTPResponse.Status
func (r PostResourcesWithGroupsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostResourcesWithGroupsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetHealthzWithResponse request returning *GetHealthzResponse
func (c *ClientWithResponses) GetHealthzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthzResponse, error) {
	rsp, err := c.GetHealthz(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthzResponse(rsp)
}

// GetHtmlBrowseWithResponse request returning *GetHtmlBrowseResponse
func (c *ClientWithResponses) GetHtmlBrowseWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHtmlBrowseResponse, error) {
	rsp, err := c.GetHtmlBrowse(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHtmlBrowseResponse(rsp)
}

// GetHtmlGroupTypesGroupTypeNameGroupsWithResponse request returning *GetHtmlGroupTypesGroupTypeNameGroupsResponse
func (c *ClientWithResponses) GetHtmlGroupTypesGroupTypeNameGroupsWithResponse(ctx context.Context, groupTypeName GroupTypeName, reqEditors ...RequestEditorFn) (*GetHtmlGroupTypesGroupTypeNameGroupsResponse, error) {
	rsp, err := c.GetHtmlGroupTypesGroupTypeNameGroups(ctx, groupTypeName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHtmlGroupTypesGroupTypeNameGroupsResponse(rsp)
}

// GetHtmlGroupsGroupTypeGroupSlugWithResponse request returning *GetHtmlGroupsGroupTypeGroupSlugResponse
func (c *ClientWithResponses) GetHtmlGroupsGroupTypeGroupSlugWithResponse(ctx context.Context, groupType GroupType, groupSlug GroupSlug, reqEditors ...RequestEditorFn) (*GetHtmlGroupsGroupTypeGroupSlugResponse, error) {
	rsp, err := c.GetHtmlGroupsGroupTypeGroupSlug(ctx, groupType, groupSlug, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHtmlGroupsGroupTypeGroupSlugResponse(rsp)
}

// GetLinksWithResponse request returning *GetLinksResponse
func (c *ClientWithResponses) GetLinksWithResponse(ctx context.Context, params *GetLinksParams, reqEditors ...RequestEditorFn) (*GetLinksResponse, error) {
	rsp, err := c.GetLinks(ctx, params, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetLinksResponse(rsp)
}

// GetHtmlMenuWithResponse request returning *GetHtmlMenuResponse
func (c *ClientWithResponses) GetHtmlMenuWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHtmlMenuResponse, error) {
	rsp, err := c.GetHtmlMenu(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHtmlMenuResponse(rsp)
}

// GetHtmlMenuMenuItemWithResponse request returning *GetHtmlMenuMenuItemResponse
func (c *ClientWithResponses) GetHtmlMenuMenuItemWithResponse(ctx context.Context, menuItem MenuItem, reqEditors ...RequestEditorFn) (*GetHtmlMenuMenuItemResponse, error) {
	rsp, err := c.GetHtmlMenuMenuItem(ctx, menuItem, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHtmlMenuMenuItemResponse(rsp)
}

// GetReadyzWithResponse request returning *GetReadyzResponse
func (c *ClientWithResponses) GetReadyzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetReadyzResponse, error) {
	rsp, err := c.GetReadyz(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetReadyzResponse(rsp)
}

// PostResourcesWithGroupsWithBodyWithResponse request with arbitrary body returning *PostResourcesWithGroupsResponse
func (c *ClientWithResponses) PostResourcesWithGroupsWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostResourcesWithGroupsResponse, error) {
	rsp, err := c.PostResourcesWithGroupsWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostResourcesWithGroupsResponse(rsp)
}

func (c *ClientWithResponses) PostResourcesWithGroupsWithResponse(ctx context.Context, body PostResourcesWithGroupsJSONRequestBody, reqEditors ...RequestEditorFn) (*PostResourcesWithGroupsResponse, error) {
	rsp, err := c.PostResourcesWithGroups(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostResourcesWithGroupsResponse(rsp)
}

// ParseGetHealthzResponse parses an HTTP response from a GetHealthzWithResponse call
func ParseGetHealthzResponse(rsp *http.Response) (*GetHealthzResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHealthzResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetHtmlBrowseResponse parses an HTTP response from a GetHtmlBrowseWithResponse call
func ParseGetHtmlBrowseResponse(rsp *http.Response) (*GetHtmlBrowseResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHtmlBrowseResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetHtmlGroupTypesGroupTypeNameGroupsResponse parses an HTTP response from a GetHtmlGroupTypesGroupTypeNameGroupsWithResponse call
func ParseGetHtmlGroupTypesGroupTypeNameGroupsResponse(rsp *http.Response) (*GetHtmlGroupTypesGroupTypeNameGroupsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHtmlGroupTypesGroupTypeNameGroupsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetHtmlGroupsGroupTypeGroupSlugResponse parses an HTTP response from a GetHtmlGroupsGroupTypeGroupSlugWithResponse call
func ParseGetHtmlGroupsGroupTypeGroupSlugResponse(rsp *http.Response) (*GetHtmlGroupsGroupTypeGroupSlugResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHtmlGroupsGroupTypeGroupSlugResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetLinksResponse parses an HTTP response from a GetLinksWithResponse call
func ParseGetLinksResponse(rsp *http.Response) (*GetLinksResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetLinksResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetHtmlMenuResponse parses an HTTP response from a GetHtmlMenuWithResponse call
func ParseGetHtmlMenuResponse(rsp *http.Response) (*GetHtmlMenuResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHtmlMenuResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetHtmlMenuMenuItemResponse parses an HTTP response from a GetHtmlMenuMenuItemWithResponse call
func ParseGetHtmlMenuMenuItemResponse(rsp *http.Response) (*GetHtmlMenuMenuItemResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHtmlMenuMenuItemResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetReadyzResponse parses an HTTP response from a GetReadyzWithResponse call
func ParseGetReadyzResponse(rsp *http.Response) (*GetReadyzResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetReadyzResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParsePostResourcesWithGroupsResponse parses an HTTP response from a PostResourcesWithGroupsWithResponse call
func ParsePostResourcesWithGroupsResponse(rsp *http.Response) (*PostResourcesWithGroupsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostResourcesWithGroupsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest []IdObjects
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest UnexpectedError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}
