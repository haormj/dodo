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
	service1 := registry.Service{
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
	service2 := registry.Service{
		Protocol:  "rest",
		Address:   "127.0.0.1:17312",
		Name:      "World",
		Version:   "0.1.0",
		Funcs:     []string{"SayWorld"},
		Codecs:    []string{"json"},
		Transport: "http",
		Side:      "provider",
		TLS:       true,
		Timestamp: 1543311057,
		Labels: map[string]string{
			"nodeID": "1",
		},
	}

	go func() {
		for {
			if err := r.Register(service1); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			if err := r.Register(service2); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second)
		}

	}()

	w, err := r.Watch()
	if err != nil {
		log.Fatalln(err)
	}

	for {
		result, err := w.Next()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(result)
	}

}
