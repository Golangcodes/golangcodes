![](https://www.notion.so/icons/cursor_green.svg)

# Maps

# Maps in Go: From Basics to Advanced

Maps are one of Go's most useful built-in data types, providing efficient key-value storage. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Map Declaration](Maps%202091ccc344b180d9860dc3b035e6b792.html)

2.  [Map Initialization](Maps%202091ccc344b180d9860dc3b035e6b792.html)

3.  [Map Operations](Maps%202091ccc344b180d9860dc3b035e6b792.html)

4.  [Map Internals](Maps%202091ccc344b180d9860dc3b035e6b792.html)

5.  [Common Operations](Maps%202091ccc344b180d9860dc3b035e6b792.html)

6.  [Advanced Patterns](Maps%202091ccc344b180d9860dc3b035e6b792.html)

7.  [Performance Considerations](Maps%202091ccc344b180d9860dc3b035e6b792.html)

8.  [Best Practices](Maps%202091ccc344b180d9860dc3b035e6b792.html)

## Basic Map Declaration

A map is an unordered collection of key-value pairs:

```Go
var m map[string]int  // nil map (can't be used directly)
m := make(map[string]int)  // initialized map
```

## Map Initialization

Multiple ways to create maps:

```Go
// 1. Using make
ages := make(map[string]int)

// 2. Literal initialization
colors := map[string]string{
    "red":   "#ff0000",
    "green": "#00ff00",
}

// 3. Nil map (must be initialized before use)
var counts map[string]int
counts = make(map[string]int)  // initialization
```

## Map Operations

### Inserting/Updating Elements

```Go
m["key"] = 42      // Insert or update
m["nonexistent"]++ // Increment (initializes to 0 if key doesn't exist)
```

### Accessing Elements

```Go
value := m["key"]  // Returns zero value if key doesn't exist
value, exists := m["key"]  // exists is bool indicating presence
```

### Deleting Elements

```Go
delete(m, "key")  // No-op if key doesn't exist
```

### Checking Existence

```Go
if value, ok := m["key"]; ok {
    // key exists
} else {
    // key doesn't exist
}
```

### Length

```Go
size := len(m)  // Number of key-value pairs
```

## Map Internals

Go maps are implemented as hash tables with:

*   Average O(1) time complexity for operations

*   Automatic growth when load factor increases

*   Non-deterministic iteration order

Important characteristics:

*   Keys must be comparable (using == and !=)

*   Slices, functions, and other maps cannot be keys

*   Values can be any type

## Common Operations

### Iteration

```Go
for key, value := range m {
    fmt.Printf("%s: %v\\n", key, value)
}

// Keys only
for key := range m {
    fmt.Println(key)
}
```

### Map as Set

```Go
set := make(map[string]bool)
set["item1"] = true
set["item2"] = true

if set["item1"] {
    fmt.Println("item1 exists")
}
```

### Default Values

```Go
// Count word frequencies
wordCount := make(map[string]int)
for _, word := range words {
    wordCount[word]++  // Automatically initializes to 0
}
```

## Advanced Patterns

### Nested Maps

```Go
graph := make(map[string]map[string]int)
graph["A"] = make(map[string]int)
graph["A"]["B"] = 5  // A -> B with weight 5
```

### Concurrent Maps

```Go
var mu sync.RWMutex
var safeMap = make(map[string]int)

// Write
mu.Lock()
safeMap["key"] = 42
mu.Unlock()

// Read
mu.RLock()
value := safeMap["key"]
mu.RUnlock()
```

### Functional Operations

```Go
func mapValues(m map[string]int, fn func(int) int) map[string]int {
    result := make(map[string]int)
    for k, v := range m {
        result[k] = fn(v)
    }
    return result
}

squared := mapValues(ages, func(x int) int { return x * x })
```

### Inverted Map

```Go
func invertMap(m map[string]int) map[int]string {
    inverted := make(map[int]string)
    for k, v := range m {
        inverted[v] = k
    }
    return inverted
}
```

## Performance Considerations

1.  **Pre-allocate size** for large maps:
    
    ```Go
    m := make(map[string]int, 1000)  // Hints initial capacity
    ```

2.  **Small maps** (≤8 elements) are more efficient when initialized with literals

3.  **Struct values** can be more efficient than pointer values:
    
    ```Go
    // Often better than map[string]*MyStruct
    m := make(map[string]MyStruct)
    ```

4.  **Concurrent access** requires synchronization (sync.Mutex or sync.RWMutex)

5.  **Memory leaks** can occur if maps grow large and never shrink (consider recreating if needed)

## Best Practices

1.  **Check existence** when zero value might be valid

2.  **Document key requirements** for custom key types

3.  **Consider alternatives** for performance-critical code (arrays, slices)

4.  **Use struct keys** for composite keys:
    
    ```Go
    type Point struct { X, Y int }
    m := make(map[Point]string)
    ```

5.  **Avoid global maps** when possible

6.  **Handle concurrent access** properly

7.  **Keep maps small** when serializing (JSON, etc.)

8.  **Use meaningful names** for map variables

9.  **Consider sync.Map** for specific concurrent patterns

10.  **Benchmark** map operations in critical paths

## Complete Example

```Go
package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func main() {
	// Basic map operations
	wordCount := make(map[string]int)
	text := "hello world hello go world go hello"

	for _, word := range strings.Fields(text) {
		wordCount[word]++
	}

	fmt.Println("Word counts:")
	for word, count := range wordCount {
		fmt.Printf("%s: %d\\n", word, count)
	}

	// Checking existence
	if count, exists := wordCount["nonexistent"]; exists {
		fmt.Println("Count:", count)
	} else {
		fmt.Println("Word not found")
	}

	// Ordered iteration
	fmt.Println("\\nSorted word counts:")
	words := make([]string, 0, len(wordCount))
	for word := range wordCount {
		words = append(words, word)
	}
	sort.Strings(words)
	for _, word := range words {
		fmt.Printf("%s: %d\\n", word, wordCount[word])
	}

	// Set implementation
	primes := map[int]bool{
		2:  true,
		3:  true,
		5:  true,
		7:  true,
		11: true,
	}
	if primes[7] {
		fmt.Println("\\n7 is prime")
	}

	// Concurrent map
	var concurrentMap = struct {
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}

	// Write
	concurrentMap.Lock()
	concurrentMap.m["counter"] = 42
	concurrentMap.Unlock()

	// Read
	concurrentMap.RLock()
	value := concurrentMap.m["counter"]
	concurrentMap.RUnlock()
	fmt.Println("\\nConcurrent map value:", value)

	// Advanced: inverted map
	original := map[string]int{
		"apple":  1,
		"banana": 2,
	}
	inverted := invertMap(original)
	fmt.Println("\\nInverted map:")
	for k, v := range inverted {
		fmt.Printf("%d: %s\\n", k, v)
	}
}

func invertMap(m map[string]int) map[int]string {
	inverted := make(map[int]string)
	for k, v := range m {
		inverted[v] = k
	}
	return inverted
}
```