package main

import (
	"fmt"
	// Local
	"github.com/ggcr/trie-go/internal/trie"
)

func main() {
	t := trie.New()
	t = t.Put("hello", 32)
	t = t.Put("hell", 15)
	val := t.Get("hell")
	fmt.Println(val)
}
