package goscope

import (
	"testing"
)

func Init() *Scope {
	s := NewScope()
	s.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com", "example2.com:8080", "*.example.*.test")
	s.AddExclude("somedomain.com", "exclude.example.com", "192.168.0.6", "192.168.0.2:8080", "exclude.example.com:443")
	s.AddTargetToScope("example.com", "example.com/", "http://example.com", "http://example.com/", "http://example.com/foo", "http://example.com/foo/", "https://example.com/robots.txt")

	return s
}

func TestIsIncluded(t *testing.T) {
	s := Init()
	if !s.IsTargetIncluded("192.168.0.2") {
		t.Fatal()
	}

	if s.IsTargetIncluded("192.168.0.7") {
		t.Fatal()
	}

	if !s.IsTargetIncluded("192.168.0.4") {
		t.Fatal()
	}

	if s.IsTargetIncluded("192.168.0.8") {
		t.Fatal()
	}

	if !s.IsTargetIncluded("192.168.10.50") {
		t.Fatal()
	}

	if s.IsTargetIncluded("192.168.11.50") {
		t.Fatal()
	}

	if !s.IsTargetIncluded("foo.example.com") {
		t.Fatal()
	}

	if s.IsTargetIncluded("bar.otherdomain.com") {
		t.Fatal()
	}

	if !s.IsTargetIncluded("example2.com:8080") {
		t.Fatal()
	}

	if s.IsTargetIncluded("example2.com:1234") {
		t.Fatal()
	}

	if !s.IsTargetIncluded("foo.example.bar.test") {
		t.Fatal()
	}

	if s.IsTargetIncluded("foo.bar.baz.test") {
		t.Fatal()
	}
}

func TestIsExcluded(t *testing.T) {
	s := Init()
	if !s.IsTargetExcluded("192.168.0.6") {
		t.Fatal()
	}

	if s.IsTargetExcluded("192.168.0.1") {
		t.Fatal()
	}

	if !s.IsTargetExcluded("exclude.example.com") {
		t.Fatal()
	}

	if !s.IsTargetExcluded("sub.somedomain.com") {
		t.Fatal()
	}

	if !s.IsTargetExcluded("192.168.0.2:8080") {
		t.Fatal()
	}

	if s.IsTargetExcluded("192.168.0.2:9090") {
		t.Fatal()
	}

	if !s.IsTargetExcluded("exclude.example.com:443") {
		t.Fatal()
	}

	if s.IsTargetExcluded("exclude.example.com:80") {
		t.Fatal()
	}
}

func TestIsInScope(t *testing.T) {
	s := Init()

	if s.IsTargetInScope("192.168.0.7") {
		t.Fatal()
	}

	if !s.IsTargetInScope("192.168.0.2") {
		t.Fatal()
	}

	if !s.IsTargetInScope("https://example.com/robots.txt") {
		t.Fatal()
	}

	if !s.IsTargetInScope("http://foo.example.com") {
		t.Fatal()
	}

	if !s.IsTargetInScope("sub.example.com") {
		t.Fatal()
	}

	if !s.IsTargetInScope("example.com") {
		t.Fatal()
	}

	if !s.IsTargetInScope("http://example.com/foo/scope.html") {
		t.Fatal()
	}
}

func TestRemoveTarget(t *testing.T) {
	s := Init()

	if err := s.RemoveTargetFromScope("example.com"); err != nil {
		t.Fatal()
	}

	if s.IsTargetInScope("example.com") {
		t.Fatal()
	}

	if err := s.RemoveTargetFromScope("example.com"); err == nil {
		t.Fatal()
	}
}
