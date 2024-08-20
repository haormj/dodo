package consumer

import (
	"crypto/tls"

	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/client/rpc"
	"github.com/haormj/dodo/codec"
	"github.com/haormj/dodo/codec/json"
	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/registry/etcd"
	"github.com/haormj/dodo/selector"
	"github.com/haormj/dodo/selector/cache"
)

type Options struct {
	Registry     registry.Registry
	Codecs       []codec.Codec
	Clients      []client.Client
	Selector     selector.Selector
	Interceptors []invoker.Interceptor

	CallOptions CallOptions
}

type Option func(*Options)

type CallOptions struct {
	// Address of remote host
	Address   string
	TLSConfig *tls.Config
	// Middleware for low level call func
	Interceptors []invoker.Interceptor
	Filters      []selector.Filter
}

type CallOption func(*CallOptions)

func newOptions(opts ...Option) Options {
	options := Options{
		Registry: etcd.NewRegistry(),
		Selector: cache.NewSelector(),
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Clients) == 0 {
		options.Clients = append(options.Clients, rpc.NewClient())
	}

	if len(options.Codecs) == 0 {
		options.Codecs = append(options.Codecs, json.NewCodec())
	}

	return options
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

func Codec(c codec.Codec) Option {
	return func(o *Options) {
		o.Codecs = append(o.Codecs, c)
	}
}

func Client(c client.Client) Option {
	return func(o *Options) {
		o.Clients = append(o.Clients, c)
	}
}

func Selector(s selector.Selector) Option {
	return func(o *Options) {
		o.Selector = s
	}
}

func Intercept(i ...invoker.Interceptor) Option {
	return func(o *Options) {
		o.Interceptors = append(o.Interceptors, i...)
	}
}

// WithAddress sets the remote address to use rather than using service discovery
func WithAddress(a string) CallOption {
	return func(o *CallOptions) {
		o.Address = a
	}
}

func WithTLSConfig(t *tls.Config) CallOption {
	return func(o *CallOptions) {
		o.TLSConfig = t
	}
}

// WithCallWrapper is a CallOption which adds to the existing CallFunc wrappers
func WithInterceptor(i ...invoker.Interceptor) CallOption {
	return func(o *CallOptions) {
		o.Interceptors = append(o.Interceptors, i...)
	}
}

func FilterLabel(k, v string) CallOption {
	return func(o *CallOptions) {
		o.Filters = append(o.Filters, selector.FilterLabel(k, v))
	}
}

func FilterVersion(version string) CallOption {
	return func(o *CallOptions) {
		o.Filters = append(o.Filters, selector.FilterVersion(version))
	}
}
