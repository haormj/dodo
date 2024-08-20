package rpc

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"runtime/debug"
	"sync"

	"github.com/haormj/dodo/codec"
	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/metadata"
	"github.com/haormj/dodo/server"
	"github.com/haormj/dodo/transport"
	"github.com/haormj/dodo/util"
)

// Server implement by rpc, in the life cycle of a rpc server instance,
// it can Start or Stop once actually,but you can repeatedly invoke Start and Stop method.
// Invoke Start counter will increase 1,and invoke Stop counter will decrease 1 each time.
// Only when counter == 1 && started == false, server will truly Start.
// Only when started == true && counter == 0 && stopped == false, server will truly Stop.
// Once the server is stopped, it cannot be started.
// You should avoid calling Start more than 2^32 times, which will cause the counter to
// overflow and cause the server to stop prematurely.
type Server struct {
	sync.RWMutex
	exit     chan chan error
	opts     server.Options
	invokers map[string]invoker.Invoker
	codecs   map[string]codec.Codec

	// graceful exit
	wg sync.WaitGroup

	started bool
	stopped bool
	// call Start +1, call Stop -1
	counter uint32
}

// NewServer implement by rpc
func NewServer(opts ...server.Option) server.Server {
	options := newOptions(opts...)
	return &Server{
		opts:     options,
		invokers: make(map[string]invoker.Invoker),
		codecs:   make(map[string]codec.Codec),
		exit:     make(chan chan error),
	}
}

func (s *Server) accept(sock transport.Socket) {
	defer func() {
		// close socket
		sock.Close()

		if r := recover(); r != nil {
			log.Error("panic recovered: ", r)
			log.Error(string(debug.Stack()))
		}
	}()

	for {
		var mi transport.Message
		if err := sock.Recv(&mi); err != nil {
			return
		}

		// add to wait group
		s.wg.Add(1)
		pi := parse(mi)
		var fn = func(pi protocol) (po protocol) {
			defer func() {
				if r := recover(); r != nil {
					log.Error(r, string(debug.Stack()))
					po.header.Error = "Internal Server Error"
				}
			}()
			po = protocol{
				header: header{
					ServiceName: pi.header.ServiceName,
					FuncName:    pi.header.FuncName,
					Codec:       pi.header.Codec,
				},
			}
			ctx := metadata.NewContext(context.Background(), pi.header.Metadata)
			inv, ok := s.invokers[pi.header.ServiceName]
			if !ok {
				s := "not find service " + pi.header.ServiceName
				po.header.Error = s
				return po
			}
			f, err := inv.Function(pi.header.FuncName)
			if err != nil {
				po.header.Error = err.Error()
				return po
			}

			c, ok := s.codecs[pi.header.Codec]
			if !ok {
				s := "not find codec " + pi.header.Codec
				po.header.Error = s
				return po
			}

			in := f.In()
			reqVal := util.InitPointer(in[1])
			rspVal := util.InitPointer(in[2])

			// decode
			if err := c.Unmarshal(mi.Body, reqVal.Addr().Interface()); err != nil {
				po.header.Error = err.Error()
				return po
			}

			mi := invoker.NewMessage()
			mi.SetFuncName(pi.header.FuncName)
			mi.SetParameters([]interface{}{ctx, reqVal.Interface(), rspVal.Interface()})
			mo, err := inv.Invoke(ctx, mi)
			if err != nil {
				po.header.Error = err.Error()
				return po
			}

			if mo.Parameters()[0] != nil {
				po.header.Error = mo.Parameters()[0].(error).Error()
			}
			// encode
			rspBytes, err := c.Marshal(rspVal.Interface())
			if err != nil {
				po.header.Error = err.Error()
				return po
			}
			po.Body = rspBytes
			return po
		}
		mo := format(fn(pi))
		if err := sock.Send(&mo); err != nil {
			log.Error(err)
			return
		}

		s.wg.Done()
	}
}

func (s *Server) Init(opts ...server.Option) error {
	s.Lock()
	defer s.Unlock()
	for _, opt := range opts {
		opt(&s.opts)
	}

	// auto get ip address and port, if address is empty
	host, port, err := net.SplitHostPort(s.opts.Address)
	if err == nil {
		addr, err := util.Address(host)
		if err == nil {
			host = addr
		}
	}
	s.opts.Address = net.JoinHostPort(host, port)

	for _, c := range s.opts.Codecs {
		s.codecs[c.String()] = c
	}

	if err := s.opts.Transport.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Options() server.Options {
	s.RLock()
	defer s.RUnlock()
	return s.opts
}

// Register invoker to server
func (s *Server) Register(inv invoker.Invoker) error {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.invokers[inv.Name()]; ok {
		s := "duplicate service name " + inv.Name()
		return errors.New(s)
	}
	s.invokers[inv.Name()] = inv
	return nil
}

// Start server when counter == 1 and started == false
func (s *Server) Start() error {
	s.Lock()
	defer s.Unlock()
	if s.counter = s.counter + 1; s.counter > 1 {
		return nil
	}
	if s.started {
		return nil
	}
	s.started = true

	var ts transport.Listener
	var err error
	if s.opts.TLSEnable {
		var config *tls.Config
		// use not give key and cert, will generate
		if len(s.opts.TLSKeyFile) == 0 || len(s.opts.TLSCertFile) == 0 {
			config, err = util.GetTLSConfigByAddr(s.opts.Address)
		} else {
			config, err = util.GetTLSConfig(s.opts.TLSKeyFile, s.opts.TLSCertFile)
		}
		if err != nil {
			return err
		}
		ts, err = s.opts.Transport.Listen(
			s.opts.Address,
			transport.WithListenTLSConfig(config),
		)
	} else {
		ts, err = s.opts.Transport.Listen(s.opts.Address)
	}
	if err != nil {
		return err
	}

	log.Infof("Listening on %s", ts.Addr())
	go ts.Accept(s.accept)

	go func() {
		// wait for exit
		ch := <-s.exit

		// wait for requests to finish
		if s.opts.Wait {
			s.wg.Wait()
		}

		// close transport listener
		ch <- ts.Close()
	}()

	return nil
}

// Stop server when started == true, counter == 0 and stopped == false
func (s *Server) Stop() error {
	s.Lock()
	defer s.Unlock()
	if !s.started {
		return nil
	}
	if s.counter > 0 {
		s.counter = s.counter - 1
	}
	if s.counter > 0 {
		return nil
	}
	if s.stopped {
		return nil
	}
	s.stopped = true

	ch := make(chan error)
	s.exit <- ch
	return <-ch
}

func (s *Server) String() string {
	return "rpc"
}
