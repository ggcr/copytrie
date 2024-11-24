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

func TestTrieCopyPutOps(t *testing.T) {
	trie := New()

	trie1 := trie.Put("hello", 16)
	trie2 := trie1.Put("hello", 32)
	trie3 := trie2.Put("hello", 64)

	// Validate that we do not modify trie by ref
	val1 := trie1.Get("hello")
	val2 := trie2.Get("hello")
	val3 := trie3.Get("hello")
	if val1 != 16 || val2 != 32 || val3 != 64 {
		t.Errorf("Values got overwritten on different tries. Got val1=%v, val2=%v, val3=%v", val1, val2, val3)
	}

	// New key should not be present in original trie instances
	trie4 := trie3.Put("test", "hello")
	val1 = trie1.Get("test")
	val2 = trie2.Get("test")
	val3 = trie3.Get("test")
	if val1 != nil || val2 != nil || val3 != nil {
		t.Errorf("Earlier tries have access to keys added in newer tries.")
	}

	// New key should be present on the new trie instance
	val4 := trie4.Get("test")
	if val4 != "hello" {
		t.Errorf("Expected 'hello', got %v", val4)
	}
}

func TestTrieCopyRemoveOps(t *testing.T) {
	trie := New()

	trie1 := trie.Put("hello", 16)
	trie2 := trie1.Put("world", 32)
	trie3 := trie2.Put("golang", 64)

	trie4 := trie3.Remove("world")

	if trie4.Get("world") != nil {
		t.Errorf("'world' should be removed in trie4.")
	}
	if trie2.Get("world") == nil {
		t.Errorf("'world' should still exist in trie2.")
	}

	val1 := trie1.Get("hello")
	val2 := trie2.Get("hello")
	val3 := trie3.Get("golang")
	if val1 != 16 || val2 != 16 || val3 != 64 {
		t.Errorf("Values got overwritten unexpectedly. val1=%v, val2=%v, val3=%v", val1, val2, val3)
	}

	if trie4.Get("golang") != 64 {
		t.Errorf("'golang' should still exist in trie4.")
	}

	trie5 := trie3.Remove("hello")

	if trie5.Get("hello") != nil {
		t.Errorf("'hello' should be removed in trie5.")
	}
	if trie3.Get("hello") == nil {
		t.Errorf("'hello' should still exist in trie3.")
	}
}
