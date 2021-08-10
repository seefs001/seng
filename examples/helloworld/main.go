package main

import (
	"fmt"
	"log"

	"github.com/seefs001/seng"
)

func main() {
	engine := seng.Default()
	//engine.GET("/:name", func(c *seng.Context) error {
	//	param, exists := c.Param("name")
	//	if !exists{
	//		return c.Text("not found")
	//	}
	//	return c.JSON(seng.Map{
	//		"name":param,
	//	})
	//})
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
	//routerGroup.Use(cors.Default())
	//routerGroup.Use(recovery.Default())
	routerGroup.GET("/mv", func(context *seng.Context) error {
		data, exists := context.Get("x")
		fmt.Println(exists)
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
