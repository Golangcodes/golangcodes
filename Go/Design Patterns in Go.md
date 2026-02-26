# Design Patterns in Go

Design patterns are reusable solutions to common software engineering problems. Go's simplicity — composition over inheritance, first-class functions, interfaces, and goroutines — makes many classic OOP patterns either unnecessary or elegantly simple to implement.

This lesson covers the most practical patterns you'll encounter in real Go codebases.

---

## 1. Functional Options Pattern

Use when you have a struct with many optional configuration fields and want a clean, extensible constructor.

```Go
type Server struct {
	host string
	port int
	tls  bool
}

type Option func(*Server)

func WithHost(h string) Option {
	return func(s *Server) { s.host = h }
}

func WithPort(p int) Option {
	return func(s *Server) { s.port = p }
}

func WithTLS(enabled bool) Option {
	return func(s *Server) { s.tls = enabled }
}

func NewServer(opts ...Option) *Server {
	s := &Server{host: "localhost", port: 8080}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Usage
srv := NewServer(WithHost("0.0.0.0"), WithTLS(true))
```

https://youtu.be/SaeYzGL3370

---

## 2. Builder Pattern

Use when constructing complex objects step-by-step with a fluent API.

```Go
type ContainerSpec struct {
	Image   string
	CPU     int
	Memory  int
	Env     map[string]string
	Ports   []int
}

type ContainerBuilder struct {
	spec *ContainerSpec
}

func NewContainerBuilder() *ContainerBuilder {
	return &ContainerBuilder{
		spec: &ContainerSpec{Env: make(map[string]string)},
	}
}

func (b *ContainerBuilder) SetImage(img string) *ContainerBuilder {
	b.spec.Image = img
	return b
}

func (b *ContainerBuilder) SetCPU(cpu int) *ContainerBuilder {
	b.spec.CPU = cpu
	return b
}

func (b *ContainerBuilder) SetMemory(mem int) *ContainerBuilder {
	b.spec.Memory = mem
	return b
}

func (b *ContainerBuilder) AddEnv(key, value string) *ContainerBuilder {
	b.spec.Env[key] = value
	return b
}

func (b *ContainerBuilder) AddPort(port int) *ContainerBuilder {
	b.spec.Ports = append(b.spec.Ports, port)
	return b
}

func (b *ContainerBuilder) Build() ContainerSpec {
	return *b.spec // return a copy for immutability
}
```

```Go
// Usage
spec := NewContainerBuilder().
	SetImage("nginx:latest").
	SetCPU(4).
	SetMemory(2048).
	AddEnv("ENV", "production").
	AddPort(8080).
	AddPort(443).
	Build()

fmt.Printf("%+v\n", spec)
```

https://youtu.be/g8TMyTUhT08

---

## 3. Worker Pool Pattern

Use when you need to process many jobs concurrently without spawning unlimited goroutines.

```Go
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		results <- j * 2
	}
}

func main() {
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// Start 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for a := 1; a <= 5; a++ {
		fmt.Println(<-results)
	}
}
```

https://youtu.be/ZWMiKQXmh9s

---

## 4. Context Pattern (Timeout & Cancellation)

Use for clean cancellation and timeouts in API calls, background workers, and HTTP handlers.

```Go
func fetch(ctx context.Context) error {
	select {
	case <-time.After(3 * time.Second):
		return nil // work completed
	case <-ctx.Done():
		return ctx.Err() // cancelled or timed out
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := fetch(ctx); err != nil {
		log.Println("fetch error:", err)
	}
}
```

---

## 5. Fan-In / Fan-Out

**Fan-Out**: Distribute work across multiple goroutines.
**Fan-In**: Merge results from multiple channels into one.

```Go
// Fan-Out: multiple workers reading from one channel
func worker(id int, jobs <-chan int, out chan<- int) {
	for j := range jobs {
		out <- j * j
	}
}

// Fan-In: merge multiple channels into one
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	wg.Add(len(cs))
	for _, c := range cs {
		go func(ch <-chan int) {
			for v := range ch {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
```

---

## 6. Middleware Chaining

Use for composable, reusable HTTP handler logic (logging, auth, CORS, etc.).

```Go
type Middleware func(http.Handler) http.Handler

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello")
	})

	http.Handle("/", Chain(handler, LoggingMiddleware))
	http.ListenAndServe(":8080", nil)
}
```

---

## 7. Command Pattern

Use to encapsulate actions as objects. Common in CLI tools, task queues, and undo/redo systems.

```Go
type Command interface {
	Execute() error
}

type DeployCommand struct{}

func (d DeployCommand) Execute() error {
	fmt.Println("Deploying...")
	return nil
}

type RollbackCommand struct{}

func (r RollbackCommand) Execute() error {
	fmt.Println("Rolling back...")
	return nil
}

func run(cmd Command) {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Usage
run(DeployCommand{})
run(RollbackCommand{})
```

---

## 8. Pub/Sub via Channels

Use for decoupled, event-driven communication between components.

```Go
type EventBus struct {
	subscribers map[string][]chan string
	lock        sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{subscribers: make(map[string][]chan string)}
}

func (eb *EventBus) Subscribe(topic string) <-chan string {
	ch := make(chan string, 10)
	eb.lock.Lock()
	eb.subscribers[topic] = append(eb.subscribers[topic], ch)
	eb.lock.Unlock()
	return ch
}

func (eb *EventBus) Publish(topic, msg string) {
	eb.lock.RLock()
	for _, ch := range eb.subscribers[topic] {
		ch <- msg
	}
	eb.lock.RUnlock()
}

// Usage
bus := NewEventBus()
ch := bus.Subscribe("alerts")

go func() {
	for msg := range ch {
		fmt.Println("Alert:", msg)
	}
}()

bus.Publish("alerts", "CPU usage at 90%")
```

---

## 9. Plugin Architecture

Use to build extensible systems where functionality can be added without modifying the core.

```Go
type Plugin interface {
	Name() string
	Run() error
}

type LoggingPlugin struct{}

func (l LoggingPlugin) Name() string { return "logger" }
func (l LoggingPlugin) Run() error   { fmt.Println("Logging..."); return nil }

type PluginManager struct {
	plugins map[string]Plugin
}

func (pm *PluginManager) Register(p Plugin) {
	if pm.plugins == nil {
		pm.plugins = make(map[string]Plugin)
	}
	pm.plugins[p.Name()] = p
}

func (pm *PluginManager) RunAll() {
	for _, p := range pm.plugins {
		p.Run()
	}
}

// Usage
pm := &PluginManager{}
pm.Register(LoggingPlugin{})
pm.RunAll()
```

---

## 10. Event Sourcing (Simplified)

Use when you need full audit trails, state replay, or versioning. Store events instead of state.

```Go
type Event struct {
	Type string
	Data string
}

type Store struct {
	events []Event
}

func (s *Store) Record(e Event) {
	s.events = append(s.events, e)
}

func (s *Store) Replay() {
	for _, e := range s.events {
		fmt.Printf("[%s] %s\n", e.Type, e.Data)
	}
}

// Usage
store := &Store{}
store.Record(Event{"CREATE", "user:123"})
store.Record(Event{"UPDATE", "user:123 name=Alice"})
store.Replay()
```

---

## Quick Reference

| Pattern | Best For |
|---------|----------|
| Functional Options | Clean constructors with optional config |
| Builder | Complex object construction with fluent API |
| Worker Pool | Bounded concurrent job processing |
| Context | Cancellation, timeouts, request-scoped values |
| Fan-In / Fan-Out | Parallelism and result aggregation |
| Middleware | Composable HTTP/RPC handler chains |
| Command | Encapsulating actions as objects |
| Pub/Sub | Decoupled event-driven communication |
| Plugin | Extensible, modular architectures |
| Event Sourcing | Audit logs, state replay, versioning |

---

## Additional Patterns in Go

These patterns are less commonly needed but worth knowing:

- **Singleton** — Use `sync.Once` to ensure a single instance
- **Factory** — Return interface types from constructor functions
- **Strategy** — Pass function types or interfaces to swap algorithms
- **Decorator** — Wrap structs or functions to add behavior
- **Adapter** — Bridge incompatible interfaces for third-party integrations
- **Facade** — Simplify complex subsystems behind a single interface
- **Proxy** — Control access to objects (caching, lazy loading, protection)
- **Iterator** — Largely replaced by Go's built-in `range` keyword