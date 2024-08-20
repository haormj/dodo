package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/haormj/dodo/invoker/receiver"
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
	s := rpc.NewServer()
	if err := s.Init(); err != nil {
		log.Fatalln(err)
	}
	if err := s.Register(inv); err != nil {
		log.Fatalln(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			if err := s.Start(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	time.Sleep(time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			if err := s.Stop(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	time.Sleep(time.Second * 5)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5000; i++ {
			if err := s.Stop(); err != nil {
				log.Fatalln(err)
			}
		}
	}()

	wg.Wait()
}
