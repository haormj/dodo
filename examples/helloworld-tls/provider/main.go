package main

import (
	"context"
	"log"

	"github.com/haormj/dodo/provider"
	"github.com/haormj/dodo/server"
	"github.com/haormj/dodo/server/rpc"
)

type Hello struct{}

func (h *Hello) SayHello(ctx context.Context, in string, out *string) error {
	*out = in + " dodo Hello"
	return nil
}

func main() {
	p := provider.NewProvider(
		new(Hello),
		provider.Label("nodeID", "1"),
		provider.Server(rpc.NewServer(server.TLSEnable(true))),
	)
	if err := p.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := p.Run(); err != nil {
		log.Fatalln(err)
	}
}
