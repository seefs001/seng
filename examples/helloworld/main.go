package main

import "github.com/seefs001/seng"

func main() {
	engine := seng.Default()
	engine.Get("/ping", func(c *seng.Context) error {
		return c.Text("pong")
	})
	engine.Get("/test", func(c *seng.Context) error {
		return c.Json(seng.Response{
			Code:    1,
			Message: "success",
		})
	})
	engine.Run(":8080")
}
