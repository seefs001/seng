package logger

import (
	"log"
	"time"

	"github.com/seefs001/seng"
)

func Default() seng.Handler {
	return func(c *seng.Context) error {
		t := time.Now()
		_ = c.Next()
		latency := time.Since(t)
		log.Printf("path: %s --- method: %s --- latency: %s", c.Path, c.Method, latency)
		return nil
	}
}
