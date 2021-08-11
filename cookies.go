package seng

import (
	"net/http"
)

// GetCookie get cookie by key
func (c *Context) GetCookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

// GetCookies get all cookies
func (c *Context) GetCookies() []*http.Cookie {
	return c.Request.Cookies()
}

// SetCookie set cookie
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Writer, cookie)
	return
}

// SetCookieWithValue set cookie with key value expires
func (c *Context) SetCookieWithValue(key, value string, expires int, httpOnly bool, secure bool) {
	cookie := &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: expires,
		// Set httponly = true cookies cannot be obtained by JS, unable
		// to use Document.cookie to play cookie content.
		HttpOnly: httpOnly,
		// If a cookie is set for secure = true, this cookie can only be sent
		// to the server with HTTPS protocol.
		Secure: secure,
		// http.SameSiteStrictMode http.SameSiteLaxMode http.SameSiteNoneMode
		// http.SameSiteNoneMode must set secure to true
		SameSite: c.engine.config.CookieSameSite,
	}
	http.SetCookie(c.Writer, cookie)
	return
}
