![](https://www.notion.so/icons/corn_green.svg)


# Pointers in Go: From Basics to Advanced

Pointers are a fundamental concept in Go that enable efficient memory management and data sharing. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Pointer Concepts](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

2.  [Pointer Declaration](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

3.  [Pointer Operations](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

4.  [Pointers with Structs](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

5.  [Pointers with Functions](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

6.  [Pointer Arithmetic (Absent in Go)](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

7.  [Advanced Pointer Patterns](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

8.  [Performance Considerations](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

9.  [Best Practices](Pointers%202091ccc344b1807caf3df37f8abfcb24.html)

## Basic Pointer Concepts

A pointer holds the memory address of a value:

*   `&` operator generates a pointer to its operand

*   operator dereferences a pointer (accesses the underlying value)

*   Zero value of a pointer is `nil`

```Go
var x int = 10
var p *int = &x  // p points to x
fmt.Println(*p)  // 10 (dereferencing)
```

## Pointer Declaration

Different ways to declare and use pointers:

```Go
// Method 1: Full declaration
var ptr1 *int
num := 42
ptr1 = &num

// Method 2: Short declaration
ptr2 := &num

// Method 3: Using new()
ptr3 := new(int)  // Allocates memory, returns pointer
*ptr3 = 100
```

## Pointer Operations

### Basic Operations

```Go
a := 10
b := &a          // b points to a
fmt.Println(*b)  // 10 (read through pointer)
*b = 20          // write through pointer
fmt.Println(a)   // 20 (original value changed)
```

### Nil Pointers

```Go
var p *int       // nil pointer
if p != nil {
    fmt.Println(*p)  // Would panic if dereferenced
}
```

### Pointer Comparison

```Go
x := 10
p1 := &x
p2 := &x
fmt.Println(p1 == p2)  // true (same memory address)
```

## Pointers with Structs

### Value vs Pointer Receivers

```Go
type Point struct {
    X, Y int
}

// Value receiver (works on copy)
func (p Point) MoveByValue(dx, dy int) {
    p.X += dx
    p.Y += dy
}

// Pointer receiver (modifies original)
func (p *Point) MoveByPointer(dx, dy int) {
    p.X += dx
    p.Y += dy
}

p := Point{1, 2}
p.MoveByValue(1, 1)  // p unchanged
p.MoveByPointer(1, 1) // p modified
```

### Automatic Dereferencing

Go automatically dereferences struct pointers:

```Go
p := &Point{1, 2}
fmt.Println(p.X)  // Same as (*p).X
```

## Pointers with Functions

### Passing Pointers to Functions

```Go
func increment(p *int) {
    *p++
}

x := 10
increment(&x)
fmt.Println(x)  // 11
```

### Returning Pointers

```Go
func createPoint(x, y int) *Point {
    return &Point{x, y}  // Safe to return address of local variable
}

p := createPoint(5, 10)
```

## Pointer Arithmetic (Absent in Go)

Unlike C/C++, Go doesn't support pointer arithmetic:

```Go
arr := [3]int{1, 2, 3}
p := &arr[0]
// p++  // Compile error: invalid operation
```

## Advanced Pointer Patterns

### Pointer to Pointer

```Go
var x int = 10
var p *int = &x
var pp **int = &p

fmt.Println(**pp)  // 10
```

### Function Pointers

```Go
var op func(int, int) int
op = func(a, b int) int { return a + b }
result := op(3, 4)  // 7
```

### Interface Pointers

```Go
var w io.Writer = os.Stdout
pw := &w
(*pw).Write([]byte("hello\\n"))
```

### Linked List Example

```Go
type Node struct {
    Value int
    Next  *Node
}

list := &Node{1, &Node{2, &Node{3, nil}}}
```

## Performance Considerations

1.  **Passing pointers** to large structs is more efficient than passing values

2.  **Pointer indirection** has a small performance cost

3.  **Escape analysis** determines if variables are allocated on heap or stack

4.  **Cache locality** can be better with values than pointers

5.  **Garbage collection** overhead increases with many heap allocations

## Best Practices

1.  **Use pointers** to modify function arguments

2.  **Avoid unnecessary pointers** - they add complexity

3.  **Document pointer ownership** - who is responsible for the data

4.  **Check for nil** before dereferencing

5.  **Prefer value receivers** for small structs

6.  **Use pointer receivers** when methods need to modify the struct

7.  **Be consistent** with receiver types across all methods of a type

8.  **Avoid returning pointers** to local variables (except when it's safe)

9.  **Consider sync.Pool** for frequently allocated pointer types

10.  **Profile before optimizing** pointer usage

## Complete Example

```Go
package main

import (
	"fmt"
	"unsafe"
)

type Person struct {
	Name string
	Age  int
}

func updateName(p *Person, name string) {
	p.Name = name
}

func (p *Person) birthday() {
	p.Age++
	fmt.Printf("%s is now %d years old\\n", p.Name, p.Age)
}

func main() {
	// Basic pointer usage
	x := 10
	ptr := &x
	fmt.Println("Value of x through pointer:", *ptr)

	// Pointer to struct
	person := Person{"Alice", 30}
	updateName(&person, "Alicia")
	fmt.Println("Updated person:", person)

	// Method with pointer receiver
	person.birthday()

	// Pointer to pointer
	pp := &ptr
	fmt.Println("Pointer to pointer value:", **pp)

	// Unsafe pointer (advanced usage)
	bytes := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}
	str := *(*string)(unsafe.Pointer(&bytes))
	fmt.Println("String from bytes (unsafe):", str)

	// Function pointer
	var op func(int, int) int = func(a, b int) int { return a * b }
	fmt.Println("Function pointer result:", op(5, 6))

	// Linked list with pointers
	type Node struct {
		Value int
		Next  *Node
	}
	list := &Node{1, &Node{2, &Node{3, nil}}}
	for n := list; n != nil; n = n.Next {
		fmt.Println("Node value:", n.Value)
	}
}
```
