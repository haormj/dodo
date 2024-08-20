package util

import (
	"net"
	"net/http"
	"strings"
)

func GetHttpClientIP(request *http.Request) string {
	xForwardedFor := request.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 && ips[0] != "" {
			rip, _, err := net.SplitHostPort(ips[0])
			if err != nil {
				rip = ips[0]
			}
			return rip
		}
	}
	if ip, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		return ip
	}
	return request.RemoteAddr
}
