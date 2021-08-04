package seng

import "time"

type Config struct {
	// 严格路由 /foo/ /foo Default: false
	StrictRouting bool `json:"strict_routing"`
	// 大小写敏感
	CaseSensitive bool `json:"case_sensitive"`
	// Body限制
	BodyLimit int `json:"body_limit"`

	GETOnly          bool         `json:"get_only"`
	ErrorHandler     ErrorHandler `json:"-"`
	NotFoundHandler  HandlerFunc  `json:"-"`
	DisableKeepalive bool         `json:"disable_keepalive"`

	Addr    string `json:"addr"`
	AppName string `json:"app_name"`

	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
}
