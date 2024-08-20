package client

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/haormj/dodo/codec"
	"github.com/haormj/dodo/transport"
)

type Options struct {
	Transport transport.Transport
	Codecs    []codec.Codec
	// Connection Pool
	PoolSize int
	PoolTTL  time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type CallOptions struct {
	Codec     string
	TLSConfig *tls.Config
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
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

// PoolSize sets the connection pool size
func PoolSize(d int) Option {
	return func(o *Options) {
		o.PoolSize = d
	}
}

// PoolSize sets the connection pool size
func PoolTTL(d time.Duration) Option {
	return func(o *Options) {
		o.PoolTTL = d
	}
}

func WithCodec(c string) CallOption {
	return func(o *CallOptions) {
		o.Codec = c
	}
}

func WithTLSConfig(t *tls.Config) CallOption {
	return func(o *CallOptions) {
		o.TLSConfig = t
	}
}
