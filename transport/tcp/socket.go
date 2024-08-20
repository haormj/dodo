package tcp

import (
	"bufio"
	"encoding/gob"
	"errors"
	"net"
	"time"

	"github.com/haormj/dodo/transport"
)

type socket struct {
	conn    net.Conn
	enc     *gob.Encoder
	dec     *gob.Decoder
	encBuf  *bufio.Writer
	timeout time.Duration
}

func (s *socket) Recv(m *transport.Message) error {
	if m == nil {
		return errors.New("message passed in is nil")
	}

	// set timeout if its greater than 0
	if s.timeout > time.Duration(0) {
		s.conn.SetDeadline(time.Now().Add(s.timeout))
	}

	return s.dec.Decode(&m)
}

func (s *socket) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if s.timeout > time.Duration(0) {
		s.conn.SetDeadline(time.Now().Add(s.timeout))
	}
	if err := s.enc.Encode(m); err != nil {
		return err
	}
	return s.encBuf.Flush()
}

func (s *socket) Close() error {
	return s.conn.Close()
}
