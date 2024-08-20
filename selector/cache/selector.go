package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/haormj/dodo/log"
	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/selector"

	"gopkg.in/yaml.v2"
)

type Selector struct {
	sync.RWMutex
	opts     selector.Options
	services map[string][]registry.Service
	exit     chan struct{}
}

func NewSelector(opts ...selector.Option) selector.Selector {
	options := selector.Options{
		Strategy: selector.Random,
	}

	for _, o := range opts {
		o(&options)
	}

	s := &Selector{
		opts:     options,
		services: make(map[string][]registry.Service),
	}

	return s
}

func (s *Selector) getCacheDir() string {
	if s.opts.Context != nil {
		if v := s.opts.Context.Value(cacheDirKey{}); v != nil {
			return v.(string)
		}
	}
	return "./.dodo/consumer/selector"
}

func (s *Selector) getConfigDir() string {
	if s.opts.Context != nil {
		if v := s.opts.Context.Value(configDirKey{}); v != nil {
			return v.(string)
		}
	}
	return "../config/selector"
}

func (s *Selector) readFromConfig() error {
	return s.readFromLocal(s.getConfigDir())
}

func (s *Selector) readFromCache() error {
	return s.readFromLocal(s.getCacheDir())
}

func (s *Selector) readFromLocal(d string) error {
	fis, err := ioutil.ReadDir(d)
	if err != nil {
		return err
	}
	for _, fi := range fis {
		// ignore dir
		if fi.IsDir() {
			continue
		}
		// only support yaml/yml config now
		if !strings.HasSuffix(fi.Name(), "yaml") && !strings.HasSuffix(fi.Name(), "yml") {
			log.Warn("ext must be yaml/yml, current is " + fi.Name())
			continue
		}
		p := filepath.Join(d, fi.Name())
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		svcs := make([]registry.Service, 0)
		if err := yaml.Unmarshal(b, &svcs); err != nil {
			return err
		}
		s.addService(svcs)
	}
	return nil
}

func (s *Selector) readFromRegistry() error {
	svcs, err := s.opts.Registry.ListServices()
	if err != nil {
		return err
	}
	m := make(map[string][]registry.Service)
	for _, svc := range svcs {
		services := m[svc.Name]
		services = append(services, svc)
		m[svc.Name] = services
	}
	s.Lock()
	s.services = m
	s.Unlock()
	return nil
}

func (s *Selector) readFromRemote() error {
	if err := s.readFromRegistry(); err != nil {
		log.Error(err)
		if err = s.readFromCache(); err != nil {
			return err
		}
	}
	s.writeToCache()

	// pull
	go func() {
		// time can use option
		t := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-t.C:
				log.Debug("hello")
				if err := s.readFromRegistry(); err != nil {
					log.Error(err)
					continue
				}
				s.writeToCache()
			case <-s.exit:
				return
			}
		}
	}()

	// watch
	go s.watch()

	return nil
}

func (s *Selector) watch() {
	defer func() {
		if r := recover(); r != nil {
			log.Error(r, string(debug.Stack()))
		}
		// can not stop watch :)
		go s.watch()
	}()
	w, err := s.opts.Registry.Watch()
	if err != nil {
		log.Error(err)
		return
	}
	for {
		select {
		case <-s.exit:
			w.Stop()
			return
		default:
			result, err := w.Next()
			if err != nil {
				log.Error(err)
				continue
			}

			// TODO
			// there we should use Action

			// get service from registry
			services, err := s.opts.Registry.GetService(result.Service.Name)
			if err != nil {
				log.Error(err)
				continue
			}
			s.Lock()
			s.services[result.Service.Name] = services
			s.Unlock()

			s.writeToCache()
		}
	}
}

func (s *Selector) writeToCache() error {
	s.RLock()
	for name, svcs := range s.services {
		b, err := yaml.Marshal(svcs)
		if err != nil {
			log.Error(err)
			continue
		}
		if err := os.MkdirAll(s.getCacheDir(), 0755); err != nil {
			log.Error(err)
			continue
		}
		p := filepath.Join(s.getCacheDir(), name+".yaml")
		f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Error(err)
			continue
		}
		if _, err := f.Write(b); err != nil {
			log.Error(err)
			f.Close()
			continue
		}
		f.Close()
	}
	s.RUnlock()
	return nil
}

func (s *Selector) read() error {
	if err := s.readFromConfig(); err != nil {
		// log.Info(err)
		if err := s.readFromRemote(); err != nil {
			return err
		}
	}
	return nil
}

func (s Selector) addService(services []registry.Service) {
	if len(services) == 0 {
		// TODO
		return
	}
	n := services[0].Name
	s.Lock()
	s.services[n] = services
	s.Unlock()
}

func (s *Selector) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(&s.opts)
	}

	if err := s.read(); err != nil {
		return err
	}
	return nil
}

func (s *Selector) Options() selector.Options {
	return s.opts
}

func (s *Selector) Select(service string, opts ...selector.SelectOption) (registry.Service, error) {
	var svc registry.Service

	sopts := selector.SelectOptions{
		Strategy: s.opts.Strategy,
	}

	for _, o := range opts {
		o(&sopts)
	}
	s.RLock()
	services, ok := s.services[service]
	s.RUnlock()
	if !ok || len(services) == 0 {
		return svc, selector.ErrNoneAvailable
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	if len(services) == 0 {
		return svc, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services)
}

func (s *Selector) Mark(service string, err error) {

}

func (s *Selector) Reset(service string) {

}

func (s *Selector) Close() error {
	close(s.exit)
	return nil
}

func (s *Selector) String() string {
	return "cache"
}
