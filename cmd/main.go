// cmd/main.go
package main

import (
	"encoding/csv"
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

		file, _, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			fileIps, err := parseFile(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			ips = append(ips, fileIps...)
		} else {
			manualIps := r.FormValue("ips")
			ips = append(ips, parseManualIPs(manualIps)...)
		}

		if len(ips) == 0 || len(ips) > 100 {
			http.Error(w, "Please provide between 1 to 100 IP addresses", http.StatusBadRequest)
			return
		}

		results, err := service.GetBatch(ips)
		if err != nil {
			log.Printf("errors during batch lookup: %v", err)
		}

		w.Header().Set("Content-Type", "text/csv")
		wantDownload := r.FormValue("download") == "1"
		if wantDownload {
			w.Header().Set("Content-Disposition", "attachment; filename=geoip_results.csv")
		}

		w.Header().Set("Content-Type", "text/csv")
		writer := csv.NewWriter(w)
		writer.Write([]string{"IP", "CountryIsoCode"})

		for ip, res := range results {
			country := ""
			if res != nil {
				country = res.CountryIsoCode
			}
			writer.Write([]string{ip, country})
		}
		writer.Flush()
	})

	fmt.Println("Listening at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func parseManualIPs(input string) []string {
	fields := strings.Split(input, ",")
	var ips []string
	for _, field := range fields {
		trimmed := strings.TrimSpace(field)
		if trimmed != "" {
			ips = append(ips, trimmed)
		}
	}
	return ips
}

func parseFile(file io.Reader) ([]string, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	var ips []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			ips = append(ips, trimmed)
		}
	}
	return ips, nil
}
