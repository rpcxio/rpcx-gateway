package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	gateway "github.com/rpcxio/rpcx-gateway"

	"github.com/smallnest/rpcx/codec"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

var (
	n = flag.Int("n", 1000, "total requests")
	c = flag.Int("c", runtime.GOMAXPROCS(-1), "concurrent")
)

func main() {
	flag.Parse()

	cc := &codec.MsgpackCodec{}

	args := &Args{
		A: 10,
		B: 20,
	}

	data, _ := cc.Encode(args)

	count := *n / (*c)

	var wg sync.WaitGroup
	wg.Add(*c)
	start := time.Now()
	for i := 0; i < *c; i++ {
		i := i
		go func() {
			client := &http.Client{
				Transport: &http.Transport{
					MaxIdleConnsPerHost: 2000,
				},
				Timeout: time.Duration(1) * time.Second,
			}
			req, err := http.NewRequest("POST", "http://127.0.0.1:9981/", bytes.NewReader(data))
			if err != nil {
				log.Fatal("failed to create request: ", err)
				return
			}

			h := req.Header
			h.Set(gateway.XMessageType, "0")
			h.Set(gateway.XSerializeType, "3")
			h.Set(gateway.XServicePath, "Arith")
			h.Set(gateway.XServiceMethod, "Mul")

			for j := 0; j < count; j++ {
				h.Set(gateway.XMessageID, strconv.Itoa(i*count+j))
				req.Body = ioutil.NopCloser(bytes.NewReader(data))
				res, err := client.Do(req)
				if err != nil {
					log.Fatal("failed to call: ", err)
				}
				defer res.Body.Close()

				// handle http response
				replyData, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal("failed to read response: ", err)
				}

				reply := &Reply{}
				err = cc.Decode(replyData, reply)
				if err != nil {
					log.Fatal("failed to decode reply: ", err)
				}

				//log.Printf("goroutine %d: %d * %d = %d", i, args.A, args.B, reply.C)
			}

			wg.Done()
		}()
	}
	wg.Wait()
	t := time.Since(start).Nanoseconds()

	log.Printf("call %d times, took %d ms, tps: %d", *n, t/1e6, uint64(*n)*1e9/uint64(t))

}
