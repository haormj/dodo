package rpc

import (
	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/codec/json"
	"github.com/haormj/dodo/transport/grpc"
)

func newOptions(opt ...client.Option) client.Options {
	options := client.Options{
		Transport: grpc.NewTransport(),
		PoolSize:  client.DefaultPoolSize,
		PoolTTL:   client.DefaultPoolTTL,
	}

	for _, o := range opt {
		o(&options)
	}

	if len(options.Codecs) == 0 {
		options.Codecs = append(options.Codecs, json.NewCodec())
	}

	return options
}
