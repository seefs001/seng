package seng

import "github.com/valyala/fasthttp"

type Header struct {
	Fasthttp *fasthttp.RequestHeader
}

func (h *Header) Get(key string) []byte {
	return h.Fasthttp.Peek(key)
}
