package seng

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net"

	"github.com/seefs001/seng/utils"
	"github.com/valyala/fasthttp"
)

type Context struct {
	Fasthttp *fasthttp.RequestCtx

	handlers []Handler
	index    int8

	fullPath   []byte
	path       []byte
	postBody   []byte
	method     []byte
	host       []byte
	remoteAddr net.Addr
	statusCode int

	queryParams []Params
	postForm    []Params

	header *Header
	engine *Engine
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *Context) set(fasthttp *fasthttp.RequestCtx) {
	c.Fasthttp = fasthttp
	c.path = c.Fasthttp.Path()
	c.fullPath = c.Fasthttp.RequestURI()
	c.postBody = c.Fasthttp.PostBody()
	c.method = c.Fasthttp.Method()
	c.host = c.Fasthttp.Host()
	c.header = &Header{Fasthttp: &c.Fasthttp.Request.Header}
	c.remoteAddr = c.Fasthttp.RemoteAddr()
	// args
	c.Fasthttp.QueryArgs().VisitAll(func(key, value []byte) {
		c.queryParams = append(c.queryParams, Params{
			Key:   utils.Bytes2String(key),
			Value: utils.Bytes2String(value),
		})
	})
	c.Fasthttp.PostArgs().VisitAll(func(key, value []byte) {
		c.postForm = append(c.postForm, Params{
			Key:   utils.Bytes2String(key),
			Value: utils.Bytes2String(value),
		})
	})
}

func (c *Context) Cookie(key string) []byte {
	return c.Fasthttp.Request.Header.Cookie(key)
}

func (c *Context) SetCookieKV(key string, value string) *Context {
	//c.Fasthttp.Request.Header.Set(key, value)
	c.Fasthttp.Response.Header.Set(key, value)

	cookie := fasthttp.AcquireCookie()
	cookie.SetKey("cookie-name")
	cookie.SetValue("cookie-value")
	c.Fasthttp.Response.Header.SetCookie(cookie)

	return c
}

func (c *Context) Set(key string, value interface{}) *Context {
	c.Fasthttp.SetUserValue(key, value)
	return c
}

func (c *Context) SetUserValueBytes(key []byte, value interface{}) *Context {
	c.Fasthttp.SetUserValueBytes(key, value)
	return c
}

func (c *Context) UserValue(key string) interface{} {
	return c.Fasthttp.UserValue(key)
}

func (c *Context) UserValueBytes(key []byte) interface{} {
	return c.Fasthttp.UserValueBytes(key)
}

func (c *Context) SetContentType(key string) *Context {
	c.Fasthttp.Request.Header.SetContentType(key)
	return c
}

func (c *Context) Path() []byte {
	return c.fullPath
}

func (c *Context) Body() []byte {
	return c.postBody
}

func (c *Context) Method() []byte {
	return c.method
}

func (c *Context) Header() *Header {
	return c.header
}

func (c *Context) SetHeaderKV(key string, value string) {
	c.Fasthttp.Response.Header.Set(key, value)
}

func (c *Context) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *Context) GetHeader(key string) []byte {
	return c.header.Get(key)
}

func (c *Context) Host() []byte {
	return c.host
}

func (c *Context) PostForm() []Params {
	return c.postForm
}

func (c *Context) QueryParams() []Params {
	return c.queryParams
}

func (c *Context) Query(key string) (string, error) {
	data := c.Fasthttp.QueryArgs().Peek(key)
	if data == nil {
		return "", ErrQueryParamNotFound
	}
	return utils.Bytes2String(data), nil
}

func (c *Context) QueryDefaultValue(key string, defaultValue string) string {
	query, err := c.Query(key)
	if err != nil {
		return defaultValue
	}
	return query
}

func (c *Context) FormValue(key string) []byte {
	return c.Fasthttp.FormValue(key)
}

func (c *Context) FormFile(key string) (*multipart.FileHeader, error) {
	return c.Fasthttp.FormFile(key)
}

func (c *Context) MultipartForm() (*multipart.Form, error) {
	return c.Fasthttp.MultipartForm()
}

func (c *Context) RemoteIP() net.IP {
	return c.Fasthttp.RemoteIP()
}

func (c *Context) String(format string) error {
	c.SetHeaderKV("Content-Type", "text/plain")
	_, err := fmt.Fprint(c.Fasthttp, format)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) JSON(data interface{}) error {
	c.SetHeaderKV("Content-Type", "application/json")
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(c.Fasthttp, utils.Bytes2String(d))
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) HTML(data interface{}) error {
	c.SetHeaderKV("Content-Type", "text/html")
	_, err := fmt.Fprint(c.Fasthttp, data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) Status(status int) *Context {
	c.statusCode = status
	c.Fasthttp.Response.SetStatusCode(status)
	return c
}

func (c *Context) SaveFile(fileheader *multipart.FileHeader, path string) error {
	return fasthttp.SaveMultipartFile(fileheader, path)
}
