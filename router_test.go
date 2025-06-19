// router_test.go
// Tests for DNS TXT-based routing (mocked)
package main

import "testing"

func TestLookupRoute(t *testing.T) {
	// This is a placeholder test. Actual DNS lookup should be mocked or tested with a test domain.
	_, err := LookupRoute("bob#emsg.dev")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
