package rpc

import (
	"strings"

	"github.com/haormj/dodo/transport"
)

type protocol struct {
	header header
	Body   []byte
}

type header struct {
	ServiceName string
	FuncName    string
	Codec       string
	Error       string
	Metadata    map[string]string
}

func parse(m transport.Message) protocol {
	p := protocol{
		header: header{
			Metadata: make(map[string]string),
		},
		Body: m.Body,
	}
	for k, v := range m.Header {
		switch k {
		case "ServiceName":
			p.header.ServiceName = v
		case "FuncName":
			p.header.FuncName = v
		case "Codec":
			p.header.Codec = v
		case "Error":
			p.header.Error = v
		default:
			k = strings.TrimPrefix(k, "Meta-")
			p.header.Metadata[k] = v
		}
	}
	return p
}

func format(p protocol) transport.Message {
	m := transport.Message{
		Header: map[string]string{
			"ServiceName": p.header.ServiceName,
			"FuncName":    p.header.FuncName,
			"Codec":       p.header.Codec,
			"Error":       p.header.Error,
		},
		Body: p.Body,
	}

	for k, v := range p.header.Metadata {
		k = "Meta-" + k
		m.Header[k] = v
	}

	return m
}
