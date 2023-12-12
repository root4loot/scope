package goscope

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/iputil"
	"github.com/root4loot/goutils/sliceutil"
)

// Targets represents a set of IPs, domains, and other targets.
type Targets struct {
	IPs     []string
	Domains []string
	Other   []string
}

// Scope represents a set of includes, excludes, and targets.
// It must be initialized using NewScope() before use.
type Scope struct {
	Includes map[string]bool
	Excludes map[string]bool
	Targets  *Targets
}

// NewScope returns a new Scope.
func NewScope() *Scope {
	return &Scope{
		Includes: make(map[string]bool),
		Excludes: make(map[string]bool),
		Targets:  &Targets{},
	}
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

// GetTargetsCIDR returns a string slice representation of the scope's Targets CIDR list.
func (s *Scope) GetTargetCIDRs() (cidrs []string, err error) {
	cidrs, err = iputil.IPsToCIDR(s.Targets.IPs)
	if err != nil {
		return nil, err
	}
	return cidrs, nil
}

// GetAllTargets returns all hosts as a string slice
func (s *Scope) GetAllTargets() (targets []string) {
	targets = append(targets, s.Targets.IPs...)
	targets = append(targets, s.Targets.Domains...)
	targets = append(targets, s.Targets.Other...)
	return targets
}

// GetTargetDomains returns all domains as a string slice
func (s *Scope) GetTargetDomains() (domains []string) {
	return s.Targets.Domains
}

// GetTargetIPs returns all IPs as a string slice
func (s *Scope) GetTargetIPs() (ips []string) {
	return s.Targets.IPs
}

// GetTargetOther returns all other targets as a string slice
func (s *Scope) GetTargetOther() (other []string) {
	return s.Targets.Other
}

// AddTargetToScope adds one or more targets to the scope's Targets list.
func (s *Scope) AddTargetToScope(targets ...string) error {
	for _, target := range targets {
		target = strings.ToLower(target)
		target = removeSchemeAndTrailSlash(target)

		if s.IsTargetExcluded(target) {
			return fmt.Errorf("target %s is excluded", target)
		}

		if s.IsTargetAdded(target) {
			return nil // Target already added
		}

		if iputil.IsIP(target) {
			s.Targets.IPs = append(s.Targets.IPs, target)
		} else if domainutil.IsDomainName(target) {
			s.Targets.Domains = append(s.Targets.Domains, target)
		} else if iputil.IsCIDR(target) {
			cidrs, err := iputil.ParseCIDR(target)
			if err != nil {
				return err
			}
			for i := range cidrs {
				s.Targets.IPs = append(s.Targets.IPs, cidrs[i].String())
			}
		} else if iputil.IsIPRange(target) {
			ips, err := iputil.ParseIPRange(target)
			if err != nil {
				return err
			}
			for i := range ips {
				s.Targets.IPs = append(s.Targets.IPs, ips[i].String())
			}
		} else {
			s.Targets.Other = append(s.Targets.Other, target)
		}

		s.AddInclude(target) // Include added target
	}
	return nil
}

// RemoveTargetFromScope removes a target from the scope's Targets list.
func (s *Scope) RemoveTargetFromScope(target string) error {
	target = strings.ToLower(target)
	target = removeSchemeAndTrailSlash(target)

	if !s.IsTargetAdded(target) {
		return fmt.Errorf("target %s is not in scope", target)
	}

	if sliceutil.Contains(s.Targets.IPs, target) {
		s.Targets.IPs = sliceutil.Remove(s.Targets.IPs, target)
		delete(s.Includes, target)
	}

	if sliceutil.Contains(s.Targets.Domains, target) {
		s.Targets.Domains = sliceutil.Remove(s.Targets.Domains, target)
		delete(s.Includes, target)
	}
	return nil
}

// AddInclude adds one or more targets to the scope's Includes list.
func (s *Scope) AddInclude(targets ...string) error {
	for _, target := range targets {
		target = removeSchemeAndTrailSlash(target)
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
		target = removeSchemeAndTrailSlash(target)
		target = strings.ToLower(target)
		if err := s.addHostToScope(target, s.Excludes); err != nil {
			return err
		}
	}
	return nil
}

// IsTargetIncluded returns true if the target is in the scope's Includes list.
func (s *Scope) IsTargetIncluded(target string) bool {
	target = removeSchemeAndTrailSlash(target)
	target = strings.ToLower(target)

	if s.Includes[target] || isParentDirectoryIncluded(target, s.Includes) {
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

// IsTargetExcluded checks if a target is excluded, considering domain hierarchies and specific conditions.
func (s *Scope) IsTargetExcluded(target string) bool {
	target = removeSchemeAndTrailSlash(target)
	target = strings.ToLower(target)

	targetParts := strings.Split(target, ":")
	targetBase := targetParts[0]
	targetPort := ""
	if len(targetParts) > 1 {
		targetPort = targetParts[1]
	}

	for exclude := range s.Excludes {
		excludeParts := strings.Split(exclude, ":")
		excludeBase := excludeParts[0]
		excludePort := ""
		if len(excludeParts) > 1 {
			excludePort = excludeParts[1]
		}

		if targetBase == excludeBase {
			if targetPort == excludePort {
				return true // Exact match including port
			}
			if targetPort != "" && excludePort == "" {
				continue // Target has port, but exclude does not
			}
			if targetPort == "" && excludePort != "" {
				continue // Target does not have port, but exclude does
			}
			if targetPort == "" && excludePort == "" {
				return true // Match without port specification
			}
		} else if strings.HasSuffix(targetBase, "."+excludeBase) {
			// Check for domain hierarchy match without considering ports
			return true
		}
	}

	return false
}

// IsTargetInScope now considers parent directories as well.
func (s *Scope) IsTargetInScope(target string) bool {
	target = strings.ToLower(target)
	target = removeSchemeAndTrailSlash(target)
	return s.IsTargetIncluded(target) && !s.IsTargetExcluded(removeSchemeAndTrailSlash(target))
}

// IsTargetAdded returns true if the target is in the scope's Targets list.
func (s *Scope) IsTargetAdded(target string) bool {
	target = strings.ToLower(target)
	target = removeSchemeAndTrailSlash(target)

	for i := range s.Targets.Domains {
		if s.Targets.Domains[i] == target {
			return true
		}
	}

	for i := range s.Targets.IPs {
		if s.Targets.IPs[i] == target {
			return true
		}
	}

	for i := range s.Targets.Other {
		if s.Targets.Other[i] == target {
			return true
		}
	}

	return false
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

// removeSchemeAndTrailSlash removes the URL scheme (e.g., "http://", "https://", "ftp://") and trailing slash from the given string.
func removeSchemeAndTrailSlash(target string) string {
	target = removeScheme(target)
	target = removeTrailSlash(target)
	return target
}

// removeTrailSlash removes the trailing slash from the given string.
func removeTrailSlash(target string) string {
	if strings.HasSuffix(target, "/") {
		return target[:len(target)-1]
	}
	return target
}

// removeScheme removes the URL scheme (e.g., "http://", "https://", "ftp://") from the given string.
func removeScheme(target string) string {
	if idx := strings.Index(target, "://"); idx != -1 {
		return target[idx+3:]
	}
	return target
}

// isWildcardMatch returns true if the wildcard notation in the pattern matches the input string
func isWildcardMatch(pattern, input string) bool {
	regexPattern := strings.ReplaceAll(pattern, ".", `\.`)     // Escape dots in the pattern
	regexPattern = strings.ReplaceAll(regexPattern, "*", ".*") // Replace asterisks with .* for wildcard matching

	matched, _ := regexp.MatchString("^"+regexPattern+"$", input)
	return matched
}

// isParentDirectoryIncluded checks if any parent directory of the target is included.
func isParentDirectoryIncluded(target string, includes map[string]bool) bool {
	for include := range includes {
		if strings.HasPrefix(target, include) {
			return true
		}
	}
	return false
}

// isParentDirectoryExcluded checks if any parent directory of the target is excluded.
func isParentDirectoryExcluded(target string, excludes map[string]bool) bool {
	for exclude := range excludes {
		if strings.HasPrefix(target, exclude) {
			return true
		}
	}
	return false
}
