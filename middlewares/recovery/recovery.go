package recovery

import (
	"fmt"
	"os"
	"runtime"

	"github.com/seefs001/seng"
)

var defaultStackTraceBufLen = 1024

func Default() seng.Handler {
	return func(c *seng.Context) (err error) {
		// Catch panics
		defer func() {
			if r := recover(); r != nil {
				defaultStackTraceHandler(r)

				var ok bool
				if err, ok = r.(error); !ok {
					// Set error that will call the global error handler
					err = fmt.Errorf("%v", r)
				}
			}
		}()

		// Return err if exist, else move to next handler
		return c.Next()
	}
}

func defaultStackTraceHandler(e interface{}) {
	buf := make([]byte, defaultStackTraceBufLen)
	buf = buf[:runtime.Stack(buf, false)]
	_, _ = os.Stderr.WriteString(fmt.Sprintf("panic: %v\n%s\n", e, buf))
}
