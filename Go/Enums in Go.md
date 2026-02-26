![](https://www.notion.so/icons/computer-chip_green.svg)

# Enums in Go

# Enums in Go: From Basics to Advanced

While Go doesn't have a built-in `enum` type like some other languages, it provides several idiomatic ways to implement enumeration patterns. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Enum Concepts](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

2.  [Iota-Based Enums](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

3.  [String Enums](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

4.  [Custom Types](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

5.  [Enum Methods](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

6.  [Validation](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

7.  [JSON Marshaling](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

8.  [Advanced Patterns](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

9.  [Best Practices](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

## Basic Enum Concepts

Go implements enums using:

*   `const` blocks with `iota`

*   Custom types with constants

*   String representations

## Iota-Based Enums

The most common approach using `iota` for auto-incrementing values:

```Go
type Weekday int

const (
    Sunday Weekday = iota  // 0
    Monday                 // 1
    Tuesday                // 2
    Wednesday              // 3
    Thursday               // 4
    Friday                 // 5
    Saturday               // 6
)
```

### Advanced Iota Usage

```Go
type BitFlag int

const (
    Active BitFlag = 1 << iota  // 1 << 0 = 1
    Send                        // 1 << 1 = 2
    Receive                     // 1 << 2 = 4
)

// Combined flags
const (
    All = Active | Send | Receive  // 7
)
```

## String Enums

For more readable representations:

```Go
type Color string

const (
    Red    Color = "RED"
    Green  Color = "GREEN"
    Blue   Color = "BLUE"
)
```

## Custom Types

Creating type-safe enums:

```Go
type Status struct {
    name string
}

var (
    Pending = Status{"PENDING"}
    Active  = Status{"ACTIVE"}
    Closed  = Status{"CLOSED"}
)
```

## Enum Methods

Adding behavior to enums:

```Go
func (d Weekday) String() string {
    return [...]string{"Sunday", "Monday", "Tuesday",
        "Wednesday", "Thursday", "Friday", "Saturday"}[d]
}

func (d Weekday) Weekend() bool {
    return d == Saturday || d == Sunday
}
```

## Validation

Ensuring valid enum values:

```Go
func (c Color) IsValid() bool {
    switch c {
    case Red, Green, Blue:
        return true
    }
    return false
}

func SetColor(c Color) error {
    if !c.IsValid() {
        return fmt.Errorf("invalid color: %v", c)
    }
    // ...
    return nil
}
```

## JSON Marshaling

Proper JSON handling for enums:

```Go
// Marshal
func (c Color) MarshalJSON() ([]byte, error) {
    return json.Marshal(string(c))
}

// Unmarshal
func (c *Color) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    *c = Color(s)
    if !c.IsValid() {
        return fmt.Errorf("invalid color: %s", s)
    }
    return nil
}
```

## Advanced Patterns

### Enum Generators

```Go
func NewStatus(name string) (Status, error) {
    switch name {
    case "PENDING":
        return Pending, nil
    case "ACTIVE":
        return Active, nil
    case "CLOSED":
        return Closed, nil
    default:
        return Status{}, fmt.Errorf("invalid status: %s", name)
    }
}
```

### Enum Collections

```Go
var AllWeekdays = []Weekday{
    Sunday, Monday, Tuesday, Wednesday,
    Thursday, Friday, Saturday,
}

func (d Weekday) Next() Weekday {
    return AllWeekdays[(d+1)%len(AllWeekdays)]
}
```

### Database Enum Mapping

```Go
func (d *Weekday) Scan(value interface{}) error {
    if v, ok := value.(int64); ok {
        *d = Weekday(v)
        return nil
    }
    return fmt.Errorf("invalid weekday value")
}

func (d Weekday) Value() (driver.Value, error) {
    return int64(d), nil
}
```

## Best Practices

1.  **Use** `**iota**` **for simple numeric enums**

2.  **Implement** `**String()**` **for readable output**

3.  **Add validation methods**

4.  **Consider string enums for JSON APIs**

5.  **Document all enum values**

6.  **Handle JSON marshaling/unmarshaling**

7.  **Add database support if needed**

8.  **Keep enum sets small and focused**

9.  **Use constants for enum values**

10.  **Consider code generation for large enums**

## Complete Example

```Go
package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Iota-based enum with string representation
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var weekdays = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func (d Weekday) String() string {
	if d < Sunday || d > Saturday {
		return "Invalid"
	}
	return weekdays[d]
}

func (d Weekday) Weekend() bool {
	return d == Saturday || d == Sunday
}

// String-based enum with validation
type Color string

const (
	Red   Color = "RED"
	Green Color = "GREEN"
	Blue  Color = "BLUE"
)

func (c Color) IsValid() bool {
	switch c {
	case Red, Green, Blue:
		return true
	}
	return false
}

func (c *Color) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = Color(s)
	if !c.IsValid() {
		return fmt.Errorf("invalid color: %s", s)
	}
	return nil
}

// Custom type enum with methods
type Status struct {
	name string
}

var (
	Pending = &Status{"PENDING"}
	Active  = &Status{"ACTIVE"}
	Closed  = &Status{"CLOSED"}
)

func (s *Status) String() string {
	return s.name
}

func (s *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.name)
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	switch name {
	case "PENDING":
		*s = *Pending
	case "ACTIVE":
		*s = *Active
	case "CLOSED":
		*s = *Closed
	default:
		return errors.New("invalid status")
	}
	return nil
}

func (s *Status) Scan(value interface{}) error {
	if name, ok := value.(string); ok {
		switch name {
		case "PENDING":
			*s = *Pending
		case "ACTIVE":
			*s = *Active
		case "CLOSED":
			*s = *Closed
		default:
			return errors.New("invalid status")
		}
		return nil
	}
	return errors.New("invalid type for status")
}

func (s *Status) Value() (driver.Value, error) {
	return s.name, nil
}

func main() {
	// Using Weekday enum
	day := Tuesday
	fmt.Printf("%s is weekend? %t\\n", day, day.Weekend())

	// Using Color enum
	color := Green
	fmt.Printf("Color: %s, valid: %t\\n", color, color.IsValid())

	// Using Status enum
	status := Active
	fmt.Printf("Status: %s\\n", status)

	// JSON marshaling
	jsonData, _ := json.Marshal(struct {
		Day    Weekday
		Color  Color
		Status *Status
	}{day, color, status})
	fmt.Println("JSON:", string(jsonData))
}
```

This comprehensive guide covers enums in Go from basic to advanced concepts. Would you like me to focus on any particular aspect in more detail?

# Enums in Go: From Basics to Advanced

While Go doesn't have a built-in `enum` type like some other languages, it provides several idiomatic ways to implement enumeration patterns. Let's explore them comprehensively.

## Table of Contents

1.  [Basic Enum Concepts](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

2.  [Iota-Based Enums](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

3.  [String Enums](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

4.  [Custom Types](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

5.  [Enum Methods](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

6.  [Validation](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

7.  [JSON Marshaling](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

8.  [Advanced Patterns](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

9.  [Best Practices](Enums%20in%20Go%202091ccc344b180aa8679f1dd3b57b930.html)

## Basic Enum Concepts

Go implements enums using:

*   `const` blocks with `iota`

*   Custom types with constants

*   String representations

## Iota-Based Enums

The most common approach using `iota` for auto-incrementing values:

```Go
type Weekday int

const (
    Sunday Weekday = iota  // 0
    Monday                 // 1
    Tuesday                // 2
    Wednesday              // 3
    Thursday               // 4
    Friday                 // 5
    Saturday               // 6
)
```

### Advanced Iota Usage

```Go
type BitFlag int

const (
    Active BitFlag = 1 << iota  // 1 << 0 = 1
    Send                        // 1 << 1 = 2
    Receive                     // 1 << 2 = 4
)

// Combined flags
const (
    All = Active | Send | Receive  // 7
)
```

## String Enums

For more readable representations:

```Go
type Color string

const (
    Red    Color = "RED"
    Green  Color = "GREEN"
    Blue   Color = "BLUE"
)
```

## Custom Types

Creating type-safe enums:

```Go
type Status struct {
    name string
}

var (
    Pending = Status{"PENDING"}
    Active  = Status{"ACTIVE"}
    Closed  = Status{"CLOSED"}
)
```

## Enum Methods

Adding behavior to enums:

```Go
func (d Weekday) String() string {
    return [...]string{"Sunday", "Monday", "Tuesday",
        "Wednesday", "Thursday", "Friday", "Saturday"}[d]
}

func (d Weekday) Weekend() bool {
    return d == Saturday || d == Sunday
}
```

## Validation

Ensuring valid enum values:

```Go
func (c Color) IsValid() bool {
    switch c {
    case Red, Green, Blue:
        return true
    }
    return false
}

func SetColor(c Color) error {
    if !c.IsValid() {
        return fmt.Errorf("invalid color: %v", c)
    }
    // ...
    return nil
}
```

## JSON Marshaling

Proper JSON handling for enums:

```Go
// Marshal
func (c Color) MarshalJSON() ([]byte, error) {
    return json.Marshal(string(c))
}

// Unmarshal
func (c *Color) UnmarshalJSON(data []byte) error {
    var s string
    if err := json.Unmarshal(data, &s); err != nil {
        return err
    }
    *c = Color(s)
    if !c.IsValid() {
        return fmt.Errorf("invalid color: %s", s)
    }
    return nil
}
```

## Advanced Patterns

### Enum Generators

```Go
func NewStatus(name string) (Status, error) {
    switch name {
    case "PENDING":
        return Pending, nil
    case "ACTIVE":
        return Active, nil
    case "CLOSED":
        return Closed, nil
    default:
        return Status{}, fmt.Errorf("invalid status: %s", name)
    }
}
```

### Enum Collections

```Go
var AllWeekdays = []Weekday{
    Sunday, Monday, Tuesday, Wednesday,
    Thursday, Friday, Saturday,
}

func (d Weekday) Next() Weekday {
    return AllWeekdays[(d+1)%len(AllWeekdays)]
}
```

### Database Enum Mapping

```Go
func (d *Weekday) Scan(value interface{}) error {
    if v, ok := value.(int64); ok {
        *d = Weekday(v)
        return nil
    }
    return fmt.Errorf("invalid weekday value")
}

func (d Weekday) Value() (driver.Value, error) {
    return int64(d), nil
}
```

## Best Practices

1.  **Use** `**iota**` **for simple numeric enums**

2.  **Implement** `**String()**` **for readable output**

3.  **Add validation methods**

4.  **Consider string enums for JSON APIs**

5.  **Document all enum values**

6.  **Handle JSON marshaling/unmarshaling**

7.  **Add database support if needed**

8.  **Keep enum sets small and focused**

9.  **Use constants for enum values**

10.  **Consider code generation for large enums**

## Complete Example

```Go
package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

// Iota-based enum with string representation
type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

var weekdays = [...]string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

func (d Weekday) String() string {
	if d < Sunday || d > Saturday {
		return "Invalid"
	}
	return weekdays[d]
}

func (d Weekday) Weekend() bool {
	return d == Saturday || d == Sunday
}

// String-based enum with validation
type Color string

const (
	Red   Color = "RED"
	Green Color = "GREEN"
	Blue  Color = "BLUE"
)

func (c Color) IsValid() bool {
	switch c {
	case Red, Green, Blue:
		return true
	}
	return false
}

func (c *Color) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*c = Color(s)
	if !c.IsValid() {
		return fmt.Errorf("invalid color: %s", s)
	}
	return nil
}

// Custom type enum with methods
type Status struct {
	name string
}

var (
	Pending = &Status{"PENDING"}
	Active  = &Status{"ACTIVE"}
	Closed  = &Status{"CLOSED"}
)

func (s *Status) String() string {
	return s.name
}

func (s *Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.name)
}

func (s *Status) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	switch name {
	case "PENDING":
		*s = *Pending
	case "ACTIVE":
		*s = *Active
	case "CLOSED":
		*s = *Closed
	default:
		return errors.New("invalid status")
	}
	return nil
}

func (s *Status) Scan(value interface{}) error {
	if name, ok := value.(string); ok {
		switch name {
		case "PENDING":
			*s = *Pending
		case "ACTIVE":
			*s = *Active
		case "CLOSED":
			*s = *Closed
		default:
			return errors.New("invalid status")
		}
		return nil
	}
	return errors.New("invalid type for status")
}

func (s *Status) Value() (driver.Value, error) {
	return s.name, nil
}

func main() {
	// Using Weekday enum
	day := Tuesday
	fmt.Printf("%s is weekend? %t\\n", day, day.Weekend())

	// Using Color enum
	color := Green
	fmt.Printf("Color: %s, valid: %t\\n", color, color.IsValid())

	// Using Status enum
	status := Active
	fmt.Printf("Status: %s\\n", status)

	// JSON marshaling
	jsonData, _ := json.Marshal(struct {
		Day    Weekday
		Color  Color
		Status *Status
	}{day, color, status})
	fmt.Println("JSON:", string(jsonData))
}
```
