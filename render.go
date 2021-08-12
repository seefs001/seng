// Package seng stole from https://github.com/masseelch/render/blob/master/render.go
package seng

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
)

type Response struct {
	Code   int         `json:"code" xml:"code"`                         // http response status code
	Status string      `json:"status" xml:"status"`                     // user-level status message
	Errors interface{} `json:"errors,omitempty" xml:"errors,omitempty"` // application-level error messages, for debugging
}

func NewResponse(code int, err interface{}) Response {
	r := Response{
		Code:   code,
		Status: http.StatusText(code),
		Errors: err,
	}
	return r
}

func RenderBadRequest(w http.ResponseWriter, r *http.Request, err interface{}) {
	resp := NewResponse(http.StatusBadRequest, err)
	Render(w, r, resp.Code, resp)
}

func RenderCreated(w http.ResponseWriter, r *http.Request, err interface{}) {
	Render(w, r, http.StatusCreated, err)
}

func RenderForbidden(w http.ResponseWriter, r *http.Request, err interface{}) {
	resp := NewResponse(http.StatusForbidden, err)
	Render(w, r, resp.Code, resp)
}

func RenderInternalServerError(w http.ResponseWriter, r *http.Request, err interface{}) {
	resp := NewResponse(http.StatusInternalServerError, err)
	Render(w, r, resp.Code, resp)
}

func RenderNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func RenderNotFound(w http.ResponseWriter, r *http.Request, err interface{}) {
	resp := NewResponse(http.StatusNotFound, err)
	Render(w, r, resp.Code, resp)
}

func RenderUnauthorized(w http.ResponseWriter, r *http.Request, err interface{}) {
	resp := NewResponse(http.StatusUnauthorized, err)
	Render(w, r, resp.Code, resp)
}

func RenderOK(w http.ResponseWriter, r *http.Request, err interface{}) {
	Render(w, r, http.StatusOK, err)
}

func RenderPartialContent(w http.ResponseWriter, r *http.Request, err interface{}) {
	Render(w, r, http.StatusPartialContent, err)
}

func Render(w http.ResponseWriter, r *http.Request, code int, d interface{}) {
	switch r.Header.Get(HeaderAccept) {
	case ContentTypeXml:
		RenderXML(w, code, d)
	default:
		RenderJSON(w, code, d)
	}
}

func RenderJSON(w http.ResponseWriter, code int, d interface{}) {
	b, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, ContentTypeJson+CharsetSuffix)
	w.WriteHeader(code)
	_, _ = w.Write(b)
}

func RenderXML(w http.ResponseWriter, code int, d interface{}) {
	buf := new(bytes.Buffer)
	if err := xml.NewEncoder(w).Encode(d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set(HeaderContentType, ContentTypeXml+CharsetSuffix)
	w.WriteHeader(code)
	_, _ = w.Write(buf.Bytes())
}

func RenderRaw(w http.ResponseWriter, code int, d []byte) {
	w.Header().Set(HeaderContentType, ContentTypeTextPlain+CharsetSuffix)
	w.WriteHeader(code)
	_, _ = w.Write(d)
}

func RenderHTML(w http.ResponseWriter, code int, d []byte) {
	w.Header().Set(HeaderContentType, ContentTypeTextHtml+CharsetSuffix)
	w.WriteHeader(code)
	_, _ = w.Write(d)
}
