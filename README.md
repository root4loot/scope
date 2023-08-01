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
	s := goscope.NewScope()

	s.AddInclude("192.168.0.1-5", "192.168.10/24")
	s.AddInclude("*.example.com")
	s.AddInclude("example2.com:8080")
	s.AddInclude("*.example.*.test")
	s.AddExclude("exclude.example.com")

	fmt.Println(s.InScope("192.168.0.2"))
	fmt.Println(s.InScope("192.168.0.6"))
	fmt.Println(s.InScope("192.168.10.50"))
	fmt.Println(s.InScope("foo.example.com"))
	fmt.Println(s.InScope("example2.com:8080"))
	fmt.Println(s.InScope("example2.com:1234"))
	fmt.Println(s.InScope("foo.example.bar.test"))
}

// Output:
// true
// false
// false
// true
// true
// false
// true
```

## Contributing

Contributions to goscope are welcome. If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.