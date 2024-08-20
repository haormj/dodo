package tcp

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"

	"github.com/haormj/dodo/transport"
)

type client struct {
	dialOpts transport.DialOptions
	conn     net.Conn
	enc      *gob.Encoder
	dec      *gob.Decoder
	encBuf   *bufio.Writer
	timeout  time.Duration
}

func (c *client) Send(m *transport.Message) error {
	// set timeout if its greater than 0
	if c.timeout > time.Duration(0) {
		c.conn.SetDeadline(time.Now().Add(c.timeout))
	}
	if err := c.enc.Encode(m); err != nil {
		return err
	}
	return c.encBuf.Flush()
}

func (c *client) Recv(m *transport.Message) error {
	// set timeout if its greater than 0
	if c.timeout > time.Duration(0) {
		c.conn.SetDeadline(time.Now().Add(c.timeout))
	}
	return c.dec.Decode(&m)
}

func (c *client) Close() error {
	return c.conn.Close()
}
