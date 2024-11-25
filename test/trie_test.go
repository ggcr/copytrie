package trie_test

import (
	"testing"
	"github.com/ggcr/copytrie/trie"
)

func TestTrieConstruct(t *testing.T) {
	tr := trie.New()
	Root := tr.Root
	if Root.Value != nil || len(Root.Children) != 0 {
		t.Errorf("trie.New a initialized without Root being <nil>.")
	}
}

func TestTrieEmptyKey(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("", 32)
	val := tr.Get("")
	if val != nil {
		t.Errorf("Expected nil got %d", val)
	}
}

func TestTrieEmptyVal(t *testing.T) {
	tr := trie.New()
	var empty interface{}
	tr = tr.Put("key", empty)
	val := tr.Get("key")
	if val != nil {
		t.Errorf("Expected nil got %d", val)
	}
}

func TestTrieNilVal(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("key", nil)
	val := tr.Get("key")
	if val != nil {
		t.Errorf("Expected nil got %d", val)
	}
}

func TestTrieBasicPut(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("hello", 32)
	tr = tr.Put("hell", 32)
	val := tr.Get("hell")
	if val != 32 {
		t.Errorf("Expected 32 got %d", val)
	}
}

func TestTrieBasicPutRev(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("hell", 32)
	tr = tr.Put("hello", 32)
	val := tr.Get("hello")
	if val != 32 {
		t.Errorf("Expected 32 got %d", val)
	}
}

func TestTriePutUpdate(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("hello", 32)
	tr = tr.Put("hello", 32)
	val := tr.Get("hello")
	if val != 32 {
		t.Errorf("Expected 32, got %d", val)
	}
}

func TestTrieNonExistentKey(t *testing.T) {
	tr := trie.New()
	val := tr.Get("nonexistent")
	if val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestTrieDifferentValueTypes(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("key1", "string Value")
	tr = tr.Put("key2", 42.42)
	tr = tr.Put("key3", true)

	val1 := tr.Get("key1")
	if val1 != "string Value" {
		t.Errorf("Expected 'string Value', got %v", val1)
	}

	val2 := tr.Get("key2")
	if val2 != 42.42 {
		t.Errorf("Expected 42.42, got %v", val2)
	}

	val3 := tr.Get("key3")
	if val3 != true {
		t.Errorf("Expected true, got %v", val3)
	}
}

func TestTrieOverwrite(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("overwrite", 100)
	tr = tr.Put("overwrite", 200)
	val := tr.Get("overwrite")
	if val != 200 {
		t.Errorf("Expected %d, got %v", 200, val)
	}
}

func TestTrieSharedPrefix(t *testing.T) {
	tr := trie.New()
	tr = tr.Put("cat", "meow")
	tr = tr.Put("car", "vroom")
	val1 := tr.Get("cat")
	val2 := tr.Get("car")
	if val1 != "meow" {
		t.Errorf("Expected 'meow', got %v", val1)
	}
	if val2 != "vroom" {
		t.Errorf("Expected 'vroom', got %v", val2)
	}
}

func isNodePresent(root *trie.Node, key string) bool {
	if len(key) == 0 { return false }
	node := root
	for i := 0; i < len(key); i++ {
		ch := rune(key[i])
		nextN, ok := node.Children[ch]
		if ok != true {
			return false
		}
		node = nextN
	}
	return true
}

func TestTrieRemoval(t *testing.T) {
	tr := trie.New()

	tr = tr.Put("cat", "meow")
	tr = tr.Put("car", "vroom")

	tr = tr.Remove("cat")

	val := tr.Get("cat")
	if val != nil { // check Value to be nil
		t.Errorf("Expected <nil>, got %v", val)
	}

	if isNodePresent(tr.Root, "cat") != false { // check that the node has been freed
		t.Errorf("Expected 'cat' to be removed, however still exists")
	}
}

func TestTrieRemovalNodePresent(t *testing.T) {
	tr := trie.New()

	tr = tr.Put("cat", "meow")
	tr = tr.Put("catacumba", "yes")
	tr = tr.Put("car", "vroom")

	tr = tr.Remove("cat")

	val := tr.Get("cat")
	if val != nil { // check Value to be nil
		t.Errorf("Expected <nil>, got %v", val)
	}

	if isNodePresent(tr.Root, "cat") == false { // check that the node has been freed
		t.Errorf("Expected 'cat' to NOT be removed, however it does not exist")
	}
}

func TestTrieCopyPutOps(t *testing.T) {
	tr := trie.New()

	tr1 := tr.Put("hello", 16)
	tr2 := tr1.Put("hello", 32)
	tr3 := tr2.Put("hello", 64)

	// Validate that we do not modify tr by ref
	val1 := tr1.Get("hello")
	val2 := tr2.Get("hello")
	val3 := tr3.Get("hello")
	if val1 != 16 || val2 != 32 || val3 != 64 {
		t.Errorf("Values got overwritten on different trs. Got val1=%v, val2=%v, val3=%v", val1, val2, val3)
	}

	// trie.New key should not be present in original tr instances
	tr4 := tr3.Put("test", "hello")
	val1 = tr1.Get("test")
	val2 = tr2.Get("test")
	val3 = tr3.Get("test")
	if val1 != nil || val2 != nil || val3 != nil {
		t.Errorf("Earlier trs have access to keys added in newer trs.")
	}

	// trie.New key should be present on the new tr instance
	val4 := tr4.Get("test")
	if val4 != "hello" {
		t.Errorf("Expected 'hello', got %v", val4)
	}
}

func TestTrieCopyRemoveOps(t *testing.T) {
	tr := trie.New()

	tr1 := tr.Put("hello", 16)
	tr2 := tr1.Put("world", 32)
	tr3 := tr2.Put("golang", 64)

	tr4 := tr3.Remove("world")

	if tr4.Get("world") != nil {
		t.Errorf("'world' should be removed in tr4.")
	}
	if tr2.Get("world") == nil {
		t.Errorf("'world' should still exist in tr2.")
	}

	val1 := tr1.Get("hello")
	val2 := tr2.Get("hello")
	val3 := tr3.Get("golang")
	if val1 != 16 || val2 != 16 || val3 != 64 {
		t.Errorf("Values got overwritten unexpectedly. val1=%v, val2=%v, val3=%v", val1, val2, val3)
	}

	if tr4.Get("golang") != 64 {
		t.Errorf("'golang' should still exist in tr4.")
	}

	tr5 := tr3.Remove("hello")

	if tr5.Get("hello") != nil {
		t.Errorf("'hello' should be removed in tr5.")
	}
	if tr3.Get("hello") == nil {
		t.Errorf("'hello' should still exist in tr3.")
	}
}

func TestTrieExample(t *testing.T) {
	t0 := trie.New()
	t1 := t0.Put("test", "hello")     // Put
	t2 := t1.Put("test", 32)          // Update
	t3 := t2.Put("testing", "works")  // trie.New key
	t4 := t3.Remove("test")           // Remove

	if val := t1.Get("test"); val != "hello" {
		t.Errorf("t1: Expected 'hello', got %v", val)
	}

	if val := t2.Get("test"); val != 32 {
		t.Errorf("t2: Expected 32, got %v", val)
	}

	if val := t3.Get("testing"); val != "works" {
		t.Errorf("t3: Expected 'works', got %v", val)
	}

	if val := t3.Get("test"); val != 32 {
		t.Errorf("t3: Expected 32 for 'test', got %v", val)
	}

	if val := t4.Get("test"); val != nil {
		t.Errorf("t4: Expected nil for removed key 'test', got %v", val)
	}

	if val := t4.Get("testing"); val != "works" {
		t.Errorf("t4: Expected 'works' for 'testing', got %v", val)
	}
}
