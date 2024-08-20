package grpc

import (
	"runtime/debug"

	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc/pb"
)

type service struct {
	fn func(transport.Socket)
}

func (s *service) Stream(ts pb.Transport_StreamServer) error {
	sock := &socket{
		stream: ts,
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error(r, string(debug.Stack()))
			sock.Close()
		}
	}()

	s.fn(sock)
	return nil
}
