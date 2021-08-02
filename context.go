package seng

import (
	"encoding/json"
	"fmt"

	"github.com/seefs001/seng/utils"
	"github.com/valyala/fasthttp"
)

type Context struct {
	Fasthttp *fasthttp.RequestCtx
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
