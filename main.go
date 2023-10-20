package goscope

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/iputil"
)

type Scope struct {
	Includes map[string]bool
	Excludes map[string]bool
}

// NewScope returns a new Scope.
func NewScope() *Scope {
	return &Scope{
		Includes: make(map[string]bool),
		Excludes: make(map[string]bool),
	}
}

// AddInclude adds hosts to the scope's Includes list.
func (s *Scope) AddInclude(hosts ...string) error {
	for _, host := range hosts {
		if err := s.addHostToScope(host, s.Includes); err != nil {
			return err
		}
	}
	return nil
}

// AddExclude adds a host to the scope's Excludes list.
func (s *Scope) AddExclude(hosts ...string) error {
	for _, host := range hosts {
		if err := s.addHostToScope(host, s.Excludes); err != nil {
			return err
		}
	}
	return nil
}

// IsIncluded returns true if the host is in the scope's Includes list.
func (s *Scope) IsIncluded(host string) bool {
	if s.Includes[host] {
		return true
	}

	for include := range s.Includes {
		if IsWildcardMatch(include, host) {
			return true
		}
	}

	if iputil.IsIP(host) {
		ip := net.ParseIP(host)
		for include := range s.Includes {
			if iputil.IsCIDR(include) {
				b, _ := iputil.IsIPInCIDR(ip.String(), include)
				if b {
					return true
				}
			}

			if iputil.IsIPRange(include) {
				b, _ := iputil.IsIPInRange(ip.String(), include)
				if b {
					return true
				}
			}
		}
	}

	return false
}

// IsExcluded returns true if the host is in the scope's Excludes list.
func (s *Scope) IsExcluded(host string) bool {
	// Direct match
	if s.Excludes[host] {
		return true
	}

	// Check for wildcard or other complex rules
	for exclude := range s.Excludes {
		if IsWildcardMatch(exclude, host) {
			return true
		}
	}

	// Check if the parent domain is excluded
	parts := strings.Split(host, ".")
	for i := 0; i < len(parts); i++ {
		parentDomain := strings.Join(parts[i:], ".")
		if s.Excludes[parentDomain] {
			return true
		}
		for exclude := range s.Excludes {
			if IsWildcardMatch(exclude, parentDomain) {
				return true
			}
		}
	}
	return false
}

// InScope returns true if the host is in the scope's Includes list and not in the Excludes list.
func (s *Scope) InScope(host string) bool {
	return s.IsIncluded(host) && !s.IsExcluded(host)
}

// IsWildcardMatch returns true if the wildcard notation in the pattern matches the input string
func IsWildcardMatch(pattern, input string) bool {
	regexPattern := strings.ReplaceAll(pattern, ".", `\.`)     // Escape dots in the pattern
	regexPattern = strings.ReplaceAll(regexPattern, "*", ".*") // Replace asterisks with .* for wildcard matching

	matched, _ := regexp.MatchString("^"+regexPattern+"$", input)
	return matched
}

// addHostToScope adds a host to the scope
func (s *Scope) addHostToScope(host string, scope map[string]bool) error {
	host, port := splitIPAndPort(host)

	if iputil.IsValidIP(host) || domainutil.IsDomainName(host) {
		// fmt.Println("Adding host to scope:", host)
		if port != "" {
			scope[host+":"+port] = true
		} else {
			scope[host] = true
		}
		return nil
	} else if iputil.IsValidCIDR(host) ||
		iputil.IsValidIPRange(host) ||
		strings.Contains(host, "*") {
		scope[host] = true
		return nil
	} else {
		return fmt.Errorf("invalid host: %s", host)
	}
}

// splitIPAndPort splits a host:port string into host and port
func splitIPAndPort(input string) (string, string) {
	host, port, err := net.SplitHostPort(input)
	if err != nil {
		// Invalid or no port found
		return input, ""
	}

	return host, port
}
