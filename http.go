package assist

import (
	"net"
	"net/http"
)

func HttpRequestIP(r *http.Request) net.IP {
	ip := r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return net.ParseIP(ip)
}
