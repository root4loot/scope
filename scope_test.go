package scope

import (
	"testing"
)

func TestDomainInclusionsWithPorts(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"example.com", true},
		{"example.com:8080", false},
		{"example.com:9090", false},
		{"sub.example.com", true},
		{"another.example.com", true},
		{"deep.sub.example.com", true},
		{"http://example.com", true},
		{"https://example.com", true},
		{"https://example.com:443", true},
		{"https://example.com:8080", false},
	}

	sc := NewScope()

	sc.AddIncludes([]string{
		"example.com",
		"example.com:8080",
		"sub.example.com",
		"another.example.com",
		"deep.sub.example.com",
		"http://example.com",
		"https://example.com",
		"https://example.com:443",
	})

	sc.AddExcludes([]string{
		"https://example.com:8080",
		"example.com:9090",
	})

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result := sc.IsInScope(test.url)
			if result != test.expected {
				t.Errorf("expected %v for URL '%s', got %v", test.expected, test.url, result)
			}
		})
	}
}
func TestIPInclusions(t *testing.T) {
	sc := NewScope()

	sc.AddIncludes([]string{
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"192.168.3.2-5",
		"192.168.2.0/24",
	})

	sc.AddExcludes([]string{
		"http://192.168.1.1",
		"192.168.2.0/24",
	})

	testCases := []struct {
		url      string
		expected bool
	}{
		{"192.168.1.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"192.168.1.1", true},
		{"192.168.3.3", true},
		{"192.168.3.5", true},
		{"192.168.1.6", false},
		{"10.0.0.2", false},
		{"172.16.0.2", false},
		{"http://192.168.1.1", false},
	}

	for _, tc := range testCases {
		t.Run(tc.url, func(t *testing.T) {
			result := sc.IsInScope(tc.url)

			if result != tc.expected {
				t.Errorf("expected %v for URL '%s', got %v", tc.expected, tc.url, result)
			}
		})
	}
}
