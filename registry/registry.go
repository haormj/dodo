// Package registry is an interface for service discovery
package registry

import "errors"

// The registry provides an interface for service discovery
// and an abstraction over varying implementations
// {consul, etcd, zookeeper, ...}
type Registry interface {
	Init(...Option) error
	Register(Service, ...RegisterOption) error
	Deregister(Service) error
	GetService(string) ([]Service, error)
	ListServices() ([]Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

type WatchOption func(*WatchOptions)

var (
	ErrNotFound = errors.New("not found")
)
