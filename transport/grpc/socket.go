package grpc

import (
	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc/pb"
)

type socket struct {
	stream pb.Transport_StreamServer
}

func (s *socket) Recv(m *transport.Message) error {
	if m == nil {
		return nil
	}

	msg, err := s.stream.Recv()
	if err != nil {
		return err
	}

	m.Header = msg.Header
	m.Body = msg.Body
	return nil
}

func (s *socket) Send(m *transport.Message) error {
	if m == nil {
		return nil
	}

	return s.stream.Send(&pb.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (s *socket) Close() error {
	return nil
}
