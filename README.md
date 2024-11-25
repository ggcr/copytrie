# TrieStore

A persistent trie (prefix tree) key-value storage system implementation in Go that supports copy-on-write and concurrent versions. In both versions an instance of a trie is returned thus ensuring the ability of back-tracking to an older state of the trie.


## Install
`go get -u github.com/ggcr/triestore`

## Copy-on-Write Trie (Sequential)

```go
package main

import (
    "fmt"
    "github.com/ggcr/triestore/copytrie" // Import
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
```

## TrieStore (Concurrent)

```go
package main

import (
    "fmt"
    "sync"
    "github.com/ggcr/triestore"
)

func main() {
    store := triestore.New()

    var wg sync.WaitGroup
    const keysPerThread = 10000
    const numThreads = 4

    for tid := 0; tid < numThreads; tid++ {
        wg.Add(1)
        go func(tid int) {
            defer wg.Done()

            // Phase 1: Put
            for i := 0; i < keysPerThread; i++ {
                key := fmt.Sprintf("%05d", i*4+tid)
                value := fmt.Sprintf("value-%08d", i*4+tid)
                store.Put(key, value)
            }

            // Phase 2: Remove
            for i := 0; i < keysPerThread; i++ {
                key := fmt.Sprintf("%05d", i*4+tid)
                store.Remove(key)
            }

            // Phase 3: Update
            for i := 0; i < keysPerThread; i++ {
                key := fmt.Sprintf("%05d", i*4+tid)
                value := fmt.Sprintf("new-value-%08d", i*4+tid)
                store.Put(key, value)
            }
        }(tid)
    }
    wg.Wait()

    fmt.Println(store.Get("01234").Value) // {TrieStore t*, "new-value-00001234"}
    fmt.Println(store.Get("nonexistent-key")) // <nil>
}
```

## Tests
You can clone it locally and run tests if needed:
```bash
git clone https://github.com/ggcr/triestore.git
cd copytrie
go test ./... -v
```

Check out the tests output:

```
?       github.com/ggcr/triestore/copytrie      [no test files]
?       github.com/ggcr/triestore/triestore     [no test files]
=== RUN   TestTrieConstruct
--- PASS: TestTrieConstruct (0.00s)
=== RUN   TestTrieEmptyKey
--- PASS: TestTrieEmptyKey (0.00s)
=== RUN   TestTrieEmptyVal
--- PASS: TestTrieEmptyVal (0.00s)
=== RUN   TestTrieNilVal
--- PASS: TestTrieNilVal (0.00s)
=== RUN   TestTrieBasicPut
--- PASS: TestTrieBasicPut (0.00s)
=== RUN   TestTrieBasicPutRev
--- PASS: TestTrieBasicPutRev (0.00s)
=== RUN   TestTriePutUpdate
--- PASS: TestTriePutUpdate (0.00s)
=== RUN   TestTrieNonExistentKey
--- PASS: TestTrieNonExistentKey (0.00s)
=== RUN   TestTrieDifferentValueTypes
--- PASS: TestTrieDifferentValueTypes (0.00s)
=== RUN   TestTrieOverwrite
--- PASS: TestTrieOverwrite (0.00s)
=== RUN   TestTrieSharedPrefix
--- PASS: TestTrieSharedPrefix (0.00s)
=== RUN   TestTrieRemoval
--- PASS: TestTrieRemoval (0.00s)
=== RUN   TestTrieRemovalNodePresent
--- PASS: TestTrieRemovalNodePresent (0.00s)
=== RUN   TestTrieCopyPutOps
--- PASS: TestTrieCopyPutOps (0.00s)
=== RUN   TestTrieCopyRemoveOps
--- PASS: TestTrieCopyRemoveOps (0.00s)
=== RUN   TestTrieExample
--- PASS: TestTrieExample (0.00s)
=== RUN   TestTrieSequential
--- PASS: TestTrieSequential (0.00s)
=== RUN   TestTrieStoreConstruct
--- PASS: TestTrieStoreConstruct (0.00s)
=== RUN   TestTrieBasic
--- PASS: TestTrieBasic (0.00s)
=== RUN   TestTrieGuard
--- PASS: TestTrieGuard (0.00s)
=== RUN   TestTrieStoreSequential
--- PASS: TestTrieStoreSequential (0.06s)
=== RUN   TestTrieStoreConcurrentPhases
--- PASS: TestTrieStoreConcurrentPhases (0.11s)
=== RUN   TestTrieStoreMixedConcurrent
--- PASS: TestTrieStoreMixedConcurrent (0.10s)
PASS
ok      github.com/ggcr/triestore/test
```
</details>
