package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	gateway "github.com/rpcx-ecosystem/rpcx-gateway"
	"github.com/smallnest/rpcx/client"
)

var (
	addr       = flag.String("addr", ":9981", "http server address")
	st         = flag.String("st", "http1", "server type: http1 or h2c")
	registry   = flag.String("registry", "peer2peer://127.0.0.1:8972", "registry address")
	basePath   = flag.String("basepath", "/rpcx", "basepath for zookeeper, etcd and consul")
	failmode   = flag.Int("failmode", int(client.Failover), "failMode, Failover in default")
	selectMode = flag.Int("selectmode", int(client.RoundRobin), "selectMode, RoundRobin in default")
)

func main() {
	flag.Parse()

	d, err := createServiceDiscovery(*registry)
	if err != nil {
		log.Fatal(err)
	}
	gw := gateway.NewGateway(*addr, gateway.ServerType(*st), d, client.FailMode(*failmode), client.SelectMode(*selectMode), client.DefaultOption)

	gw.Serve()
}

func createServiceDiscovery(regAddr string) (client.ServiceDiscovery, error) {
	i := strings.Index(regAddr, "://")
	if i < 0 {
		return nil, errors.New("wrong format registry address. The right fotmat is [registry_type://address]")
	}

	regType := regAddr[:i]
	regAddr = regAddr[i+3:]

	switch regType {
	case "peer2peer": //peer2peer://127.0.0.1:8972
		return client.NewPeer2PeerDiscovery("tcp@"+regAddr, ""), nil
	case "multiple":
		var pairs []*client.KVPair
		pp := strings.Split(regAddr, ",")
		for _, v := range pp {
			pairs = append(pairs, &client.KVPair{Key: v})
		}
		return client.NewMultipleServersDiscovery(pairs), nil
	case "zookeeper":
		return client.NewZookeeperDiscoveryTemplate(*basePath, []string{regAddr}, nil), nil
	case "etcd":
		return client.NewEtcdDiscoveryTemplate(*basePath, []string{regAddr}, nil), nil
	case "consul":
		return client.NewConsulDiscoveryTemplate(*basePath, []string{regAddr}, nil), nil
	case "mdns":
		return client.NewMDNSDiscoveryTemplate(10*time.Second, 10*time.Second, ""), nil
	default:
		return nil, fmt.Errorf("wrong registry type %s. only support peer2peer,multiple, zookeeper, etcd, consul and mdns", regType)
	}
}
