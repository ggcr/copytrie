package triestore_test

import (
	"testing"
	"fmt"
	"github.com/ggcr/triestore/triestore"
	"sync"
)

func TestTrieStoreConstruct(t *testing.T) {
	store := triestore.New()
	Root := store.Root
	if Root.Value != nil || len(Root.Children) != 0 {
		t.Errorf("triestore.New a initialized without Root being <nil>.")
	}
}

func TestTrieBasic(t *testing.T) {
	store := triestore.New()
	val := store.Get("233")
	if val != nil {
		t.Errorf("Expected Trie to be empty, got value %v instead.", val)
	}
	store.Put("233", 2333)
	val = store.Get("233")
	if val == nil || val.Value != 2333 {
		t.Errorf("Expected key '233' to contain 2333, got %v instead.", val)
	}
	store.Remove("233")
	val = store.Get("233")
	if val != nil {
		t.Errorf("Expected Trie to be empty after removal, got value %v instead.", val)
	}
}

func TestTrieGuard(t *testing.T) {
	store := triestore.New()
	val := store.Get("233")
	if val != nil {
		t.Errorf("Expected Trie to be empty, got value %v instead.", val)
	}
	store.Put("233", 2333)
	vg1 := store.Get("233")
	if vg1 == nil || vg1.Value != 2333 {
		t.Errorf("Expected key '233' to contain 2333, got guard value %+v instead.", vg1)
	}
	store.Remove("233")
	vg2 := store.Get("233")
	if vg2 != nil {
		t.Errorf("Expected Trie to be empty after removal, got guard value %+v instead.", vg2)
	}
	// we can still access previous instances of the trie through our valueguard
	if vg1.Value != 2333 {
		t.Errorf("Expected guard to retain value 2333 after removal, got %+v instead.", vg1.Value)
	}
}

func TestTrieStoreSequential(t *testing.T) {
	store := triestore.New()
	totalTh := 67000
	
	for i := 0; i < totalTh; i++ { // Phase 1: puts
		key := fmt.Sprintf("%05d", i)
		value := fmt.Sprintf("value-%08d", i)
		store.Put(key, value)
	}
	for i := 0; i < totalTh; i+=2 { // Phase 2: update even keys
		key := fmt.Sprintf("%05d", i)
		value := fmt.Sprintf("new-value-%08d", i)
		store.Put(key, value)
	}
	for i := 0; i < totalTh; i+=3 { // Phase 3: remove triplet keys
		key := fmt.Sprintf("%05d", i)
		store.Remove(key)
	}

	// verify
	for i := 0; i < totalTh; i++ {
		key := fmt.Sprintf("%05d", i)
		if i % 3 == 0 { // key should not be on the trie
			if vg1 := store.Get(key); vg1 != nil {
				t.Errorf("Expected Trie to be empty after removal, got guard value %+v instead.", vg1)
			}
		} else if i % 2 == 0 { // key should contain the updated version
			exp := fmt.Sprintf("new-value-%08d", i)
			if vg2 := store.Get(key); vg2 == nil || vg2.Value != exp {
				t.Errorf("Expected %+v, got guard value %+v instead.", exp, vg2)
			}
		} else { // key should contain the original value
			exp := fmt.Sprintf("value-%08d", i)
			if vg3 := store.Get(key); vg3 == nil || vg3.Value != exp {
				t.Errorf("Expected %+v, got guard value %+v instead.", exp, vg3)
			}
		}
	}
}

func TestTrieStoreConcurrentPhases(t *testing.T) {
	store := triestore.New()
	var wg sync.WaitGroup
	totalTh := 67000

	// Phase 1: put values
	for i := 0; i < totalTh; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%05d", i)
			value := fmt.Sprintf("value-%08d", i)
			store.Put(key, value)
		}(i)
	}
	wg.Wait()

	// Phase 2: update even keys
	for i := 0; i < totalTh; i += 2 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%05d", i)
			value := fmt.Sprintf("new-value-%08d", i)
			store.Put(key, value)
		}(i)
	}
	wg.Wait()

	// Phase 3: remove triplet keys
	for i := 0; i < totalTh; i += 3 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%05d", i)
			store.Remove(key)
		}(i)
	}
	wg.Wait()

	// Verify trie 
	for i := 0; i < totalTh; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("%05d", i)
			if i%3 == 0 { // key should not be in the trie
				if vg1 := store.Get(key); vg1 != nil {
					t.Errorf("Expected Trie to be empty after removal, got guard value %+v instead.", vg1)
				}
			} else if i%2 == 0 { // key should contain the updated version
				exp := fmt.Sprintf("new-value-%08d", i)
				if vg2 := store.Get(key); vg2 == nil || vg2.Value != exp {
					t.Errorf("Expected %+v, got guard value %+v instead.", exp, vg2)
				}
			} else { // key should contain the original value
				exp := fmt.Sprintf("value-%08d", i)
				if vg3 := store.Get(key); vg3 == nil || vg3.Value != exp {
					t.Errorf("Expected %+v, got guard value %+v instead.", exp, vg3)
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestTrieStoreMixedConcurrent(t *testing.T) {
	store := triestore.New()

	var wg sync.WaitGroup
	const keysPerThread = 1000
	const numThreads = 67

	for tid := 0; tid < numThreads; tid++ {
		wg.Add(1)
		go func(tid int) {
			defer wg.Done()

			// Phase 1: Put
			for i := 0; i < keysPerThread; i++ {
				key := fmt.Sprintf("%05d", i*numThreads+tid)
				value := fmt.Sprintf("value-%08d", i*numThreads+tid)
				store.Put(key, value)
			}

			// Phase 2: Remove
			for i := 0; i < keysPerThread; i++ {
				key := fmt.Sprintf("%05d", i*numThreads+tid)
				store.Remove(key)
			}

			// Phase 3: Update
			for i := 0; i < keysPerThread; i++ {
				key := fmt.Sprintf("%05d", i*numThreads+tid)
				value := fmt.Sprintf("new-value-%08d", i*numThreads+tid)
				store.Put(key, value)
			}
		}(tid)
	}

	wg.Wait()

	fmt.Println(store.Get("01234").Value)

	for i := 0; i < keysPerThread*numThreads; i++ {
		key := fmt.Sprintf("%05d", i)
		exp := fmt.Sprintf("new-value-%08d", i)
		vg := store.Get(key)
		if vg == nil || vg.Value != exp {
			t.Errorf("Expected %s, but got %+v", exp, vg)
		}
	}
}
