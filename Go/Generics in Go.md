![](https://www.notion.so/icons/flag-pennant_green.svg)

# Generics in Go

# Generics in Go: From Basics to Advanced

Generics were introduced in Go 1.18, bringing parametric polymorphism to the language. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Generic Concepts](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

2.  [Type Parameters](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

3.  [Generic Functions](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

4.  [Generic Types](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

5.  [Type Constraints](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

6.  [Type Inference](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

7.  [Advanced Patterns](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

8.  [Performance Considerations](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

9.  [Best Practices](Generics%20in%20Go%202091ccc344b1804988d2c11ce5a475fb.html)

## Basic Generic Concepts

Generics allow writing code that works with multiple types while maintaining type safety:

```Go
func PrintSlice[T any](s []T) {
    for _, v := range s {
        fmt.Println(v)
    }
}
```

Key concepts:

*   Type parameters (`[T any]`)

*   Type constraints (`any`, `comparable`, custom constraints)

*   Instantiation (compile-time generation of concrete functions)

## Type Parameters

Declared in square brackets before function parameters:

```Go
// Single type parameter
func Identity[T any](t T) T {
    return t
}

// Multiple type parameters
func MapKeys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}
```

## Generic Functions

Functions that work with multiple types:

```Go
// Basic generic function
func Swap[T any](a, b T) (T, T) {
    return b, a
}

// Using with different types
a, b := Swap(1, 2)           // int
x, y := Swap("hello", "world") // string
```

## Generic Types

Types that can be parameterized:

```Go
// Generic stack type
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
    if len(s.items) == 0 {
        panic("empty stack")
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item
}

// Usage
intStack := Stack[int]{}
intStack.Push(1)
intStack.Push(2)
fmt.Println(intStack.Pop()) // 2
```

## Type Constraints

Define what operations are available on type parameters:

### Built-in Constraints

*   `any` - any type (equivalent to `interface{}`)

*   `comparable` - types that can be compared with == and !=

### Custom Constraints

```Go
type Number interface {
    ~int | ~float64 | ~uint // Using type approximation (~)
}

func Sum[T Number](numbers []T) T {
    var total T
    for _, n := range numbers {
        total += n
    }
    return total
}
```

### Method Constraints

```Go
type Stringer interface {
    String() string
}

func PrintString[T Stringer](t T) {
    fmt.Println(t.String())
}
```

## Type Inference

Go can often infer type parameters:

```Go
// Type inferred from arguments
fmt.Println(Sum([]int{1, 2, 3}))       // T inferred as int
fmt.Println(Sum([]float64{1.1, 2.2})) // T inferred as float64

// Explicit type parameters
fmt.Println(Sum[float64]([]int{1, 2, 3})) // Convert int to float64
```

## Advanced Patterns

### Generic Interfaces

```Go
type Container[T any] interface {
    Add(item T)
    Get(index int) T
}

type GenericSlice[T any] []T

func (g *GenericSlice[T]) Add(item T) {
    *g = append(*g, item)
}

func (g GenericSlice[T]) Get(index int) T {
    return g[index]
}
```

### Recursive Generic Types

```Go
type TreeNode[T any] struct {
    Value  T
    Left   *TreeNode[T]
    Right  *TreeNode[T]
}
```

### Higher-Order Generic Functions

```Go
func Filter[T any](slice []T, predicate func(T) bool) []T {
    var result []T
    for _, v := range slice {
        if predicate(v) {
            result = append(result, v)
        }
    }
    return result
}

// Usage
numbers := []int{1, 2, 3, 4, 5}
evens := Filter(numbers, func(n int) bool { return n%2 == 0 })
```

### Generic Methods (with type parameters on receivers)

```Go
type Pair[A, B any] struct {
    First  A
    Second B
}

func (p Pair[A, B]) Swap() Pair[B, A] {
    return Pair[B, A]{p.Second, p.First}
}
```

## Performance Considerations

1.  **No runtime overhead** - code is generated at compile time

2.  **Binary size increase** - due to multiple instantiations

3.  **Compile time impact** - generics can slow down compilation

4.  **Memory usage** - each instantiation creates specialized code

## Best Practices

1.  **Start simple** - don't overuse generics prematurely

2.  **Use descriptive type parameter names** (T for simple cases, K/V for maps, etc.)

3.  **Prefer constraints** over `any` when possible

4.  **Document generic functions** thoroughly

5.  **Consider performance implications** for hot code paths

6.  **Avoid complex type hierarchies** - keep it Go-like

7.  **Test with different types** to ensure flexibility

8.  **Use type inference** where it improves readability

9.  **Be cautious with method sets** - generic methods have limitations

10.  **Watch for error messages** - they can be complex with generics

## Complete Example

```Go
package main

import (
	"fmt"
	"golang.org/x/exp/constraints"
)

// Basic generic function
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Generic type
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() T {
	if len(s.items) == 0 {
		panic("stack is empty")
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

// Custom constraint
type Number interface {
	constraints.Integer | constraints.Float
}

func Average[T Number](numbers []T) T {
	if len(numbers) == 0 {
		return 0
	}
	var sum T
	for _, n := range numbers {
		sum += n
	}
	return sum / T(len(numbers))
}

// Higher-order generic function
func Map[T, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

func main() {
	// Using Max function
	fmt.Println("Max int:", Max(3, 7))
	fmt.Println("Max float:", Max(3.14, 2.71))

	// Using Stack type
	intStack := Stack[int]{}
	intStack.Push(1)
	intStack.Push(2)
	fmt.Println("Popped:", intStack.Pop())

	// Using Average function
	fmt.Println("Average:", Average([]float64{1.0, 2.0, 3.0}))

	// Using Map function
	numbers := []int{1, 2, 3}
	squares := Map(numbers, func(n int) int { return n * n })
	fmt.Println("Squares:", squares)

	// Type inference in action
	fmt.Println("Max inferred:", Max(3.5, 7.2)) // float64
}
```

