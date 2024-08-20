package tcp

import (
	"bufio"
	"encoding/gob"
	"net"
	"time"

	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/transport"
)

type listener struct {
	ln      net.Listener
	timeout time.Duration
}

func (l *listener) Addr() string {
	return l.ln.Addr().String()
}

func (l *listener) Close() error {
	return l.ln.Close()
}

func (l *listener) Accept(fn func(transport.Socket)) error {
	var tempDelay time.Duration

	for {
		c, err := l.ln.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Errorf("http: Accept error: %v; retrying in %v\n", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		encBuf := bufio.NewWriter(c)
		sock := &socket{
			timeout: l.timeout,
			conn:    c,
			encBuf:  encBuf,
			enc:     gob.NewEncoder(encBuf),
			dec:     gob.NewDecoder(c),
		}

		go func() {
			// TODO: think of a better error response strategy
			defer func() {
				if r := recover(); r != nil {
					sock.Close()
				}
			}()

			fn(sock)
		}()
	}
}
