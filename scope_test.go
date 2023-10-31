package goscope

import (
	"testing"
)

func TestScope(t *testing.T) {
	s := NewScope()
	s.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com", "example2.com:8080", "*.example.*.test")
	s.AddExclude("somedomain.com", "exclude.example.com", "192.168.0.6", "192.168.0.2:8080", "exclude.example.com:443")

	// Test single IPs
	if !s.IsIncluded("192.168.0.2") {
		t.Errorf("Expected true for IsIncluded, got false")
	}

	if s.IsIncluded("192.168.0.7") {
		t.Errorf("Expected false for IsIncluded, got true")
	}

	if !s.IsExcluded("192.168.0.6") {
		t.Errorf("Expected true for IsExcluded, got false")
	}

	if s.IsExcluded("192.168.0.1") {
		t.Errorf("Expected false for IsExcluded, got true")
	}

	// Test IP range
	if !s.IsIncluded("192.168.0.4") {
		t.Errorf("Expected true for IsIncluded, got false")
	}

	if s.IsIncluded("192.168.0.8") {
		t.Errorf("Expected false for IsIncluded, got true")
	}

	// Test CIDR notation
	if !s.IsIncluded("192.168.10.50") {
		t.Errorf("Expected true for IsIncluded, got false")
	}

	if s.IsIncluded("192.168.11.50") {
		t.Errorf("Expected false for IsIncluded, got true")
	}

	// Test Domains
	if !s.IsIncluded("foo.example.com") {
		t.Errorf("Expected true for IsIncluded, got false")
	}

	if s.IsIncluded("bar.otherdomain.com") {
		t.Errorf("Expected false for IsIncluded, got true")
	}

	if !s.IsExcluded("exclude.example.com") {
		t.Errorf("Expected true for IsExcluded, got false")
	}

	// Test Subdomains
	if !s.IsExcluded("sub.somedomain.com") {
		t.Errorf("Expected true for IsExcluded, got false")
	}

	if s.IsExcluded("sub.otherdomain.com") {
		t.Errorf("Expected false for IsExcluded, got true")
	}

	// Test port numbers
	if !s.IsIncluded("example2.com:8080") {
		t.Errorf("Expected true for IsIncluded, got false")
	}

	if s.IsIncluded("example2.com:1234") {
		t.Errorf("Expected false for IsIncluded, got true")
	}

	// Test wildcard
	if !s.IsIncluded("foo.example.bar.test") {
		t.Errorf("Expected true for IsIncluded, got false")
	}

	if s.IsIncluded("foo.bar.baz.test") {
		t.Errorf("Expected false for IsIncluded, got true")
	}

	// Test IP and port number
	if !s.IsExcluded("192.168.0.2:8080") {
		t.Errorf("Expected true for IsExcluded, got false")
	}

	if s.IsExcluded("192.168.0.2:9090") {
		t.Errorf("Expected false for IsExcluded, got true")
	}

	// Test domain and port number
	if !s.IsExcluded("exclude.example.com:443") {
		t.Errorf("Expected true for IsExcluded, got false")
	}

	if s.IsExcluded("exclude.example.com:80") {
		t.Errorf("Expected false for IsExcluded, got true")
	}

	// Test InScope
	if !s.InScope("192.168.0.2") {
		t.Errorf("Expected true for InScope, got false")
	}

	if s.InScope("192.168.0.7") {
		t.Errorf("Expected false for InScope, got true")
	}

	// Test Scheme
	if !s.InScope("http://foo.example.com") {
		t.Errorf("Expected true for InScope with scheme, got false")
	}

	if s.InScope("https://bar.otherdomain.com") {
		t.Errorf("Expected false for InScope with scheme, got true")
	}

	// Test excluded subdomain due to parent domain
	if s.InScope("sub.somedomain.com") {
		t.Errorf("Expected false for InScope, got true")
	}

	if s.InScope("sub.otherdomain.com") {
		t.Errorf("Expected false for InScope, got true")
	}
}
