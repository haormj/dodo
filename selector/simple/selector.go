package simple

import (
	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/selector"
)

type Selector struct {
	so selector.Options
}

func NewSelector(opts ...selector.Option) selector.Selector {
	sopts := selector.Options{
		Strategy: selector.Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	return &Selector{
		so: sopts,
	}
}

func (r *Selector) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(&r.so)
	}
	return nil
}

func (r *Selector) Options() selector.Options {
	return r.so
}

func (r *Selector) Select(service string, opts ...selector.SelectOption) (registry.Service, error) {
	var svc registry.Service
	sopts := selector.SelectOptions{
		Strategy: r.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	services, err := r.so.Registry.GetService(service)
	if err != nil {
		return svc, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return svc, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services)
}

func (r *Selector) Mark(service string, err error) {
	return
}

func (r *Selector) Reset(service string) {
	return
}

func (r *Selector) Close() error {
	return nil
}

func (r *Selector) String() string {
	return "simple"
}
