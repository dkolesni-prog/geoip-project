// cmd/main.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	registerRoutes(r, service)

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
