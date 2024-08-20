// Package tcp provides a TCP transport
package tcp

import (
	"bufio"
	"crypto/tls"
	"encoding/gob"
	"net"

	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/util"
)

type Transport struct {
	opts transport.Options
}

func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}
	return &Transport{opts: options}
}

func (t *Transport) Init(opts ...transport.Option) error {
	for _, o := range opts {
		o(&t.opts)
	}
	return nil
}

func (t *Transport) Options() transport.Options {
	return t.opts
}

func (t *Transport) Dial(addr string, opts ...transport.DialOption) (transport.Client, error) {
	dopts := transport.DialOptions{
		Timeout: transport.DefaultDialTimeout,
	}

	for _, opt := range opts {
		opt(&dopts)
	}

	var conn net.Conn
	var err error

	if dopts.TLSConfig != nil {
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: dopts.Timeout}, "tcp", addr,
			dopts.TLSConfig)
	} else {
		conn, err = net.DialTimeout("tcp", addr, dopts.Timeout)
	}

	if err != nil {
		return nil, err
	}

	encBuf := bufio.NewWriter(conn)

	return &client{
		dialOpts: dopts,
		conn:     conn,
		encBuf:   encBuf,
		enc:      gob.NewEncoder(encBuf),
		dec:      gob.NewDecoder(conn),
		timeout:  t.opts.Timeout,
	}, nil
}

func (t *Transport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	var l net.Listener
	var err error

	// TODO: support use of listen options
	if options.TLSConfig != nil {
		fn := func(addr string) (net.Listener, error) {
			return tls.Listen("tcp", addr, options.TLSConfig)
		}

		l, err = util.Listen(addr, fn)
	} else {
		fn := func(addr string) (net.Listener, error) {
			return net.Listen("tcp", addr)
		}

		l, err = util.Listen(addr, fn)
	}

	if err != nil {
		return nil, err
	}

	return &listener{
		timeout: t.opts.Timeout,
		ln:      l,
	}, nil
}

func (t *Transport) String() string {
	return "tcp"
}
