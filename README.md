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
	// Initialize a new Scope instance
	sc := scope.NewScope()

	// Add includes
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
	fmt.Println(sc.IsInScope("10.0.0.1"))           // Output: false

	// Get active scope
	activeScope := sc.GetScope()
	fmt.Printf("Active scope: %v\n", activeScope)
}
```

For more detailed usage, see [example.go](https://github.com/root4loot/scope/blob/main/example/example.go)

## Contributing

Contributions to goscope are welcome. If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.