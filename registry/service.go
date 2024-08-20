package registry

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// Service register info
type Service struct {
	Protocol  string
	Address   string
	Name      string
	Version   string
	Funcs     []string
	Codecs    []string
	Transport string
	Side      string
	TLS       bool
	Timestamp int64
	Labels    map[string]string
}

// Parse string to Service
func Parse(str string) (Service, error) {
	service := Service{
		Labels: make(map[string]string),
	}
	str, err := url.QueryUnescape(str)
	if err != nil {
		return service, err
	}
	u, err := url.Parse(str)
	if err != nil {
		return service, err
	}
	service.Protocol = u.Scheme
	if len(service.Protocol) == 0 {
		return service, errors.New("protocol is empty")
	}
	service.Address = u.Host
	if len(service.Address) == 0 {
		return service, errors.New("address is empty")
	}
	service.Name = strings.Trim(u.Path, "/")
	if len(service.Name) == 0 {
		return service, errors.New("service name is empty")
	}
	values, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return service, err
	}
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		switch k {
		case "version":
			service.Version = values[k][0]
		case "funcs":
			service.Funcs = strings.Split(values[k][0], ",")
		case "codecs":
			service.Codecs = strings.Split(values[k][0], ",")
		case "transport":
			service.Transport = values[k][0]
		case "side":
			service.Side = values[k][0]
		case "tls":
			b, err := strconv.ParseBool(values[k][0])
			if err != nil {
				return service, err
			}
			service.TLS = b
		case "timestamp":
			i, err := strconv.ParseInt(values[k][0], 10, 64)
			if err != nil {
				return service, err
			}
			service.Timestamp = i
		default:
			service.Labels[k] = values[k][0]
		}
	}

	if len(service.Codecs) == 0 {
		return service, errors.New("codec is empty")
	}

	if len(service.Funcs) == 0 {
		return service, errors.New("funcs is empty")
	}

	if len(service.Side) == 0 {
		return service, errors.New("side is empty")
	}

	if len(service.Version) == 0 {
		return service, errors.New("version is empty")
	}

	if service.Timestamp == 0 {
		return service, errors.New("timestamp invalid")
	}

	return service, nil
}

// Format Service to string and url encode
func Format(s Service) string {
	str := fmt.Sprintf("%s://%s/%s?version=%s&funcs=%s&codecs=%s&transport=%s&side=%s&tls=%t&timestamp=%d",
		s.Protocol,
		s.Address,
		s.Name,
		s.Version,
		strings.Join(s.Funcs, ","),
		strings.Join(s.Codecs, ","),
		s.Transport,
		s.Side,
		s.TLS,
		s.Timestamp,
	)
	var keys []string
	for k := range s.Labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		str += fmt.Sprintf("&%s=%s", k, s.Labels[k])
	}
	return url.QueryEscape(str)
}

// Equal judge s1 == s2
func Equal(s1, s2 Service) bool {
	if Format(s1) == Format(s2) {
		return true
	}
	return false
}

// Contains judge is service in services
func Contains(services []Service, service Service) bool {
	for _, s := range services {
		if Equal(s, service) {
			return true
		}
	}
	return false
}
