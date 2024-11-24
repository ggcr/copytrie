package trie

import (
	"github.com/brunoga/deep"
)

type Node struct {
	value interface{}
	children map[rune]*Node
}

type Trie struct {
	root *Node
}

func New() *Trie {
	return &Trie{root: &Node{children: make(map[rune]*Node)}}
}


// Main API

func (t *Trie) Put(key string, value interface{}) *Trie {
	if len(key) == 0 || value == nil { return t }
	newroot := deep.MustCopy(t.root)
	node := newroot
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.children[ch]
		if ok != true {
			nextN = &Node{children: make(map[rune]*Node)}
			node.children[ch] = nextN
		}
		node = nextN
	}
	node.value = value
	return &Trie{root: newroot}
}

func (t *Trie) Get(key string) interface{} {
	if len(key) == 0 { return nil }
	node := t.root
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.children[ch]
		if ok != true {
			return nil
		}
		node = nextN
	}
	return node.value
}

func (t *Trie) Remove(key string) *Trie {
	cp_t := &Trie{root: deep.MustCopy(t.root)}
	cp_t.removeHelper(cp_t.root, key)
	return cp_t
}


// Helper functions

func (t *Trie) removeHelper(node *Node, key string) bool {
	if len(key) == 0 {
		node.value = nil
		return len(node.children) == 0
	} else {
		ch := rune(key[0])
		if nextN, ok := node.children[ch]; ok {
			ret := t.removeHelper(nextN, key[1:])
			if ret == true {
				delete(node.children, ch)
			}
			return ret == true && len(node.children) == 0
		}
		return false
	}
}

