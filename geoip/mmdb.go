package geoip

import (
	"net"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func GenerateMMDB(results map[string]*Response, filepath string) error {
	writer, err := mmdbwriter.New(mmdbwriter.Options{
		DatabaseType: "GeoIP2-Country",
		RecordSize:   24,
	})
	if err != nil {
		return err
	}

	for ipStr, res := range results {
		if res == nil {
			continue
		}
		ip := net.ParseIP(ipStr)
		if ip == nil {
			continue
		}

		// Convert net.IP to *net.IPNet with /32 for IPv4 or /128 for IPv6
		var ipNet *net.IPNet
		if ip.To4() != nil {
			_, ipNet, _ = net.ParseCIDR(ip.String() + "/32")
		} else {
			_, ipNet, _ = net.ParseCIDR(ip.String() + "/128")
		}

		entry := mmdbtype.Map{
			"country": mmdbtype.Map{
				"iso_code": mmdbtype.String(res.CountryIsoCode),
				"names": mmdbtype.Map{
					"en": mmdbtype.String(res.CountryName),
				},
			},
		}

		err := writer.Insert(ipNet, entry)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = writer.WriteTo(file)
	return err
}
