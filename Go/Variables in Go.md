# Variables in Go

---

## Declaration

```Go
// Explicit type
var name string
var age int
var isActive bool

// With initialization
var name string = "Alice"
var age int = 30

// Grouped declaration
var (
    host    string = "localhost"
    port    int    = 8080
    verbose bool   = true
)
```

---

## Type Inference

Go infers the type from the value:

```Go
var name = "Alice"  // string
var age = 30        // int
var rate = 3.14     // float64
```

---

## Short Declaration (`:=`)

Inside functions only. Also called the walrus operator.

```Go
name := "Alice"
age := 30
x, y := 10, 20
```

Useful in `if` and `for` statements:

```Go
if err := doSomething(); err != nil {
    log.Fatal(err)
}
```

**Cannot be used at package level** — use `var` instead.

---

## Zero Values

Uninitialized variables get their **zero value**:

| Type | Zero Value |
|------|-----------|
| `int`, `float64` | `0` |
| `bool` | `false` |
| `string` | `""` |
| Pointers, slices, maps, channels, interfaces | `nil` |

---

## Constants

```Go
const Pi = 3.14159
const MaxRetries = 3

// Grouped
const (
    StatusOK      = 200
    StatusCreated = 201
)

// iota — auto-incrementing constant generator
const (
    Read    = 1 << iota // 1
    Write               // 2
    Execute             // 4
)
```

Constants must be known at compile time. They can be untyped, which gives them more flexibility.

---

## Scope

```Go
var global = "package level" // visible throughout the package

func main() {
    local := "function level" // visible within main

    if true {
        block := "block level" // visible only in this if-block
        fmt.Println(block)
    }
    // block is undefined here
}
```

---

## Type Conversion

Go requires **explicit** conversion — no implicit casting.

```Go
i := 42
f := float64(i)     // int → float64
u := uint(f)         // float64 → uint

// String conversions
s := strconv.Itoa(123)        // int → string: "123"
n, err := strconv.Atoi("123") // string → int: 123
r := string(rune(65))         // rune → string: "A"
```

---

## Type Aliases

Create distinct types from existing ones:

```Go
type Celsius float64
type Fahrenheit float64

var c Celsius = 100.0
var f Fahrenheit = 212.0
// c = f  // compile error — different types
```

---

## Blank Identifier (`_`)

Discard values you don't need:

```Go
_, err := os.Open("file.txt")
```

---

## Variable Shadowing

A common pitfall — `:=` in an inner scope creates a **new** variable that hides the outer one.

```Go
x := 10
if true {
    x := 5          // new x — shadows outer x
    fmt.Println(x)  // 5
}
fmt.Println(x)      // 10 (unchanged)
```

Watch out for shadowing `err`:

```Go
result, err := step1()
if err := step2(result); err != nil {
    return err // this is step2's err, not step1's
}
```

Catch shadowing bugs with:

```bash
go vet -vettool=$(which shadow) ./...
```

---

## Best Practices

1. Use `:=` inside functions, `var` at package level
2. Declare variables **close to where they're used**
3. Use `const` for values that never change
4. Use meaningful names — avoid single letters except loop indices
5. Be explicit with type conversions
6. Watch for **variable shadowing** with `:=` in inner scopes