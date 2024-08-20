package cache

import (
	"context"

	"github.com/haormj/dodo/selector"
)

type cacheDirKey struct{}
type configDirKey struct{}

func CacheDir(d string) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, cacheDirKey{}, d)
	}
}

func ConfigDir(d string) selector.Option {
	return func(o *selector.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, configDirKey{}, d)
	}
}

func newOptions(opts ...selector.Option) selector.Options {
	options := selector.Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}
