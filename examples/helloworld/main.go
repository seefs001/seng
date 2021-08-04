package main

import (
	"github.com/seefs001/seng"
)

func main() {
	engine := seng.Default()
	engine.Get("/api/aaaa", func(c *seng.Context) error {
		return c.String("aaa")
	})

	engine.Get("/api/aaaa/test", func(c *seng.Context) error {
		return c.JSON(seng.Response{
			Code: 0,
			Msg:  "test",
		})
	})
	engine.Get("/api/*", func(c *seng.Context) error {
		return c.JSON(seng.Response{
			Code: 1,
			Msg:  "wild",
		})
	})
	engine.Get("/abc", func(c *seng.Context) error {
		query := c.QueryDefaultValue("test", "testdefault")
		c.SetCookieKV("test", "cookietest")
		c.SetHeaderKV("test-header", "xxxxx")
		return c.JSON(seng.Response{
			Code: 200,
			Msg:  query,
		})
	})
	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
