// Package server_oapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package server_oapi

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// Error defines model for Error.
type Error struct {
	// Code Error code
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Message defines model for Message.
type Message struct {
	// Id Message ID
	Id      int32  `json:"id"`
	Message string `json:"message"`
}

// NewMessage defines model for NewMessage.
type NewMessage struct {
	// Message message to create
	Message string `json:"message"`
}

// PostV1MessageJSONRequestBody defines body for PostV1Message for application/json ContentType.
type PostV1MessageJSONRequestBody = NewMessage

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
	// PostV1Message request with any body
	PostV1MessageWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostV1Message(ctx context.Context, body PostV1MessageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetV1MessageId request
	GetV1MessageId(ctx context.Context, id int32, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) PostV1MessageWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostV1MessageRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostV1Message(ctx context.Context, body PostV1MessageJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostV1MessageRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetV1MessageId(ctx context.Context, id int32, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetV1MessageIdRequest(c.Server, id)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewPostV1MessageRequest calls the generic PostV1Message builder with application/json body
func NewPostV1MessageRequest(server string, body PostV1MessageJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostV1MessageRequestWithBody(server, "application/json", bodyReader)
}

// NewPostV1MessageRequestWithBody generates requests for PostV1Message with any type of body
func NewPostV1MessageRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/message")
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

// NewGetV1MessageIdRequest generates requests for GetV1MessageId
func NewGetV1MessageIdRequest(server string, id int32) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "id", runtime.ParamLocationPath, id)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v1/message/%s", pathParam0)
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
	// PostV1Message request with any body
	PostV1MessageWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostV1MessageResponse, error)

	PostV1MessageWithResponse(ctx context.Context, body PostV1MessageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostV1MessageResponse, error)

	// GetV1MessageId request
	GetV1MessageIdWithResponse(ctx context.Context, id int32, reqEditors ...RequestEditorFn) (*GetV1MessageIdResponse, error)
}

type PostV1MessageResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Message
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r PostV1MessageResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostV1MessageResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetV1MessageIdResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Message
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r GetV1MessageIdResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetV1MessageIdResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// PostV1MessageWithBodyWithResponse request with arbitrary body returning *PostV1MessageResponse
func (c *ClientWithResponses) PostV1MessageWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostV1MessageResponse, error) {
	rsp, err := c.PostV1MessageWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostV1MessageResponse(rsp)
}

func (c *ClientWithResponses) PostV1MessageWithResponse(ctx context.Context, body PostV1MessageJSONRequestBody, reqEditors ...RequestEditorFn) (*PostV1MessageResponse, error) {
	rsp, err := c.PostV1Message(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostV1MessageResponse(rsp)
}

// GetV1MessageIdWithResponse request returning *GetV1MessageIdResponse
func (c *ClientWithResponses) GetV1MessageIdWithResponse(ctx context.Context, id int32, reqEditors ...RequestEditorFn) (*GetV1MessageIdResponse, error) {
	rsp, err := c.GetV1MessageId(ctx, id, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetV1MessageIdResponse(rsp)
}

// ParsePostV1MessageResponse parses an HTTP response from a PostV1MessageWithResponse call
func ParsePostV1MessageResponse(rsp *http.Response) (*PostV1MessageResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostV1MessageResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Message
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseGetV1MessageIdResponse parses an HTTP response from a GetV1MessageIdWithResponse call
func ParseGetV1MessageIdResponse(rsp *http.Response) (*GetV1MessageIdResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetV1MessageIdResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Message
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /v1/message)
	PostV1Message(ctx echo.Context) error

	// (GET /v1/message/{id})
	GetV1MessageId(ctx echo.Context, id int32) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostV1Message converts echo context to params.
func (w *ServerInterfaceWrapper) PostV1Message(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostV1Message(ctx)
	return err
}

// GetV1MessageId converts echo context to params.
func (w *ServerInterfaceWrapper) GetV1MessageId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetV1MessageId(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/v1/message", wrapper.PostV1Message)
	router.GET(baseURL+"/v1/message/:id", wrapper.GetV1MessageId)

}

type PostV1MessageRequestObject struct {
	Body *PostV1MessageJSONRequestBody
}

type PostV1MessageResponseObject interface {
	VisitPostV1MessageResponse(w http.ResponseWriter) error
}

type PostV1Message200JSONResponse Message

func (response PostV1Message200JSONResponse) VisitPostV1MessageResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostV1MessagedefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response PostV1MessagedefaultJSONResponse) VisitPostV1MessageResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type GetV1MessageIdRequestObject struct {
	Id int32 `json:"id"`
}

type GetV1MessageIdResponseObject interface {
	VisitGetV1MessageIdResponse(w http.ResponseWriter) error
}

type GetV1MessageId200JSONResponse Message

func (response GetV1MessageId200JSONResponse) VisitGetV1MessageIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetV1MessageIddefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response GetV1MessageIddefaultJSONResponse) VisitGetV1MessageIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {

	// (POST /v1/message)
	PostV1Message(ctx context.Context, request PostV1MessageRequestObject) (PostV1MessageResponseObject, error)

	// (GET /v1/message/{id})
	GetV1MessageId(ctx context.Context, request GetV1MessageIdRequestObject) (GetV1MessageIdResponseObject, error)
}

type StrictHandlerFunc func(ctx echo.Context, args interface{}) (interface{}, error)

type StrictMiddlewareFunc func(f StrictHandlerFunc, operationID string) StrictHandlerFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// PostV1Message operation middleware
func (sh *strictHandler) PostV1Message(ctx echo.Context) error {
	var request PostV1MessageRequestObject

	var body PostV1MessageJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.PostV1Message(ctx.Request().Context(), request.(PostV1MessageRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostV1Message")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(PostV1MessageResponseObject); ok {
		return validResponse.VisitPostV1MessageResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// GetV1MessageId operation middleware
func (sh *strictHandler) GetV1MessageId(ctx echo.Context, id int32) error {
	var request GetV1MessageIdRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.GetV1MessageId(ctx.Request().Context(), request.(GetV1MessageIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetV1MessageId")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(GetV1MessageIdResponseObject); ok {
		return validResponse.VisitGetV1MessageIdResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("Unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RUwXLTMBD9FbFw1NRpe/MRypQMU+DAwKHpQbU2iRhbEqt12oxH/85ItrGTtEMPvXFJ",
	"bO3qvd19+9xB5RrvLFoOUHYQqi02Kj9+JHKUHjw5j8QG83HlNKZ/jaEi49k4C2WfLHJMwtpRoxhKMJYv",
	"L0AC7z32r7hBgiihwRDUJgMNwcBk7AZilED4uzWEGspbGCDH/Lso4Wa6e1ia0aeFDcliefXKhRl9VNYX",
	"fHi2shnqYXlDQLATFaFinIp6hneiTBFj1+4UVeNatTVfzQ4lsOEap+D3/CphhxT6W4uz8zQB59Eqb6CE",
	"y7N0JMEr3uY2it15MWvFu8Cn5B9yH0IJiw9izM6wpFLKUkMJ31zgH+c3f6OpQwz83ul9v2SW0WZw5X1t",
	"qnyz+BUSQwf4qBpf96NdpJ+dqtvUW7ca57OCcgWfsK6d+Omo1m9WEPMs+w1Pl94RrqGEt8VkgWLY/2Im",
	"Zoy9BME7G3rOi8XiBVW+jGlGI48m+fUz5LOs16sR9r5+gq61+OixYtQCx5wo56IXndEx4W/wCeGvkYUa",
	"FRf3+95zh7pf4yT7UufdItUgIwUob48Rl1fCrcXMI4TcUlpmk8JpL0GCVU02sYa5U5hanKv9T+/Hu/9X",
	"5eQLpN2oQks1lLBl9mVRdFsXOM04FumzIGGnyKj7wX5jsF+HoQeoXaXqFErod/FPAAAA//+1HpmZZgYA",
	"AA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
