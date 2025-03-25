// cmd/endpoints.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/dkolesni-prog/whois/geoip"
	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, service *geoip.Service) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/build/index.html")
	})

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("frontend/build"))))

	r.Post("/check_ips", func(w http.ResponseWriter, r *http.Request) {
		var ips []string

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// Handle uploaded file OR manual input
		file, header, fileErr := r.FormFile("file")
		if fileErr == nil {
			defer file.Close()
			switch {
			case strings.HasSuffix(header.Filename, ".json"):
				ips, err = parseJSONFile(file)
			case strings.HasSuffix(header.Filename, ".txt"):
				ips, err = parseTxtFile(file)
			case strings.HasSuffix(header.Filename, ".csv"):
				http.Error(w, "CSV uploads are not supported. Use .txt or .json", http.StatusBadRequest)
				return
			default:
				http.Error(w, "Unsupported file type. Upload .txt or .json", http.StatusBadRequest)
				return
			}
			if err != nil {
				http.Error(w, "Failed to parse file: "+err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			ips = parseManualIPs(r.FormValue("ips"))
		}

		if len(ips) == 0 || len(ips) > 100 {
			http.Error(w, "Please provide between 1 to 100 IP addresses", http.StatusBadRequest)
			return
		}

		results, err := service.GetBatch(ips)
		if err != nil {
			fmt.Printf("batch lookup error: %v\n", err)
		}

		export := r.FormValue("export")

		switch export {
		case "mmdb":
			w.Header().Set("Content-Disposition", "attachment; filename=geoip_results.mmdb")
			w.Header().Set("Content-Type", "application/octet-stream")

			tmpFile, err := os.CreateTemp("", "geoip_*.mmdb")
			if err != nil {
				http.Error(w, "could not create mmdb temp file", http.StatusInternalServerError)
				return
			}
			defer os.Remove(tmpFile.Name())
			defer tmpFile.Close()

			if err := geoip.GenerateMMDB(results, tmpFile.Name()); err != nil {
				http.Error(w, "failed to generate mmdb: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.ServeFile(w, r, tmpFile.Name())
			return

		case "csv":
			w.Header().Set("Content-Disposition", "attachment; filename=geoip_results.csv")
			fallthrough
		default:
			w.Header().Set("Content-Type", "text/csv")
			writer := csv.NewWriter(w)
			if err := writer.Write([]string{"IP", "CountryIsoCode", "CountryName"}); err != nil {
				http.Error(w, "Failed to write CSV header", http.StatusInternalServerError)
				return
			}
			for ip, res := range results {
				code, name := "", ""
				if res != nil {
					code = res.CountryIsoCode
					name = res.CountryName
				}
				if err := writer.Write([]string{ip, code, name}); err != nil {
					fmt.Printf("error writing row for %s: %v\n", ip, err)
				}
			}
			writer.Flush()
		}
	})
}

func parseManualIPs(input string) []string {
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return r == ',' || r == '\n'
	})
	var ips []string
	for _, ip := range parts {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			ips = append(ips, ip)
		}
	}
	return ips
}

func parseTxtFile(r io.Reader) ([]string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	raw := strings.ReplaceAll(string(content), "\n", ",")
	return parseManualIPs(raw), nil
}

func parseJSONFile(r io.Reader) ([]string, error) {
	var list []string
	err := json.NewDecoder(r).Decode(&list)
	return list, err
}
