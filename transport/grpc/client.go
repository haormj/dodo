package grpc

import (
	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/transport/grpc/pb"

	"google.golang.org/grpc"
)

type client struct {
	conn   *grpc.ClientConn
	stream pb.Transport_StreamClient
}

func (c *client) Recv(m *transport.Message) error {
	if m == nil {
		return nil
	}

	msg, err := c.stream.Recv()
	if err != nil {
		return err
	}
	m.Header = msg.Header
	m.Body = msg.Body
	return nil
}

func (c *client) Send(m *transport.Message) error {
	if m == nil {
		return nil
	}

	return c.stream.Send(&pb.Message{
		Header: m.Header,
		Body:   m.Body,
	})
}

func (c *client) Close() error {
	return c.conn.Close()
}
