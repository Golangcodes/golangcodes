# Structs in Go

Structs are Go's way of defining custom data types — a collection of named fields grouped together. Think of them as blueprints for organizing related data.

https://www.youtube.com/watch?v=fXZJu_JuH0A

---

## Declaring Structs

```Go
type Person struct {
    FirstName string
    LastName  string
    Age       int
}
```

---

## Creating Instances

```Go
// Zero value (all fields set to defaults)
var p1 Person

// Named fields (preferred)
p2 := Person{
    FirstName: "Jane",
    LastName:  "Smith",
    Age:       25,
}

// Positional (fragile — avoid for exported types)
p3 := Person{"John", "Doe", 30}

// Pointer
p4 := &Person{FirstName: "Alice", Age: 35}
```

---

## Accessing Fields

```Go
fmt.Println(p2.FirstName) // "Jane"
p2.Age = 26               // update

// Pointers auto-dereference
p := &Person{FirstName: "Bob"}
fmt.Println(p.FirstName) // no need for (*p).FirstName
```

---

## Methods

Attach functions to a struct with a **receiver**.

### Value Receiver (read-only — works on a copy)

```Go
func (p Person) FullName() string {
    return p.FirstName + " " + p.LastName
}
```

### Pointer Receiver (can modify the original)

```Go
func (p *Person) Birthday() {
    p.Age++
}

p2.Birthday()
fmt.Println(p2.Age) // 26
```

**Rule of thumb:** Use pointer receivers when you need to mutate the struct or it's large (>64 bytes). Use value receivers for small, read-only methods.

---

## Constructor Functions

Go doesn't have constructors — use factory functions instead.

```Go
func NewPerson(first, last string, age int) *Person {
    return &Person{
        FirstName: first,
        LastName:  last,
        Age:       age,
    }
}
```

---

## Anonymous Structs

For one-off data shapes you don't need to reuse.

```Go
point := struct {
    X, Y int
}{10, 20}
```

---

## Embedded Structs (Composition)

Go uses composition instead of inheritance. Embed one struct inside another.

```Go
type Address struct {
    Street string
    City   string
}

type Employee struct {
    Person           // embedded — fields are "promoted"
    Position string
    Address          // also embedded
}

emp := Employee{
    Person:   Person{FirstName: "John", LastName: "Doe", Age: 30},
    Position: "Developer",
    Address:  Address{City: "New York"},
}

// Access promoted fields directly
fmt.Println(emp.FirstName) // "John"
fmt.Println(emp.City)      // "New York"
```

Methods are also promoted — `emp.FullName()` works if `Person` has that method.

---

## Struct Tags

Metadata attached to fields, used by packages like `encoding/json` and ORMs.

```Go
type User struct {
    ID       int    `json:"id" db:"user_id"`
    Username string `json:"username"`
    Password string `json:"-"` // excluded from JSON
    Email    string `json:"email,omitempty"`
}
```

```Go
user := User{ID: 1, Username: "jdoe"}
data, _ := json.Marshal(user)
// {"id":1,"username":"jdoe"}
```

---

## Method Overriding via Embedding

Embedded methods can be "overridden" by defining the same method on the outer struct.

```Go
type Animal struct{ Name string }

func (a Animal) Speak() string {
    return a.Name + " makes a sound"
}

type Dog struct {
    Animal
    Breed string
}

func (d Dog) Speak() string {
    return d.Name + " says: Woof!"
}

dog := Dog{Animal: Animal{Name: "Rex"}, Breed: "Labrador"}
fmt.Println(dog.Speak())        // "Rex says: Woof!"
fmt.Println(dog.Animal.Speak()) // "Rex makes a sound"
```

---

## Functional Options Pattern

A clean way to configure structs with many optional fields.

```Go
type Server struct {
    Host    string
    Port    int
    Timeout time.Duration
}

type Option func(*Server)

func WithHost(h string) Option {
    return func(s *Server) { s.Host = h }
}

func WithTimeout(t time.Duration) Option {
    return func(s *Server) { s.Timeout = t }
}

func NewServer(port int, opts ...Option) *Server {
    s := &Server{Host: "localhost", Port: port, Timeout: 10 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(8080, WithHost("0.0.0.0"), WithTimeout(30*time.Second))
```

---

## Memory Layout & Field Ordering

Structs are contiguous blocks of memory. Field order affects padding and size.

```Go
// Wasteful — 24 bytes (padding between fields)
type Bad struct {
    Active bool    // 1 byte + 7 padding
    ID     int64   // 8 bytes
    Age    int32   // 4 bytes + 4 padding
}

// Efficient — 16 bytes
type Good struct {
    ID     int64   // 8 bytes
    Age    int32   // 4 bytes
    Active bool    // 1 byte + 3 padding
}
```

Order fields from largest to smallest to minimize wasted space.

---

## Common Pitfalls

1. **Accidental copies** — value receivers don't modify the original
2. **Non-comparable structs** — structs with maps or slices can't use `==`
3. **Unexported fields** — lowercase fields aren't visible outside the package
4. **Nil pointer** — calling methods on a nil pointer panics

---

## Best Practices

1. Use **pointer receivers** for methods that mutate or for large structs
2. Use **constructor functions** (`NewX`) for complex initialization
3. Make **zero values useful** — design structs so the default state is valid
4. Keep structs **small and focused** — one responsibility per struct
5. Use **struct tags** consistently for JSON, DB, and validation
6. Order fields **largest to smallest** for memory efficiency
7. Prefer **composition** (embedding) over deep nesting