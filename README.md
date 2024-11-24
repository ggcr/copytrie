# Copy-Trie

A persistent trie (prefix tree) implementation in Go that uses copy-on-write. Each modification operation returns a new trie instance while preserving the original structure.

## Install
`go get -u github.com/ggcr/copytrie/trie@v0.1.0`

## Usage

```go
package main

import (
	 "fmt"
	 "github.com/ggcr/copytrie/trie" // Import
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

## Tests
You can clone it locally and run tests if needed:
```bash
git clone https://github.com/yourusername/copytrie.git
cd copytrie
go test ./... -v
```

Check out the tests output:

  ```bash
go test ./... -v
?       github.com/ggcr/copytrie/cmd    [no test files]
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
PASS
ok      github.com/ggcr/copytrie/trie
```
</details>
