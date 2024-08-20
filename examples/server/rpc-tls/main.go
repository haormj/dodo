package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/haormj/dodo/invoker/receiver"
	"github.com/haormj/dodo/server"
	"github.com/haormj/dodo/server/rpc"
)

type Hello struct{}

func (*Hello) SayHello(ctx context.Context, req string, rsp *string) error {
	*rsp = req + " dodo"
	return nil
}

func main() {
	inv := receiver.NewInvoker(new(Hello))
	if err := inv.Init(); err != nil {
		log.Fatalln(err)
	}
	s := rpc.NewServer(
		server.TLSEnable(true),
	)
	if err := s.Init(); err != nil {
		log.Fatalln(err)
	}
	if err := s.Register(inv); err != nil {
		log.Fatalln(err)
	}

	if err := s.Start(); err != nil {
		log.Fatalln(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	var sig os.Signal
	select {
	case sig = <-ch:
		log.Printf("receive signal %s\n", sig.String())
		// stop to receive signal
		signal.Stop(ch)
	}

	if err := s.Stop(); err != nil {
		log.Fatalln(err)
	}
}
