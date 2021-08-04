package seng

type ErrorHandler func(c *Context) error
type HandlerFunc func(c *Context) error
