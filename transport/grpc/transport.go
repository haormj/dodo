package grpc

import (
	"context"
	"net"

	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc/pb"
	"github.com/haormj/dodo/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Transport implement by grpc
type Transport struct {
	opts transport.Options
}

// NewTransport -
func NewTransport(opts ...transport.Option) transport.Transport {
	var options transport.Options
	for _, o := range opts {
		o(&options)
	}

	t := &Transport{
		opts: options,
	}

	return t
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

	for _, o := range opts {
		o(&dopts)
	}

	options := []grpc.DialOption{
		grpc.WithTimeout(dopts.Timeout),
	}

	if dopts.TLSConfig != nil {
		creds := credentials.NewTLS(dopts.TLSConfig)
		options = append(options, grpc.WithTransportCredentials(creds))
	} else {
		options = append(options, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(addr, options...)
	if err != nil {
		return nil, err
	}

	stream, err := pb.NewTransportClient(conn).Stream(context.Background())
	if err != nil {
		return nil, err
	}

	return &client{
		conn:   conn,
		stream: stream,
	}, nil

}

func (t *Transport) Listen(addr string, opts ...transport.ListenOption) (transport.Listener, error) {
	var options transport.ListenOptions
	for _, o := range opts {
		o(&options)
	}

	ln, err := util.Listen(addr, func(addr string) (net.Listener, error) {
		return net.Listen("tcp", addr)
	})
	if err != nil {
		return nil, err
	}
	return &listener{
		ln:  ln,
		tls: options.TLSConfig,
	}, nil
}

func (t *Transport) String() string {
	return "grpc"
}
