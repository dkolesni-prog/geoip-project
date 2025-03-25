package geoip

import (
	"net"
	"strings"
)

// NormalizeIP ensures consistent format, compressing IPv6, etc.
// fox example 2001:0db8:85a3:0000:0000:8a2e:0370:7334  --> 2001:db8:85a3::8a2e:370:7334
func NormalizeIP(ip string) string {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return ip
	}
	return parsed.String()
}
