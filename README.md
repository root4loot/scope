# scope

Simple library for managing scopes of hosts and IPv4 addresses, ideal for penetration testing tools and other network-related Go applications.

## Installation

```bash
go get github.com/root4loot/scope@latest
```

## Simple Usage

```go
package main

import (
	"fmt"
	"github.com/root4loot/scope"
)

func main() {
	// Initialize a new scope instance
	sc := scope.NewScope()

	// Explicit scope mode: Only defined includes are in scope
	sc.AddInclude("example.com")
	sc.AddInclude("192.168.1.1")

	// Add excludes
	sc.AddExclude("example.com:8080")
	sc.AddExclude("http://192.168.1.1")

	// Check if a domain is in scope
	fmt.Println(sc.IsInScope("example.com"))       // Output: true
	fmt.Println(sc.IsInScope("example.com:8080"))  // Output: false

	// Check if an IP is in scope
	fmt.Println(sc.IsInScope("192.168.1.1"))        // Output: false
	fmt.Println(sc.IsInScope("10.0.0.1"))           // Output: true

	// Implicit scope mode: Everything is in scope unless explicitly excluded
	sc2 := scope.NewScope()
	sc2.AddExclude("blocked.com")
	fmt.Println(sc2.IsInScope("random.com"))       // Output: true
	fmt.Println(sc2.IsInScope("blocked.com"))     // Output: false

	// Get active scope
	activeScope := sc.GetScope()
	fmt.Printf("Active scope: %v\n", activeScope)
}
```

### **Scope Behavior**
- If no includes are defined, **everything is in scope** unless explicitly excluded.
- If includes are defined, only those targets are in scope, and everything else is out of scope by default.
- Exclusions always take priority over inclusions.