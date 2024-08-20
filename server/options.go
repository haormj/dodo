package server

import (
	"context"

	"github.com/haormj/dodo/codec"
	"github.com/haormj/dodo/transport"
)

type Options struct {
	Address     string
	Wait        bool
	TLSEnable   bool
	TLSKeyFile  string
	TLSCertFile string

	Transport transport.Transport
	Codecs    []codec.Codec

	Context context.Context
}

type Option func(*Options)

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

func TLSEnable(e bool) Option {
	return func(o *Options) {
		o.TLSEnable = e
	}
}

func TLSKeyFile(k string) Option {
	return func(o *Options) {
		o.TLSKeyFile = k
	}
}

func TLSCertFile(c string) Option {
	return func(o *Options) {
		o.TLSCertFile = c
	}
}

// Transport mechanism for communication e.g http, rabbitmq, etc
func Transport(t transport.Transport) Option {
	return func(o *Options) {
		o.Transport = t
	}
}

func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codecs = append(o.Codecs, c)
	}
}

// Wait tells the server to wait for requests to finish before exiting
func Wait(b bool) Option {
	return func(o *Options) {
		o.Wait = b
	}
}
