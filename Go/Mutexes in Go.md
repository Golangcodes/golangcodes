![](https://www.notion.so/icons/aries_green.svg)


# Mutexes in Go: From Basics to Advanced

Mutexes (mutual exclusion locks) are synchronization primitives that protect shared resources in concurrent Go programs. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Mutex Concepts](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

2.  [sync.Mutex](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

3.  [sync.RWMutex](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

4.  [Mutex Patterns](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

5.  [Deadlocks](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

6.  [Mutex vs Channels](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

7.  [Advanced Patterns](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

8.  [Performance Considerations](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

9.  [Best Practices](Mutexes%20in%20Go%202091ccc344b180a4865bd6ba69023269.html)

## Basic Mutex Concepts

*   Protects shared resources from concurrent access

*   Only one goroutine can hold the lock at a time

*   Other goroutines block until lock is released

*   Prevents race conditions

```Go
var count int
var mu sync.Mutex

func increment() {
    mu.Lock()
    count++
    mu.Unlock()
}
```

## sync.Mutex

Basic mutual exclusion lock:

```Go
var mu sync.Mutex
var sharedData int

func update() {
    mu.Lock()         // Acquire lock
    sharedData = 42   // Critical section
    mu.Unlock()       // Release lock
}
```

Important methods:

*   `Lock()` - Acquires the mutex (blocks if already locked)

*   `Unlock()` - Releases the mutex

*   `TryLock()` (Go 1.18+) - Non-blocking lock attempt

## sync.RWMutex

Reader/writer mutual exclusion lock:

*   Multiple readers can hold lock simultaneously

*   Only one writer can hold lock (exclusive)

```Go
var rwmu sync.RWMutex
var config map[string]string

func readConfig(key string) string {
    rwmu.RLock()         // Reader lock
    defer rwmu.RUnlock() // Ensure unlock happens
    return config[key]
}

func updateConfig(key, value string) {
    rwmu.Lock()         // Writer lock
    defer rwmu.Unlock()
    config[key] = value
}
```

Methods:

*   `RLock()` - Acquire read lock

*   `RUnlock()` - Release read lock

*   `Lock()` - Acquire write lock

*   `Unlock()` - Release write lock

*   `TryLock()`, `TryRLock()` (Go 1.18+) - Non-blocking attempts

## Mutex Patterns

### Protecting Shared Data

```Go
type SafeCounter struct {
    mu    sync.Mutex
    count int
}

func (c *SafeCounter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

### Lazy Initialization

```Go
var resource *Resource
var once sync.Once

func GetResource() *Resource {
    once.Do(func() {
        resource = &Resource{}
    })
    return resource
}
```

## Deadlocks

Common causes:

1.  Locking mutex twice in same goroutine

2.  Circular wait between goroutines

3.  Forgetting to unlock

Example deadlock:

```Go
var mu sync.Mutex

func main() {
    mu.Lock()
    mu.Lock() // Deadlock - blocks forever
}
```

Debugging:

*   Use Go's deadlock detector (`go run -race`)

*   Avoid holding locks while doing I/O

*   Use `defer` for unlocks where possible

## Mutex vs Channels

When to use mutexes:

*   Protecting simple state

*   When performance is critical

*   For cache implementations

*   When you need read/write distinctions

When to use channels:

*   Transferring ownership of data

*   Communicating between goroutines

*   Implementing pipelines

*   Complex coordination patterns

## Advanced Patterns

### Mutex with Timeout

```Go
func tryWithTimeout(mu *sync.Mutex, timeout time.Duration) bool {
    ch := make(chan struct{})
    go func() {
        mu.Lock()
        close(ch)
    }()

    select {
    case <-ch:
        return true
    case <-time.After(timeout):
        return false
    }
}
```

### Scoped Lock

```Go
func withLock(mu *sync.Mutex, f func()) {
    mu.Lock()
    defer mu.Unlock()
    f()
}

// Usage:
withLock(&mu, func() {
    // Critical section
})
```

### Recursive Mutex

```Go
type RecursiveMutex struct {
    mu     sync.Mutex
    owner  int64 // goroutine id
    count  int
}

func (m *RecursiveMutex) Lock() {
    gid := goid.Get() // hypothetical goroutine ID
    if atomic.LoadInt64(&m.owner) == gid {
        m.count++
        return
    }
    m.mu.Lock()
    atomic.StoreInt64(&m.owner, gid)
    m.count = 1
}

func (m *RecursiveMutex) Unlock() {
    gid := goid.Get()
    if atomic.LoadInt64(&m.owner) != gid {
        panic("unlock of unlocked mutex")
    }
    m.count--
    if m.count == 0 {
        atomic.StoreInt64(&m.owner, 0)
        m.mu.Unlock()
    }
}
```

## Performance Considerations

1.  **Mutex contention** is expensive (minimize locked sections)

2.  **RWMutex** is better for read-heavy loads

3.  **Channel overhead** is higher than mutex for simple cases

4.  **Sync.Pool** can help with allocation pressure

5.  **Atomic operations** are faster for simple counters

Benchmark example:

```Go
func BenchmarkMutex(b *testing.B) {
    var mu sync.Mutex
    var count int

    for i := 0; i < b.N; i++ {
        mu.Lock()
        count++
        mu.Unlock()
    }
}

func BenchmarkRWMutex(b *testing.B) {
    var mu sync.RWMutex
    var count int

    for i := 0; i < b.N; i++ {
        mu.Lock()
        count++
        mu.Unlock()
    }
}
```

## Best Practices

1.  **Use defer** to ensure mutexes are unlocked

2.  **Keep critical sections** as small as possible

3.  **Document** which mutex protects which data

4.  **Avoid nested locks** where possible

5.  **Prefer RWMutex** for read-heavy scenarios

6.  **Consider channels** for complex synchronization

7.  **Use the race detector** (`go run -race`)

8.  **Profile contention** in performance-critical code

9.  **Avoid mutex copies** - pass by pointer

10.  **Zero-value mutexes** are usable (no initialization needed)

## Complete Example

```Go
package main

import (
	"fmt"
	"sync"
	"time"
)

type BankAccount struct {
	balance int
	mu      sync.Mutex
}

func (a *BankAccount) Deposit(amount int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.balance += amount
}

func (a *BankAccount) Withdraw(amount int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.balance >= amount {
		a.balance -= amount
		return true
	}
	return false
}

func (a *BankAccount) Balance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func main() {
	account := &BankAccount{balance: 1000}

	var wg sync.WaitGroup
	wg.Add(2)

	// Depositor
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			account.Deposit(100)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Withdrawer
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			if account.Withdraw(150) {
				fmt.Println("Withdrawal successful")
			} else {
				fmt.Println("Withdrawal failed")
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	wg.Wait()
	fmt.Println("Final balance:", account.Balance())

	// RWMutex example
	var cache = struct {
		sync.RWMutex
		items map[string]string
	}{items: make(map[string]string)}

	// Writer
	cache.Lock()
	cache.items["key"] = "value"
	cache.Unlock()

	// Readers
	var wg2 sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			cache.RLock()
			fmt.Println("Cache value:", cache.items["key"])
			cache.RUnlock()
		}()
	}
	wg2.Wait()
}
```
