package seng

import (
	"fmt"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	r := NewRouterTrie()
	r.add("/test/ttttt", func(c *Context) error {
		return c.String("1")
	})
	r.add("/get/test/2", func(c *Context) error {
		return c.String("2")
	})
	node, path, _ := r.findNode([]string{"get", "test", "2"})
	join := strings.Join(path, "x")
	fmt.Println(join)
	fmt.Println(node.Key)
}
