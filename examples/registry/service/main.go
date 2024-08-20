package main

import (
	"log"

	"github.com/haormj/dodo/registry"
)

func main() {
	str := "rpc%3A%2F%2F172.18.1.131%3A17312%2FHello%3Fversion%3D0.1.0%26funcs%3DSayHello%26codecs%3Djson%26transport%3Dgrpc%26side%3Dprovider%26tls%3Dfalse%26timestamp%3D1543303208%26nodeID%3D1"
	s1, err := registry.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}
	s2, err := registry.Parse(str)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(registry.Equal(s1, s2))
}
