# goscope

`goscope` is a lightweight Go library for managing scopes of hosts and IPv4 addresses, ideal for penetration testing tools and similar Go programs.

## Installation

```bash
go get github.com/root4loot/goscope@master
```

## Example

```go
package main

import (
	"fmt"

	"github.com/root4loot/goscope"
)

func main() {
	// Create a new scope instance
	myScope := goscope.NewScope()

	// Define includes and excludes in the scope
	_ = myScope.AddInclude("192.168.0.1-5", "192.168.10.0/24", "*.example.com")
	_ = myScope.AddExclude("somedomain.com", "192.168.0.6")

	// Evaluate a single IP
	fmt.Println("IsIncluded:", myScope.IsIncluded("192.168.0.2"))   // Expect true
	fmt.Println("IsExcluded:", myScope.IsExcluded("192.168.0.6"))   // Expect true
	fmt.Println("InScope:", myScope.IsTargetInScope("192.168.0.2")) // Expect true

	// Evaluate an IP within a range
	fmt.Println("IsIncluded:", myScope.IsIncluded("192.168.0.4")) // Expect true

	// Evaluate an IP within a CIDR
	fmt.Println("IsIncluded:", myScope.IsIncluded("192.168.10.50")) // Expect true

	// Evaluate a domain
	fmt.Println("IsIncluded:", myScope.IsIncluded("foo.example.com")) // Expect true
	fmt.Println("IsExcluded:", myScope.IsExcluded("somedomain.com"))  // Expect true

	// Add a new host to the scope
	fmt.Println("Adding a new host to the scope")
	if err := myScope.AddTargetToScope("newHost.com"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Successfully added newHost.com to scope")
	}

	// Loop includes
	fmt.Println("Current Includes:")
	for _, include := range myScope.GetIncludes() {
		fmt.Println(" - ", include)
	}

	// Loop excludes
	fmt.Println("Current Excludes:")
	for _, exclude := range myScope.GetExcludes() {
		fmt.Println(" - ", exclude)
	}

	// Loop hosts
	fmt.Println("Current Hosts:")
	for _, host := range myScope.GetTargets() {
		fmt.Println(" - ", host)
	}

	// Loop hosts and their types
	for host, hostType := range myScope.GetTargetsAndTypeMap() {
		fmt.Printf("Host: %s, Type: %s\n", host, hostType.String()) // Host: newHost.com, Type: Domain
	}

	// Remove host from scope
	fmt.Println("Removing a host from the scope")
	if err := myScope.RemoveTargetFromScope("newHost.com"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Successfully removed newHost.com from scope")
	}
}

```

## Contributing

Contributions to goscope are welcome. If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.