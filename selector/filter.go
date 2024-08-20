package selector

import (
	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/registry"
	"github.com/haormj/dodo/util"
)

// FilterLabel is a label based Select Filter which will
// only return services with the label specified.
func FilterLabel(key, val string) Filter {
	return func(old []registry.Service) []registry.Service {
		var services []registry.Service

		for _, service := range old {
			if service.Labels == nil {
				continue
			}
			if service.Labels[key] == val {
				services = append(services, service)
			}
		}

		return services
	}
}

// FilterVersion is a version based Select Filter which will
// only return services with the version specified.
func FilterVersion(version string) Filter {
	return func(old []registry.Service) []registry.Service {
		var services []registry.Service

		for _, service := range old {
			if service.Version == version {
				services = append(services, service)
			}
		}

		return services
	}
}

func FilterTLS() Filter {
	return func(old []registry.Service) []registry.Service {
		var services []registry.Service

		for _, service := range old {
			if service.TLS == true {
				services = append(services, service)
			}
		}

		return services
	}
}

func FilterClient(clis []client.Client) Filter {
	return func(old []registry.Service) []registry.Service {
		var services []registry.Service

		for _, cli := range clis {
			t := make([]string, 0)
			for _, cdc := range cli.Options().Codecs {
				t = append(t, cdc.String())
			}
			var transport, protocol string
			protocol = cli.String()
			opts := cli.Options()
			if opts.Transport != nil {
				transport = opts.Transport.String()
			}

			for _, service := range old {
				if service.Protocol != protocol {
					continue
				}
				if service.Transport != transport {
					continue
				}
				if len(util.ArrayIntersectString(service.Codecs, t)) == 0 {
					continue
				}
				if !registry.Contains(services, service) {
					services = append(services, service)
				}
			}
		}

		return services
	}
}
