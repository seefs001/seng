package main

import (
	"github.com/seefs001/seng"
)

func main() {
	engine := seng.Default()
	engine.Get("/aaaa", func(c *seng.Context) error {
		return c.String("aaa")
	})
	engine.Get("/abc", func(c *seng.Context) error {
		return c.JSON(seng.Response{
			Code: 200,
			Msg:  "success",
		})
	})
	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}
