// Package selector is a way to load balance service nodes
package selector

import (
	"errors"

	"github.com/haormj/dodo/registry"
)

// Selector builds on the registry as a mechanism to pick nodes
// and mark their status. This allows host pools and other things
// to be built using various algorithms.
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	// Select returns a function which should return the next node
	Select(service string, opts ...SelectOption) (registry.Service, error)
	// Mark sets the success/error against a node
	Mark(service string, err error)
	// Reset returns state back to zero for a service
	Reset(service string)
	// Close renders the selector unusable
	Close() error
	// Name of the selector
	String() string
}

// Filter is used to filter a service during the selection process
type Filter func([]registry.Service) []registry.Service

// Strategy is a selection strategy e.g random, round robin
type Strategy func([]registry.Service) (registry.Service, error)

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)
