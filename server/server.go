package server

import (
	"github.com/haormj/dodo/invoker"
)

type Server interface {
	Init(...Option) error
	Options() Options
	Register(invoker.Invoker) error
	Start() error
	Stop() error
	String() string
}
