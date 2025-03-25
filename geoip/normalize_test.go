// geoip/normalize_test.go
package geoip

import "testing"

func TestNormalizeIP(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:db8:85a3::8a2e:370:7334"},
		{"8.8.8.8", "8.8.8.8"},
		{"   1.1.1.1  ", "1.1.1.1"},
		{"", ""},
		{"not-an-ip", "not-an-ip"}, // invalid input, fallback expected
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			out := NormalizeIP(tt.input)
			if out != tt.expected {
				t.Errorf("NormalizeIP(%q) = %q; want %q", tt.input, out, tt.expected)
			}
		})
	}
}
