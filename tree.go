package seng

import "strings"

type RouteTrie struct {
	handler       Handler
	Key           string
	RouteChildren map[string]*RouteTrie
	ParamChildren map[string]*RouteTrie
}

func NewRouterTrie() *RouteTrie {
	return &RouteTrie{
		handler: func(c *Context) error {
			return c.String("")
		},
		RouteChildren: make(map[string]*RouteTrie, 10),
		ParamChildren: make(map[string]*RouteTrie, 10),
	}
}

func (r *RouteTrie) add(path string, handler Handler) {
	arr := strings.Split(path[1:], "/")
	p := r
	for index, item := range arr {
		if index == len(arr)-1 {
			// the last one, insert directly
			newNode := &RouteTrie{
				handler:       handler,
				Key:           item,
				RouteChildren: make(map[string]*RouteTrie),
				ParamChildren: make(map[string]*RouteTrie),
			}
			if item[0] == ':' {
				if t, ok := p.ParamChildren[item]; ok {
					t.handler = handler
				} else {
					p.ParamChildren[item] = newNode
				}
			} else {
				if t, ok := p.RouteChildren[item]; ok {
					t.handler = handler
				} else {
					p.RouteChildren[item] = newNode
				}
			}
		} else {
			// maybe the node not exists
			newNode := &RouteTrie{
				handler:       nil,
				Key:           item,
				RouteChildren: make(map[string]*RouteTrie),
				ParamChildren: make(map[string]*RouteTrie),
			}
			if item[0] == ':' {
				if _, ok := p.ParamChildren[item]; !ok {
					p.ParamChildren[item] = newNode
				}
				p = p.ParamChildren[item]
			} else {
				if _, ok := p.RouteChildren[item]; !ok {
					p.RouteChildren[item] = newNode
				}
				p = p.RouteChildren[item]
			}
		}
	}
}

// 匹配到的节点、Params 的 Key、Params 的 Value
func (r *RouteTrie) findNode(arr []string) (*RouteTrie, []string, Handler) {
	item := arr[0]
	// 2. the last one
	if len(arr) == 1 {
		if target, ok := r.RouteChildren[item]; ok {
			return target, []string{}, target.handler
		}
		for _, target := range r.ParamChildren {
			return target, []string{target.Key}, target.handler
		}
		return nil, []string{}, nil
	}
	// 1. use RouteChildren
	if target, ok := r.RouteChildren[item]; ok {
		return target.findNode(arr[1:])
	}
	// 3. search in ParamChildren
	for _, target := range r.ParamChildren {
		if t, keys, _ := target.findNode(arr[1:]); t != nil {
			return t, append(keys, target.Key), t.handler
		}
	}
	return nil, []string{}, nil
}
