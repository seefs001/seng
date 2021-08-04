package seng

type RouterGroup struct {
	prefix      string
	middleWares []Handler    // 中间件
	parent      *RouterGroup // 支持多级分组
	engine      *Engine
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	engine := r.engine
	newGroup := &RouterGroup{
		prefix: r.prefix + prefix,
		parent: r,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (r *RouterGroup) Get(path string, handler Handler) {
	r.add(MethodGet, path, handler)
}

func (r *RouterGroup) Post(path string, handler Handler) {
	r.add(MethodPost, path, handler)
}

func (r *RouterGroup) Put(path string, handler Handler) {
	r.add(MethodPut, path, handler)
}

func (r *RouterGroup) Delete(path string, handler Handler) {
	r.add(MethodDelete, path, handler)
}
func (r *RouterGroup) Head(path string, handler Handler) {
	r.add(MethodHead, path, handler)
}

func (r *RouterGroup) Patch(path string, handler Handler) {
	r.add(MethodPatch, path, handler)
}

func (r *RouterGroup) Connect(path string, handler Handler) {
	r.add(MethodConnect, path, handler)
}

func (r *RouterGroup) Trace(path string, handler Handler) {
	r.add(MethodTrace, path, handler)
}

func (r *RouterGroup) Options(path string, handler Handler) {
	r.add(MethodOptions, path, handler)
}

func (r *RouterGroup) add(method string, path string, handler Handler) {
	path = r.prefix + path
	r.engine.router.add(method, path, handler)
}
