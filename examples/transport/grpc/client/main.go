package main

import (
	"log"
	"time"

	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc"
)

func main() {
	t := grpc.NewTransport()
	c, err := t.Dial("127.0.0.1:8888", transport.WithTimeout(time.Second*2))
	if err != nil {
		log.Fatalln(err)
	}

	m := transport.Message{
		Header: map[string]string{"hello": "world"},
		Body:   nil,
	}
	var mm transport.Message

	for i := 0; i < 10; i++ {
		if err := c.Send(&m); err != nil {
			log.Fatalln(err)
		}

		if err := c.Recv(&mm); err != nil {
			log.Fatalln(err)
		}
		log.Println(mm)
	}
	if err := c.Close(); err != nil {
		log.Fatalln(err)
	}
}
