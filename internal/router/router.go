// router.go
// DNS-based routing for EMSG Daemon
package router

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

// LookupRoute parses 'user#domain.com' and fetches DNS TXT records from '_emsg.domain.com'.
func LookupRoute(address string) (string, error) {
	parts := strings.Split(address, "#")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid address format: %s", address)
	}
	_, domain := parts[0], parts[1]
	txtDomain := "_emsg." + domain
	txts, err := net.LookupTXT(txtDomain)
	if err != nil {
		return "", fmt.Errorf("DNS TXT lookup failed for %s: %w", txtDomain, err)
	}
	if len(txts) == 0 {
		return "", fmt.Errorf("no TXT records found for %s", txtDomain)
	}
	// For now, return the first TXT record (could be routing info, pubkey, etc.)
	return txts[0], nil
}

// RouteInfo represents routing information from DNS TXT records
type RouteInfo struct {
	Server    string `json:"server"`  // EMSG server endpoint
	PublicKey string `json:"pubkey"`  // Domain's public key
	Version   string `json:"version"` // EMSG protocol version
	TTL       int    `json:"ttl"`     // Cache TTL in seconds
}

// ParseRouteInfo parses a TXT record into structured routing information
func ParseRouteInfo(txtRecord string) (*RouteInfo, error) {
	// Try to parse as JSON first
	var routeInfo RouteInfo
	if err := json.Unmarshal([]byte(txtRecord), &routeInfo); err == nil {
		return &routeInfo, nil
	}

	// Fallback: parse as simple server URL
	if strings.HasPrefix(txtRecord, "http://") || strings.HasPrefix(txtRecord, "https://") {
		return &RouteInfo{
			Server:  txtRecord,
			Version: "1.0",
			TTL:     3600, // 1 hour default
		}, nil
	}

	return nil, fmt.Errorf("unable to parse route info: %s", txtRecord)
}

// GetRouteInfo gets comprehensive routing information for an address
func GetRouteInfo(address string) (*RouteInfo, error) {
	txtRecord, err := LookupRoute(address)
	if err != nil {
		return nil, err
	}

	return ParseRouteInfo(txtRecord)
}

// ValidateAddress checks if an EMSG address is valid
func ValidateAddress(address string) error {
	parts := strings.Split(address, "#")
	if len(parts) != 2 {
		return fmt.Errorf("invalid address format: must be user#domain.com")
	}

	user, domain := parts[0], parts[1]
	if user == "" {
		return fmt.Errorf("user part cannot be empty")
	}
	if domain == "" {
		return fmt.Errorf("domain part cannot be empty")
	}

	// Basic domain validation
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("domain must contain at least one dot")
	}

	return nil
}

// RouteMessage determines where to route a message based on recipient addresses
func RouteMessage(recipients []string) (map[string][]string, error) {
	routes := make(map[string][]string) // server -> list of recipients

	for _, recipient := range recipients {
		if err := ValidateAddress(recipient); err != nil {
			return nil, fmt.Errorf("invalid recipient %s: %w", recipient, err)
		}

		routeInfo, err := GetRouteInfo(recipient)
		if err != nil {
			return nil, fmt.Errorf("failed to get route info for %s: %w", recipient, err)
		}

		if routeInfo.Server == "" {
			return nil, fmt.Errorf("no server found for %s", recipient)
		}

		routes[routeInfo.Server] = append(routes[routeInfo.Server], recipient)
	}

	return routes, nil
}

// IsLocalDomain checks if a domain should be handled locally
func IsLocalDomain(domain string, localDomains []string) bool {
	for _, localDomain := range localDomains {
		if domain == localDomain {
			return true
		}
	}
	return false
}
