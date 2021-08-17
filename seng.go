package seng

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Version version of seng
const Version = "0.0.3"

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
	// Default: unlimited
	ReadHeaderTimeout time.Duration `json:"read_header_timeout"`
	// Default: false
	GETOnly bool `json:"get_only"`
	// print routes
	// Default: true
	Debug bool `json:"debug"`
	// Cookie
	// http.SameSiteStrictMode http.SameSiteLaxMode http.SameSiteNoneMode
	// http.SameSiteNoneMode must set secure to true
	CookieSameSite http.SameSite `json:"cookie_same_site"`
	// Default: false
	DisableKeepalive bool `json:"disable_keepalive"`
	// ErrorHandler Default: DefaultErrorHandler
	ErrorHandler ErrorHandler `json:"-"`
	// NotFoundHandler Default: DefaultNotFoundErrorHandler
	NotFoundErrorHandler Handler `json:"-"`
	// seng version
	SengVersion string      `json:"seng_version"`
	Logger      *log.Logger `json:"logger"`
}

// Default Config values
const (
	DefaultBodyLimit       = 4 * 1024 * 1024
	DefaultReadBufferSize  = 4096
	DefaultWriteBufferSize = 4096
	DefaultCookieSameSite  = http.SameSiteLaxMode
	DefaultListenAddr      = ":8080"
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
	Logger:               log.Default(),
	StrictRouting:        false,
	BodyLimit:            DefaultBodyLimit,
	GETOnly:              false,
	DisableKeepalive:     false,
	Debug:                true,
	CookieSameSite:       DefaultCookieSameSite,
	ErrorHandler:         DefaultErrorHandler,
	NotFoundErrorHandler: DefaultNotFoundErrorHandler,
}

// Engine struct
type Engine struct {
	mutex sync.Mutex
	// router
	*RouterGroup
	config Config
	Logger *log.Logger
	// Ctx pool
	ctxPool sync.Pool
	// Validator Pool
	validatorPool sync.Pool
	router        *Router
	groups        []*RouterGroup
	// template
	htmlTemplates *template.Template
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
		validatorPool: sync.Pool{New: func() interface{} {
			return new(Validator)
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
	if engine.config.CookieSameSite == 0 {
		engine.config.CookieSameSite = http.SameSiteLaxMode
	}
	if engine.Logger == nil {
		logger := log.Default()
		engine.config.Logger = logger
		engine.Logger = logger
	}
	engine.config.SengVersion = Version
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

	// debug print info
	if e.config.Debug {
		e.Logger.Printf("Seng version: %s", e.config.SengVersion)
		e.Logger.Printf("Disable keepalive: %t", e.config.DisableKeepalive)
		e.Logger.Printf("GETOnly: %t", e.config.GETOnly)
	}
	e.RouterGroup = &RouterGroup{engine: e}
	e.groups = []*RouterGroup{e.RouterGroup}
}

// Default engine with default middlewares and config
func Default() *Engine {
	engine := New(defaultConfig)

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
func (e *Engine) Listen(address ...string) (err error) {
	if len(address) == 0 {
		e.config.Addr = DefaultListenAddr
	} else {
		// set addr to config
		e.config.Addr = address[0]
	}
	s := http.Server{
		Addr:              e.config.Addr,
		Handler:           e,
		ReadTimeout:       e.config.ReadTimeout,
		ReadHeaderTimeout: e.config.ReadHeaderTimeout,
		WriteTimeout:      e.config.WriteTimeout,
		IdleTimeout:       e.config.IdleTimeout,
	}
	// disable keepalive
	s.SetKeepAlivesEnabled(!e.config.DisableKeepalive)
	if e.config.Debug {
		e.Logger.Printf("Listening on %s", e.config.Addr)
	}
	// http serve
	return s.ListenAndServe()
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
	// TODO graceful down
	return
}
