package geoip

import (
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

func GenerateMMDB(results map[string]*Response, outPath string) error {
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

		if err := writer.Insert(ipNet, entry); err != nil {
			return err
		}
	}

	safePath := filepath.Clean(outPath)
	file, err := os.Create(safePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("couldnt close file: %s", err)
		}
	}(file)

	_, err = writer.WriteTo(file)
	return err
}
