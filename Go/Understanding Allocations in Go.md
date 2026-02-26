# Understanding Allocations in Go

Go manages memory automatically, but understanding where variables are allocated — **stack** vs **heap** — is key to writing performant code.

---

## Stack vs Heap

### The Stack

Each goroutine has its own stack (starting at 2KB, grows dynamically). Stack allocation is fast — just moving a pointer. When a function returns, its stack frame is automatically freed.

```
main()          ← stack frame
  └─ calculate() ← stack frame (freed when function returns)
       └─ helper() ← stack frame
```

### The Heap

Shared memory across all goroutines. Heap allocation is slower because the runtime must find free space, and the **garbage collector** must later clean it up.

**Stack = fast, automatic cleanup. Heap = slower, needs GC.**

---

## What Goes Where?

The Go compiler runs **escape analysis** to decide if a variable can stay on the stack or must move to the heap.

### Stays on the Stack

Variables that don't outlive their function:

```Go
func stackIt() int {
    y := 2
    return y * 2 // y is copied — stays on stack
}
```

```
BenchmarkStackIt  680439016  1.52 ns/op  0 B/op  0 allocs/op
```

### Escapes to the Heap

Variables whose reference outlives the function:

```Go
func heapIt() *int {
    y := 2
    res := y * 2
    return &res // res must survive after heapIt returns → heap
}
```

```
BenchmarkHeapIt  70922517  16.0 ns/op  8 B/op  1 allocs/op
```

10x slower because of heap allocation + GC overhead.

### Passing Pointers Down Is Fine

```Go
func main() {
    y := 2
    _ = process(&y) // pointer goes DOWN the stack → no allocation
}

func process(y *int) int {
    return *y * 2
}
```

```
BenchmarkProcess  705347884  1.62 ns/op  0 B/op  0 allocs/op
```

**Key rule:** Sharing pointers **up** the stack causes heap allocations. Sharing pointers **down** the stack does not.

---

## Viewing Escape Analysis

Use `gcflags` to see what the compiler decides:

```bash
go build -gcflags '-m -l' .
```

Output shows which variables escape:

```
./main.go:10:2: moved to heap: res
./main.go:18:14: y does not escape
```

---

## Why Heap Allocations Matter

Heap allocations trigger garbage collection. With many allocations, GC can consume significant CPU:

```
+-------------------+--------+---------+
|                   | Copy   | Pointer |
+-------------------+--------+---------+
| Time (20M ops)    | 0.28s  | 2.22s   |
| ns/op             | 5.20   | 52.6    |
| B/op              | 0      | 80      |
| allocs/op         | 0      | 1       |
| GC pauses (STW)   | 0      | 397     |
+-------------------+--------+---------+
```

Returning a struct by **value** (copy) kept everything on the stack — no GC activity. Returning by **pointer** caused heap allocations and ~400 stop-the-world GC pauses.

---

## Common Causes of Heap Allocations

| Cause | Example |
|-------|---------|
| Returning a pointer to a local variable | `return &result` |
| Storing in an interface | `var i interface{} = x` |
| Closures capturing variables | `go func() { use(x) }()` |
| Large stack objects | Very large arrays/structs |
| Slice/map growth | `append()` may allocate |

---

## Reducing Allocations

### 1. Return values instead of pointers when possible

```Go
// Causes heap allocation
func newUser() *User {
    u := User{Name: "Alice"}
    return &u
}

// Stays on stack (caller's frame)
func newUser() User {
    return User{Name: "Alice"}
}
```

### 2. Pre-allocate slices

```Go
// Many small allocations from append
var items []string
for _, v := range data {
    items = append(items, v)
}

// One allocation upfront
items := make([]string, 0, len(data))
for _, v := range data {
    items = append(items, v)
}
```

### 3. Use `sync.Pool` for frequently allocated objects

```Go
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func process() {
    buf := bufPool.Get().(*bytes.Buffer)
    defer bufPool.Put(buf)
    buf.Reset()
    // use buf...
}
```

### 4. Avoid unnecessary interfaces

```Go
// Forces heap allocation (value stored in interface)
func log(v any) { fmt.Println(v) }

// No allocation
func log(s string) { fmt.Println(s) }
```

---

## Benchmarking Allocations

Use `-benchmem` to track allocations:

```bash
go test -bench . -benchmem
```

Output:

```
BenchmarkFunc-8  67836464  16.0 ns/op  8 B/op  1 allocs/op
```

- **B/op** — bytes allocated per operation
- **allocs/op** — heap allocations per operation

---

## Go Memory Architecture (Overview)

Go's allocator is inspired by TCMalloc and uses a multi-level caching strategy:

| Level | Name | Purpose |
|-------|------|---------|
| Per-P cache | `mcache` | Lock-free allocation for small objects (≤32KB) |
| Central cache | `mcentral` | Shared pool of spans per size class |
| Global heap | `mheap` | Manages all pages, talks to the OS |

- Objects ≤ 16B → tiny allocator in `mcache`
- Objects 16B–32KB → size-class spans from `mcache`
- Objects > 32KB → allocated directly from `mheap`

Memory is requested from the OS in large chunks called **arenas** (~64MB) to amortize syscall overhead.

---

## Summary

1. **Stack allocation is fast and free** — no GC involved
2. **Heap allocation is expensive** — requires GC to clean up
3. **Pointers shared up the stack escape to the heap**
4. Use `go build -gcflags '-m'` to see escape analysis decisions
5. Use `go test -bench -benchmem` to measure allocations
6. **Correctness first** — only optimize allocations when you've identified a bottleneck

## References

- [Understanding Allocations in Go — James Kirk](https://medium.com/eureka-engineering/understanding-allocations-in-go-stack-heap-memory-9a2631b5035d)
- [Visual Guide to Go Memory Allocator — Ankur Anand](https://medium.com/@ankur_anand/a-visual-guide-to-golang-memory-allocator-from-ground-up-e132258453ed)
- [Go FAQ: Stack or Heap?](https://golang.org/doc/faq#stack_or_heap)
- [GC in Go — Ardan Labs](https://www.ardanlabs.com/blog/2018/12/garbage-collection-in-go-part1-semantics.html)