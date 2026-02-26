# Error Handling in Go

Go uses explicit error returns instead of exceptions. Errors are values — they implement the `error` interface and are returned alongside results. This forces you to handle failures at every step.

```Go
type error interface {
    Error() string
}
```

---

## Basic Error Handling

Functions return errors as the last return value. Always check them.

```Go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := divide(10, 0)
if err != nil {
    log.Printf("Error: %v", err)
}
```

---

## Creating Errors

### Simple Errors

```Go
var ErrNotFound = errors.New("not found")
var ErrTimeout  = errors.New("operation timed out")
```

### Custom Error Types

Use when you need to carry structured information with the error.

```Go
type InputError struct {
    Field   string
    Message string
}

func (e *InputError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func validateInput(input string) error {
    if len(input) < 5 {
        return &InputError{
            Field:   "input",
            Message: "must be at least 5 characters",
        }
    }
    return nil
}
```

---

## Error Wrapping (Go 1.13+)

Use `%w` in `fmt.Errorf` to wrap errors with context while preserving the original error for inspection.

```Go
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("failed to open %s: %w", path, err)
    }
    defer f.Close()
    // ...
    return nil
}
```

Wrapped errors build a chain of context:

```
"failed to process request"
  └─ "failed to open config.json"
     └─ "no such file or directory"
```

---

## Error Inspection

### `errors.Is` — Compare with sentinel errors

```Go
if errors.Is(err, ErrNotFound) {
    // handle not found case
}

if errors.Is(err, os.ErrNotExist) {
    // file doesn't exist
}
```

### `errors.As` — Extract a specific error type

```Go
var inputErr *InputError
if errors.As(err, &inputErr) {
    fmt.Println("Invalid field:", inputErr.Field)
}
```

### `errors.Unwrap` — Get the underlying error

```Go
if unwrapped := errors.Unwrap(err); unwrapped != nil {
    fmt.Println("Original error:", unwrapped)
}
```

---

## Sentinel Errors

Predefined error values for common, expected conditions. Think of them like HTTP status codes.

```Go
var (
    ErrInvalidInput = errors.New("invalid input")
    ErrNotFound     = errors.New("not found")
    ErrTimeout      = errors.New("operation timed out")
)

func fetchData() error {
    return ErrTimeout
}

if err := fetchData(); errors.Is(err, ErrTimeout) {
    log.Println("retrying after timeout...")
}
```

---

## Panic and Recover

`panic` stops normal execution. `recover` catches it inside a deferred function. Use panic only for truly unrecoverable errors (programmer bugs, impossible states).

```Go
func criticalOperation() {
    panic("invariant violated: negative balance")
}

func safeOperation() (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()

    criticalOperation()
    return nil
}
```

---

## Error Aggregation

Collect multiple errors when validating or processing batches.

```Go
type MultiError struct {
    Errors []error
}

func (m *MultiError) Error() string {
    var msgs []string
    for _, err := range m.Errors {
        msgs = append(msgs, err.Error())
    }
    return strings.Join(msgs, "; ")
}

func (m *MultiError) Add(err error) {
    if err != nil {
        m.Errors = append(m.Errors, err)
    }
}

func validate(input string) error {
    var merr MultiError
    if len(input) < 5 {
        merr.Add(errors.New("too short"))
    }
    if !strings.ContainsAny(input, "!@#$") {
        merr.Add(errors.New("missing special character"))
    }
    if len(merr.Errors) > 0 {
        return &merr
    }
    return nil
}
```

---

## Complete Example

```Go
package main

import (
    "errors"
    "fmt"
    "log"
    "os"
)

var ErrNotFound = errors.New("file not found")

type FileError struct {
    Path   string
    Op     string
    Reason string
    Cause  error
}

func (e *FileError) Error() string {
    return fmt.Sprintf("%s %s: %s", e.Op, e.Path, e.Reason)
}

func (e *FileError) Unwrap() error {
    return e.Cause
}

func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return &FileError{
            Path:   path,
            Op:     "open",
            Reason: "failed to access file",
            Cause:  err,
        }
    }
    defer f.Close()
    return nil
}

func main() {
    err := processFile("missing.txt")
    if err != nil {
        var ferr *FileError
        if errors.As(err, &ferr) {
            log.Printf("File operation failed: %v", ferr)
            if errors.Is(ferr.Cause, os.ErrNotExist) {
                log.Println("The file does not exist")
            }
        }
    }
}
```

---

## Best Practices

1. Always check errors — never use `_` to discard them
2. Wrap errors with context using `fmt.Errorf("doing X: %w", err)`
3. Use `errors.Is` and `errors.As` for inspection (not `==` or type assertions)
4. Define sentinel errors as package-level variables
5. Use custom error types when you need structured error data
6. Use `panic` only for programmer errors, never for expected failures
7. Return `error` as the last return value
8. Keep error messages lowercase and without trailing punctuation
9. Each layer should add context, not repeat information
10. Design your errors — they're part of your API