package server

import (
	"github.com/dkolesni-prog/whois/geoip"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func NewRouter(service *geoip.Service) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", serveFrontendIndex)
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("frontend/build"))))
	r.Post("/check_ips", CheckIPsHandler(service))

	return r
}
