package geoip

import (
	"encoding/json"
	"io"
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

func ParseTxtFile(r io.Reader) ([]string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	raw := string(content)
	raw = strings.ReplaceAll(raw, "\n", ",")
	parts := strings.Split(raw, ",")
	var ips []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			ips = append(ips, p)
		}
	}
	return ips, nil
}

func ParseJSONFile(r io.Reader) ([]string, error) {
	var list []string
	err := json.NewDecoder(r).Decode(&list)
	return list, err
}
