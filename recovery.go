package recovery

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gocontrib/context"
)

// Config for recovery middleware
type Config struct {
	Log func(string)
}

// New returns middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible.
//
// Prints a request ID if one is provided.
func New(config Config) func(http.Handler) http.Handler {
	if config.Log == nil {
		config.Log = defaultLog
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := context.GetRequestID(r)

			defer func() {
				if err := recover(); err != nil {
					// print panic
					var prefix = ""
					if len(reqID) > 0 {
						prefix = fmt.Sprintf("[%s] ", reqID)
					}
					config.Log(fmt.Sprintf("%spanic: %+v", prefix, err))

					debug.PrintStack()

					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func defaultLog(s string) {
	log.Print(s)
}
