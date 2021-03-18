module github.com/rpcxio/rpcx-gateway

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/julienschmidt/httprouter v1.3.0
	github.com/labstack/echo v3.3.10+incompatible
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/rpcxio/rpcx-etcd v0.0.0-20201229103411-8317fc934fbb
	github.com/smallnest/rpcx v0.0.0-20210313032152-d969a5f620f8
	github.com/xtaci/lossyconn v0.0.0-20190602105132-8df528c0c9ae // indirect
)

replace (
	github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.3
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
