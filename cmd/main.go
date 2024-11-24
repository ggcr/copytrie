package main

import (
	"fmt"
	// Local
	"github.com/ggcr/copytrie/internal/trie"
)

func main() {
    t0 := trie.New()

    t1 := t0.Put("test", "hello")     // Put
    t2 := t1.Put("test", 32)          // Update
    t3 := t2.Put("testing", "works")  // New key
    t4 := t3.Remove("test")           // Remove

    // Original values are preserved
    fmt.Println(t1.Get("test"))       // "hello"
    fmt.Println(t2.Get("test"))       // 32
    fmt.Println(t3.Get("testing"))    // "works"
    fmt.Println(t4.Get("test"))       // nil
	 fmt.Println(t4.Get("testing"))    // "works"
}
