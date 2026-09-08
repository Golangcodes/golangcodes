package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// pageNames are the top-level pages. Each is parsed together with layout.html
// and every component in templates/components, so a page can call any
// component without importing anything.
var pageNames = []string{
	"landing.html",
	"steps.html",
	"components.html",
	"404.html",
}

// funcs are the helpers components are built on.
//
// Go templates take a single argument, so components are called with dict:
//
//	{{template "ui/button" dict "Label" "Run" "Variant" "primary"}}
//
// def supplies a fallback for a key the caller left out, and opt does the same
// for a lookup into a variant table, so an unknown variant degrades to the
// default rather than rendering a nil.
var funcs = template.FuncMap{
	"dict":  dict,
	"slice": func(items ...any) []any { return items },
	"def":   def,
	"opt":   opt,
	"safe":  func(s string) template.HTML { return template.HTML(s) },
	"attrs": func(s string) template.HTMLAttr {
		return template.HTMLAttr(s)
	},
	"pct":      pct,
	"paginate": paginate,
	"add":      func(a, b int) int { return a + b },
	"sub":      func(a, b int) int { return a - b },
	"inc":      func(i int) int { return i + 1 },
}

// dict builds the argument for a component from alternating keys and values.
func dict(kv ...any) (map[string]any, error) {
	if len(kv)%2 != 0 {
		return nil, fmt.Errorf("dict: got %d arguments, want an even number", len(kv))
	}
	m := make(map[string]any, len(kv)/2)
	for i := 0; i < len(kv); i += 2 {
		k, ok := kv[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict: key %d is %T, want string", i/2, kv[i])
		}
		m[k] = kv[i+1]
	}
	return m, nil
}

// def returns value unless it is missing or empty, in which case fallback.
func def(fallback, value any) any {
	switch v := value.(type) {
	case nil:
		return fallback
	case string:
		if v == "" {
			return fallback
		}
	case int:
		if v == 0 {
			return fallback
		}
	}
	return value
}

// opt looks key up in a variant table, falling back to the named default both
// when the key is empty and when it is not a variant the table knows.
func opt(table map[string]any, key any, fallback string) string {
	k := fallback
	if s, ok := key.(string); ok && s != "" {
		k = s
	}
	if v, ok := table[k]; ok {
		return fmt.Sprint(v)
	}
	return fmt.Sprint(table[fallback])
}

// pct is value as a percentage of max, clamped, for progress bars.
func pct(value, max any) float64 {
	v, m := toFloat(value), toFloat(max)
	if m <= 0 {
		return 0
	}
	p := v / m * 100
	switch {
	case p < 0:
		return 0
	case p > 100:
		return 100
	}
	return p
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case int:
		return float64(n)
	case int64:
		return float64(n)
	case float64:
		return n
	}
	return 0
}

// paginate returns the page numbers to render around current, with 0 standing
// in for a gap. It keeps 1000 pages down to a handful of links.
func paginate(current, total int) []int {
	if total < 1 {
		return nil
	}
	if current < 1 {
		current = 1
	}
	if current > total {
		current = total
	}

	const window = 2
	keep := map[int]bool{1: true, total: true}
	for i := current - window; i <= current+window; i++ {
		if i >= 1 && i <= total {
			keep[i] = true
		}
	}

	var out []int
	prev := 0
	for i := 1; i <= total; i++ {
		if !keep[i] {
			continue
		}
		if prev != 0 && i-prev > 1 {
			out = append(out, 0)
		}
		out = append(out, i)
		prev = i
	}
	return out
}

var pages map[string]*template.Template

// LoadTemplates parses every page once at startup so a broken template fails
// the process rather than a request.
func LoadTemplates(dir string) error {
	components := filepath.Join(dir, "components", "*.html")

	parsed := make(map[string]*template.Template, len(pageNames))
	for _, name := range pageNames {
		t, err := template.New(name).Funcs(funcs).ParseGlob(components)
		if err != nil {
			return fmt.Errorf("parsing components for %s: %w", name, err)
		}
		t, err = t.ParseFiles(
			filepath.Join(dir, "layout.html"),
			filepath.Join(dir, name),
		)
		if err != nil {
			return fmt.Errorf("parsing %s: %w", name, err)
		}
		parsed[name] = t
	}
	pages = parsed
	return nil
}

// page carries what layout.html needs on every request.
type page struct {
	Title       string
	Description string
	Canonical   string
	Nav         string
	BodyClass   string
	FixedNav    bool
	Dev         bool
	Year        int
}

func newPage(title, description, canonical, nav string) page {
	return page{
		Title:       title,
		Description: description,
		Canonical:   canonical,
		Nav:         nav,
		BodyClass:   "flex flex-col min-h-screen",
		FixedNav:    true,
		Dev:         Dev,
		Year:        time.Now().Year(),
	}
}

// render writes a page through a buffer, so a mid-execution template error
// cannot leave a half-written response behind.
func render(w http.ResponseWriter, name string, data any) {
	t, ok := pages[name]
	if !ok {
		log.Printf("render: unknown template %q", name)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "layout", data); err != nil {
		log.Printf("render %s: %v", name, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("write %s: %v", name, err)
	}
}
