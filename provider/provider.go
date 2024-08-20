package provider

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/haormj/dodo/invoker"
	"github.com/haormj/dodo/invoker/function"
	"github.com/haormj/dodo/invoker/receiver"
	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/util"
)

type Provider struct {
	name    string
	service interface{}
	opts    Options
	inv     invoker.Invoker
	codecs  []string
	exit    chan struct{}
}

func NewProvider(service interface{}, opts ...Option) *Provider {
	options := newOptions(opts...)
	p := &Provider{
		service: service,
		opts:    options,
		exit:    make(chan struct{}),
	}
	return p
}

func (p *Provider) getServices() []registry.Service {
	services := make([]registry.Service, 0)
	funcNames := make([]string, 0)
	for _, f := range p.inv.Functions() {
		funcNames = append(funcNames, f.FuncName())
	}

	service := registry.Service{
		Name:      p.name,
		Version:   p.opts.Version,
		Funcs:     funcNames,
		Labels:    p.opts.Labels,
		Timestamp: time.Now().Unix(),
		Side:      "provider",
	}
	for _, s := range p.opts.Servers {
		svc := service
		sopts := s.Options()
		svc.TLS = sopts.TLSEnable
		svc.Address = sopts.Address
		if sopts.Transport != nil {
			svc.Transport = sopts.Transport.String()
		}
		svc.Protocol = s.String()
		codecs := make([]string, 0)
		for _, c := range sopts.Codecs {
			if !util.ArrayContainsString(codecs, c.String()) {
				codecs = append(codecs, c.String())
			}
		}
		svc.Codecs = codecs
		services = append(services, svc)
	}
	return services
}

func (p *Provider) start(services []registry.Service) error {
	for _, s := range p.opts.Servers {
		if err := s.Start(); err != nil {
			return err
		}
		if err := s.Register(p.inv); err != nil {
			return err
		}
	}
	if err := p.register(services); err != nil {
		return err
	}
	return nil
}

func (p *Provider) stop(services []registry.Service) error {
	if err := p.deregister(services); err != nil {
		return err
	}
	for _, s := range p.opts.Servers {
		if err := s.Stop(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Provider) register(services []registry.Service) error {
	for _, r := range p.opts.Registries {
		for _, s := range services {
			if err := r.Register(s, registry.RegisterTTL(p.opts.RegisterTTL)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Provider) registerLoop(services []registry.Service) {
	if p.opts.RegisterInterval <= time.Duration(0) {
		return
	}
	go func() {
		t := time.NewTicker(p.opts.RegisterInterval)
		for {
			select {
			case <-t.C:
				for _, r := range p.opts.Registries {
					for _, s := range services {
						err := r.Register(s, registry.RegisterTTL(p.opts.RegisterTTL))
						if err != nil {
							log.Error(err)
						}
					}
				}
			case <-p.exit:
				t.Stop()
				return
			}
		}
	}()
}

func (p *Provider) deregister(services []registry.Service) error {
	for _, r := range p.opts.Registries {
		for _, s := range services {
			if err := r.Deregister(s); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Provider) Init(opts ...Option) error {
	for _, o := range opts {
		o(&p.opts)
	}
	if reflect.TypeOf(p.service).Kind() == reflect.Func {
		p.inv = function.NewInvoker(p.service)
	} else {
		p.inv = receiver.NewInvoker(p.service)
	}
	if err := p.inv.Init(invoker.Name(p.opts.Name)); err != nil {
		return err
	}
	p.name = p.inv.Name()

	for _, s := range p.opts.Servers {
		if err := s.Init(); err != nil {
			log.Error(err)
			return err
		}
	}

	for _, r := range p.opts.Registries {
		if err := r.Init(); err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
}

func (p *Provider) Run() error {
	services := p.getServices()

	if err := p.start(services); err != nil {
		log.Error(err)
		return err
	}

	p.registerLoop(services)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	var sig os.Signal
	select {
	case sig = <-ch:
		log.Infof("receive signal %s", sig.String())
		// stop to receive signal
		signal.Stop(ch)
		// close exit channel
		close(p.exit)
	}

	if err := p.stop(services); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
