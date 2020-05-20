package std

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Logger is a middleware that prints http logs.
func Logger(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		ts := time.Now()
		h(w, r, params)
		te := time.Since(ts)
		fmt.Printf("\x1b[%d;%dm[GW]\x1b[0m  %s |  %d |  %v |  %s |  %s  %s\n", 44, 37,
			time.Now().Format("2006-01-02 - 15:04:05"), 200, te, ip, r.Method, r.URL.Path)
	}
}
