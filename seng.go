package seng

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// Engine app
type Engine struct {
	sync.Mutex
	*RouterGroup

	groups []*RouterGroup
	router *Router
	config Config
	server *fasthttp.Server
}

// Config 配置
type Config struct {
	// 服务器名
	Name string `json:"name"`
	// 监听地址
	Addr string `json:"addr"`
	// NotFoundHandler
	NotFoundHandler Handler `json:"error_handler"`
}

// New 新建Engine实例
func New(config ...Config) *Engine {

	router := NewRouter(config[0].NotFoundHandler)

	server := &fasthttp.Server{
		Name:    config[0].Name,
		Handler: router.RequestHandler,
	}

	engine := &Engine{
		config: config[0],
		server: server,
		router: router,
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func Default() *Engine {

	defaultConfig := Config{
		Name: "My Web Application",
	}

	engine := New(defaultConfig)
	engine.server.ErrorHandler = func(ctx *fasthttp.RequestCtx, err error) {
		ctx.Error("xxxx", 200)
	}
	// TODO Logger和Revocer
	return engine
}

func (e *Engine) Run(addr string) error {
	return e.server.ListenAndServe(addr)
}

func (e *Engine) Get(path string, handler Handler) *Engine {
	e.router.Get(path, handler)
	return e
}

func (e *Engine) Post(path string, handler Handler) *Engine {
	e.router.add(MethodPost, path, handler)
	return e
}

func (e *Engine) Put(path string, handler Handler) *Engine {
	e.router.add(MethodPut, path, handler)
	return e
}

func (e *Engine) Delete(path string, handler Handler) *Engine {
	e.router.add(MethodDelete, path, handler)
	return e
}
func (e *Engine) Head(path string, handler Handler) *Engine {
	e.router.add(MethodHead, path, handler)
	return e
}

func (e *Engine) Patch(path string, handler Handler) *Engine {
	e.router.add(MethodPatch, path, handler)
	return e
}

func (e *Engine) Connect(path string, handler Handler) *Engine {
	e.router.add(MethodConnect, path, handler)
	return e
}

func (e *Engine) Trace(path string, handler Handler) *Engine {
	e.router.add(MethodTrace, path, handler)
	return e
}

func (e *Engine) Options(path string, handler Handler) *Engine {
	e.router.add(MethodOptions, path, handler)
	return e
}
