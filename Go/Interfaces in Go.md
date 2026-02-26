# Interfaces in Go

Interfaces define behavior — a set of method signatures that a type must have. Any type that implements all the methods of an interface automatically satisfies it. No explicit declaration needed.

```Go
type Stringer interface {
    String() string
}
```

Any type with a `String() string` method satisfies `Stringer` — the compiler checks this implicitly.

https://youtu.be/rH0bpx7I2Dk

---

## Declaring Interfaces

An interface lists method signatures. Types satisfy it by implementing those methods.

```Go
type Shape interface {
    Area() float64
    Perimeter() float64
}
```

---

## Implementing Interfaces

There's no `implements` keyword. Just define the methods.

```Go
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

// Circle now satisfies Shape
```

```Go
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Rectangle also satisfies Shape
```

Both can be used anywhere a `Shape` is expected:

```Go
func printArea(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
}

printArea(Circle{Radius: 5})
printArea(Rectangle{Width: 3, Height: 4})
```

---

## The Empty Interface (`any`)

The empty interface has no methods, so every type satisfies it.

```Go
func printAnything(v any) {
    fmt.Println(v)
}

printAnything(42)
printAnything("hello")
printAnything([]int{1, 2, 3})
```

`any` is an alias for `interface{}` (Go 1.18+). Prefer `any` in modern code.

Use sparingly — you lose type safety. Prefer concrete types or defined interfaces.

---

## Type Assertions

Extract the concrete type from an interface value.

```Go
var i any = "hello"

// With check (safe)
s, ok := i.(string)
fmt.Println(s, ok) // "hello" true

f, ok := i.(float64)
fmt.Println(f, ok) // 0 false

// Without check (panics if wrong type)
s = i.(string) // works
// f = i.(float64) // panics!
```

---

## Type Switches

Handle multiple types cleanly.

```Go
func describe(i any) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case Shape:
        fmt.Printf("Shape with area %.2f\n", v.Area())
    default:
        fmt.Printf("Unknown: %T\n", v)
    }
}
```

---

## Interface Composition

Build larger interfaces from smaller ones.

```Go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

---

## Reducing Boilerplate with Interfaces

Without interfaces, you'd write separate functions for every destination:

```Go
func writeToFile(c *Customer, f *os.File) error { ... }
func writeToBuffer(c *Customer, b *bytes.Buffer) error { ... }
```

With interfaces, one function works for both — because `os.File` and `bytes.Buffer` both satisfy `io.Writer`:

```Go
func (c *Customer) WriteJSON(w io.Writer) error {
    data, err := json.Marshal(c)
    if err != nil {
        return err
    }
    _, err = w.Write(data)
    return err
}

// Works with files
f, _ := os.Create("customer.json")
customer.WriteJSON(f)

// Works with buffers
var buf bytes.Buffer
customer.WriteJSON(&buf)

// Works with HTTP responses
func handler(w http.ResponseWriter, r *http.Request) {
    customer.WriteJSON(w)
}
```

---

## Mocking for Tests

Interfaces make testing easy. Define an interface, create a mock.

**Production code:**

```Go
type ShopModel interface {
    CountCustomers(time.Time) (int, error)
    CountSales(time.Time) (int, error)
}

func calculateSalesRate(sm ShopModel) (string, error) {
    since := time.Now().Add(-24 * time.Hour)
    sales, err := sm.CountSales(since)
    if err != nil {
        return "", err
    }
    customers, err := sm.CountCustomers(since)
    if err != nil {
        return "", err
    }
    rate := float64(sales) / float64(customers)
    return fmt.Sprintf("%.2f", rate), nil
}
```

**Test code:**

```Go
type MockShopDB struct{}

func (m *MockShopDB) CountCustomers(_ time.Time) (int, error) {
    return 1000, nil
}

func (m *MockShopDB) CountSales(_ time.Time) (int, error) {
    return 333, nil
}

func TestCalculateSalesRate(t *testing.T) {
    mock := &MockShopDB{}
    rate, err := calculateSalesRate(mock)
    if err != nil {
        t.Fatal(err)
    }
    if rate != "0.33" {
        t.Fatalf("got %s, want 0.33", rate)
    }
}
```

No database needed. The function doesn't care what type it gets — only that it has the right methods.

---

## Compile-Time Interface Check

Verify a type satisfies an interface at compile time:

```Go
var _ Shape = (*Circle)(nil)     // fails to compile if Circle doesn't satisfy Shape
var _ io.Writer = (*MyType)(nil) // same pattern for standard library interfaces
```

---

## Common Standard Library Interfaces

| Interface | Methods | Used For |
|-----------|---------|----------|
| `fmt.Stringer` | `String() string` | Custom string representation |
| `error` | `Error() string` | Error values |
| `io.Reader` | `Read([]byte) (int, error)` | Reading data |
| `io.Writer` | `Write([]byte) (int, error)` | Writing data |
| `io.Closer` | `Close() error` | Closing resources |
| `io.ReadWriter` | `Read` + `Write` | Bidirectional I/O |
| `http.Handler` | `ServeHTTP(ResponseWriter, *Request)` | HTTP request handling |
| `sort.Interface` | `Len`, `Less`, `Swap` | Custom sorting |

---

## Best Practices

1. **Keep interfaces small** — 1-3 methods is ideal. `io.Reader` has one method and powers half the standard library
2. **Accept interfaces, return structs** — be flexible in inputs, specific in outputs
3. **Don't create interfaces prematurely** — write concrete code first, extract interfaces when you need decoupling or testing
4. **Interfaces belong with the consumer** — define them where they're used, not where they're implemented
5. **Name by behavior** — use `-er` suffix: `Reader`, `Writer`, `Closer`, `Formatter`
6. **Avoid `any` when possible** — prefer defined interfaces or concrete types for type safety
7. **Use composition** — build complex interfaces from simple ones

Reference: [Interfaces Explained by Alex Edwards](https://www.alexedwards.net/blog/interfaces-explained)