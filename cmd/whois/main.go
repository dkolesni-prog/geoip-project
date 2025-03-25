// cmd/main.go
package main

import (
	"github.com/dkolesni-prog/whois/internal/server"
	"log"
	"net/http"
	"time"

	"github.com/dkolesni-prog/whois/geoip"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func run() error {
	client := geoip.NewHTTPClient()
	service := geoip.NewService(client, "")

	r := server.NewRouter(service)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Listening at :8080")
	return srv.ListenAndServe()
}
