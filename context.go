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
}

func (c *Context) Cookie(key string) []byte {
	return c.Fasthttp.Request.Header.Cookie(key)
}

func (c *Context) SetCookie(key string, value string) {
	//c.Fasthttp.Request.Header.Set(key, value)
	c.Fasthttp.Response.Header.Set(key, value)

	var cookie fasthttp.Cookie
	cookie.SetKey("cookie-name")
	cookie.SetValue("cookie-value")
	c.Fasthttp.Response.Header.SetCookie(&cookie)
}

func (c *Context) SetContentType(key string) {
	c.Fasthttp.Request.Header.SetContentType(key)
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

func (c *Context) SetHeader(key string, value string) {
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

func (c *Context) PostForm() *fasthttp.Args {
	return c.Fasthttp.PostArgs()
}

func (c *Context) QueryParams() *fasthttp.Args {
	return c.Fasthttp.QueryArgs()
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

func (c *Context) String(msg string) error {
	_, err := fmt.Fprint(c.Fasthttp, msg)
	if err != nil {
		return err
	}
	return nil
}

func (c *Context) JSON(data interface{}) error {
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
