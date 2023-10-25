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
func (s *Scope) AddInclude(targets ...string) error {
	for _, target := range targets {
		target = removeScheme(target)
		if err := s.addHostToScope(target, s.Includes); err != nil {
			return err
		}
	}
	return nil
}

// AddExclude adds a host to the scope's Excludes list.
func (s *Scope) AddExclude(targets ...string) error {
	for _, target := range targets {
		target = removeScheme(target)
		if err := s.addHostToScope(target, s.Excludes); err != nil {
			return err
		}
	}
	return nil
}

// IsIncluded returns true if the target is in the scope's Includes list.
func (s *Scope) IsIncluded(target string) bool {
	host := removeScheme(target)
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

func (s *Scope) IsExcluded(target string) bool {
	for exclude := range s.Excludes {
		if target == exclude {
			return true
		}

		// Splitting the target and exclude strings by ':'
		excludeParts := strings.Split(exclude, ":")
		targetParts := strings.Split(target, ":")

		// Comparing base parts
		excludeBase := excludeParts[0]
		targetBase := targetParts[0]

		if targetBase == excludeBase || strings.HasSuffix(targetBase, "."+excludeBase) {
			// If ports are specified, they must match. Otherwise, it's a match.
			if len(excludeParts) > 1 && len(targetParts) > 1 {
				if excludeParts[1] == targetParts[1] {
					return true
				}
			} else if len(excludeParts) == 1 && len(targetParts) == 1 {
				return true
			}
		}
	}
	return false
}

// InScope returns true if the target is in the scope's Includes list and not in the Excludes list.
func (s *Scope) InScope(target string) bool {
	host := removeScheme(target)
	return s.IsIncluded(host) && !s.IsExcluded(removeScheme(host))
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

// removeScheme removes the URL scheme (e.g., "http://", "https://", "ftp://") from the given string.
func removeScheme(host string) string {
	if idx := strings.Index(host, "://"); idx != -1 {
		return host[idx+3:]
	}
	return host
}
