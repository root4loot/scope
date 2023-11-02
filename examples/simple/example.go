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
	fmt.Println("IsIncluded:", myScope.IsTargetExcluded("192.168.0.2")) // Expect true
	fmt.Println("IsExcluded:", myScope.IsTargetExcluded("192.168.0.6")) // Expect true
	fmt.Println("InScope:", myScope.IsTargetInScope("192.168.0.2"))     // Expect true

	// Evaluate an IP within a range
	fmt.Println("IsIncluded:", myScope.IsTargetIncluded("192.168.0.4")) // Expect true

	// Evaluate an IP within a CIDR
	fmt.Println("IsIncluded:", myScope.IsTargetIncluded("192.168.10.50")) // Expect true

	// Evaluate a domain
	fmt.Println("IsIncluded:", myScope.IsTargetIncluded("foo.example.com")) // Expect true
	fmt.Println("IsExcluded:", myScope.IsTargetExcluded("somedomain.com"))  // Expect true

	// Add a new target to the scope
	fmt.Println("Adding a new target to the scope")
	if err := myScope.AddTargetToScope("newhost.com"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Successfully added newhost.com to scope")
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
	if err := myScope.RemoveTargetFromScope("newhost.com"); err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Successfully removed newhost.com from scope")
	}
}
