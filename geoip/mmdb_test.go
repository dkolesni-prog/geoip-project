// geoip/mmdb_test.go
package geoip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateMMDB(t *testing.T) {
	results := map[string]*Response{
		"8.8.8.8": {
			CountryIsoCode: "US",
			CountryName:    "United States",
		},
		"217.150.32.5": {
			CountryIsoCode: "RU",
			CountryName:    "Russia",
		},
	}

	tmpFile := filepath.Join(os.TempDir(), "test.mmdb")
	defer os.Remove(tmpFile)

	err := GenerateMMDB(results, tmpFile)
	if err != nil {
		t.Fatalf("GenerateMMDB failed: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("MMDB file does not exist: %v", err)
	}
	if info.Size() == 0 {
		t.Errorf("MMDB file is empty")
	}
}
