package gateway

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/smallnest/rpcx/client"
)

// ServiceHandler converts http.Request into rpcx.Request and send it to rpcx service,
// and then converts the result and writes it into http.Response.
// You should get the http.Request and servicePath in your web handler.
type ServiceHandler func(*http.Request, string) (map[string]string, []byte, error)

// HTTPServer is a golang web interface。
// You can use echo, gin, iris or other go web frameworks to implement it.
// You must wrap ServiceHandler into your handler of your selected web framework and add it into router.
type HTTPServer interface {
	RegisterHandler(base string, handler ServiceHandler)
	Serve() error
}

// Gateway is a rpcx gateway which can convert http invoke into rpcx invoke.
type Gateway struct {
	base       string
	httpserver HTTPServer

	serviceDiscovery client.ServiceDiscovery
	FailMode         client.FailMode
	SelectMode       client.SelectMode
	Option           client.Option

	mu       sync.RWMutex
	xclients map[string]client.XClient

	seq uint64
}

// NewGateway returns a new gateway.
func NewGateway(base string, hs HTTPServer, sd client.ServiceDiscovery, failMode client.FailMode, selectMode client.SelectMode, option client.Option) *Gateway {
	// base is empty or like /abc/
	if base == "" {
		base = "/"
	}

	if base[0] != '/' {
		base = "/" + base
	}

	g := &Gateway{
		base:             base,
		httpserver:       hs,
		serviceDiscovery: sd,
		FailMode:         failMode,
		SelectMode:       selectMode,
		Option:           option,
		xclients:         make(map[string]client.XClient),
	}

	hs.RegisterHandler(base, g.handler)
	return g
}

// Serve listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
func (g *Gateway) Serve() error {
	return g.httpserver.Serve()
}

func (g *Gateway) handler(r *http.Request, servicePath string) (meta map[string]string, payload []byte, err error) {
	req, err := HttpRequest2RpcxRequest(r)
	if err != nil {
		return nil, nil, err
	}
	seq := atomic.AddUint64(&g.seq, 1)
	req.SetSeq(seq)

	var xc client.XClient
	g.mu.Lock()
	xc, err = getXClient(g, servicePath)
	g.mu.Unlock()

	if err != nil {
		return nil, nil, err
	}

	return xc.SendRaw(context.Background(), req)
}

func getXClient(g *Gateway, servicePath string) (xc client.XClient, err error) {
	defer func() {
		if e := recover(); e != nil {
			if ee, ok := e.(error); ok {
				err = ee
				return
			}

			err = fmt.Errorf("failed to get xclient: %v", e)
		}
	}()

	if g.xclients[servicePath] == nil {
		d, err := g.serviceDiscovery.Clone(servicePath)
		if err != nil {
			return nil, err
		}
		g.xclients[servicePath] = client.NewXClient(servicePath, g.FailMode, g.SelectMode, d, g.Option)
	}
	xc = g.xclients[servicePath]

	return xc, err
}
