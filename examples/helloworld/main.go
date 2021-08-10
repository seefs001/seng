package main

import (
	"log"

	"github.com/seefs001/seng"
)

func main() {
	engine := seng.Default()
	group := engine.Group("/api")
	group.GET("/", func(context *seng.Context) error {
		return context.JSON(seng.Map{
			"xx": " xxx",
		})
	})
	routerGroup := group.Group("/test")
	routerGroup.Use(func(context *seng.Context) error {
		context.Set("x", "xxx")
		log.Default().Println("mv")
		return nil
	})
	routerGroup.GET("/mv", func(context *seng.Context) error {
		data, exists := context.Get("x")
		if !exists {
			return context.Text("err")
		}
		return context.JSON(seng.Map{
			"x":   "mv",
			"ctx": data,
		})
	})
	log.Fatal(engine.Listen(":8080"))
}
