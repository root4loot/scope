package goscope

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/iputil"
)

type TargetType int

const (
	IP TargetType = iota
	Domain
	CIDR
	IPRange
	Other
)

type Scope struct {
	Includes map[string]bool
	Excludes map[string]bool
	Targets  map[string]TargetType
}

// NewScope returns a new Scope.
func NewScope() *Scope {
	return &Scope{
		Includes: make(map[string]bool),
		Excludes: make(map[string]bool),
		Targets:  make(map[string]TargetType),
	}
}

// String returns the string representation of the HostType
func (targetType TargetType) String() string {
	return [...]string{"IP", "Domain", "CIDR", "IPRange", "Other"}[targetType]
}

// GetIncludes returns a string slice representation of the scope's Includes list.
func (s *Scope) GetIncludes() (includes []string) {
	for include := range s.Includes {
		includes = append(includes, include)
	}
	return
}

// GetExcludes returns a string slice representation of the scope's Excludes list.
func (s *Scope) GetExcludes() (excludes []string) {
	for exclude := range s.Excludes {
		excludes = append(excludes, exclude)
	}
	return
}

// GetTargets returns all hosts as a string slice
func (s *Scope) GetTargets() (targets []string) {
	for target := range s.Targets {
		targets = append(targets, target)
	}
	return targets
}

// GetTargetsAndTypeMap returns all targets and their types as a map
func (s *Scope) GetTargetsAndTypeMap() map[string]TargetType {
	return s.Targets
}

// AddTargetToScope adds one or more targets to the scope's Targets list.
func (s *Scope) AddTargetToScope(targets ...string) error {
	for _, target := range targets {
		target = strings.ToLower(target)
		target = removeScheme(target)

		if s.IsTargetExcluded(target) {
			return fmt.Errorf("target %s is excluded", target)
		}
		hostType := categorizeHost(target)
		s.Targets[target] = hostType
		s.AddInclude(target) // Automatically add to Includes
	}
	return nil
}

// RemoveTargetFromScope removes a target from the scope's Targets list.
func (s *Scope) RemoveTargetFromScope(target string) error {
	if _, exists := s.Targets[target]; !exists {
		return fmt.Errorf("target %s does not exist in scope", target)
	}

	delete(s.Targets, target)
	delete(s.Includes, target)
	return nil
}

// AddInclude adds one or more targets to the scope's Includes list.
func (s *Scope) AddInclude(targets ...string) error {
	for _, target := range targets {
		target = removeScheme(target)
		target = strings.ToLower(target)
		if err := s.addHostToScope(target, s.Includes); err != nil {
			return err
		}
	}
	return nil
}

// AddExclude adds one or more targets to the scope's Excludes list.
func (s *Scope) AddExclude(targets ...string) error {
	for _, target := range targets {
		target = removeScheme(target)
		target = strings.ToLower(target)
		if err := s.addHostToScope(target, s.Excludes); err != nil {
			return err
		}
	}
	return nil
}

// IsTargetIncluded returns true if the target is in the scope's Includes list.
func (s *Scope) IsTargetIncluded(target string) bool {
	target = removeScheme(target)
	target = strings.ToLower(target)

	if s.Includes[target] {
		return true
	}

	for include := range s.Includes {
		if isWildcardMatch(include, target) {
			return true
		}
	}

	if iputil.IsIP(target) {
		ip := net.ParseIP(target)
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

// IsTargetExcluded returns true if the target is in the scope's Excludes list.
func (s *Scope) IsTargetExcluded(target string) bool {
	target = removeScheme(target)
	target = strings.ToLower(target)

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

// IsTargetInScope returns true if the target is in the scope's Includes list and not in the Excludes list.
func (s *Scope) IsTargetInScope(target string) bool {
	target = strings.ToLower(target)
	target = removeScheme(target)
	return s.IsTargetIncluded(target) && !s.IsTargetExcluded(removeScheme(target))
}

// addHostToScope adds a target to the scope
func (s *Scope) addHostToScope(target string, scope map[string]bool) error {
	target, port := splitIPAndPort(target)
	additional := ""
	if port != "" {
		additional = ":" + port
	}

	if iputil.IsValidIP(target) || domainutil.IsDomainName(target) {
		scope[target+additional] = true
		return nil
	}

	if iputil.IsValidCIDR(target) || iputil.IsValidIPRange(target) || strings.Contains(target, "*") {
		scope[target] = true
		return nil
	}

	return fmt.Errorf("invalid host: %s", target)
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
func removeScheme(target string) string {
	if idx := strings.Index(target, "://"); idx != -1 {
		return target[idx+3:]
	}
	return target
}

// categorizeHost categorizes the target into IP or Domain or other.
func categorizeHost(target string) TargetType {
	if iputil.IsIP(target) {
		return IP
	}
	if domainutil.IsDomainName(target) {
		return Domain
	}
	if iputil.IsCIDR(target) {
		return CIDR
	}
	if iputil.IsIPRange(target) {
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
