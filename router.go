package seng

import (
	"sync"

	"github.com/seefs001/seng/utils"
	"github.com/valyala/fasthttp"
)

var ctxPool *sync.Pool

func init() {
	ctxPool = &sync.Pool{
		New: func() interface{} {
			return new(Context)
		},
	}
}

type Router struct {
	Routes       map[string]Handler
	ErrorHandler Handler
}

func (r *Router) RequestHandler(c *fasthttp.RequestCtx) {
	ctx := ctxPool.Get().(*Context)
	ctx.set(c)
	handler, err := r.match(utils.Bytes2String(c.Path()), utils.Bytes2String(c.Method()))
	if err != nil {
		// TODO error handler
		if err == ErrNotFoundRoute {
			err := r.ErrorHandler(ctx)
			if err != nil {
				// TODO error handler
				return
			}
		}
		return
	}
	err = handler(ctx)
	if err != nil {
		// TODO error handler
		ctxPool.Put(ctx)
	}
	ctxPool.Put(ctx)
}

func (r *Router) Get(path string, handler Handler) {
	r.add(MethodGet, path, handler)
}

func (r *Router) Post(path string, handler Handler) {
	r.add(MethodPost, path, handler)
}

func (r *Router) Put(path string, handler Handler) {
	r.add(MethodPut, path, handler)
}

func (r *Router) Delete(path string, handler Handler) {
	r.add(MethodDelete, path, handler)
}
func (r *Router) Head(path string, handler Handler) {
	r.add(MethodHead, path, handler)
}

func (r *Router) Patch(path string, handler Handler) {
	r.add(MethodPatch, path, handler)
}

func (r *Router) Connect(path string, handler Handler) {
	r.add(MethodConnect, path, handler)
}

func (r *Router) Trace(path string, handler Handler) {
	r.add(MethodTrace, path, handler)
}

func (r *Router) Options(path string, handler Handler) {
	r.add(MethodOptions, path, handler)
}

func (r *Router) add(method string, path string, handler Handler) {
	key := method + "-" + path
	r.Routes[key] = handler
}

func (r *Router) match(path string, method string) (Handler, error) {
	// TODO 模糊匹配
	key := method + "-" + path
	if handler, ok := r.Routes[key]; !ok {
		return nil, ErrNotFoundRoute
	} else {
		return handler, nil
	}
}
