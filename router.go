package seng

import (
	"strings"
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
	handlers        map[string]Handler
	root            map[string]*RouteNode
	Params          map[string]string
	NotFoundHandler Handler
}

func NewRouter(notFoundHandler ...Handler) *Router {

	handler := func(c *Context) error {
		return c.JSON(Response{
			Code: 404,
			Msg:  "Not Found",
		})
	}
	if notFoundHandler[0] != nil {
		handler = notFoundHandler[0]
	}

	return &Router{
		handlers:        make(map[string]Handler),
		root:            make(map[string]*RouteNode),
		NotFoundHandler: handler,
	}
}

func parsePath(path string) []string {
	parts := strings.Split(path, "/")

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

func (r *Router) RequestHandler(c *fasthttp.RequestCtx) {
	ctx := ctxPool.Get().(*Context)
	ctx.set(c)
	key := combineKey(utils.Bytes2String(c.Method()), utils.Bytes2String(c.Path()))
	routeNode, params := r.getRoute(utils.Bytes2String(c.Method()), utils.Bytes2String(c.Path()))
	r.Params = params
	// 未查找到
	if routeNode == nil {
		err := r.NotFoundHandler(ctx)
		if err != nil {
			// TODO error handler
			ctxPool.Put(ctx)
			return
		}
		ctxPool.Put(ctx)
		return
	}
	if handler, ok := r.handlers[key]; ok {
		err := handler(ctx)
		if err != nil {
			ctxPool.Put(ctx)
			return
		}
		ctxPool.Put(ctx)
		return
	} else {
		err := r.NotFoundHandler(ctx)
		if err != nil {
			// TODO error handler
			ctxPool.Put(ctx)
			return
		}
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
	parts := parsePath(path)
	key := combineKey(method, path)
	// 获取根节点
	_, ok := r.root[method]
	if !ok { // 不存在则创建
		r.root[method] = &RouteNode{}
	}
	// 往根节点中插入
	r.root[method].insert(path, parts, 0)
	r.handlers[key] = handler
}

func combineKey(method string, path string) string {
	return method + "+" + path
}

func (r *Router) getRoute(method string, path string) (*RouteNode, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string)

	root, ok := r.root[method] // 获取根节点
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil { // 匹配节点非空
		parts := parsePath(n.part)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 { // *开头，且不只有*
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
