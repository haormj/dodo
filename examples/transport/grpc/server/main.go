package main

import (
	"log"
	"runtime/debug"

	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc"
)

func main() {
	t := grpc.NewTransport()
	ln, err := t.Listen(":8888")
	if err != nil {
		log.Fatalln(err)
	}
	ln.Accept(func(s transport.Socket) {
		defer func() {
			s.Close()
			if r := recover(); r != nil {
				log.Println(r)
				log.Println(string(debug.Stack()))
			}
		}()
		for {
			var m transport.Message
			if err := s.Recv(&m); err != nil {
				// if client close
				// rpc error: code = Canceled desc = context canceled
				// ignore this error
				return
			}
			log.Println(m)
			if err := s.Send(&m); err != nil {
				log.Println(err)
				return
			}
		}
	})
}
