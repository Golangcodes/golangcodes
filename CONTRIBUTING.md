# Contributing to Golang Codes

Thank you for your interest in contributing! This project is community-driven, and we appreciate every contribution — from fixing a typo to adding entire lessons.

---

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Getting Started](#getting-started)
- [Lesson Guidelines](#lesson-guidelines)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Reporting Bugs](#reporting-bugs)

---

## Code of Conduct

This project follows our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree to uphold a respectful and welcoming environment for everyone.

---

## How Can I Contribute?

### 📖 Improve Lessons
- Fix typos, grammar, or unclear explanations
- Add runnable Go code examples
- Add real-world use cases
- Improve the structure or flow of existing lessons

### 📝 Add New Lessons
- Propose new topics via a GitHub Issue first
- Write lessons in Markdown format (see [Lesson Guidelines](#lesson-guidelines))
- Place lesson files in the `Go/` directory

### 🐛 Fix Bugs
- Check the Issues tab for open bugs
- Fix UI issues, broken links, or JavaScript errors

### ✨ Add Features
- Improve the code editor experience
- Enhance the gamification system
- Add accessibility improvements
- Improve mobile responsiveness

---

## Getting Started

### 1. Fork and Clone

```bash
git clone https://github.com/<your-username>/golangcodes.git
cd golang.codes
```

### 2. Install Dependencies

```bash
# Install Air for hot-reloading
go install github.com/air-verse/air@v1.52.3
```

### 3. Run Locally

```bash
make dev
```

Visit `http://localhost:8080` to see the app running.

### 4. Create a Branch

```bash
git checkout -b feature/your-feature-name
```

---

## Lesson Guidelines

All lessons live in the `Go/` directory as Markdown (`.md`) files.

### File Naming
- Use descriptive names: `Variables in Go.md`, `Context Package.md`
- The filename becomes the lesson title in the sidebar

### Structure
Each lesson should follow this general structure:

```markdown
# Lesson Title

Brief introduction explaining the topic and why it matters.

---

## Section 1

Explanation with examples.

\```Go
// Runnable code example
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
\```

---

## Section 2

More content...

---

## Best Practices

Numbered list of practical tips.
```

### Guidelines
- **Include runnable code examples** — Wrap standalone examples in `package main` with a `main()` function
- **Keep explanations concise** — Systems engineers prefer clarity over verbosity
- **Use real-world context** — Show how concepts apply in production systems
- **Avoid AI-generated filler** — No conversational text, summaries tables, or "Would you like me to..." endings
- **Use standard Markdown** — Headers, code blocks, bold, lists. No custom HTML

### Adding to the Sidebar Order

If you add a new lesson, also add it to the `linearOrder` slice in `handlers/pages.go`:

```go
linearOrder := []string{
    "Introduction",
    "Variables in Go",
    // ... existing lessons ...
    "Your New Lesson",  // Add here
    "Community Resources",
}
```

---

## Code Style

### Go
- Follow standard `gofmt` formatting
- Keep handlers clean and focused
- Use meaningful variable names

### JavaScript
- Use the `StorageAPI` abstraction in `static/storage.js` for all persistence
- Never use `localStorage` directly in `main.js` or inline scripts

### HTML/CSS
- Use Tailwind CSS utility classes
- Support both light and dark modes
- Ensure all interactive elements have descriptive `title` attributes

---

## Submitting Changes

### Pull Request Process

1. **Ensure your code works** — Test locally with `make dev`
2. **Write a clear PR title** — e.g., "Add Goroutines lesson" or "Fix broken sidebar link"
3. **Describe your changes** — Explain what you changed and why
4. **Keep PRs focused** — One feature or fix per PR
5. **Update the sidebar order** if adding a new lesson

### Commit Messages

Use clear, descriptive commit messages:

```
Add: Goroutines and Concurrency lesson
Fix: Broken click-to-edit on Variables page
Update: Improve error handling examples in Errors lesson
```

---

## Reporting Bugs

Open a GitHub Issue with:

1. **Description** — What's broken?
2. **Steps to Reproduce** — How can we see the bug?
3. **Expected Behavior** — What should happen?
4. **Screenshots** — If it's a UI issue, include a screenshot
5. **Environment** — Browser, OS, screen size

---

## 🙏 Thank You

Every contribution makes this a better learning resource for the Go community. We review all PRs promptly and appreciate your time and effort!
