package main

import (
	"context"
	"log"

	"github.com/haormj/dodo/provider"
)

func SayHello(ctx context.Context, in string, out *string) error {
	*out = in + " dodo this is function example"
	return nil
}

func main() {
	p1 := provider.NewProvider(
		SayHello,
		provider.Name("Hello"),
	)
	if err := p1.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := p1.Run(); err != nil {
		log.Fatalln(err)
	}
}
