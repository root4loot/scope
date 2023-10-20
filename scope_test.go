package goscope

import (
	"testing"
)

func TestScope(t *testing.T) {
	tests := []struct {
		host       string
		isIncluded bool
		isExcluded bool
		isInScope  bool
	}{
		{"192.168.0.1", true, false, true},
		{"192.168.0.6", false, false, false},
		{"foo.example.com", true, false, true},
		{"wildcard.example.com", true, false, true},
		{"excluded.example.com", false, true, false},
		{"sub.excluded.com", false, true, false}, // New test for subdomain exclusion
		{"test.example2.com", true, false, true},
		{"example.com:8080", true, false, true},
		{"abc.example.xyz.test", true, false, true},
		{"notinscope.com", false, false, false},
	}

	s := NewScope()
	s.AddInclude("192.168.0.1-5", "*.example.com", "example.com:8080", "test.example2.com", "*.example.*.test")
	s.AddExclude("excluded.example.com")
	s.AddExclude("excluded.com")

	for _, test := range tests {
		inScope := s.InScope(test.host)
		if inScope != test.isInScope {
			t.Errorf("InScope(%s) = %t, expected %t", test.host, inScope, test.isInScope)
		}
	}
}
