package grpc

import (
	"crypto/tls"
	"net"

	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type listener struct {
	ln  net.Listener
	tls *tls.Config
}

func (l *listener) Addr() string {
	return l.ln.Addr().String()
}

func (l *listener) Close() error {
	return l.ln.Close()
}

func (l *listener) Accept(fn func(transport.Socket)) error {
	var opts []grpc.ServerOption

	if l.tls != nil {
		creds := credentials.NewTLS(l.tls)
		opts = append(opts, grpc.Creds(creds))
	}

	srv := grpc.NewServer(opts...)

	// register service
	pb.RegisterTransportServer(srv, &service{fn: fn})

	return srv.Serve(l.ln)
}
