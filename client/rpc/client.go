package rpc

import (
	"context"
	"errors"
	"sync"

	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/codec"
	"github.com/haormj/dodo/metadata"
	"github.com/haormj/dodo/transport"
)

type Client struct {
	sync.RWMutex
	opts   client.Options
	pool   *pool
	codecs map[string]codec.Codec
}

func NewClient(opts ...client.Option) client.Client {
	options := newOptions(opts...)

	c := &Client{
		opts:   options,
		pool:   newPool(options.PoolSize, options.PoolTTL),
		codecs: make(map[string]codec.Codec),
	}

	return c
}

func (c *Client) Init(opts ...client.Option) error {
	c.Lock()
	for _, o := range opts {
		o(&c.opts)
	}
	c.Unlock()

	for _, cdc := range c.opts.Codecs {
		c.codecs[cdc.String()] = cdc
	}

	return nil
}

func (c *Client) Options() client.Options {
	c.RLock()
	opts := c.opts
	c.RUnlock()
	return opts
}

func (c *Client) Call(ctx context.Context, address string, serviceName string, funcName string,
	req interface{}, rsp interface{}, opts ...client.CallOption) error {
	callOptions := client.CallOptions{
		Codec: "json",
	}
	for _, o := range opts {
		o(&callOptions)
	}

	cdc, ok := c.codecs[callOptions.Codec]
	if !ok {
		// TODO better handle error
		s := "not find codec " + callOptions.Codec
		return errors.New(s)
	}

	reqBytes, err := cdc.Marshal(req)
	if err != nil {
		return err
	}
	md, _ := metadata.FromContext(ctx)
	pi := protocol{
		header: header{
			ServiceName: serviceName,
			FuncName:    funcName,
			Codec:       callOptions.Codec,
			Metadata:    md,
		},
		Body: reqBytes,
	}
	// transport dail option
	var topts []transport.DialOption
	if callOptions.TLSConfig != nil {
		topts = append(topts, transport.WithDailTLSConfig(callOptions.TLSConfig))
	}

	conn, err := c.pool.getConn(address, c.opts.Transport, topts...)
	if err != nil {
		return err
	}
	var cerr error
	defer func() {
		c.pool.release(address, conn, cerr)
	}()
	mi := format(pi)
	if err := conn.Send(&mi); err != nil {
		cerr = err
		return err
	}
	var mo transport.Message
	if err := conn.Recv(&mo); err != nil {
		cerr = err
		return err
	}
	po := parse(mo)
	if len(po.header.Error) != 0 {
		return errors.New(po.header.Error)
	}
	if err := cdc.Unmarshal(po.Body, rsp); err != nil {
		return err
	}
	return nil
}

func (c *Client) String() string {
	return "rpc"
}
