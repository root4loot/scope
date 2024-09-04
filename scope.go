package scope

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/root4loot/goutils/domainutil"
	"github.com/root4loot/goutils/iputil"
	"github.com/root4loot/goutils/urlutil"
)

// ScopeDefinition holds the original string and its compiled regex
type ScopeDefinition struct {
	Original string
	Regex    *regexp.Regexp
}

// Scope holds the include and exclude lists
type Scope struct {
	includes []ScopeDefinition
	excludes []ScopeDefinition
}

// NewScope initializes and returns a new Scope instance
func NewScope() *Scope {
	return &Scope{
		includes: []ScopeDefinition{},
		excludes: []ScopeDefinition{},
	}
}

// AddInclude adds a new single include scope definition
func (s *Scope) AddInclude(definition string) error {
	var regex *regexp.Regexp
	var err error

	if iputil.IsIP(definition) || iputil.IsIPRange(definition) {
		s.includes = append(s.includes, ScopeDefinition{
			Original: definition,
			Regex:    nil,
		})
		return nil
	}

	regex, err = convertToRegex(definition)
	if err != nil {
		return fmt.Errorf("failed to add include '%s': %w", definition, err)
	}

	s.includes = append(s.includes, ScopeDefinition{
		Original: definition,
		Regex:    regex,
	})
	return nil
}

// AddIncludes adds multiple include scope definitions from a slice
func (s *Scope) AddIncludes(definitions []string) error {
	for _, def := range definitions {
		if err := s.AddInclude(def); err != nil {
			return fmt.Errorf("failed to add include '%s': %w", def, err)
		}
	}
	return nil
}

// AddExclude adds a new single exclude scope definition
func (s *Scope) AddExclude(definition string) error {
	var regex *regexp.Regexp
	var err error

	if iputil.IsIP(definition) || iputil.IsIPRange(definition) {
		s.excludes = append(s.excludes, ScopeDefinition{
			Original: definition,
			Regex:    nil,
		})
		return nil
	}

	regex, err = convertToRegex(definition)
	if err != nil {
		return fmt.Errorf("failed to add exclude '%s': %w", definition, err)
	}

	s.excludes = append(s.excludes, ScopeDefinition{
		Original: definition,
		Regex:    regex,
	})
	return nil
}

// AddExcludes adds multiple exclude scope definitions from a slice
func (s *Scope) AddExcludes(definitions []string) error {
	for _, def := range definitions {
		if err := s.AddExclude(def); err != nil {
			return fmt.Errorf("failed to add exclude '%s': %w", def, err)
		}
	}
	return nil
}

// IsInScope checks if a given URL or domain is in scope
func (s *Scope) IsInScope(target string) bool {
	if domainutil.IsDomainName(target) || urlutil.IsURL(target) {
		return s.inScopeURL(target)
	}

	if iputil.IsIP(target) {
		return s.inScopeIP(target)
	}

	return false
}

// GetScope returns the active inclusions, removing any that are excluded
func (s *Scope) GetScope() []string {
	var result []string
	for _, include := range s.includes {
		excluded := false
		for _, exclude := range s.excludes {
			if exclude.Regex.MatchString(include.Original) {
				excluded = true
				break
			}
		}
		if !excluded {
			result = append(result, include.Original)
		}
	}
	return result
}

func (s *Scope) inScopeIP(ip string) bool {
	checkMatch := func(definitions []ScopeDefinition, shouldMatch bool) bool {
		for _, def := range definitions {
			if def.Regex == nil {
				if ip == def.Original {
					return shouldMatch
				}
				if iputil.IsIPRange(def.Original) && iputil.IsIPInRange(ip, def.Original) {
					return shouldMatch
				}
				if iputil.IsCIDR(def.Original) && iputil.IsIPInCIDR(ip, def.Original) {
					return shouldMatch
				}
			} else if def.Regex.MatchString(ip) {
				return shouldMatch
			}
		}
		return !shouldMatch
	}

	if !checkMatch(s.excludes, false) {
		return false
	}
	return checkMatch(s.includes, true)
}

func (s *Scope) inScopeURL(url string) bool {
	checkMatch := func(definitions []ScopeDefinition, shouldMatch bool) bool {
		for _, def := range definitions {
			if def.Regex.MatchString(url) {
				return shouldMatch
			}
		}
		return !shouldMatch
	}

	if !checkMatch(s.excludes, false) {
		return false
	}
	return checkMatch(s.includes, true)
}

func convertToRegex(definition string) (*regexp.Regexp, error) {
	var regexPattern string

	switch {
	case strings.HasPrefix(definition, "http://"):
		definition = strings.TrimPrefix(definition, "http://")
		regexPattern = `^http://` + regexp.QuoteMeta(definition) + `(:\d+)?$`
	case strings.HasPrefix(definition, "https://"):
		definition = strings.TrimPrefix(definition, "https://")
		regexPattern = `^https://` + regexp.QuoteMeta(definition) + `(:\d+)?$`
	default:
		if strings.HasPrefix(definition, "*.") {
			definition = strings.TrimPrefix(definition, "*.")
			regexPattern = `^.*\.` + regexp.QuoteMeta(definition) + `(:\d+)?$`
		} else {
			regexPattern = `^` + regexp.QuoteMeta(definition) + `(:\d+)?$`
		}
	}

	compiledRegex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	return compiledRegex, nil
}
