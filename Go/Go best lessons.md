# Go Best Practices & Patterns

A quick-reference catalog of common Go patterns, organized by category. Each pattern is covered in detail in its respective lesson.

---

## Concurrency Patterns

- **Worker Pool** — Distribute jobs across a fixed number of goroutines
- **Fan-Out / Fan-In** — Split work among workers, merge results into one channel
- **Pipeline** — Chain processing stages connected by channels
- **Generator** — Function that returns a channel producing values on demand
- **Timeout** — Use `select` with `time.After` for operation deadlines
- **Quit Channel** — Signal goroutines to stop via `chan struct{}`
- **Semaphore** — Limit concurrent access using a buffered channel
- **Pub/Sub** — Broadcast messages to multiple subscribers

---

## Error Handling Patterns

- **Error Wrapping** — Add context with `fmt.Errorf("doing X: %w", err)`
- **Sentinel Errors** — Predefined error values (`var ErrNotFound = errors.New(...)`)
- **Custom Error Types** — Structs implementing `error` with extra fields
- **Error Aggregation** — Collect multiple errors into a `MultiError`
- **Retry with Backoff** — Reattempt failed operations with exponential delay
- **Circuit Breaker** — Stop retrying after repeated failures

---

## Structural Patterns

- **Functional Options** — Clean constructors with `Option` functions
- **Builder** — Fluent API for step-by-step object construction
- **Middleware** — Chain processing layers (common in HTTP servers)
- **Dependency Injection** — Pass dependencies explicitly via interfaces
- **Decorator** — Wrap functions to add behavior (logging, timing, etc.)

---

## Data Patterns

- **Map as Set** — `map[T]struct{}` for O(1) membership checks
- **Slice Filtering** — `s = append(s[:i], s[i+1:]...)`
- **Pre-allocation** — `make([]T, 0, capacity)` to avoid repeated allocations
- **Memoization** — Cache expensive function results in a map
- **Sharded Map** — Split a map across shards for concurrent access
- **Object Pool** — `sync.Pool` for reusing expensive objects

---

## Interface Patterns

- **Interface-Based Polymorphism** — Different types satisfying the same interface
- **Embedding** — Composition over inheritance via struct embedding
- **Strategy** — Swap algorithms at runtime via interfaces
- **Null Object** — Default no-op implementation (e.g., `io.Discard`)

```Go
// Singleton with sync.Once
var (
    instance *DB
    once     sync.Once
)

func GetDB() *DB {
    once.Do(func() { instance = &DB{} })
    return instance
}
```

---

## Testing Patterns

- **Table-Driven Tests** — Test cases in a slice, loop through them

```Go
func TestAdd(t *testing.T) {
    cases := []struct{ a, b, want int }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }
    for _, c := range cases {
        got := Add(c.a, c.b)
        if got != c.want {
            t.Errorf("Add(%d, %d) = %d, want %d", c.a, c.b, got, c.want)
        }
    }
}
```

- **Interface Mocks** — Fake implementations for testing without real dependencies
- **Golden Files** — Compare output against saved reference files
- **Test Fixtures** — Setup and teardown state for each test

---

## Generic Patterns (Go 1.18+)

- **Generic Data Structures** — Reusable containers like `Stack[T]`, `Set[T]`
- **Type Constraints** — Define allowed operations via interface constraints
- **Type-Safe Functions** — Algorithms that work across multiple types

---

## API Design Patterns

- **Return `(value, error)`** — Always return errors as the last value
- **Context Propagation** — Pass `context.Context` as the first parameter
- **Options Pattern** — Configurable initialization without breaking changes
- **`chan struct{}` for signaling** — Zero-allocation done/quit channels
