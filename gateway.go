package gateway

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/smallnest/rpcx/client"
	"golang.org/x/net/http2"
)

type ServerType string

const (
	HTTP1  ServerType = "http1"
	HTTP2             = "http2"
	HTTP2c            = "h2c"
)

type Gateway struct {
	// http listen address
	Addr       string
	ServerType ServerType

	serviceDiscovery client.ServiceDiscovery
	FailMode         client.FailMode
	SelectMode       client.SelectMode
	Option           client.Option

	mu       sync.RWMutex
	xclients map[string]client.XClient
}

func NewGateway(addr string, st ServerType, sd client.ServiceDiscovery, failMode client.FailMode, selectMode client.SelectMode, option client.Option) *Gateway {

	return &Gateway{
		Addr:             addr,
		ServerType:       st,
		serviceDiscovery: sd,
		FailMode:         failMode,
		SelectMode:       selectMode,
		Option:           option,
		xclients:         make(map[string]client.XClient),
	}
}

func (g *Gateway) Serve() {
	router := httprouter.New()
	router.POST("/*servicePath", g.handleRequest)
	router.GET("/*servicePath", g.handleRequest)
	router.PUT("/*servicePath", g.handleRequest)

	switch g.ServerType {
	case HTTP2c:
		g.startH2c(router)
	case HTTP2:
		panic("unsupported")
	default:
		g.startHttp1(router)
	}
}

func (g *Gateway) startHttp1(handler http.Handler) {
	if err := http.ListenAndServe(g.Addr, handler); err != nil {
		log.Fatalf("error in ListenAndServe: %s", err)
	}
}

func (g *Gateway) startH2c(handler http.Handler) {
	server := &http.Server{
		Addr:         g.Addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	//http2.Server.ServeConn()
	s2 := &http2.Server{
		IdleTimeout: 1 * time.Minute,
	}
	http2.ConfigureServer(server, s2)
	l, _ := net.Listen("tcp", g.Addr)
	defer l.Close()
	log.Println("Start server...")
	for {
		rwc, err := l.Accept()
		if err != nil {
			log.Println("accept err:", err)
			continue
		}
		go s2.ServeConn(rwc, &http2.ServeConnOpts{BaseConfig: server})
	}
}

func (g *Gateway) handleRequest(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Header.Get(XServicePath) == "" {
		servicePath := params.ByName("servicePath")
		if strings.HasPrefix(servicePath, "/") {
			servicePath = servicePath[1:]
		}

		r.Header.Set(XServicePath, servicePath)
	}

	servicePath := r.Header.Get(XServicePath)

	wh := w.Header()
	req, err := HttpRequest2RpcxRequest(r)
	if err != nil {
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

	var xc client.XClient
	g.mu.Lock()
	xc, err = getXClient(g, servicePath)
	g.mu.Unlock()

	if err != nil {
		wh.Set(XMessageStatusType, "Error")
		wh.Set(XErrorMessage, err.Error())
		return
	}

	m, payload, err := xc.SendRaw(context.Background(), req)
	for k, v := range m {
		wh.Set(k, v)
	}
	if err != nil {
		wh.Set(XMessageStatusType, "Error")
		wh.Set(XErrorMessage, err.Error())
		return
	}

	w.Write(payload)

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
		g.xclients[servicePath] = client.NewXClient(servicePath, g.FailMode, g.SelectMode, g.serviceDiscovery.Clone(servicePath), g.Option)
	}
	xc = g.xclients[servicePath]

	return xc, err
}
