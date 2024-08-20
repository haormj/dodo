package main

import (
	"log"
	"time"

	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/registry/etcd"
)

func main() {
	r := etcd.NewRegistry()
	if err := r.Init(); err != nil {
		log.Fatalln(err)
	}
	service := registry.Service{
		Protocol:  "rpc",
		Address:   "127.0.0.1:17312",
		Name:      "Hello",
		Version:   "0.1.0",
		Funcs:     []string{"SayHello"},
		Codecs:    []string{"json"},
		Transport: "grpc",
		Side:      "provider",
		TLS:       true,
		Timestamp: 1543311057,
		Labels: map[string]string{
			"nodeID": "1",
		},
	}
	if err := r.Register(service, registry.RegisterTTL(time.Second*10)); err != nil {
		log.Fatalln(err)
	}
}
