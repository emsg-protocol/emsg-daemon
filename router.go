// router.go
// DNS-based routing for EMSG Daemon
package main

import (
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
