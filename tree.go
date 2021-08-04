package seng

import (
	"strings"
)

type RouteNode struct {
	path          string
	part          string
	RouteChildren []*RouteNode
	ParamChildren []*RouteNode
}

func (r *RouteNode) matchChild(part string) *RouteNode {
	for _, child := range r.RouteChildren {
		if child.part == part {
			return child
		}
	}
	if r.ParamChildren != nil {
		return r.ParamChildren[0]
	}
	return nil
}

func (r *RouteNode) matchChildren(part string) []*RouteNode {
	nodes := make([]*RouteNode, 0)
	for _, child := range r.RouteChildren {
		if child.part == part {
			nodes = append(nodes, child)
		}
	}
	for _, child := range r.ParamChildren {
		if child.part == "*" {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (r *RouteNode) insert(pattern string, parts []string, depth int) {
	if len(parts) == depth { // 当深度到达目标层时
		r.path = pattern // 设置当前节点的 pattern 为注册的 pattern
		return
	}

	part := parts[depth]
	child := r.matchChild(part)
	if child == nil {
		if part[0] == ':' || part[0] == '*' {
			child = &RouteNode{
				part: part,
			}
			r.ParamChildren = append(r.ParamChildren, child)
		} else {
			child = &RouteNode{
				part: part,
			}
			r.RouteChildren = append(r.RouteChildren, child)
		}
	}
	// 递归调用
	child.insert(pattern, parts, depth+1)
}

func (r *RouteNode) search(parts []string, depth int) *RouteNode {
	if len(parts) == depth || strings.HasPrefix(r.part, "*") { // 到达最底层或者当前为 * 的模糊匹配
		if r.path == "" { // pattern 为空，说明未到最底层，查找失败
			return nil
		}
		return r
	}

	part := parts[depth]
	children := r.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, depth+1)
		if result != nil {
			return result
		}
	}
	return nil
}
