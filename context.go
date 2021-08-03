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

func (c *Context) SetCookieKV(key string, value string) {
	//c.Fasthttp.Request.Header.Set(key, value)
	c.Fasthttp.Response.Header.Set(key, value)

	cookie := fasthttp.AcquireCookie()
	cookie.SetKey("cookie-name")
	cookie.SetValue("cookie-value")
	c.Fasthttp.Response.Header.SetCookie(cookie)
}

func (c *Context) Set(key string, value interface{}) {
	c.Fasthttp.SetUserValue(key, value)
}

func (c *Context) SetUserValueBytes(key []byte, value interface{}) {
	c.Fasthttp.SetUserValueBytes(key, value)
}

func (c *Context) UserValue(key string) interface{} {
	return c.Fasthttp.UserValue(key)
}

func (c *Context) UserValueBytes(key []byte) interface{} {
	return c.Fasthttp.UserValueBytes(key)
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
