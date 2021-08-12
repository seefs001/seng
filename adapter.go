package seng

import "net/http"

// AdapterHandlerFunc adapter http.HandlerFunc -> seng.Handler
func AdapterHandlerFunc(handlerFunc http.HandlerFunc) Handler {
	return func(c *Context) (err error) {
		handlerFunc.ServeHTTP(c.Writer, c.Request)
		return
	}
}

// AdapterHandler adapter http.Handler -> seng.Handler
func AdapterHandler(handler http.Handler) Handler {
	return func(c *Context) (err error) {
		handler.ServeHTTP(c.Writer, c.Request)
		return
	}
}
