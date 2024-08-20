package selector

import (
	"math/rand"
	"time"

	"github.com/haormj/dodo/registry"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random is a random strategy algorithm for node selection
func Random(services []registry.Service) (registry.Service, error) {
	var svc registry.Service
	if len(services) == 0 {
		return svc, ErrNoneAvailable
	}

	i := rand.Int() % len(services)
	return services[i], nil
}
