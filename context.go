package seng

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
)

// Context represents the Context which hold the HTTP request and response.
type Context struct {
	// reference to engine
	engine *Engine
	// reference to router
	router *Router
	// origin objects
	// http Writer
	Writer http.ResponseWriter
	// http Request
	Request *http.Request
	// request info
	Method   string
	HostName string
	// path -> /ping
	Path string
	// params routeParams
	Params map[string]string
	// status code
	StatusCode int
	// middleWares
	handlers []Handler
	// index of handler
	indexHandler int
	// context value
	Values map[string]interface{}
	// user context
	userContext context.Context
}

// NewContext new context with default
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	ctx := &Context{}
	ctx.init(w, req)
	return ctx
}

// init Context
func (c *Context) init(w http.ResponseWriter, req *http.Request) {
	c.Writer = w
	c.Request = req
	c.Method = req.Method
	c.Path = req.URL.Path
	c.HostName = req.Host
	c.Values = make(map[string]interface{})
	c.indexHandler = -1
}

// ReSet context from w,req
func (c *Context) ReSet(w http.ResponseWriter, req *http.Request) *Context {
	// init
	c.init(w, req)
	return c
}

// Engine return the engine reference
func (c *Context) Engine() *Engine {
	return c.engine
}

// Router return the router reference
func (c *Context) Router() *Router {
	return c.router
}

// UserContext return user context
func (c *Context) UserContext() context.Context {
	return c.userContext
}

// SetUserContext set user context
func (c *Context) SetUserContext(u context.Context) {
	c.userContext = u
}

// FormFile returns the first file by key from a MultipartForm.
func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

// defaultString returns the value or a default value if it is set
func defaultString(value string, defaultValue []string) string {
	if len(value) == 0 && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// FormValue returns the value
func (c *Context) FormValue(key string, defaultValue ...string) string {
	return defaultString(c.Request.FormValue(key), defaultValue)
}

// GetHeader get value from header
func (c *Context) GetHeader(key string, defaultValue ...string) string {
	return defaultString(c.Request.Header.Get(key), defaultValue)
}

// DelHeader del value from header
func (c *Context) DelHeader(key string) {
	c.Request.Header.Del(key)
	c.Writer.Header().Del(key)
	return
}

// GetContentType get content type
func (c *Context) GetContentType() string {
	return c.GetHeader(HeaderContentType)
}

// SetHeader set header
func (c *Context) SetHeader(key string, value string) {
	c.Request.Header.Set(key, value)
	c.Writer.Header().Set(key, value)
}

// RemoveHeader remove specific key header
func (c *Context) RemoveHeader(key string) {
	c.Request.Header.Del(key)
	c.Writer.Header().Del(key)
}

// GetHeaderValues get header values from header
func (c *Context) GetHeaderValues(key string) []string {
	return c.Writer.Header().Values(key)
}

// IP get remote ip address
func (c *Context) IP() string {
	return c.Request.RemoteAddr
}

// Set set user value
func (c *Context) Set(key string, value interface{}) {
	c.Values[key] = value
}

// Get get user value
func (c *Context) Get(key string) (data interface{}, exists bool) {
	if data, ok := c.Values[key]; ok {
		return data, true
	}
	return nil, false
}

// MultipartForm request multipart form
func (c *Context) MultipartForm() *multipart.Form {
	return c.Request.MultipartForm
}

// Refer returns refer url
func (c *Context) Refer() string {
	return c.Request.Referer()
}

// String ...
func (c *Context) String() string {
	// TODO
	return fmt.Sprintf(
		"%s <-> %s seng version： %s",
		c.Request.Host,
		c.Request.RemoteAddr,
		c.engine.config.SengVersion,
	)
}

// PostForm get value from FormValue
func (c *Context) PostForm(key string, defaultValue ...string) string {
	return defaultString(c.Request.FormValue(key), defaultValue)
}

// Query get value from query
func (c *Context) Query(key string, defaultValue ...string) string {
	return defaultString(c.Request.URL.Query().Get(key), defaultValue)
}

// Status set response status code
func (c *Context) Status(code int) *Context {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
	return c
}

// Text return text
func (c *Context) Text(format string, values ...interface{}) (err error) {
	c.SetHeader(HeaderContentType, MIMETextPlainCharsetUTF8)
	_, err = c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
	if err != nil {
		return err
	}
	return
}

// JSON return JSON
func (c *Context) JSON(obj interface{}) (err error) {
	c.SetHeader(HeaderContentType, MINEApplicationJSON)
	encoder := json.NewEncoder(c.Writer)
	if err = encoder.Encode(obj); err != nil {
		return err
	}
	return
}

// Protobuf return protobuf
func (c *Context) Protobuf(data []byte) error {
	c.SetHeader(HeaderContentType, MINEApplicationProtobuf)
	return c.Data(data)
}

// Data return []byte Data
func (c *Context) Data(data []byte) (err error) {
	_, err = c.Writer.Write(data)
	if err != nil {
		return err
	}
	return
}

// Fail return fail message
func (c *Context) Fail(code int, err string) error {
	c.indexHandler = len(c.handlers)
	return c.JSON(NewError(code, err))
}

// HTML renders HTML
func (c *Context) HTML(name string, data interface{}) (err error) {
	c.SetHeader(HeaderContentType, MINETextHTML)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		return c.Fail(http.StatusInternalServerError, err.Error())
	}
	return
}

// Param get values from route parameters
func (c *Context) Param(key string, defaultValue ...string) (string, bool) {
	value, ok := c.Params[key]
	if !ok {
		if len(defaultValue) == 0 {
			return "", false
		} else {
			return defaultValue[0], true
		}
	}
	return value, true
}

// ParamsInt param -> int
func (c *Context) ParamsInt(key string, defaultValue ...int) (int, bool) {
	param, ok := c.Param(key)
	if !ok {
		if len(defaultValue) == 0 {
			return 0, false
		} else {
			return defaultValue[0], true
		}
	}
	paramInt, err := strconv.Atoi(param)
	if err != nil {
		return 0, false
	}
	return paramInt, true
}

// ParamsInt64 param -> int64
func (c *Context) ParamsInt64(key string, defaultValue ...int64) (int64, bool) {
	param, ok := c.Param(key)
	if !ok {
		if len(defaultValue) == 0 {
			return 0, false
		} else {
			return defaultValue[0], true
		}
	}
	paramInt64, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, false
	}
	return paramInt64, true
}

// Next middleware
// All handlers are traversed here, because not all handlers will manually call c.Next()
// For handlers that only act before the request, you can omit c.Next()
func (c *Context) Next() error {
	c.indexHandler++
	size := len(c.handlers)
	for ; c.indexHandler < size; c.indexHandler++ {
		err := c.handlers[c.indexHandler](c)
		if err != nil {
			return err
		}
	}
	return nil
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// JSONResponse return json response with status 200
func (c *Context) JSONResponse(code int, msg string, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}
