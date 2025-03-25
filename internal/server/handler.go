package server

import (
	"encoding/csv"
	"net/http"
	_ "os"

	"github.com/dkolesni-prog/whois/geoip"
)

func CheckIPsHandler(service *geoip.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ips []string

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		file, header, fileErr := r.FormFile("file")
		if fileErr == nil {
			defer file.Close()
			var err error
			ips, err = parseUploadedFile(file, header.Filename)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
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
			http.Error(w, "Failed to get batch: "+err.Error(), http.StatusInternalServerError)
			return
		}

		switch r.FormValue("export") {
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
			_ = writer.Write([]string{"IP", "CountryIsoCode", "CountryName"})
			for ip, res := range results {
				code, name := "", ""
				if res != nil {
					code = res.CountryIsoCode
					name = res.CountryName
				}
				_ = writer.Write([]string{ip, code, name})
			}
			writer.Flush()
		}
	}
}

func serveFrontendIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/build/index.html")
}
