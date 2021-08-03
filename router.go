package seng

import (
	"github.com/valyala/fasthttp"
)

type Router struct {
	Routes map[string]Handler
}

func (r *Router) RequestHandler(c *fasthttp.RequestCtx) {
	if handler, ok := r.Routes[string(c.Path())]; ok {
		ctx := &Context{}
		ctx.set(c)
		err := handler(ctx)
		if err != nil {
			// TODO error handler
		}
	}
}

func (r *Router) Get(path string, handler Handler) {
	r.Routes[path] = handler
}
