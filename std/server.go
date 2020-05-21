package std

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	. "github.com/rpcxio/rpcx-gateway"
)

// Server implements gateway.HTTPServer by using julienschmidt/httprouter.
type Server struct {
	addr       string
	router     *httprouter.Router
	Middleware func(next httprouter.Handle) httprouter.Handle
}

// New returns a server.
func New(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// NewWithRouter returns a server with preconfigured router.
func NewWithRouter(addr string, router *httprouter.Router) *Server {
	return &Server{
		addr:   addr,
		router: router,
	}
}

// RegisterHandler configures the handler to handle http rpcx invoke.
// It wraps ServiceHandler into httprouter.Handle.
func (s *Server) RegisterHandler(base string, handler ServiceHandler) {
	router := s.router
	if router == nil {
		router = httprouter.New()
	}
	h := wrapServiceHandler(handler)
	if s.Middleware != nil {
		h = s.Middleware(h)
	}

	router.POST(base, h)
	router.GET(base, h)
	router.PUT(base, h)
	s.router = router
}

func wrapServiceHandler(handler ServiceHandler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		if r.Header.Get(XServicePath) == "" {
			servicePath := params.ByName("servicePath")
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
			w.Write(payload)
			return
		}

		rh := r.Header
		for k, v := range rh {
			if strings.HasPrefix(k, "X-RPCX-") && len(v) > 0 {
				wh.Set(k, v[0])
			}
		}

		wh.Set(XMessageStatusType, "Error")
		wh.Set(XErrorMessage, err.Error())
		return
	}
}

func (s *Server) Serve() error {
	return http.ListenAndServe(s.addr, s.router)
}
