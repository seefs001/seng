package seng

import (
	"strings"
)

// Router handlers
type Router struct {
	// node
	roots map[string]*node
	// Amount of registered routes
	routesCount uint32
	// Amount of registered handlers
	handlerCount uint32
	handlers     map[string]Handler
}

// NewRouter create a new Router instance
func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]Handler),
	}
}

// parsePattern
// /static/*filepath ->  [static]
// /static/:name/doc ->  [static, :name, doc]
func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")

	results := make([]string, 0)

	for _, part := range parts {
		if part != "" {
			results = append(results, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return results
}

// combineRouteKey method-key
func combineRouteKey(method string, pattern string) string {
	return method + "-" + pattern
}

// addRoute add route to router
func (r *Router) addRoute(method string, pattern string, handler Handler) {
	parts := parsePattern(pattern)
	key := combineRouteKey(method, pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	// insert node
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// getRoute match route
func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	// wild match
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			// TODO wild match :,*
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			// example: *aaa
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// handle function to handle request
func (r *Router) handle(c *Context) error {
	// get params
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		// set to r.Params
		c.Params = params
		key := combineRouteKey(c.Method, n.pattern)
		if handler, ok := r.handlers[key]; ok {
			c.handlers = append(c.handlers, handler)
		} else {
			c.handlers = append(c.handlers, c.engine.config.NotFoundErrorHandler)
		}
	}
	// not found
	return c.Next()
}
