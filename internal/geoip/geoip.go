// Package geoip provides IP-to-location lookup using free GeoIP APIs.
package geoip

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// Location holds the geo-location result for an IP address.
type Location struct {
	Country string
	City    string
}

// Service performs GeoIP lookups with in-memory caching.
type Service struct {
	client *http.Client
	cache  sync.Map // map[string]Location
}

// NewService creates a new GeoIP lookup service.
func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Lookup returns the country and city for the given IP address.
// Results are cached in memory to avoid repeated API calls.
// Returns empty Location (no error) for private/localhost IPs.
func (s *Service) Lookup(ip string) Location {
	// Skip private / loopback IPs
	if isPrivateIP(ip) {
		return Location{Country: "Local", City: "Localhost"}
	}

	// Check cache first
	if cached, ok := s.cache.Load(ip); ok {
		return cached.(Location)
	}

	// Try ip-api.com (free, no API key, 45 req/min)
	loc := s.lookupIPAPI(ip)

	// Cache the result regardless of success
	s.cache.Store(ip, loc)
	return loc
}

// LookupAsync performs a GeoIP lookup in a goroutine and calls the callback with the result.
// Use this to avoid blocking the HTTP response.
func (s *Service) LookupAsync(ip string, callback func(Location)) {
	go func() {
		loc := s.Lookup(ip)
		if callback != nil {
			callback(loc)
		}
	}()
}

// ── ip-api.com ────────────────────────────────────────────────

type ipAPIResponse struct {
	Status  string `json:"status"`
	Country string `json:"country"`
	City    string `json:"city"`
}

func (s *Service) lookupIPAPI(ip string) Location {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,city", ip)

	resp, err := s.client.Get(url)
	if err != nil {
		log.Printf("⚠️ GeoIP lookup failed for %s: %v", ip, err)
		return Location{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ GeoIP lookup returned %d for %s", resp.StatusCode, ip)
		return Location{}
	}

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("⚠️ GeoIP decode failed for %s: %v", ip, err)
		return Location{}
	}

	if result.Status != "success" {
		return Location{}
	}

	return Location{
		Country: result.Country,
		City:    result.City,
	}
}

// ── Helpers ───────────────────────────────────────────────────

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return true // unparseable → treat as private
	}

	// Loopback
	if ip.IsLoopback() {
		return true
	}

	// Private ranges
	privateRanges := []struct {
		start net.IP
		end   net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
	}

	for _, r := range privateRanges {
		if bytesInRange(ip.To4(), r.start.To4(), r.end.To4()) {
			return true
		}
	}
	return false
}

func bytesInRange(ip, start, end net.IP) bool {
	if ip == nil || start == nil || end == nil {
		return false
	}
	for i := 0; i < len(ip); i++ {
		if ip[i] < start[i] {
			return false
		}
		if ip[i] > end[i] {
			return false
		}
	}
	return true
}
