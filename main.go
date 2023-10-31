package goscope

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/iputil"
)

type HostType int

const (
	IP HostType = iota
	Domain
	CIDR
	IPRange
	Other
)

type Scope struct {
	Includes map[string]bool
	Excludes map[string]bool
	Hosts    map[string]HostType
}

// NewScope returns a new Scope.
func NewScope() *Scope {
	return &Scope{
		Includes: make(map[string]bool),
		Excludes: make(map[string]bool),
		Hosts:    make(map[string]HostType),
	}
}

// String returns the string representation of the HostType
func (h HostType) String() string {
	return [...]string{"IP", "Domain", "CIDR", "IPRange", "Other"}[h]
}

// IncludeSlice returns a string slice representation of the scope's Includes list.
func (s *Scope) IncludeSlice() (includes []string) {
	for include := range s.Includes {
		includes = append(includes, include)
	}
	return
}

// ExcludeSlice returns a string slice representation of the scope's Excludes list.
func (s *Scope) ExcludeSlice() (excludes []string) {
	for exclude := range s.Excludes {
		excludes = append(excludes, exclude)
	}
	return
}

// HostSlice returns all hosts as a string slice
func (s *Scope) HostSlice() (hosts []string) {
	for host := range s.Hosts {
		hosts = append(hosts, host)
	}
	return hosts
}

// HostAndTypes returns all hosts and their types as a map
func (s *Scope) HostAndTypes() map[string]HostType {
	return s.Hosts
}

// AddToScope adds a host to the scope's Hosts list, with error checking against Excludes.
func (s *Scope) AddToScope(hosts ...string) error {
	for _, host := range hosts {
		if s.IsExcluded(host) {
			return fmt.Errorf("host %s is excluded", host)
		}
		hostType := categorizeHost(host)
		s.Hosts[host] = hostType
		s.AddInclude(host) // Automatically add to Includes
	}
	return nil
}

// RemoveFromScope removes a host from the scope's Hosts and Includes list
func (s *Scope) RemoveFromScope(host string) error {
	if _, exists := s.Hosts[host]; !exists {
		return fmt.Errorf("host %s does not exist in scope", host)
	}

	delete(s.Hosts, host)
	delete(s.Includes, host)
	return nil
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
		if isWildcardMatch(include, host) {
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

// IsExcluded returns true if the target is in the scope's Excludes list.
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

// addHostToScope adds a host to the scope
func (s *Scope) addHostToScope(host string, scope map[string]bool) error {
	host, port := splitIPAndPort(host)
	additional := ""
	if port != "" {
		additional = ":" + port
	}

	if iputil.IsValidIP(host) || domainutil.IsDomainName(host) {
		scope[host+additional] = true
		return nil
	}

	if iputil.IsValidCIDR(host) || iputil.IsValidIPRange(host) || strings.Contains(host, "*") {
		scope[host] = true
		return nil
	}

	return fmt.Errorf("invalid host: %s", host)
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

// categorizeHost categorizes the host into IP or Domain or other.
func categorizeHost(host string) HostType {
	if iputil.IsIP(host) {
		return IP
	}
	if domainutil.IsDomainName(host) {
		return Domain
	}
	if iputil.IsCIDR(host) {
		return CIDR
	}
	if iputil.IsIPRange(host) {
		return IPRange
	}
	return Other
}

// isWildcardMatch returns true if the wildcard notation in the pattern matches the input string
func isWildcardMatch(pattern, input string) bool {
	regexPattern := strings.ReplaceAll(pattern, ".", `\.`)     // Escape dots in the pattern
	regexPattern = strings.ReplaceAll(regexPattern, "*", ".*") // Replace asterisks with .* for wildcard matching

	matched, _ := regexp.MatchString("^"+regexPattern+"$", input)
	return matched
}
