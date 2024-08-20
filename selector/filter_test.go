package selector

import (
	"reflect"
	"testing"

	"github.com/haormj/dodo/client"
	"github.com/haormj/dodo/client/rpc"
	"github.com/haormj/dodo/registry"
)

func TestFilterLabel(t *testing.T) {
	type args struct {
		key      string
		val      string
		services []registry.Service
	}
	tests := []struct {
		name string
		args args
		want []registry.Service
	}{
		{
			name: "FilterLabel",
			args: args{
				key: "nodeID",
				val: "1",
				services: []registry.Service{
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "2",
						},
					},
				},
			},
			want: []registry.Service{
				{
					Protocol:  "rpc",
					Address:   "127.0.0.1:17312",
					Name:      "Hello",
					Version:   "0.1.0",
					Funcs:     []string{"SayHello"},
					Codecs:    []string{"json"},
					Transport: "grpc",
					Side:      "provider",
					TLS:       true,
					Timestamp: 1543311057,
					Labels: map[string]string{
						"nodeID": "1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterLabel(tt.args.key, tt.args.val)(tt.args.services)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterLabel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterVersion(t *testing.T) {
	type args struct {
		version  string
		services []registry.Service
	}
	tests := []struct {
		name string
		args args
		want []registry.Service
	}{
		{
			name: "FilterVersion",
			args: args{
				version: "0.1.0",
				services: []registry.Service{
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "1.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
				},
			},
			want: []registry.Service{
				{
					Protocol:  "rpc",
					Address:   "127.0.0.1:17312",
					Name:      "Hello",
					Version:   "0.1.0",
					Funcs:     []string{"SayHello"},
					Codecs:    []string{"json"},
					Transport: "grpc",
					Side:      "provider",
					TLS:       true,
					Timestamp: 1543311057,
					Labels: map[string]string{
						"nodeID": "1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterVersion(tt.args.version)(tt.args.services)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterTLS(t *testing.T) {
	type args struct {
		services []registry.Service
	}
	tests := []struct {
		name string
		args args
		want []registry.Service
	}{
		{
			name: "FilterTLS",
			args: args{
				services: []registry.Service{
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       false,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
				},
			},
			want: []registry.Service{
				{
					Protocol:  "rpc",
					Address:   "127.0.0.1:17312",
					Name:      "Hello",
					Version:   "0.1.0",
					Funcs:     []string{"SayHello"},
					Codecs:    []string{"json"},
					Transport: "grpc",
					Side:      "provider",
					TLS:       true,
					Timestamp: 1543311057,
					Labels: map[string]string{
						"nodeID": "1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterTLS()(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterTLS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterClient(t *testing.T) {
	type args struct {
		clis     []client.Client
		services []registry.Service
	}
	tests := []struct {
		name string
		args args
		want []registry.Service
	}{
		{
			name: "FilterClient",
			args: args{
				clis: []client.Client{rpc.NewClient()},
				services: []registry.Service{
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "grpc",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "tcp",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
					{
						Protocol:  "rest",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"json"},
						Transport: "tcp",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
					{
						Protocol:  "rpc",
						Address:   "127.0.0.1:17312",
						Name:      "Hello",
						Version:   "0.1.0",
						Funcs:     []string{"SayHello"},
						Codecs:    []string{"gob"},
						Transport: "tcp",
						Side:      "provider",
						TLS:       true,
						Timestamp: 1543311057,
						Labels: map[string]string{
							"nodeID": "1",
						},
					},
				},
			},
			want: []registry.Service{
				{
					Protocol:  "rpc",
					Address:   "127.0.0.1:17312",
					Name:      "Hello",
					Version:   "0.1.0",
					Funcs:     []string{"SayHello"},
					Codecs:    []string{"json"},
					Transport: "grpc",
					Side:      "provider",
					TLS:       true,
					Timestamp: 1543311057,
					Labels: map[string]string{
						"nodeID": "1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterClient(tt.args.clis)(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FilterClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
