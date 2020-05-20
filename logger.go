package gateway

import (
  "fmt"
  "net"
  "net/http"
  "time"
)

/**
网关日志
 */


func Logger(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    HttpCode = 200
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    ts := time.Now()
    next.ServeHTTP(w, r)
    te := time.Since(ts)
    fmt.Printf("\x1b[%d;%dm[GW]\x1b[0m  %s |  %d |  %v |  %s |  %s  %s\n", 44, 37,
        time.Now().Format("2006-01-02 - 15:04:05"), HttpCode, te, ip, r.Method, r.URL.Path)
  })
}




