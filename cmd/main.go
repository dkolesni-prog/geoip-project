// cmd/main.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/dkolesni-prog/whois/geoip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	client := geoip.NewHTTPClient()
	service := geoip.NewService(client, "")

	r := chi.NewRouter()
	r.Use(middleware.Logger)

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

		// 🔍 Handle uploaded file OR manual input
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
			log.Printf("batch lookup error: %v", err)
		}

		export := r.FormValue("export")

		switch export {
		case "mmdb":
			w.Header().Set("Content-Disposition", "attachment; filename=geoip_results.mmdb")
			w.Header().Set("Content-Type", "application/octet-stream")

			tmpFile := "geoip_results.mmdb"
			if err := geoip.GenerateMMDB(results, tmpFile); err != nil {
				http.Error(w, "failed to generate mmdb: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.ServeFile(w, r, tmpFile)
			return

		case "csv":
			w.Header().Set("Content-Disposition", "attachment; filename=geoip_results.csv")
			fallthrough

		default:
			w.Header().Set("Content-Type", "text/csv")
			writer := csv.NewWriter(w)
			writer.Write([]string{"IP", "CountryIsoCode", "CountryName"})
			for ip, res := range results {
				code, name := "", ""
				if res != nil {
					code = res.CountryIsoCode
					name = res.CountryName
				}
				writer.Write([]string{ip, code, name})
			}
			writer.Flush()
		}
	})

	fmt.Println("Listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Parses comma-separated or newline-separated IPs from a textarea input
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

// Parses .txt files with newline or comma-separated IPs
func parseTxtFile(r io.Reader) ([]string, error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	raw := string(content)
	raw = strings.ReplaceAll(raw, "\n", ",")
	parts := strings.Split(raw, ",")
	var ips []string
	for _, ip := range parts {
		ip = strings.TrimSpace(ip)
		if ip != "" {
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

// Parses .json file with array of IPs
func parseJSONFile(r io.Reader) ([]string, error) {
	var list []string
	err := json.NewDecoder(r).Decode(&list)
	return list, err
}
