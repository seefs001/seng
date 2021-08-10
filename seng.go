package seng

import (
	"html/template"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/seefs001/seng/utils"
)

// Version version of seng
const Version = "0.0.1"

// Handler defines a function to serve HTTP requests.
type Handler func(c *Context) error

// ErrorHandler defines a function to handle errors
type ErrorHandler func(c *Context, err error) error

// Map a shortcut for map[string]interface{}
type Map map[string]interface{}

type Config struct {
	// Addr "ip:port"
	Addr string `json:"addr"`
	// When set to true, the Router treats "/foo" and "/foo/" as different.
	// Default: false
	StrictRouting bool `json:"strict_routing"`
	// Default: 4 * 1024 * 1024
	BodyLimit int `json:"body_limit"`
	// Default: unlimited
	ReadTimeout time.Duration `json:"read_timeout"`
	// Default: unlimited
	WriteTimeout time.Duration `json:"write_timeout"`
	// Default: unlimited
	IdleTimeout time.Duration `json:"idle_timeout"`
	// Default: false
	GETOnly bool `json:"get_only"`
	// print routes
	// Default: true
	Debug bool `json:"debug"`
	// Default: false
	DisableKeepalive bool `json:"disable_keepalive"`
	// ErrorHandler Default: DefaultErrorHandler
	ErrorHandler ErrorHandler `json:"-"`
	// NotFoundHandler Default: DefaultNotFoundErrorHandler
	NotFoundErrorHandler Handler `json:"-"`
}

// Default Config values
const (
	DefaultBodyLimit       = 4 * 1024 * 1024
	DefaultReadBufferSize  = 4096
	DefaultWriteBufferSize = 4096
)

// DefaultErrorHandler default error handler.
var DefaultErrorHandler = func(c *Context, err error) error {
	code := http.StatusInternalServerError
	if e, ok := err.(*Error); ok {
		code = e.Code
	}
	return c.Status(code).Text(err.Error())
}

var DefaultNotFoundErrorHandler = func(c *Context) error {
	code := http.StatusNotFound
	return c.Status(code).Text("404 not found")
}

// defaultConfig default engine config
var defaultConfig = Config{
	StrictRouting:    false,
	BodyLimit:        DefaultBodyLimit,
	GETOnly:          false,
	DisableKeepalive: false,
	Debug:            true,
	ErrorHandler:     DefaultErrorHandler,
}

// Engine struct
type Engine struct {
	mutex sync.Mutex

	*RouterGroup
	config Config
	// Ctx pool
	ctxPool       sync.Pool
	router        *Router
	groups        []*RouterGroup     // 存储所有分组
	htmlTemplates *template.Template // 用于 html 渲染
	funcMap       template.FuncMap
}

// New create a new instance of Engine
func New(config ...Config) *Engine {
	engine := &Engine{
		router: NewRouter(),
		// context pool
		ctxPool: sync.Pool{New: func() interface{} {
			return new(Context)
		}},
		config: Config{},
	}

	if len(config) > 0 {
		engine.config = config[0]
	}

	// Override default values
	if engine.config.BodyLimit == 0 {
		engine.config.BodyLimit = DefaultBodyLimit
	}
	if engine.config.ErrorHandler == nil {
		engine.config.ErrorHandler = DefaultErrorHandler
	}
	if engine.config.NotFoundErrorHandler == nil {
		engine.config.NotFoundErrorHandler = DefaultNotFoundErrorHandler
	}
	if engine.config.Debug == false {
		engine.config.Debug = true
	}
	// init Engine
	engine.init()
	return engine
}

// SetMode set server mode
func (e *Engine) SetMode(mode bool) {
	e.config.Debug = mode
}

// SetDebugMode set debug mode
func (e *Engine) SetDebugMode() {
	e.config.Debug = ModeDebug
}

// SetReleaseMode set release mode
func (e *Engine) SetReleaseMode() {
	e.config.Debug = ModeRelease
}

// init engine
func (e *Engine) init() {
	// Lock
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}
}

// Default engine with default middlewares and config
func Default() *Engine {
	engine := New(defaultConfig)

	// apply middlewares
	// TODO
	utils.ApplyDefaultMiddlewares(engine)
	return engine
}

// ServeHTTP implements http.Handler
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middleWares []Handler
	// add middlewares
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middleWares = append(middleWares, group.middleWares...)
		}
	}
	// ctxPool
	ctx := e.AcquireCtx().ReSet(w, req)
	//ctx := NewContext(w, req)
	ctx.engine = e
	ctx.router = e.router
	ctx.handlers = middleWares
	// handle request
	if err := e.router.handle(ctx); err != nil {
		err := e.config.ErrorHandler(ctx, err)
		if err != nil {
			// release
			e.ReleaseCtx(ctx)
			return
		}
	}
	// release
	e.ReleaseCtx(ctx)
}

// AcquireCtx acquired context from ctxPool
func (e *Engine) AcquireCtx() *Context {
	return e.ctxPool.Get().(*Context)
}

func (e *Engine) ReleaseCtx(ctx *Context) {
	// clean
	ctx.Writer = nil
	ctx.Request = nil
	ctx.handlers = nil
	// put to ctxPool
	e.ctxPool.Put(ctx)
	return
}

// Listen serve seng instance
func (e *Engine) Listen(addr string) (err error) {
	// set addr to config
	e.config.Addr = addr
	// http serve
	return http.ListenAndServe(addr, e)
}

// Config get engine config
func (e *Engine) Config() Config {
	return e.config
}

// ShutDown clean resources
func (e *Engine) ShutDown() (err error) {
	// Lock
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return
}
