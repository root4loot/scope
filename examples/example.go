package main

import (
	"fmt"

	"github.com/root4loot/goscope"
)

func main() {
	s := goscope.NewScope()

	// Adding includes
	s.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com", "example2.com:8080", "*.example.*.test")

	// Adding excludes
	s.AddExclude("somedomain.com", "exclude.example.com", "192.168.0.6")

	// Testing with single IP
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.2")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.7")) // Returns false
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.6")) // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.1")) // Returns false
	fmt.Println("InScope:", s.InScope("192.168.0.2"))       // Returns true
	fmt.Println("InScope:", s.InScope("192.168.0.7"))       // Returns false

	// Testing with IP-range
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.4")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.0.8")) // Returns false
	fmt.Println("IsExcluded:", s.IsExcluded("192.168.0.6")) // Returns true
	fmt.Println("InScope:", s.InScope("192.168.0.4"))       // Returns true
	fmt.Println("InScope:", s.InScope("192.168.0.8"))       // Returns false

	// Testing with IP-CIDR
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.10.50")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("192.168.11.50")) // Returns false
	fmt.Println("InScope:", s.InScope("192.168.10.50"))       // Returns true
	fmt.Println("InScope:", s.InScope("192.168.11.50"))       // Returns false

	// Testing with domain
	fmt.Println("IsIncluded:", s.IsIncluded("foo.example.com"))     // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("bar.otherdomain.com")) // Returns false
	fmt.Println("IsExcluded:", s.IsExcluded("exclude.example.com")) // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("include.example.com")) // Returns false
	fmt.Println("InScope:", s.InScope("foo.example.com"))           // Returns true
	fmt.Println("InScope:", s.InScope("bar.otherdomain.com"))       // Returns false

	// Testing with domain and scheme
	fmt.Println("InScope with scheme:", s.InScope("http://foo.example.com"))      // Returns true
	fmt.Println("InScope with scheme:", s.InScope("https://bar.otherdomain.com")) // Returns false

	// Testing if subdomain is excluded due to parent domain being in excludes list
	fmt.Println("IsExcluded:", s.IsExcluded("sub.somedomain.com"))  // Returns true
	fmt.Println("IsExcluded:", s.IsExcluded("sub.otherdomain.com")) // Returns false
	fmt.Println("InScope:", s.InScope("sub.somedomain.com"))        // Returns false
	fmt.Println("InScope:", s.InScope("sub.otherdomain.com"))       // Returns false

	// Testing with port number
	fmt.Println("IsIncluded:", s.IsIncluded("example2.com:8080")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("example2.com:1234")) // Returns false
	fmt.Println("InScope:", s.InScope("example2.com:8080"))       // Returns true
	fmt.Println("InScope:", s.InScope("example2.com:1234"))       // Returns false

	// Testing against wildcard
	fmt.Println("IsIncluded:", s.IsIncluded("foo.example.bar.test")) // Returns true
	fmt.Println("IsIncluded:", s.IsIncluded("foo.bar.baz.test"))     // Returns false
	fmt.Println("InScope:", s.InScope("foo.example.bar.test"))       // Returns true
	fmt.Println("InScope:", s.InScope("foo.bar.baz.test"))           // Returns false
}
