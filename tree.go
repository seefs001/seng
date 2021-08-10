package seng

import "strings"

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否模糊匹配，part 开头为 : 或 * 时为true
}

// 子节点中第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children { // 遍历所有子节点
		// 满足子节点的 part 与当前 part 相等 or 当前子节点为模糊匹配
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的子节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children { // 遍历所有子节点
		// 满足子节点的 part 与当前 part 相等 or 当前子节点为模糊匹配
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 注册路由时，插入结点
func (n *node) insert(pattern string, parts []string, depth int) {
	if len(parts) == depth { // 当深度到达目标层时
		n.pattern = pattern // 设置当前节点的 pattern 为注册的 pattern
		return
	}

	part := parts[depth]
	child := n.matchChild(part)
	if child == nil {
		child = &node{ //  非最底层的节点不设置 pattern
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		// 将匹配到的第一个子节点加入当前节点的 children 列表中
		n.children = append(n.children, child)
	}
	// 递归调用
	child.insert(pattern, parts, depth+1)
}

// 匹配路由时，需要查找满足条件的节点
func (n *node) search(parts []string, depth int) *node {
	if len(parts) == depth || strings.HasPrefix(n.part, "*") { // 到达最底层或者当前为 * 的模糊匹配
		if n.pattern == "" { // pattern 为空，说明未到最底层，查找失败
			return nil
		}
		return n
	}

	part := parts[depth]
	children := n.matchChildren(part)

	for _, child := range children { // 遍历所有满足条件的子节点
		result := child.search(parts, depth+1) // 递归调用
		if result != nil {
			return result
		}
	}
	return nil
}
