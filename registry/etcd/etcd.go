// Package etcd provides an etcd registry
package etcd

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/registry"

	etcd "github.com/coreos/etcd/client"
)

var (
	prefix = "/dodo"
)

type etcdRegistry struct {
	client  etcd.KeysAPI
	options registry.Options
}

func NewRegistry(opts ...registry.Option) registry.Registry {
	e := &etcdRegistry{
		options: registry.Options{},
	}
	configure(e, opts...)
	return e
}

func configure(e *etcdRegistry, opts ...registry.Option) error {
	config := etcd.Config{
		Endpoints: []string{"http://127.0.0.1:2379"},
	}

	for _, o := range opts {
		o(&e.options)
	}

	if e.options.Timeout == 0 {
		e.options.Timeout = etcd.DefaultRequestTimeout
	}

	if e.options.Secure || e.options.TLSConfig != nil {
		tlsConfig := e.options.TLSConfig
		if tlsConfig == nil {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}

		// for InsecureSkipVerify
		t := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     tlsConfig,
		}

		runtime.SetFinalizer(&t, func(tr **http.Transport) {
			(*tr).CloseIdleConnections()
		})

		config.Transport = t

		// default secure address
		config.Endpoints = []string{"https://127.0.0.1:2379"}
	}

	var cAddrs []string

	for _, addr := range e.options.Addrs {
		if len(addr) == 0 {
			continue
		}

		if e.options.Secure {
			// replace http:// with https:// if its there
			addr = strings.Replace(addr, "http://", "https://", 1)

			// has the prefix? no... ok add it
			if !strings.HasPrefix(addr, "https://") {
				addr = "https://" + addr
			}
		}

		cAddrs = append(cAddrs, addr)
	}

	// if we got addrs then we'll update
	if len(cAddrs) > 0 {
		config.Endpoints = cAddrs
	}

	c, err := etcd.New(config)
	if err != nil {
		return err
	}
	e.client = etcd.NewKeysAPI(c)
	return nil
}

func servicePath(s string) string {
	return path.Join(prefix, s, "providers")
}

func (e *etcdRegistry) Init(opts ...registry.Option) error {
	return configure(e, opts...)
}

func (e *etcdRegistry) Options() registry.Options {
	return e.options
}

func (e *etcdRegistry) Deregister(s registry.Service) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	key := path.Join(prefix, s.Name, s.Side+"s", registry.Format(s))
	if _, err := e.client.Delete(ctx, key, &etcd.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

func (e *etcdRegistry) Register(s registry.Service, opts ...registry.RegisterOption) error {
	var options registry.RegisterOptions
	for _, o := range opts {
		o(&options)
	}

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	key := path.Join(prefix, s.Name, s.Side+"s", registry.Format(s))
	_, err := e.client.Set(ctx, key, "", &etcd.SetOptions{TTL: options.TTL})
	if err != nil {
		return err
	}

	return nil
}

func (e *etcdRegistry) GetService(name string) ([]registry.Service, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, servicePath(name), &etcd.GetOptions{})
	if err != nil && !strings.HasPrefix(err.Error(), "100: Key not found") {
		return nil, err
	}

	if rsp == nil {
		return nil, registry.ErrNotFound
	}

	services := make([]registry.Service, 0)
	nodes := e.get(rsp.Node)
	for _, n := range nodes {
		splits := strings.Split(n.Key, "/")
		if len(splits) != 5 {
			log.Warn(n.Key + " invalid service")
			continue
		}
		service, err := registry.Parse(splits[4])
		if err != nil {
			log.Errorf("key:%v,err:%v", n.Key, err)
			continue
		}
		services = append(services, service)
	}

	return services, nil
}

func (e *etcdRegistry) get(n *etcd.Node) []*etcd.Node {
	if len(n.Nodes) == 0 {
		return []*etcd.Node{n}
	}
	var nodes []*etcd.Node
	for _, node := range n.Nodes {
		nodes = append(nodes, e.get(node)...)
	}
	return nodes
}

func (e *etcdRegistry) ListServices() ([]registry.Service, error) {
	services := make([]registry.Service, 0)

	ctx, cancel := context.WithTimeout(context.Background(), e.options.Timeout)
	defer cancel()

	rsp, err := e.client.Get(ctx, prefix, &etcd.GetOptions{Recursive: true, Sort: true})
	if err != nil && !strings.HasPrefix(err.Error(), "100: Key not found") {
		return nil, err
	}

	if rsp == nil {
		return services, nil
	}

	nodes := e.get(rsp.Node)
	for _, n := range nodes {
		splits := strings.Split(n.Key, "/")
		if len(splits) != 5 {
			log.Warn(n.Key + " invalid service")
			continue
		}
		service, err := registry.Parse(splits[4])
		if err != nil {
			log.Errorf("key:%v,err:%v", n.Key, err)
			continue
		}
		services = append(services, service)
	}

	return services, nil
}

func (e *etcdRegistry) Watch(opts ...registry.WatchOption) (registry.Watcher, error) {
	return newEtcdWatcher(e, opts...)
}

func (e *etcdRegistry) String() string {
	return "etcd"
}
