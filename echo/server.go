package echo

import (
	"net/http"
	"strings"

	"github.com/labstack/echo"
	. "github.com/rpcxio/rpcx-gateway"
)

type Server struct {
	addr string
	e    *echo.Echo
}

// New returns a server.
func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// NewWithEcho returns a server with preconfigured echo.
func NewWithEcho(addr string, e *echo.Echo) *Server {
	return &Server{
		addr: addr,
		e:    e,
	}
}

// RegisterHandler configures the handler to handle http rpcx invoke.
// It wraps ServiceHandler into httprouter.Handle.
func (s *Server) RegisterHandler(base string, handler ServiceHandler) {
	e := s.e
	if e == nil {
		e = echo.New()
	}
	h := wrapServiceHandler(handler)

	e.POST(base, h)
	e.GET(base, h)
	e.PUT(base, h)
	s.e = e
}

func wrapServiceHandler(handler ServiceHandler) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		r := ctx.Request()
		w := ctx.Response()
		if r.Header.Get(XServicePath) == "" {
			servicePath := ctx.Param("servicePath")
			if strings.HasPrefix(servicePath, "/") {
				servicePath = servicePath[1:]
			}

			r.Header.Set(XServicePath, servicePath)
		}

		servicePath := r.Header.Get(XServicePath)
		messageID := r.Header.Get(XMessageID)
		wh := w.Header()
		if messageID != "" {
			wh.Set(XMessageID, messageID)
		}

		meta, payload, err := handler(r, servicePath)
		for k, v := range meta {
			wh.Set(k, v)
		}

		if err == nil {
			ctx.Blob(http.StatusOK, "application/octet-stream", payload)
			return nil
		}

		rh := r.Header
		for k, v := range rh {
			if strings.HasPrefix(k, "X-RPCX-") && len(v) > 0 {
				wh.Set(k, v[0])
			}
		}

		wh.Set(XMessageStatusType, "Error")
		wh.Set(XErrorMessage, err.Error())
		ctx.String(http.StatusOK, err.Error())
		return nil
	}
}

func (s *Server) Serve() error {
	return s.e.Start(s.addr)
}
