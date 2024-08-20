package consumer

import (
	"context"
	"crypto/tls"
	"errors"

	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/selector"
	"github.com/haormj/dodo/util"
)

type Consumer struct {
	opts           Options
	invokerManager *InvokerManager
	clients        map[string]client.Client
}

func NewConsumer(opts ...Option) *Consumer {
	options := newOptions(opts...)
	c := &Consumer{
		opts:           options,
		invokerManager: NewInvokerManager(),
		clients:        make(map[string]client.Client),
	}
	return c
}

func (c *Consumer) Init(opts ...Option) error {
	for _, o := range opts {
		o(&c.opts)
	}
	if err := c.opts.Registry.Init(); err != nil {
		return err
	}
	if err := c.opts.Selector.Init(selector.Registry(c.opts.Registry)); err != nil {
		return err
	}
	for _, cli := range c.opts.Clients {
		if err := cli.Init(); err != nil {
			return err
		}
		c.clients[cli.String()] = cli
	}

	return nil
}

func (c *Consumer) Call(ctx context.Context, serviceName string, funcName string,
	in interface{}, out interface{}, opts ...CallOption) error {
	callOpts := c.opts.CallOptions
	for _, o := range opts {
		o(&callOpts)
	}
	// TLS filter
	if callOpts.TLSConfig != nil {
		callOpts.Filters = append(callOpts.Filters, selector.FilterTLS())
	}
	// codec filter
	callOpts.Filters = append(callOpts.Filters, selector.FilterClient(c.opts.Clients))

	var service registry.Service
	if len(callOpts.Address) != 0 {
		service.Address = callOpts.Address
		service.Name = serviceName
	} else {
		var err error
		service, err = c.opts.Selector.Select(serviceName, selector.WithFilter(callOpts.Filters...))
		if err != nil {
			return err
		}
	}

	// get client by protocol
	cli, _ := c.clients[service.Protocol]
	// TODO this is a shit
	t := make([]string, 0)
	for _, cdc := range cli.Options().Codecs {
		t = append(t, cdc.String())
	}

	var fn = func(invoker.InvokeFunc) invoker.InvokeFunc {
		return func(ctx context.Context, mi invoker.Message,
			opts ...invoker.InvokeOption) (invoker.Message, error) {
			mo := invoker.NewMessage()
			params := mi.Parameters()
			var copts []client.CallOption
			copts = append(copts, client.WithCodec(util.ArrayIntersectString(t, service.Codecs)[0]))
			if service.TLS {
				if callOpts.TLSConfig != nil {
					copts = append(copts, client.WithTLSConfig(callOpts.TLSConfig))
				} else {
					copts = append(copts, client.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}))
				}
			}
			err := cli.Call(ctx, service.Address, serviceName, funcName,
				params[1], params[2], copts...)

			mo.SetParameters([]interface{}{err})
			return mo, nil
		}
	}

	inv := c.invokerManager.Get(serviceName)
	mi := invoker.NewMessage()
	mi.SetFuncName(funcName)
	mi.SetParameters([]interface{}{ctx, in, out})

	var i []invoker.Interceptor
	i = append(i, callOpts.Interceptors...)
	i = append(i, c.opts.Interceptors...)
	i = append(i, fn)
	mo, err := inv.Invoke(ctx, mi, invoker.WithInterceptor(i...))
	if err != nil {
		return err
	}

	if mo.Parameters()[0] != nil {
		err, ok := mo.Parameters()[0].(error)
		if !ok {
			return errors.New("message out parameter is not error")
		}
		return err
	}

	return nil
}

func (c *Consumer) Close() error {
	return nil
}
