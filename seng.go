package seng

import (
	"net/http"
	"sync"
)

type Engine struct {
	sync.Mutex

	CtxPool sync.Pool

	routes    map[string][]*Route
	treeStack []map[string][]*Route
	config    Config
}

func New(config Config) *Engine {

	e := &Engine{
		CtxPool: sync.Pool{
			New: func() interface{} {
				return NewContext(nil, nil)
			},
		},
		routes: make(map[string][]*Route),
		config: config,
	}

	return e
}

func Default() *Engine {

	defaultConfig := Config{
		StrictRouting: false,
		CaseSensitive: true,
		GETOnly:       false,
		ErrorHandler: func(ctx *Context) error {
			return ctx.Status(http.StatusInternalServerError).Text("server error")
		},
		NotFoundHandler: func(c *Context) error {
			return c.Status(http.StatusNotFound).Text("not found")
		},
		DisableKeepalive: false,
		AppName:          "Seng Web Application",
		Addr:             "0.0.0.0:8080",
	}

	e := New(defaultConfig)

	return e
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := e.CtxPool.Get().(*Context).ReSet(r, w)
	handlerFuncs, b := e.matchRoute(ctx.method, ctx.path)
	if !b {
		err := e.config.NotFoundHandler(ctx)
		if err != nil {
			return
		}
	}
	for _, h := range handlerFuncs {
		h(ctx)
	}
}

func (e *Engine) matchRoute(method string, pattern string) ([]HandlerFunc, bool) {
	if routes, ok := e.routes[method]; !ok {
		return nil, false
	} else {
		for _, route := range routes {
			if route.Path == pattern {
				return route.Handlers, true
			}
		}
		return nil, false
	}
}

func (e *Engine) Get(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodGet, pattern, handler...)
}

func (e *Engine) Head(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodHead, pattern, handler...)
}

func (e *Engine) Post(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodPost, pattern, handler...)
}

func (e *Engine) Put(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodPut, pattern, handler...)
}

func (e *Engine) Patch(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodPatch, pattern, handler...)
}

func (e *Engine) Delete(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodDelete, pattern, handler...)
}

func (e *Engine) Connect(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodConnect, pattern, handler...)
}

func (e *Engine) Trace(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodTrace, pattern, handler...)
}

func (e *Engine) Options(pattern string, handler ...HandlerFunc) {
	e.addRoute(MethodOptions, pattern, handler...)
}

func (e *Engine) addRoute(method string, pattern string, handler ...HandlerFunc) {
	e.routes[method] = append(e.routes[method], &Route{
		Method:   method,
		Path:     pattern,
		Handlers: handler,
	})
}

func (e *Engine) Run(addr ...string) error {
	if addr == nil {
		addr[0] = e.config.Addr
	}
	return http.ListenAndServe(addr[0], e)
}
