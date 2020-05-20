package gateway

import "net/http"

// 定义404 handler
func XNotFound(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  w.WriteHeader(http.StatusNotFound) // StatusNotFound = 404
  w.Write([]byte("页面未找到"))
  HttpCode = 404        // 重写httpcode
  return
}

