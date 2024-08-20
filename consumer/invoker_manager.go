package consumer

import (
	"sync"

	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/invoker/dummy"
)

type InvokerManager struct {
	sync.RWMutex
	invokers map[string]invoker.Invoker
}

func NewInvokerManager() *InvokerManager {
	i := &InvokerManager{
		invokers: make(map[string]invoker.Invoker),
	}
	return i
}

func (i *InvokerManager) Get(n string) invoker.Invoker {
	i.RLock()
	inv, ok := i.invokers[n]
	i.RUnlock()
	if ok {
		return inv
	}

	i.Lock()
	inv = dummy.NewInvoker(invoker.Name(n))
	inv.Init()
	i.invokers[inv.Name()] = inv
	i.Unlock()

	return inv
}
