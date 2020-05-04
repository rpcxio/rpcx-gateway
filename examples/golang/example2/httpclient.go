package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

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

func main() {
	cc := &codec.MsgpackCodec{}

	args := &Args{
		A: 10,
		B: 20,
	}

	data, _ := cc.Encode(args)

	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		i := i

		go func() {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				req, err := http.NewRequest("POST", "http://127.0.0.1:9981/", bytes.NewReader(data))
				if err != nil {
					log.Fatal("failed to create request: ", err)
					return
				}

				h := req.Header
				// h.Set(gateway.XMessageID, strconv.Itoa(rand.Int()))
				h.Set(gateway.XMessageType, "0")
				h.Set(gateway.XSerializeType, "3")
				h.Set(gateway.XServicePath, "Arith")
				h.Set(gateway.XServiceMethod, "Mul")

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					log.Fatal("#%d failed to call: ", i, err)
				}
				defer res.Body.Close()

				// handle http response
				replyData, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal("#%d failed to read response: ", i, err)
				}

				reply := &Reply{}
				err = cc.Decode(replyData, reply)
				if err != nil {
					log.Fatal("#%d failed to decode reply: ", i, err)
				}

				log.Printf("#%d: %d * %d = %d", i, args.A, args.B, reply.C)
			}
		}()
	}

	wg.Wait()
}
