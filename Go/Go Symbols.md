# Go Symbols & Keywords

A complete reference of Go's keywords, operators, and built-in identifiers.

---

## Keywords (25)

### Declaration

| Keyword | Purpose |
|---------|---------|
| `var` | Declare variables |
| `const` | Declare constants |
| `type` | Define new types or aliases |
| `func` | Define functions |
| `struct` | Define struct types |
| `interface` | Define interface types |
| `map` | Declare map types |

### Control Flow

| Keyword | Purpose |
|---------|---------|
| `if` / `else` | Conditional branching |
| `switch` / `case` / `default` | Multi-way branching |
| `for` | Loops (Go's only loop keyword) |
| `range` | Iterate over slices, maps, channels, strings |
| `break` | Exit a loop or switch |
| `continue` | Skip to next iteration |
| `goto` | Jump to a labeled statement |
| `fallthrough` | Continue to next case in switch |
| `return` | Return from a function |

### Concurrency

| Keyword | Purpose |
|---------|---------|
| `go` | Launch a goroutine |
| `chan` | Declare a channel type |
| `select` | Wait on multiple channel operations |

### Error Handling & Defer

| Keyword | Purpose |
|---------|---------|
| `defer` | Schedule a call to run when the function returns |
| `panic` | Stop normal execution |
| `recover` | Catch a panic inside a deferred function |

### Packages

| Keyword | Purpose |
|---------|---------|
| `package` | Declare which package this file belongs to |
| `import` | Import other packages |

---

## Operators

### Arithmetic

```
+    addition
-    subtraction
*    multiplication
/    division
%    remainder (modulo)
++   increment (postfix only, statement not expression)
--   decrement (postfix only, statement not expression)
```

### Comparison

```
==   equal
!=   not equal
<    less than
<=   less or equal
>    greater than
>=   greater or equal
```

### Logical

```
&&   AND
||   OR
!    NOT
```

### Bitwise

```
&    AND
|    OR
^    XOR (binary) / NOT (unary)
&^   AND NOT (bit clear)
<<   left shift
>>   right shift
```

### Assignment

```
=    assign
:=   short variable declaration
+=   add and assign
-=   subtract and assign
*=   multiply and assign
/=   divide and assign
%=   modulo and assign
&=   bitwise AND assign
|=   bitwise OR assign
^=   bitwise XOR assign
<<=  left shift assign
>>=  right shift assign
```

### Pointer & Address

```
&    address of (get pointer)
*    dereference (get value from pointer)
```

### Channel

```
<-   send to or receive from channel
     ch <- v     // send v to ch
     v := <-ch   // receive from ch
```

### Other

```
...  variadic parameters / slice unpacking
.    field/method selector
_    blank identifier (discard a value)
```

---

## Punctuation

| Symbol | Purpose |
|--------|---------|
| `()` | Function calls, grouping |
| `[]` | Array/slice indexing, type parameters |
| `{}` | Blocks, composite literals |
| `,` | Separator in lists |
| `;` | Statement terminator (usually implicit) |
| `"` | Interpreted string literal |
| `` ` `` | Raw string literal |
| `'` | Rune (character) literal |
| `//` | Line comment |
| `/* */` | Block comment |

---

## Predeclared Identifiers

### Constants

| Name | Purpose |
|------|---------|
| `true` / `false` | Boolean values |
| `iota` | Auto-incrementing constant generator |
| `nil` | Zero value for pointers, slices, maps, channels, functions, interfaces |

### Types

| Type | Description |
|------|-------------|
| `bool` | Boolean |
| `string` | String |
| `int`, `int8`, `int16`, `int32`, `int64` | Signed integers |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | Unsigned integers |
| `uintptr` | Integer large enough to hold a pointer |
| `float32`, `float64` | Floating-point numbers |
| `complex64`, `complex128` | Complex numbers |
| `byte` | Alias for `uint8` |
| `rune` | Alias for `int32` (Unicode code point) |
| `any` | Alias for `interface{}` (Go 1.18+) |
| `error` | Built-in error interface |
| `comparable` | Constraint for types supporting `==` (Go 1.18+) |

### Built-in Functions

| Function | Purpose |
|----------|---------|
| `make(T, ...)` | Allocate and initialize slices, maps, channels |
| `new(T)` | Allocate zeroed memory, return pointer |
| `len(v)` | Length of string, slice, array, map, channel |
| `cap(v)` | Capacity of slice, array, channel |
| `append(s, ...)` | Append elements to a slice |
| `copy(dst, src)` | Copy elements between slices |
| `delete(m, key)` | Remove entry from a map |
| `close(ch)` | Close a channel |
| `panic(v)` | Stop normal execution |
| `recover()` | Catch a panic in a deferred function |
| `print(...)` | Built-in print (use `fmt` instead) |
| `println(...)` | Built-in println (use `fmt` instead) |
| `real(c)` | Real part of a complex number |
| `imag(c)` | Imaginary part of a complex number |
| `complex(r, i)` | Construct a complex number |
| `clear(v)` | Clear a map or zero a slice (Go 1.21+) |
| `min(x, ...)` | Return minimum value (Go 1.21+) |
| `max(x, ...)` | Return maximum value (Go 1.21+) |

---

## Type Assertion & Type Switch

```Go
// Type assertion
val, ok := x.(string)

// Type switch
switch v := x.(type) {
case string:
    fmt.Println("string:", v)
case int:
    fmt.Println("int:", v)
default:
    fmt.Printf("other: %T\n", v)
}
```
