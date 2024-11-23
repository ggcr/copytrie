package trie

import (
	"testing"
)

func TestTrieConstruct(t *testing.T) {
	trie := New()
	root := trie.root
	if root.value != nil || len(root.children) != 0 {
		t.Errorf("New trie initialized without root being <nil>.")
	}
}

func TestTrieEmptyKey(t *testing.T) {
	trie := New()
	trie = trie.Put("", 32)
	val := trie.Get("")
	if val != nil {
		t.Errorf("Expected nil got %d", val)
	}
}

func TestTrieEmptyVal(t *testing.T) {
	trie := New()
	var empty interface{}
	trie = trie.Put("key", empty)
	val := trie.Get("key")
	if val != nil {
		t.Errorf("Expected nil got %d", val)
	}
}

func TestTrieNilVal(t *testing.T) {
	trie := New()
	trie = trie.Put("key", nil)
	val := trie.Get("key")
	if val != nil {
		t.Errorf("Expected nil got %d", val)
	}
}

func TestTrieBasicPut(t *testing.T) {
	trie := New()
	trie = trie.Put("hello", 32)
	trie = trie.Put("hell", 32)
	val := trie.Get("hell")
	if val != 32 {
		t.Errorf("Expected 32 got %d", val)
	}
}

func TestTrieBasicPutRev(t *testing.T) {
	trie := New()
	trie = trie.Put("hell", 32)
	trie = trie.Put("hello", 32)
	val := trie.Get("hello")
	if val != 32 {
		t.Errorf("Expected 32 got %d", val)
	}
}

func TestTriePutUpdate(t *testing.T) {
	trie := New()
	trie = trie.Put("hello", 32)
	trie = trie.Put("hello", 32)
	val := trie.Get("hello")
	if val != 32 {
		t.Errorf("Expected 32, got %d", val)
	}
}

func TestTrieNonExistentKey(t *testing.T) {
	trie := New()
	val := trie.Get("nonexistent")
	if val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestTrieDifferentValueTypes(t *testing.T) {
	trie := New()
	trie = trie.Put("key1", "string value")
	trie = trie.Put("key2", 42.42)
	trie = trie.Put("key3", true)

	val1 := trie.Get("key1")
	if val1 != "string value" {
		t.Errorf("Expected 'string value', got %v", val1)
	}

	val2 := trie.Get("key2")
	if val2 != 42.42 {
		t.Errorf("Expected 42.42, got %v", val2)
	}

	val3 := trie.Get("key3")
	if val3 != true {
		t.Errorf("Expected true, got %v", val3)
	}
}

func TestTrieOverwrite(t *testing.T) {
	trie := New()
	trie = trie.Put("overwrite", 100)
	trie = trie.Put("overwrite", 200)
	val := trie.Get("overwrite")
	if val != 200 {
		t.Errorf("Expected %d, got %v", 200, val)
	}
}

func TestTrieSharedPrefix(t *testing.T) {
	trie := New()
	trie = trie.Put("cat", "meow")
	trie = trie.Put("car", "vroom")
	val1 := trie.Get("cat")
	val2 := trie.Get("car")
	if val1 != "meow" {
		t.Errorf("Expected 'meow', got %v", val1)
	}
	if val2 != "vroom" {
		t.Errorf("Expected 'vroom', got %v", val2)
	}
}

func (t *Trie) isNodePresent(key string) bool {
	if len(key) == 0 { return false }
	node := t.root
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.children[ch]
		if ok != true {
			return false
		}
		node = nextN
	}
	return true
}

func TestTrieRemoval(t *testing.T) {
	trie := New()

	trie = trie.Put("cat", "meow")
	trie = trie.Put("car", "vroom")

	trie = trie.Remove("cat")

	val := trie.Get("cat")
	if val != nil { // check value to be nil
		t.Errorf("Expected <nil>, got %v", val)
	}

	if trie.isNodePresent("cat") != false { // check that the node has been freed
		t.Errorf("Expected 'cat' to be removed, however still exists")
	}
}

func TestTrieRemovalNodePresent(t *testing.T) {
	trie := New()

	trie = trie.Put("cat", "meow")
	trie = trie.Put("catacumba", "yes")
	trie = trie.Put("car", "vroom")

	trie = trie.Remove("cat")

	val := trie.Get("cat")
	if val != nil { // check value to be nil
		t.Errorf("Expected <nil>, got %v", val)
	}

	if trie.isNodePresent("cat") == false { // check that the node has been freed
		t.Errorf("Expected 'cat' to NOT be removed, however it does not exist")
	}
}
