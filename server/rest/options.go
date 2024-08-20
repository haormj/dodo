package rest

import (
	"github.com/haormj/dodo/server"
)

func newOptions(opts ...server.Option) server.Options {
	options := server.Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}
