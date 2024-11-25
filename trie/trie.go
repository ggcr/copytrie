package trie

import (
	"github.com/brunoga/deep"
)

type Node struct {
	Value interface{}
	Children map[rune]*Node
}

type Trie struct {
	Root *Node
}

func New() *Trie {
	return &Trie{Root: &Node{Children: make(map[rune]*Node)}}
}


// Main API

func (t *Trie) Put(key string, value interface{}) *Trie {
	if len(key) == 0 || value == nil { return t }
	newroot := deep.MustCopy(t.Root)
	node := newroot
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.Children[ch]
		if ok != true {
			nextN = &Node{Children: make(map[rune]*Node)}
			node.Children[ch] = nextN
		}
		node = nextN
	}
	node.Value = value
	return &Trie{Root: newroot}
}

func (t *Trie) Get(key string) interface{} {
	if len(key) == 0 { return nil }
	node := t.Root
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.Children[ch]
		if ok != true {
			return nil
		}
		node = nextN
	}
	return node.Value
}

func (t *Trie) Remove(key string) *Trie {
	cp_t := &Trie{Root: deep.MustCopy(t.Root)}
	cp_t.removeHelper(cp_t.Root, key)
	return cp_t
}


// Helper functions

func (t *Trie) removeHelper(node *Node, key string) bool {
	if len(key) == 0 {
		node.Value = nil
		return len(node.Children) == 0
	} else {
		ch := rune(key[0])
		if nextN, ok := node.Children[ch]; ok {
			ret := t.removeHelper(nextN, key[1:])
			if ret == true {
				delete(node.Children, ch)
			}
			return ret == true && len(node.Children) == 0
		}
		return false
	}
}

