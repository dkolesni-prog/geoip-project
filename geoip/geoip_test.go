// geoip/geoip_test.go
package geoip

import (
	"testing"
)

func TestService_Get(t *testing.T) {
	client := NewHTTPClient()
	service := NewService(client, "")

	testIPs := []string{
		"8.8.8.8",
		"1.1.1.1", // May return 404
		"217.150.32.5",
		"2001:4860:4860::8888",
	}

	for _, ip := range testIPs {
		t.Run(ip, func(t *testing.T) {
			res, err := service.Get(ip)
			if err != nil {
				t.Errorf("Failed to fetch %s: %v", ip, err)
				return
			}
			if res == nil {
				if ip == "1.1.1.1" {
					t.Logf("IP %s not found (expected nil)", ip)
					return
				}
				t.Errorf("Expected result for IP %s, got nil", ip)
				return
			}
			if res.CountryIsoCode == "" {
				t.Errorf("Empty CountryIsoCode for IP %s", ip)
			}
		})
	}
}

func TestService_GetBatch(t *testing.T) {
	client := NewHTTPClient()
	service := NewService(client, "")

	ips := []string{
		"8.8.8.8",
		"1.1.1.1", // Might return 404
		"217.150.32.5",
		"invalid-ip", // Should fail gracefully
	}

	results, err := service.GetBatch(ips)
	if err != nil {
		t.Logf("GetBatch returned non-fatal errors: %v", err)
	}

	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			res, ok := results[ip]
			if !ok {
				t.Errorf("Missing result for IP: %s", ip)
				return
			}

			if res == nil {
				// ❗ Important: we treat these IPs as acceptable nils
				if ip == "1.1.1.1" || ip == "invalid-ip" {
					t.Logf("IP %s not found (expected nil)", ip)
					return
				}
				t.Errorf("Expected result for IP %s, got nil", ip)
				return
			}

			// Now it's safe to check CountryIsoCode
			if res.CountryIsoCode == "" {
				t.Errorf("No country code returned for valid IP: %s", ip)
			}
		})
	}
}
