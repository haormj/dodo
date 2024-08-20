// Package rest implement http server
package rest

import (
	"errors"
	"sync"

	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/server"
)

type Server struct {
	sync.RWMutex
	opts     server.Options
	invokers map[string]invoker.Invoker
}

func NewServer(opts ...server.Option) server.Server {
	options := newOptions(opts...)

	s := &Server{
		opts:     options,
		invokers: make(map[string]invoker.Invoker),
	}

	return s
}

func (s *Server) Init(opts ...server.Option) error {
	s.Lock()
	for _, o := range opts {
		o(&s.opts)
	}
	s.Unlock()

	return nil
}

func (s *Server) Options() server.Options {
	s.RLock()
	opts := s.opts
	s.RUnlock()
	return opts
}

func (s *Server) Register(inv invoker.Invoker) error {
	if _, ok := s.invokers[inv.Name()]; ok {
		s := "duplicate service name " + inv.Name()
		log.Error(s)
		return errors.New(s)
	}
	s.invokers[inv.Name()] = inv
	return nil
}

func (s *Server) Start() error {
	return nil
}

func (s *Server) Stop() error {
	return nil
}

func (s *Server) String() string {
	return "rest"
}
