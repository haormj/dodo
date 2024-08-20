package main

import (
	"context"
	"log"

	"github.com/haormj/dodo/consumer"
)

func main() {
	c := consumer.NewConsumer()
	if err := c.Init(); err != nil {
		log.Fatalln(err)
	}

	var rsp string
	if err := c.Call(context.Background(), "Hello", "SayHello", "Hello", &rsp); err != nil {
		log.Fatalln(err)
	}
	log.Println(rsp)
}
