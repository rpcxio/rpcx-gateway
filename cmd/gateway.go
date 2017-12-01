package main

import (
	"flag"

	gateway "github.com/rpcx-ecosystem/rpcx-gateway"
	"github.com/smallnest/rpcx/client"
)

var (
	addr        = flag.String("addr", ":9981", "http server address")
	serviceAddr = flag.String("saddr", "127.0.0.1:8972", "service address for Peer2PeerDiscovery")
	failmode    = flag.Int("failmode", int(client.Failover), "failMode, Failover in default")
	selectMode  = flag.Int("selectmode", int(client.RoundRobin), "selectMode, RoundRobin in default")
)

func main() {
	flag.Parse()

	d := client.NewPeer2PeerDiscovery("tcp@"+*serviceAddr, "")
	gw := gateway.NewGateway(*addr, d, client.FailMode(*failmode), client.SelectMode(*selectMode), client.DefaultOption)

	gw.Serve()
}
