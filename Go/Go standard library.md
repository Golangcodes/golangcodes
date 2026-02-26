# Go Standard Library

A practical guide to the most important packages in Go's standard library, with examples.

---

## `fmt` — Formatted I/O

```Go
// Output
fmt.Print("no newline")
fmt.Println("with newline")
fmt.Printf("Hello %s, you are %d\n", "Alice", 30)

// Return as string
s := fmt.Sprintf("formatted: %v", 42)

// Write to any io.Writer
fmt.Fprintf(os.Stderr, "error: %v\n", err)

// Read input
var name string
fmt.Scan(&name)
fmt.Scanf("%s", &name)
```

**Common format verbs:** `%v` (default), `%+v` (struct fields), `%T` (type), `%d` (int), `%s` (string), `%f` (float), `%t` (bool), `%x` (hex), `%p` (pointer)

---

## `os` — Operating System Interface

### Files

```Go
// Read entire file
data, err := os.ReadFile("config.json")

// Create and write
f, err := os.Create("output.txt")
defer f.Close()
f.WriteString("Hello, World!\n")

// Open for reading
f, err := os.Open("input.txt")
defer f.Close()
```

### Directories

```Go
os.Mkdir("mydir", 0755)
os.MkdirAll("path/to/dir", 0755)
os.Remove("file.txt")
os.RemoveAll("dir/")

entries, _ := os.ReadDir(".")
for _, e := range entries {
    fmt.Println(e.Name())
}
```

### Environment & Args

```Go
os.Setenv("KEY", "value")
val := os.Getenv("KEY")
args := os.Args[1:] // command-line arguments
os.Exit(1)          // exit with code
```

---

## `io` — Core I/O Interfaces

The `Reader` and `Writer` interfaces are used everywhere in Go.

```Go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```

### Useful functions

```Go
io.Copy(dst, src)                    // copy all data
io.ReadAll(r)                        // read everything into []byte
io.WriteString(w, "hello")           // write string to writer
io.MultiWriter(os.Stdout, file)      // write to multiple destinations
io.LimitReader(r, 1024)              // read at most N bytes
io.TeeReader(r, w)                   // read from r, copy to w
```

---

## `bufio` — Buffered I/O

Reduces system calls by buffering reads/writes.

```Go
// Buffered reading
reader := bufio.NewReader(file)
line, err := reader.ReadString('\n')

// Scanner (line-by-line)
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    fmt.Println(scanner.Text())
}

// Scanner by words
scanner.Split(bufio.ScanWords)

// Buffered writing
writer := bufio.NewWriter(file)
writer.WriteString("buffered output\n")
writer.Flush() // always flush when done
```

---

## `strings` — String Manipulation

```Go
strings.Contains("hello world", "world")  // true
strings.HasPrefix("hello", "he")          // true
strings.HasSuffix("hello", "lo")          // true
strings.Index("hello", "ll")             // 2
strings.Count("cheese", "e")             // 3

strings.ToUpper("hello")                 // "HELLO"
strings.ToLower("HELLO")                 // "hello"
strings.TrimSpace("  hi  ")             // "hi"
strings.Trim("!!hi!!", "!")             // "hi"

strings.Replace("aaa", "a", "b", 2)     // "bba"
strings.ReplaceAll("aaa", "a", "b")     // "bbb"

strings.Split("a,b,c", ",")             // ["a", "b", "c"]
strings.Join([]string{"a","b"}, "-")    // "a-b"

strings.Fields("  foo  bar  ")           // ["foo", "bar"]
strings.Repeat("ab", 3)                 // "ababab"
```

### String Builder (efficient concatenation)

```Go
var b strings.Builder
b.WriteString("Hello")
b.WriteString(" World")
result := b.String()
```

---

## `bytes` — Byte Slice Manipulation

Same API as `strings` but works on `[]byte`. Also provides `bytes.Buffer`:

```Go
var buf bytes.Buffer
buf.WriteString("Hello ")
buf.Write([]byte("World"))
fmt.Println(buf.String()) // "Hello World"
buf.Reset()
```

---

## `strconv` — String Conversions

```Go
// String ↔ Integer
i, err := strconv.Atoi("42")       // string → int
s := strconv.Itoa(42)              // int → string

// String ↔ Float
f, err := strconv.ParseFloat("3.14", 64)
s := strconv.FormatFloat(3.14, 'f', 2, 64)

// String ↔ Bool
b, err := strconv.ParseBool("true")
s := strconv.FormatBool(true)

// Different bases
n, _ := strconv.ParseInt("FF", 16, 64)   // 255
s := strconv.FormatInt(255, 2)            // "11111111"
```

---

## `errors` — Error Handling

See the dedicated [Errors lesson](/go/Erros) for full coverage.

```Go
// Create errors
err := errors.New("something failed")
err := fmt.Errorf("open %s: %w", path, cause) // wrap with context

// Inspect errors
errors.Is(err, os.ErrNotExist)  // check sentinel
errors.As(err, &target)         // extract type
errors.Unwrap(err)              // get underlying
```

---

## `log` — Simple Logging

```Go
log.Println("info message")
log.Printf("user %s logged in", username)
log.Fatal("critical error")   // logs then os.Exit(1)
log.Panic("impossible state") // logs then panic()

// Custom logger
logger := log.New(os.Stderr, "APP: ", log.Ldate|log.Ltime|log.Lshortfile)
logger.Println("custom log")

// Log to file
f, _ := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
log.SetOutput(f)
```

**Flags:** `log.Ldate`, `log.Ltime`, `log.Lmicroseconds`, `log.Lshortfile`, `log.Llongfile`

---

## `unicode` — Character Classification

```Go
unicode.IsLetter('A')  // true
unicode.IsDigit('9')   // true
unicode.IsSpace(' ')   // true
unicode.IsUpper('A')   // true
unicode.IsLower('a')   // true
unicode.IsPunct('!')   // true

unicode.ToUpper('a')   // 'A'
unicode.ToLower('A')   // 'a'

// Check script
unicode.Is(unicode.Han, '世')     // true (Chinese)
unicode.Is(unicode.Arabic, 'ا')  // true
```

---

## `reflect` — Runtime Reflection

Use sparingly — it's slower than direct code and bypasses compile-time checks.

```Go
// Type and value inspection
t := reflect.TypeOf(myVar)
v := reflect.ValueOf(myVar)
fmt.Println(t.Kind()) // struct, slice, map, etc.

// Iterate struct fields
for i := 0; i < t.NumField(); i++ {
    fmt.Println(t.Field(i).Name, v.Field(i))
}

// Modify values (requires pointer)
v := reflect.ValueOf(&x).Elem()
v.SetFloat(2.71)

// Call functions dynamically
fn := reflect.ValueOf(Add)
result := fn.Call([]reflect.Value{reflect.ValueOf(1), reflect.ValueOf(2)})
```

---

## Standard Library at a Glance

### Encoding & Data

| Package | Purpose |
|---------|---------|
| `encoding/json` | JSON encoding/decoding |
| `encoding/xml` | XML encoding/decoding |
| `encoding/csv` | CSV file handling |
| `encoding/binary` | Binary serialization |
| `encoding/gob` | Go-specific binary format |

### Networking & Web

| Package | Purpose |
|---------|---------|
| `net/http` | HTTP client and server |
| `net/url` | URL parsing |
| `net` | Low-level TCP, UDP, DNS |
| `net/http/httptest` | HTTP testing utilities |
| `html/template` | Safe HTML templating |

### Concurrency & Runtime

| Package | Purpose |
|---------|---------|
| `sync` | Mutexes, WaitGroups |
| `sync/atomic` | Atomic operations |
| `context` | Cancellation and deadlines |
| `runtime` | GC, goroutine info |

### Crypto & Hashing

| Package | Purpose |
|---------|---------|
| `crypto/sha256` | SHA-256 hashing |
| `crypto/rand` | Secure random numbers |
| `crypto/tls` | TLS/SSL |
| `crypto/aes` | AES encryption |

### Files & OS

| Package | Purpose |
|---------|---------|
| `path/filepath` | OS-aware path handling |
| `os/exec` | Run external commands |
| `regexp` | Regular expressions |
| `time` | Time and date operations |
| `database/sql` | SQL database interface |

### Testing

| Package | Purpose |
|---------|---------|
| `testing` | Test framework |
| `testing/quick` | Property-based testing |
| `runtime/pprof` | Profiling |
| `runtime/trace` | Execution tracing |

### Compression

| Package | Purpose |
|---------|---------|
| `compress/gzip` | GZIP |
| `compress/zlib` | ZLIB |
| `archive/tar` | TAR archives |
| `archive/zip` | ZIP archives |

Full reference: [pkg.go.dev/std](https://pkg.go.dev/std)