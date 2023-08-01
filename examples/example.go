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
