![](https://www.notion.so/icons/wall_green.svg)


# Slices in Go: From Basics to Advanced

Slices are one of Go's most important and versatile data structures. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Slice Declaration](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

2.  [Slice Creation](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

3.  [Slice Operations](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

4.  [Slice Internals](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

5.  [Common Operations](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

6.  [Performance Considerations](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

7.  [Advanced Patterns](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

8.  [Best Practices](Slice%202091ccc344b180f3a99cd9d826ef5400.html)

## Basic Slice Declaration

A slice is a dynamically-sized, flexible view into an array:

```Go
var s []int          // nil slice
s := []int{1, 2, 3}  // slice literal
```

## Slice Creation

Multiple ways to create slices:

```Go
// 1. Literal
letters := []string{"a", "b", "c"}

// 2. Make with length and capacity
s := make([]int, 5)     // len=5, cap=5
s := make([]int, 5, 10) // len=5, cap=10

// 3. From array
arr := [5]int{1, 2, 3, 4, 5}
s := arr[1:4]           // [2, 3, 4]

// 4. From another slice
s2 := s[1:3]            // [3, 4]
```

## Slice Operations

### Accessing Elements

```Go
s := []int{10, 20, 30}
fmt.Println(s[1]) // 20
```

### Modifying Elements

```Go
s[1] = 25
fmt.Println(s) // [10, 25, 30]
```

### Appending Elements

```Go
s = append(s, 40, 50)
fmt.Println(s) // [10, 25, 30, 40, 50]
```

### Length and Capacity

```Go
fmt.Println(len(s)) // 5
fmt.Println(cap(s)) // 6 (capacity may grow)
```

### Copying Slices

```Go
src := []int{1, 2, 3}
dst := make([]int, 2)
copy(dst, src)     // copies min(len(src), len(dst)) elements
fmt.Println(dst)   // [1, 2]
```

## Slice Internals

A slice header contains:

*   Pointer to underlying array

*   Length (number of elements)

*   Capacity (maximum length)

```Go
type sliceHeader struct {
    Length        int
    Capacity      int
    ZerothElement *byte
}
```

When you pass a slice to a function, the header is copied (value semantics), but it still references the same array.

## Common Operations

### Filtering

```Go
func filter(numbers []int, condition func(int) bool) []int {
    var result []int
    for _, num := range numbers {
        if condition(num) {
            result = append(result, num)
        }
    }
    return result
}

evens := filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
```

### Mapping

```Go
func mapSlice(numbers []int, mapper func(int) int) []int {
    result := make([]int, len(numbers))
    for i, num := range numbers {
        result[i] = mapper(num)
    }
    return result
}

squares := mapSlice([]int{1, 2, 3}, func(n int) int { return n * n })
```

### Reducing

```Go
func reduce(numbers []int, reducer func(int, int) int, initial int) int {
    result := initial
    for _, num := range numbers {
        result = reducer(result, num)
    }
    return result
}

sum := reduce([]int{1, 2, 3}, func(a, b int) int { return a + b }, 0)
```

## Performance Considerations

1.  **Pre-allocate capacity** when possible:
    
    ```Go
    // Better for large slices
    result := make([]int, 0, 1000)
    for i := 0; i < 1000; i++ {
        result = append(result, i)
    }
    ```

2.  **Avoid unnecessary allocations** by reusing slices:
    
    ```Go
    var buffer []byte
    for {
        buffer = buffer[:0] // reset length
        // ... fill buffer ...
    }
    ```

3.  **Be aware of memory leaks** from slice references:
    
    ```Go
    var bigSlice []byte
    // ... fill with large data ...
    
    // Keep just first 10 elements (but underlying array remains)
    smallSlice := bigSlice[:10]
    
    // Better - copy what you need
    smallSlice := make([]byte, 10)
    copy(smallSlice, bigSlice)
    ```

## Advanced Patterns

### Slice Tricks (from Go Wiki)

```Go
// Delete element at index i
s = append(s[:i], s[i+1:]...)

// Delete without preserving order
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// Reverse a slice
for i := len(s)/2 - 1; i >= 0; i-- {
    opp := len(s) - 1 - i
    s[i], s[opp] = s[opp], s[i]
}
```

### Slice as Stack

```Go
stack := []int{}

// Push
stack = append(stack, 1)

// Pop
value := stack[len(stack)-1]
stack = stack[:len(stack)-1]
```

### Slice as Queue

```Go
queue := []int{}

// Enqueue
queue = append(queue, 1)

// Dequeue
value := queue[0]
queue = queue[1:]
```

### Multi-dimensional Slices

```Go
// Create a 2D slice
rows, cols := 3, 4
matrix := make([][]int, rows)
for i := range matrix {
    matrix[i] = make([]int, cols)
}

// Access elements
matrix[1][2] = 5
```

## Best Practices

1.  **Prefer slices over arrays** for most use cases

2.  **Specify capacity** when you know the final size

3.  **Be careful with sub-slices** - they share memory

4.  **Use copy()** when you need independent data

5.  **Consider nil slices** as valid empty slices

6.  **Document ownership** when slices are shared

7.  **Avoid modifying slices** passed as function parameters

8.  **Use** `**range**` for iteration when index isn't needed

9.  **Benchmark slice operations** in performance-critical code

10.  **Consider custom types** for complex slice usage

## Complete Example

```Go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Basic slice operations
	numbers := []int{1, 2, 3, 4, 5}
	fmt.Println("Original:", numbers)

	// Append
	numbers = append(numbers, 6)
	fmt.Println("After append:", numbers)

	// Slice operations
	sub := numbers[1:4]
	fmt.Println("Sub-slice:", sub)

	// Modification affects both
	sub[0] = 99
	fmt.Println("After modification:")
	fmt.Println("Original:", numbers)
	fmt.Println("Sub-slice:", sub)

	// Copy to avoid sharing
	independent := make([]int, len(sub))
	copy(independent, sub)
	independent[0] = 100
	fmt.Println("After copy and modification:")
	fmt.Println("Original:", numbers)
	fmt.Println("Independent:", independent)

	// Performance demonstration
	start := time.Now()
	var s []int
	for i := 0; i < 1000000; i++ {
		s = append(s, i)
	}
	fmt.Printf("Append without pre-allocation: %v\\n", time.Since(start))

	start = time.Now()
	s = make([]int, 0, 1000000)
	for i := 0; i < 1000000; i++ {
		s = append(s, i)
	}
	fmt.Printf("Append with pre-allocation: %v\\n", time.Since(start))

	// Advanced pattern: batch processing
	data := make([]int, 100)
	for i := range data {
		data[i] = rand.Intn(1000)
	}

	const batchSize = 10
	for i := 0; i < len(data); i += batchSize {
		end := i + batchSize
		if end > len(data) {
			end = len(data)
		}
		batch := data[i:end]
		fmt.Printf("Batch %d: %v\\n", i/batchSize+1, batch)
	}
}
```
