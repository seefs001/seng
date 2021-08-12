package main

import (
	"fmt"
	"log"

	"github.com/seefs001/seng"
	"github.com/seefs001/seng/examples/helloworld/service"
	"github.com/seefs001/seng/middlewares/cors"
	"github.com/seefs001/seng/middlewares/logger"
	"github.com/seefs001/seng/middlewares/recovery"
)

func main() {
	engine := seng.Default()
	engine.Use(logger.Default())
	engine.Use(cors.Default())
	engine.Use(recovery.Default())
	engine.GET("/param/:name", func(c *seng.Context) error {
		param, exists := c.Param("name")
		if !exists {
			return c.Text("not found")
		}
		return c.JSON(seng.Map{
			"name": param,
		})
	})

	testService := service.NewTestService()
	handler := NewTestHandler(testService)
	engine.GET("/pets", seng.AdapterHandlerFunc(handler.AddPet))

	engine.GET("/header", func(c *seng.Context) error {
		c.SetHeader("X-token", "token value")
		header := c.GetHeader("X-token")
		return c.JSON(seng.Map{
			"token": header,
		})
	})
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
		return context.Next()
	})
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
	type TestRequest struct {
		Username string `json:"username" validate:"required#please input username.|min=1#please enter the legal parameters"`
		Password string `json:"password"`
	}
	routerGroup.POST("/parser", func(c *seng.Context) error {
		req := new(TestRequest)
		err := c.BodyParser(req)
		if err != nil {
			return err
		}
		err = c.Validate(*req)
		if err != nil {
			return c.JSON(seng.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(seng.Map{
			"username": req.Username,
			"password": req.Password,
		})
	})
	routerGroup.GET("/cookies", func(c *seng.Context) error {
		c.SetCookieWithValue("hello", "world", 300, false, false)
		cookie, err := c.GetCookie("hello")
		if err != nil {
			return err
		}
		return c.JSON(seng.Map{
			"cookie": cookie,
		})
	})
	log.Fatal(engine.Listen(":8080"))
}
