package etcd

import (
	"context"
	"strings"
	"sync"

	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/registry"

	etcd "github.com/coreos/etcd/client"
)

type etcdWatcher struct {
	ctx  context.Context
	once *sync.Once
	stop chan bool
	w    etcd.Watcher
}

func newEtcdWatcher(r *etcdRegistry, opts ...registry.WatchOption) (registry.Watcher, error) {
	var wo registry.WatchOptions
	for _, o := range opts {
		o(&wo)
	}

	var once sync.Once
	ctx, cancel := context.WithCancel(context.Background())
	stop := make(chan bool, 1)

	go func() {
		<-stop
		cancel()
	}()

	// watch everything by default
	watchPath := prefix
	// watch a service
	if len(wo.Service) > 0 {
		watchPath = servicePath(wo.Service)
	}

	return &etcdWatcher{
		ctx:  ctx,
		w:    r.client.Watcher(watchPath, &etcd.WatcherOptions{AfterIndex: 0, Recursive: true}),
		once: &once,
		stop: stop,
	}, nil
}

func (ew *etcdWatcher) Next() (registry.Result, error) {
	var result registry.Result
	for {
		rsp, err := ew.w.Next(ew.ctx)
		if err != nil && ew.ctx.Err() != nil {
			return result, err
		}

		var service registry.Service

		key := rsp.Node.Key
		splits := strings.Split(key, "/")
		if len(splits) != 5 {
			log.Warn(key + " invalid service")
			continue
		}
		service, err = registry.Parse(splits[4])
		if err != nil {
			log.Errorf("key:%v,err:%v", rsp.Node.Key, err)
			continue
		}

		switch rsp.Action {
		case "set", "delete", "create", "update":
			if rsp.Action == "set" {
				rsp.Action = "update"
			}
			return registry.Result{
				Action:  rsp.Action,
				Service: service,
			}, nil
		default:
			continue
		}

	}
}

func (ew *etcdWatcher) Stop() {
	ew.once.Do(func() {
		ew.stop <- true
	})
}
