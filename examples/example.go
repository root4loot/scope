package main

import (
	"fmt"

	"github.com/root4loot/goscope"
)

func main() {
	// Initialize a new scope object
	s := goscope.NewScope()

	// Add IP ranges, CIDRs, and domains to the include list
	s.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com", "example2.com:8080", "*.example.*.test")

	// Add IP ranges, CIDRs, and domains to the exclude list
	s.AddExclude("somedomain.com", "exclude.example.com", "192.168.0.6", "192.168.0.2:8080", "exclude.example.com:443")

	// Determine if single IPs are in the include list
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("192.168.0.2")) // Returns true
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("192.168.0.7")) // Returns false

	// Determine if single IPs are in the exclude list
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("192.168.0.6")) // Returns true
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("192.168.0.1")) // Returns false

	// Determine if single IPs are in scope (in includes but not in excludes)
	fmt.Println("InScope:", s.IsTargetInScope("192.168.0.2")) // Returns true
	fmt.Println("InScope:", s.IsTargetInScope("192.168.0.7")) // Returns false

	// Determine if IPs in a range are in the include list
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("192.168.0.4")) // Returns true
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("192.168.0.8")) // Returns false

	// Determine if IPs in a CIDR notation are in the include list
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("192.168.10.50")) // Returns true
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("192.168.11.50")) // Returns false

	// Determine if domains are in the include list
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("foo.example.com"))     // Returns true
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("bar.otherdomain.com")) // Returns false

	// Determine if domains are in the exclude list
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("exclude.example.com")) // Returns true
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("include.example.com")) // Returns false

	// Determine if subdomains are excluded due to their parent domain being in the excludes list
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("sub.somedomain.com"))  // Returns true
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("sub.otherdomain.com")) // Returns false

	// Determine if domains with specific ports are in the include list
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("example2.com:8080")) // Returns true
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("example2.com:1234")) // Returns false

	// Determine if wildcard domains are in the include list
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("foo.example.bar.test")) // Returns true
	fmt.Println("IsTargetIncluded:", s.IsTargetIncluded("foo.bar.baz.test"))     // Returns false

	// Determine if IPs with specific ports are in the exclude list
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("192.168.0.2:8080")) // Returns true
	fmt.Println("IsTargetExcluded:", s.IsTargetExcluded("192.168.0.2:9090")) // Returns false

	// Add specific domains and IPs to the scope
	s.AddTargetToScope("new.example.com", "192.168.15.1")
	fmt.Println("AddToScope:", s.IsTargetInScope("new.example.com")) // Should return true

	// Remove specific domains and IPs from the scope
	s.RemoveTargetFromScope("new.example.com")
	fmt.Println("RemoveFromScope:", s.IsTargetInScope("new.example.com")) // Should return false

	// Retrieve all targets in the scope as a slice
	targets := s.GetTargets()
	fmt.Println("GetHostsAsSlice:", targets)

	// Retrieve all domains in the scope as a slice
	domains := s.GetTargetDomains()
	fmt.Println("GetDomainsAsSlice:", domains)

	// Retrieve all IPs in the scope as a slice
	ips := s.GetTargetIPs()
	fmt.Println("GetIPsAsSlice:", ips)

	// Print scope details as strings
	fmt.Println(s.GetIncludes())
	fmt.Println(s.GetExcludes())
	fmt.Println(s.GetTargets())
}
