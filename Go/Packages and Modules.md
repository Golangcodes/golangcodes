![](https://www.notion.so/icons/priority-low_green.svg)


# Packages and Modules in Go: From Basics to Advanced

Go's package and module system is designed for dependency management and code organization. Let's explore it comprehensively.

## Table of Contents

1.  [Basic Package Concepts](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

2.  [Creating Packages](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

3.  [Importing Packages](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

4.  [Visibility Rules](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

5.  [Init Functions](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

6.  [Go Modules](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

7.  [Version Management](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

8.  [Vendor Directory](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

9.  [Advanced Patterns](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

10.  [Best Practices](Packages%20and%20Modules%202091ccc344b1808fa17af38b772f9bcd.html)

## Basic Package Concepts

*   Every Go file belongs to a package

*   `main` package is special (creates executable)

*   Other packages are reusable libraries

*   Packages are organized in directories

```Go
// File: greeter/greet.go
package greeter

func Hello() string {
    return "Hello, world!"
}
```

## Creating Packages

1.  Create a directory for your package

2.  Add Go files with same package declaration

3.  Export functionality using capitalized names

Example structure:

```Plain
mylib/
├── mathutil/
│   ├── math.go
│   └── stats.go
└── stringutil/
    └── strings.go
```

## Importing Packages

Basic imports:

```Go
import (
    "fmt"
    "math/rand"
    "github.com/user/mylib/mathutil"
)
```

Aliased imports:

```Go
import (
    m "math"
    "github.com/user/mylib/mathutil"
)
```

Dot imports (avoid in production):

```Go
import . "fmt" // Now can use Println directly
```

Blank imports (for side effects):

```Go
import _ "image/png" // registers PNG decoder
```

## Visibility Rules

*   Uppercase names: exported (visible outside package)

*   Lowercase names: unexported (package-private)

```Go
package mylib

var privateVar int    // package-private
var PublicVar int     // exported

func privateFunc() {} // package-private
func PublicFunc() {}  // exported
```

## Init Functions

*   Special function that runs when package is imported

*   Can have multiple init() per package

*   Used for initialization tasks

```Go
package mylib

var config map[string]string

func init() {
    config = loadConfig()
}

func loadConfig() map[string]string {
    return map[string]string{"key": "value"}
}
```

## Go Modules

Modules are Go's dependency management system:

1.  Initialize a module:

```Shell
go mod init github.com/user/mymodule
```

1.  Resulting `go.mod` file:

```Plain
module github.com/user/mymodule

go 1.21

require (
    github.com/some/dependency v1.2.3
)
```

1.  Key commands:

```Shell
go mod tidy      # Add missing and remove unused modules
go get -u        # Update dependencies
go list -m all   # List all dependencies
```

## Version Management

*   Semantic versioning (vMAJOR.MINOR.PATCH)

*   Version suffixes:
    
    *   `v1.2.3` - exact version
    
    *   `v1.2` - latest patch
    
    *   `v1` - latest minor
    
    *   `master` - branch name

Example `go get` commands:

```Shell
go get github.com/user/repo@v1.2.3
go get github.com/user/repo@latest
go get github.com/user/repo@master
```

## Vendor Directory

Local copy of dependencies:

```Shell
go mod vendor
```

This creates a `vendor` directory with all dependencies.

To use vendor:

```Shell
go build -mod=vendor
```

## Advanced Patterns

### Internal Packages

Packages in `internal/` are only importable by parent directory:

```Plain
project/
├── internal/
│   └── utils/  # Only importable by project/ code
└── main.go
```

### Submodules

Create nested modules:

```Plain
project/
├── go.mod      # root module
├── submodule/
│   └── go.mod  # submodule
```

### Plugin Architecture

```Go
// main.go
package main

import "plugin"

func loadPlugin(path string) {
    p, err := plugin.Open(path)
    if err != nil {
        panic(err)
    }
    sym, err := p.Lookup("Handler")
    if err != nil {
        panic(err)
    }
    handler := sym.(func(string))
    handler("Hello from plugin")
}
```

### Build Constraints

Control compilation with file suffixes:

```Go
// File: app_linux.go
//go:build linux

package main

func init() {
    // Linux-specific initialization
}
```

## Best Practices

1.  **Keep packages focused** - single responsibility

2.  **Avoid circular dependencies**

3.  **Use semantic versioning** properly

4.  **Document exported elements** with doc comments

5.  **Minimize package-level variables**

6.  **Use internal/ for private packages**

7.  **Pin important dependencies** to exact versions

8.  **Regularly update dependencies** (`go get -u`)

9.  **Vendor dependencies** for production deployments

10.  **Use go.sum** for reproducible builds

## Complete Example

### Module Structure

```Plain
myapp/
├── go.mod
├── go.sum
├── main.go
└── internal/
    └── calculator/
        ├── add.go
        └── subtract.go
```

### go.mod

```Plain
module github.com/user/myapp

go 1.21

require (
    github.com/sirupsen/logrus v1.9.3
)
```

### calculator/add.go

```Go
package calculator

// Add returns the sum of two integers
func Add(a, b int) int {
    return a + b
}
```

### main.go

```Go
package main

import (
    "fmt"
    "github.com/sirupsen/logrus"
    "github.com/user/myapp/internal/calculator"
)

func main() {
    logrus.Info("Starting application")
    sum := calculator.Add(5, 3)
    fmt.Printf("5 + 3 = %d\\n", sum)
}
```
