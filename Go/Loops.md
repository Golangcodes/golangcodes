# Loops in Go

Go has only one loop keyword: `for`. It covers traditional loops, while-style loops, infinite loops, and range-based iteration.

---

## Basic For Loop

```Go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

Three components: **init** (`i := 0`), **condition** (`i < 5`), **post** (`i++`).

---

## While-Style Loop

Omit init and post — just a condition:

```Go
count := 0
for count < 5 {
    fmt.Println(count)
    count++
}
```

---

## Infinite Loop

Omit everything:

```Go
for {
    fmt.Println("running...")
    break // exit when needed
}
```

Common for servers, event loops, and REPLs.

---

## Range Loop

Iterate over slices, maps, strings, and channels.

```Go
// Slices — index and value
nums := []int{10, 20, 30}
for i, v := range nums {
    fmt.Println(i, v)
}

// Maps — key and value
ages := map[string]int{"Alice": 25, "Bob": 30}
for name, age := range ages {
    fmt.Println(name, age)
}

// Strings — byte position and rune
for pos, char := range "Go 日本語" {
    fmt.Printf("%d: %c\n", pos, char)
}

// Channels — until closed
for msg := range ch {
    fmt.Println(msg)
}
```

Use `_` to discard a value:

```Go
for _, v := range nums { fmt.Println(v) }  // skip index
for k := range ages { fmt.Println(k) }     // skip value
```

---

## Break, Continue, Labels

```Go
// break — exit the loop
for i := 0; i < 10; i++ {
    if i == 5 { break }
    fmt.Println(i)
}

// continue — skip to next iteration
for i := 0; i < 10; i++ {
    if i%2 == 0 { continue }
    fmt.Println(i) // odd numbers only
}

// labeled break — exit an outer loop from inside a nested loop
outer:
for i := 0; i < 5; i++ {
    for j := 0; j < 5; j++ {
        if i*j > 6 {
            break outer
        }
        fmt.Println(i, j)
    }
}
```

---

## Common Patterns

### Filtering

```Go
var evens []int
for _, n := range numbers {
    if n%2 == 0 {
        evens = append(evens, n)
    }
}
```

### Batch Processing

```Go
batch := 3
for i := 0; i < len(items); i += batch {
    end := i + batch
    if end > len(items) {
        end = len(items)
    }
    process(items[i:end])
}
```

### Generator (Channel)

```Go
func count(max int) <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < max; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

for n := range count(5) {
    fmt.Println(n)
}
```

---

## Performance Tips

1. **Pre-allocate slices** when you know the size:
   ```Go
   result := make([]int, 0, len(input))
   ```

2. **Range copies values** — use index for large structs:
   ```Go
   for i := range bigStructs {
       bigStructs[i].Process() // avoids copying
   }
   ```

3. **Be careful with closures** in goroutines:
   ```Go
   for _, v := range items {
       v := v // capture loop variable
       go func() { process(v) }()
   }
   ```