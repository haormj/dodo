// Package client is an interface for an RPC client
package client

import (
	"context"
	"time"
)

// Client is the interface used to make requests to services.
// It supports Request/Response via Transport and Publishing via the Broker.
// It also supports bidiectional streaming of requests.
type Client interface {
	Init(...Option) error
	Options() Options
	Call(ctx context.Context, address string, serviceName string, funcName string,
		req interface{}, rsp interface{}, opts ...CallOption) error
	String() string
}

// Option used by the Client
type Option func(*Options)

// CallOption used by Call or Stream
type CallOption func(*CallOptions)

var (
	// DefaultPoolSize sets the connection pool size
	DefaultPoolSize = 1
	// DefaultPoolTTL sets the connection pool ttl
	DefaultPoolTTL = time.Minute
)
