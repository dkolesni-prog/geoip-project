package geoip

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"go.uber.org/multierr"
)

const defaultGeoIPURL = "https://geoip.noc.gov.ru/api/geoip"

type Service struct {
	client *http.Client
	url    string
}

func NewService(client *http.Client, url string) *Service {
	if url == "" {
		url = defaultGeoIPURL
	}
	return &Service{client: client, url: url}
}

func (s *Service) Get(ip string) (*Response, error) {
	ip = NormalizeIP(ip)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?ip=%s", s.url, ip), nil)
	if err != nil {
		return nil, fmt.Errorf("geoip: creating request failed: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("geoip: HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("geoip: reading response body failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geoip: unexpected status %d: %s", resp.StatusCode, body)
	}

	// Check if the body is actually an error message (API inconsistency)
	if string(body) == `{"errorCode":"404","errorMessage":"Не найдено"}` {
		return nil, nil
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("geoip: decoding failed: %w", err)
	}

	return &response, nil
}

func (s *Service) GetBatch(ips []string) (map[string]*Response, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results = make(map[string]*Response, len(ips))
		errs    error
	)

	for _, ip := range ips {
		wg.Add(1)
		ip := ip
		go func() {
			defer wg.Done()

			resp, err := s.Get(ip)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errs = multierr.Append(errs, fmt.Errorf("geoip: IP %s: %w", ip, err))
				results[ip] = nil
				return
			}

			results[ip] = resp // can be nil if IP was not found (404)
		}()
	}

	wg.Wait()
	return results, errs
}
