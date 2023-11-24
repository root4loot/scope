package goscope

import (
	"testing"
)

func TestScope(t *testing.T) {
	s := NewScope()
	s.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com", "example2.com:8080", "*.example.*.test")
	s.AddExclude("somedomain.com", "exclude.example.com", "192.168.0.6", "192.168.0.2:8080", "exclude.example.com:443")

	// Test single IPs
	if !s.IsTargetIncluded("192.168.0.2") {
		t.Errorf("Expected true for IsTargetIncluded, got false")
	}

	if s.IsTargetIncluded("192.168.0.7") {
		t.Errorf("Expected false for IsTargetIncluded, got true")
	}

	if !s.IsTargetExcluded("192.168.0.6") {
		t.Errorf("Expected true for IsTargetExcluded, got false")
	}

	if s.IsTargetExcluded("192.168.0.1") {
		t.Errorf("Expected false for IsTargetExcluded, got true")
	}

	// Test IP range
	if !s.IsTargetIncluded("192.168.0.4") {
		t.Errorf("Expected true for IsTargetIncluded, got false")
	}

	if s.IsTargetIncluded("192.168.0.8") {
		t.Errorf("Expected false for IsTargetIncluded, got true")
	}

	// Test CIDR notation
	if !s.IsTargetIncluded("192.168.10.50") {
		t.Errorf("Expected true for IsTargetIncluded, got false")
	}

	if s.IsTargetIncluded("192.168.11.50") {
		t.Errorf("Expected false for IsTargetIncluded, got true")
	}

	// Test Domains
	if !s.IsTargetIncluded("foo.example.com") {
		t.Errorf("Expected true for IsTargetIncluded, got false")
	}

	if s.IsTargetIncluded("bar.otherdomain.com") {
		t.Errorf("Expected false for IsTargetIncluded, got true")
	}

	if !s.IsTargetExcluded("exclude.example.com") {
		t.Errorf("Expected true for IsTargetExcluded, got false")
	}

	// Test Subdomains
	if !s.IsTargetExcluded("sub.somedomain.com") {
		t.Errorf("Expected true for IsTargetExcluded, got false")
	}

	if s.IsTargetExcluded("sub.otherdomain.com") {
		t.Errorf("Expected false for IsTargetExcluded, got true")
	}

	// Test port numbers
	if !s.IsTargetIncluded("example2.com:8080") {
		t.Errorf("Expected true for IsTargetIncluded, got false")
	}

	if s.IsTargetIncluded("example2.com:1234") {
		t.Errorf("Expected false for IsTargetIncluded, got true")
	}

	// Test wildcard
	if !s.IsTargetIncluded("foo.example.bar.test") {
		t.Errorf("Expected true for IsTargetIncluded, got false")
	}

	if s.IsTargetIncluded("foo.bar.baz.test") {
		t.Errorf("Expected false for IsTargetIncluded, got true")
	}

	// Test IP and port number
	if !s.IsTargetExcluded("192.168.0.2:8080") {
		t.Errorf("Expected true for IsTargetExcluded, got false")
	}

	if s.IsTargetExcluded("192.168.0.2:9090") {
		t.Errorf("Expected false for IsTargetExcluded, got true")
	}

	// Test domain and port number
	if !s.IsTargetExcluded("exclude.example.com:443") {
		t.Errorf("Expected true for IsTargetExcluded, got false")
	}

	if s.IsTargetExcluded("exclude.example.com:80") {
		t.Errorf("Expected false for IsTargetExcluded, got true")
	}

	// Test InScope
	if !s.IsTargetInScope("192.168.0.2") {
		t.Errorf("Expected true for InScope, got false")
	}

	if s.IsTargetInScope("192.168.0.7") {
		t.Errorf("Expected false for InScope, got true")
	}

	// Test Scheme
	if !s.IsTargetInScope("http://foo.example.com") {
		t.Errorf("Expected true for InScope with scheme, got false")
	}

	if s.IsTargetInScope("https://bar.otherdomain.com") {
		t.Errorf("Expected false for InScope with scheme, got true")
	}

	// Test excluded subdomain due to parent domain
	if s.IsTargetInScope("sub.somedomain.com") {
		t.Errorf("Expected false for InScope, got true")
	}

	if s.IsTargetInScope("sub.otherdomain.com") {
		t.Errorf("Expected false for InScope, got true")
	}
}

func TestScopeModifications(t *testing.T) {
	s := NewScope()
	s.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com")

	// Test AddTargetToScope
	if err := s.AddTargetToScope("newHost.com"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// With trailing slash
	if err := s.AddTargetToScope("newhostwithslash.com/"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !s.IsTargetIncluded("newHost.com") {
		t.Errorf("Expected true for IsTargetIncluded after AddToScope, got false")
	}

	if !s.IsTargetIncluded("newhostwithslash.com") {
		t.Errorf("Expected true for target with slash, got false")
	}

	if !s.IsTargetInScope("newHost.com") {
		t.Errorf("Expected true for InScope after AddToScope, got false")
	}

	if !s.IsTargetInScope("newhostwithslash.com") {
		t.Errorf("Expected true for target with slash, got false")
	}

	// Test RemoveFromScope
	if err := s.RemoveTargetFromScope("newHost.com"); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if s.IsTargetIncluded("newHost.com") {
		t.Errorf("Expected false for IsTargetIncluded after RemoveFromScope, got true")
	}

	if s.IsTargetInScope("newHost.com") {
		t.Errorf("Expected false for InScope after RemoveFromScope, got true")
	}
}
