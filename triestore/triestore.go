package triestore

import (
	"sync"
)

type ValueGuard struct {
	Value interface{}
	TrieStore *TrieStore
}

type Node struct {
	Value interface{}
	Children map[rune]*Node
}

type TrieStore struct {
	mutex sync.RWMutex
	Root *Node
}

func New() *TrieStore {
	return &TrieStore{Root: &Node{Children: make(map[rune]*Node)}}
}


// Main API

func (t *TrieStore) Put(key string, value interface{}) {
	if len(key) == 0 || value == nil { return }

	// exclusive write (wrlock)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	node := t.Root
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
}

func (t *TrieStore) Get(key string) *ValueGuard {
	if len(key) == 0 { return &ValueGuard{TrieStore: t} }

	// multiple reads (rlock)
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	node := t.Root
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.Children[ch]
		if ok != true {
			return nil
		}
		node = nextN
	}
	return &ValueGuard{Value: node.Value, TrieStore: t}
}

func (t *TrieStore) Remove(key string) {
	if len(key) == 0 { return }
	// exclusive write (wrlock)
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.removeHelper(t.Root, key)
}


// Helper functions

func (t *TrieStore) removeHelper(node *Node, key string) bool {
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
