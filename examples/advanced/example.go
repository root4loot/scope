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
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.2")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.7")) // Returns false

	// Determine if single IPs are in the exclude list
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.6")) // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.1")) // Returns false

	// Determine if single IPs are in scope (in includes but not in excludes)
	fmt.Println("InScope:", s.InScope("192.168.0.2")) // Returns true
	fmt.Println("InScope:", s.InScope("192.168.0.7")) // Returns false

	// Determine if IPs in a range are in the include list
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.4")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.8")) // Returns false

	// Determine if IPs in a CIDR notation are in the include list
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.10.50")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.11.50")) // Returns false

	// Determine if domains are in the include list
	fmt.Println("IsIncluded:", s.IsIncluded("foo.example.com"))     // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("bar.otherdomain.com")) // Returns false

	// Determine if domains are in the exclude list
	fmt.Println("IsExcluded:", s.IsExcluded("exclude.example.com")) // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("include.example.com")) // Returns false

	// Determine if subdomains are excluded due to their parent domain being in the excludes list
	fmt.Println("IsExcluded:", s.IsExcluded("sub.somedomain.com"))  // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("sub.otherdomain.com")) // Returns false

	// Determine if domains with specific ports are in the include list
	fmt.Println("IsIncluded:", s.IsIncluded("example2.com:8080")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("example2.com:1234")) // Returns false

	// Determine if wildcard domains are in the include list
	fmt.Println("IsIncluded:", s.IsIncluded("foo.example.bar.test")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("foo.bar.baz.test"))     // Returns false

	// Determine if IPs with specific ports are in the exclude list
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.2:8080")) // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.2:9090")) // Returns false

	// Add specific domains and IPs to the scope
	s.AddToScope("new.example.com", "192.168.15.1")
	fmt.Println("AddToScope:", s.InScope("new.example.com")) // Should return true

	// Remove specific domains and IPs from the scope
	s.RemoveFromScope("new.example.com")
	fmt.Println("RemoveFromScope:", s.InScope("new.example.com")) // Should return false

	// Retrieve all hosts in the scope as a slice
	hosts := s.HostSlice()
	fmt.Println("GetHostsAsSlice:", hosts)

	// Print scope details as strings
	fmt.Println(s.IncludeSlice())
	fmt.Println(s.ExcludeSlice())
	fmt.Println(s.HostSlice())
}
