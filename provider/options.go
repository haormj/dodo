package provider

import (
	"time"

	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/registry/etcd"
	"github.com/haormj/dodo/server"
	"github.com/haormj/dodo/server/rpc"
)

type Options struct {
	Name             string
	Version          string
	Labels           map[string]string
	Servers          []server.Server
	Registries       []registry.Registry
	RegisterTTL      time.Duration
	RegisterInterval time.Duration
}

type Option func(*Options)

var (
	DefaultVersion          = "0.1.0"
	DefaultRegisterInterval = time.Second * 20
	DefaultRegisterTTL      = time.Second * 60
)

func newOptions(opts ...Option) Options {
	options := Options{
		Labels:           make(map[string]string),
		Version:          DefaultVersion,
		RegisterInterval: DefaultRegisterInterval,
		RegisterTTL:      DefaultRegisterTTL,
	}

	for _, o := range opts {
		o(&options)
	}

	if len(options.Servers) == 0 {
		options.Servers = append(options.Servers, rpc.NewServer())
	}

	if len(options.Registries) == 0 {
		options.Registries = append(options.Registries, etcd.NewRegistry())
	}

	return options
}

func Server(s server.Server) Option {
	return func(o *Options) {
		o.Servers = append(o.Servers, s)
	}
}

func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registries = append(o.Registries, r)
	}
}

// RegisterTTL specifies the TTL to use when registering the service
func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

// RegisterInterval specifies the interval on which to re-register
func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

func Label(k, v string) Option {
	return func(o *Options) {
		o.Labels[k] = v
	}
}
