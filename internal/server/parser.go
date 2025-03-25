package server

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func parseUploadedFile(r io.Reader, filename string) ([]string, error) {
	switch {
	case strings.HasSuffix(filename, ".json"):
		return parseJSON(r)
	case strings.HasSuffix(filename, ".txt"):
		return parseTxt(r)
	case strings.HasSuffix(filename, ".csv"):
		return nil, fmt.Errorf("CSV uploads are not supported. Use .txt or .json")
	default:
		return nil, fmt.Errorf("unsupported file type: %s", filename)
	}
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

func parseTxt(r io.Reader) ([]string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	raw := strings.ReplaceAll(string(data), "\n", ",")
	return parseManualIPs(raw), nil
}

func parseJSON(r io.Reader) ([]string, error) {
	var list []string
	if err := json.NewDecoder(r).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}
