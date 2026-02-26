# Data Structures

Here some common data structures in Go.

* * *

##  1. **Array**

**Fixed-size ordered collection of elements of the same type.**

```Go
package main

import "fmt"

func main() {
    var nums = [5]int{1, 2, 3, 4, 5}
    fmt.Println(nums[2]) // 3
}
```

 Use when the size is known at compile time.

* * *

##  2. **Slice**

**Dynamic, resizable view into an underlying array.**

```Go
nums := []int{1, 2, 3}
nums = append(nums, 4, 5)
fmt.Println(nums) // [1 2 3 4 5]
```

Core Go data structure — used everywhere.

* * *

## 3. **Map (Hash Map / Dictionary)**

**Key-value store with O(1) average lookup.**

```Go
ages := map[string]int{"Alice": 25, "Bob": 30}
ages["Charlie"] = 35
fmt.Println(ages["Alice"])
```

 Ideal for fast lookups, caches, and frequency counting.

* * *

##  4. **Struct**

**Collection of fields that group related data.**

```Go
type User struct {
    Name string
    Age  int
}
u := User{"Yussuf", 28}
fmt.Println(u.Name)
```

Foundation of Go’s data modeling — replaces classes.

* * *

##  5. **Linked List**

**Sequential nodes linked via pointers.**

```Go
type Node struct {
    Value int
    Next  *Node
}
head := &Node{1, &Node{2, &Node{3, nil}}}
for n := head; n != nil; n = n.Next {
    fmt.Println(n.Value)
}
```

Useful in custom memory management or queue structures.

* * *

##  6. **Doubly Linked List**

**Nodes have pointers to both next and previous nodes.**

```Go
type DNode struct {
    Val  int
    Prev *DNode
    Next *DNode
}
```

Used when you need bidirectional traversal or quick deletions.

* * *

## 7. **Stack**

**LIFO — last in, first out.**

```Go
type Stack []int
func (s *Stack) Push(v int) { *s = append(*s, v) }
func (s *Stack) Pop() int {
    l := len(*s)
    v := (*s)[l-1]
    *s = (*s)[:l-1]
    return v
}
```

Use for recursion, undo functionality, parsing.

* * *

##  8. **Queue**

**FIFO — first in, first out.**

```Go
type Queue []int
func (q *Queue) Enqueue(v int) { *q = append(*q, v) }
func (q *Queue) Dequeue() int {
    v := (*q)[0]
    *q = (*q)[1:]
    return v
}
```

Use in task scheduling, BFS, pipelines.

* * *

##  9. **Deque (Double-Ended Queue)**

**Push/pop from both ends.**

```Go
type Deque struct{ data []int }
```

Useful in sliding window problems, caching, etc.

* * *

## 10. **Binary Tree**

**Each node has up to two children.**

```Go
type Tree struct {
    Val   int
    Left  *Tree
    Right *Tree
}
root := &Tree{1, &Tree{2, nil, nil}, &Tree{3, nil, nil}}
```

Hierarchical data, syntax trees, search trees.

* * *

## 🔍 11. **Binary Search Tree (BST)**

**Ordered binary tree with O(log n) search/insert.**

```Go
func insert(root *Tree, val int) *Tree {
    if root == nil { return &Tree{Val: val} }
    if val < root.Val {
        root.Left = insert(root.Left, val)
    } else {
        root.Right = insert(root.Right, val)
    }
    return root
}
```

Good for ordered data with frequent inserts/lookups.

* * *

##  12. **Heap (Priority Queue)**

**Min-heap or max-heap using container/heap.**

```Go
import (
    "container/heap"
    "fmt"
)

type IntHeap []int
func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() any {
    old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}
func main() {
    h := &IntHeap{3, 1, 4}
    heap.Init(h)
    heap.Push(h, 2)
    fmt.Println(heap.Pop(h)) // 1
}
```

 Scheduling, shortest path, event simulation.

* * *

## 13. **Graph**

**Nodes connected by edges (directed or undirected).**

```Go
graph := map[string][]string{
    "A": {"B", "C"},
    "B": {"A", "D"},
    "C": {"A", "D"},
    "D": {"B", "C"},
}
```

 Use in networks, routing, dependency graphs.

* * *

##  14. **Trie (Prefix Tree)**

**Used for fast prefix lookups (autocomplete).**

```Go
type TrieNode struct {
    Children map[rune]*TrieNode
    End bool
}
func NewNode() *TrieNode {
    return &TrieNode{Children: make(map[rune]*TrieNode)}
}
```

 Ideal for dictionaries, search engines, autocompletion.

* * *

## 15. **Hash Set**

**Map with empty struct values for uniqueness.**

```Go
set := make(map[string]struct{})
set["apple"] = struct{}{}
if _, ok := set["apple"]; ok {
    fmt.Println("Exists")
}
```

Use for deduplication, fast membership checks.

* * *

##  16. **Matrix (2D Array)**

**Grid or table of elements.**

```Go
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
}
fmt.Println(matrix[1][2]) // 6
```

 Use in pathfinding, game boards, image processing.

* * *

## 17. **Circular Queue**

**Queue that wraps around (fixed buffer).**

```Go
type CircularQueue struct {
    data []int
    head, tail, size, capacity int
}
```

 Useful in streaming, buffering, networking.

* * *

##  18. **Bitset / Bitmap**

**Stores boolean flags compactly.**

```Go
var flags uint8
flags |= 1 << 2 // set bit 2
fmt.Println(flags&(1<<2) != 0) // true
```

 Great for memory-efficient flags or sets.

* * *

## 📇 19. **Bloom Filter**

**Probabilistic data structure for membership checks.**

(Usually implemented using bit arrays + hash functions.)

Used in caches, distributed systems to avoid expensive lookups.

* * *

## 20. **Disjoint Set (Union-Find)**

**Tracks connected components efficiently.**

```Go
type DSU struct{ parent []int }
func NewDSU(n int) *DSU {
    p := make([]int, n)
    for i := range p { p[i] = i }
    return &DSU{p}
}
func (d *DSU) Find(x int) int {
    if d.parent[x] != x { d.parent[x] = d.Find(d.parent[x]) }
    return d.parent[x]
}
func (d *DSU) Union(a, b int) { d.parent[d.Find(a)] = d.Find(b) }
```

 Used in Kruskal’s algorithm, clustering, networks.

* * *

## 21. **Segment Tree**

**Tree structure for range queries.**

Use for range sums, min/max, or lazy propagation in algorithms.

* * *

##  22. **Fenwick Tree (Binary Indexed Tree)**

**Efficient prefix sum structure.**

Similar to Segment Tree but easier to implement.

* * *

##  23. **Ring Buffer**

**Fixed-size circular buffer.**

```Go
import "container/ring"

r := ring.New(3)
for i := 0; i < r.Len(); i++ {
    r.Value = i
    r = r.Next()
}
```

 Great for log streaming, buffering, rate limiting.

* * *

##  24. **LRU Cache**

**Cache with least-recently-used eviction.**

```Go
import "container/list"

type LRU struct {
    cap  int
    list *list.List
    data map[string]*list.Element
}
```

 Used in systems, databases, and caching layers.

* * *

##  25. **Queue Channel (Concurrency Primitive)**

**Go’s built-in concurrent queue.**

```Go
jobs := make(chan int, 5)
go func() {
    for j := range jobs {
        fmt.Println("Processing:", j)
    }
}()
jobs <- 1
jobs <- 2
close(jobs)
```

 Concurrency-safe queue; ideal for pipelines and workers.

* * *

## 🚀 Contribute to this Syllabus!

The best way to master Data Structures and Algorithms is by implementing them yourself. 

We highly encourage the community to contribute practical, runnable Go examples for these data structures! If you've written a clean implementation of a Linked List, a Trie, or a Graph in Go, click the **"Edit this page on GitHub"** button at the bottom of the screen.

Submit a Pull Request adding your code snippet to the relevant section above, and we'll review and merge it for everyone to learn from!