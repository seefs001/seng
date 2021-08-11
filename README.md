# seng

Web framework based on net/http, only for learning purposes

**_Warning: Do not use it for the production environment！！！_**

# Example

### Basic Example
```go
e := seng.Default()
e.GET("/ping", func(c *seng.Context) error {
	return c.JSON(seng.Map{
		"ping": pong,
	})
})
```

## Group

```go
g := engine.Group("/api")
g.GET("/ping", func(c *seng.Context) error {
	return c.JSON(seng.Map{
		"ping": " pong",
	})
})
```

## Path Params

```go
e.GET("/param/:name", func(c *seng.Context) error {
   param, exists := c.Param("name")
   if !exists {
      return c.Text("not found")
   }
   return c.JSON(seng.Map{
      "name": param,
   })
})
```

## Middleware & c.Next() & c.Get()/Set()

```go
g := group.Group("/test")
g.Use(func(context *seng.Context) error {
	context.Set("hello", "world")
	log.Default().Println("mv")
	return c.Next
})
g.GET("/mv", func(context *seng.Context) error {
	data, exists := context.Get("hello")
	if !exists {
		return context.Text("err")
	}
	return context.JSON(seng.Map{
		"ctx": data,
	})
})
```

internal middleware

```go
g.Use(cors.Default())
g.Use(recovery.Default())
g.Use(logger.Default())
```

Tips: If you need other middleware, please write it yourself.

## Cookies

```go
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
```

## BodyParser && Validator

```go
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
```

## Header

```go
e.GET("/header", func(c *seng.Context) error {
   c.SetHeader("X-token","token value")
   header := c.GetHeader("X-token")
   return c.JSON(seng.Map{
      "token":header,
   })
})
```

## Seng Config

```go
type Config struct {
   // Addr "ip:port"
   Addr string `json:"addr"`
   // When set to true, the Router treats "/foo" and "/foo/" as different.
   // Default: false
   StrictRouting bool `json:"strict_routing"`
   // Default: 4 * 1024 * 1024
   BodyLimit int `json:"body_limit"`
   // Default: unlimited
   ReadTimeout time.Duration `json:"read_timeout"`
   // Default: unlimited
   WriteTimeout time.Duration `json:"write_timeout"`
   // Default: unlimited
   IdleTimeout time.Duration `json:"idle_timeout"`
   // Default: false
   GETOnly bool `json:"get_only"`
   // print routes
   // Default: true
   Debug bool `json:"debug"`
   // Cookie
   // http.SameSiteStrictMode http.SameSiteLaxMode http.SameSiteNoneMode
   // http.SameSiteNoneMode must set secure to true
   CookieSameSite http.SameSite `json:"cookie_same_site"`
   // Default: false
   DisableKeepalive bool `json:"disable_keepalive"`
   // ErrorHandler Default: DefaultErrorHandler
   ErrorHandler ErrorHandler `json:"-"`
   // NotFoundHandler Default: DefaultNotFoundErrorHandler
   NotFoundErrorHandler Handler `json:"-"`
}
```

 If you have to see more examples, please see [examples](examples)