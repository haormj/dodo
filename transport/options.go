package transport

import (
	"context"
	"crypto/tls"
	"time"
)

type Options struct {
	// Timeout sets the timeout for Send/Recv
	Timeout time.Duration
	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type DialOptions struct {
	Timeout   time.Duration
	TLSConfig *tls.Config

	// TODO: add tls options when dialling
	// Currently set in global options

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

type ListenOptions struct {
	TLSConfig *tls.Config

	// Other options for implementations of the interface
	// can be stored in a context
	Context context.Context
}

// Timeout sets the timeout for Send/Recv execution
func Timeout(t time.Duration) Option {
	return func(o *Options) {
		o.Timeout = t
	}
}

// Timeout used when dialling the remote side
func WithTimeout(d time.Duration) DialOption {
	return func(o *DialOptions) {
		o.Timeout = d
	}
}

// WithDailTLSConfig support Dail by TLS
func WithDailTLSConfig(t *tls.Config) DialOption {
	return func(o *DialOptions) {
		o.TLSConfig = t
	}
}

// WithListenTLSConfig support Listen by TLS
func WithListenTLSConfig(t *tls.Config) ListenOption {
	return func(o *ListenOptions) {
		o.TLSConfig = t
	}
}
