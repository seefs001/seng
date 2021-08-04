package seng

import (
	"context"
	"net/http"
	"net/url"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

const maxParams = 30

type Context struct {
	sync.Mutex

	engine *Engine

	request        *http.Request
	responseWriter http.ResponseWriter

	status      int
	baseURL     string
	path        string
	method      string
	routeParams [maxParams]string
	queryParams url.Values

	UserContext context.Context
}

func NewContext(w *http.Request, r http.ResponseWriter) *Context {
	return &Context{
		request:        w,
		responseWriter: r,
		UserContext:    context.Background(),
	}
}

func (c *Context) ReSet(r *http.Request, w http.ResponseWriter) *Context {
	c.request = r
	c.responseWriter = w

	c.path = r.URL.Path
	c.method = r.Method
	return c
}

func (c *Context) Text(format string) (err error) {
	_, err = c.responseWriter.Write([]byte(format))
	if err != nil {
		return
	}
	return
}

func (c *Context) Status(code int) *Context {
	c.responseWriter.WriteHeader(code)
	return c
}

func (c *Context) Json(data interface{}) (err error) {
	marshal, err := jsoniter.Marshal(data)
	if err != nil {
		return
	}
	_, err = c.responseWriter.Write(marshal)
	if err != nil {
		return
	}
	return
}
