package main

import (
	"fmt"

	"github.com/root4loot/scope"
)

func main() {
	sc := scope.NewScope()

	err := sc.AddIncludes([]string{
		"example.com",
		"example.com:8080",
		"sub.example.com",
		"another.example.com",
		"deep.sub.example.com",
		"http://example.com",
		"https://example.com",
		"https://example.com:443",
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"192.168.3.2-5",
		"192.168.2.0/24",
	})
	if err != nil {
		fmt.Printf("Error adding includes: %v\n", err)
		return
	}

	err = sc.AddExcludes([]string{
		"https://example.com:8080",
		"example.com:9090",
		"http://192.168.1.1",
		"192.168.2.0/24",
	})
	if err != nil {
		fmt.Printf("Error adding excludes: %v\n", err)
		return
	}

	testURLs := []string{
		"example.com",
		"example.com:8080",
		"example.com:9090",
		"sub.example.com",
		"another.example.com",
		"deep.sub.example.com",
		"http://example.com",
		"https://example.com",
		"https://example.com:443",
		"https://example.com:8080",
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"192.168.3.3",
	}

	for _, url := range testURLs {
		inScope := sc.IsInScope(url)
		fmt.Printf("URL '%s' is in scope: %v\n", url, inScope)
	}

	activeScope := sc.GetScope()
	fmt.Printf("Active scope: %v\n", activeScope)
}
