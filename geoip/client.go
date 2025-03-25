package geoip

import (
	"net/http"
	"time"
)

const defaultTimeout = 5 * time.Second

func NewHTTPClient() *http.Client {
	return &http.Client{
		Timeout: defaultTimeout,
	}
}
