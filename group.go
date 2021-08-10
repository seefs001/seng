package seng

import (
	"errors"
	"log"
	"net/http"
	"path"
)

// RouterGroup router group
// route -> /api/ping
// g := r.Group("/api")
// g.GET("/ping",handler)
type RouterGroup struct {
	prefix string
	// middlewares
	middleWares []Handler
	// tree
	parent *RouterGroup
	// reference to engine
	engine *Engine
}

// Group new router group
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute add route to router
func (g *RouterGroup) addRoute(method string, pattern string, handler Handler) {
	pattern = g.prefix + pattern
	if g.engine.config.Debug {
		log.Printf("Route %4s - %s", method, pattern)
	}
	g.engine.router.addRoute(method, pattern, handler)
}

func (g *RouterGroup) GET(pattern string, handler Handler) {
	g.addRoute(http.MethodGet, pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler Handler) {
	g.addRoute(http.MethodPost, pattern, handler)
}

func (g *RouterGroup) HEAD(pattern string, handler Handler) {
	g.addRoute(http.MethodHead, pattern, handler)
}

func (g *RouterGroup) PUT(pattern string, handler Handler) {
	g.addRoute(http.MethodPut, pattern, handler)
}

func (g *RouterGroup) DELETE(pattern string, handler Handler) {
	g.addRoute(http.MethodDelete, pattern, handler)
}

func (g *RouterGroup) TRACE(pattern string, handler Handler) {
	g.addRoute(http.MethodTrace, pattern, handler)
}

func (g *RouterGroup) CONNECT(pattern string, handler Handler) {
	g.addRoute(http.MethodConnect, pattern, handler)
}

func (g *RouterGroup) OPTIONS(pattern string, handler Handler) {
	g.addRoute(http.MethodOptions, pattern, handler)
}

func (g *RouterGroup) Use(middleWares ...Handler) {
	g.middleWares = append(g.middleWares, middleWares...)
}

func (g *RouterGroup) CreateStaticHandler(relativePath string, fs http.FileSystem) Handler {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(ctx *Context) error {
		file, exists := ctx.Param("filepath")
		if !exists {
			return errors.New("filepath is not exists")
		}
		// Determine whether the file exists or have permission to process the file
		if _, err := fs.Open(file); err != nil {
			ctx.Status(http.StatusNotFound)
			return nil
		}
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		return nil
	}
}

// Static Map the root path on the hard disk to the routing relativePath
func (g *RouterGroup) Static(relativePath string, root string) {
	handler := g.CreateStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	g.GET(urlPattern, handler)
}
