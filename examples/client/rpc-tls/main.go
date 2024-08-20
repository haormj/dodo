package main

import (
	"context"
	"crypto/tls"
	"log"

	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/client/rpc"
)

func main() {
	c := rpc.NewClient()
	if err := c.Init(); err != nil {
		log.Fatalln(err)
	}
	req := "hello"
	rsp := ""
	if err := c.Call(context.Background(), "172.18.1.131:17312", "Hello", "SayHello",
		&req, &rsp, client.WithTLSConfig(&tls.Config{
			InsecureSkipVerify: true,
		})); err != nil {
		log.Fatalln(err)
	}
	log.Println(rsp)
}
