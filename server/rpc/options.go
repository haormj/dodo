package rpc

import (
	"github.com/haormj/dodo/codec/json"
	"github.com/haormj/dodo/server"
	"github.com/haormj/dodo/transport/grpc"
)

func newOptions(opt ...server.Option) server.Options {
	options := server.Options{
		Address:   ":17312",
		Transport: grpc.NewTransport(),
		Wait:      true,
	}

	for _, o := range opt {
		o(&options)
	}

	if len(options.Codecs) == 0 {
		options.Codecs = append(options.Codecs, json.NewCodec())
	}

	return options
}
