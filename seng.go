package seng

import (
	"sync"

	"github.com/valyala/fasthttp"
)

// Engine app
type Engine struct {
	sync.Mutex

	route  *Router
	config Config
	server *fasthttp.Server
}

// Config 配置
type Config struct {
	// 服务器名
	Name string `json:"name"`
	// 监听地址
	Addr string `json:"addr"`
}

// New 新建Engine实例
func New(config ...Config) *Engine {

	router := &Router{Routes: make(map[string]Handler)}

	server := &fasthttp.Server{
		Name:    config[0].Name,
		Handler: router.RequestHandler,
	}

	engine := &Engine{
		config: config[0],
		server: server,
		route:  router,
	}
	return engine
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//func requestHandler(ctx *fasthttp.RequestCtx) {
//	fmt.Println(string(ctx.Path()))
//	if handler,ok := Routes[string(ctx.Path())];ok {
//		handler(&Context{
//			ctx,
//		})
//	}
//
//	//fmt.Fprintf(ctx, "Hello, world!\n\n")
//	//
//	//fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
//	//fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
//	//fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
//	//fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
//	//fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
//	//fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
//	//fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
//	//fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
//	//fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
//	//fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())
//	//
//	//fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)
//	//
//	//ctx.SetContentType("text/plain; charset=utf8")
//	//
//	//// Set arbitrary headers
//	//ctx.Response.Header.Set("X-My-Header", "my-header-value")
//	//
//	//// Set cookies
//	//var c fasthttp.Cookie
//	//c.SetKey("cookie-name")
//	//c.SetValue("cookie-value")
//	//ctx.Response.Header.SetCookie(&c)
//}

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

func (e *Engine) Get(path string, handler Handler) {
	e.route.Get(path, handler)
}
