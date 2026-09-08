package handlers

import (
	"reflect"
	"strings"
	"testing"
)

func TestDictRejectsBadArguments(t *testing.T) {
	if _, err := dict("a", 1, "b"); err == nil {
		t.Error("odd argument count should be an error")
	}
	if _, err := dict(1, "a"); err == nil {
		t.Error("non-string key should be an error")
	}
	m, err := dict("Label", "Run", "Size", 2)
	if err != nil {
		t.Fatalf("dict: %v", err)
	}
	if m["Label"] != "Run" || m["Size"] != 2 {
		t.Errorf("dict = %v", m)
	}
}

func TestDefFallsBackOnMissingAndEmpty(t *testing.T) {
	cases := []struct {
		name  string
		value any
		want  any
	}{
		{"missing key", nil, "md"},
		{"empty string", "", "md"},
		{"zero int", 0, "md"},
		{"set", "lg", "lg"},
		{"false stays false", false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := def("md", tc.value); got != tc.want {
				t.Errorf("def(md, %v) = %v, want %v", tc.value, got, tc.want)
			}
		})
	}
}

// An unknown variant must degrade to the default. Rendering "%!s(<nil>)" into
// a class attribute is the failure this guards against.
func TestOptFallsBackOnUnknownVariant(t *testing.T) {
	table := map[string]any{"primary": "blue", "danger": "red"}

	for _, tc := range []struct{ key, want string }{
		{"danger", "red"},
		{"", "blue"},
		{"nonsense", "blue"},
	} {
		if got := opt(table, tc.key, "primary"); got != tc.want {
			t.Errorf("opt(%q) = %q, want %q", tc.key, got, tc.want)
		}
	}
}

func TestPctClampsAndHandlesZeroMax(t *testing.T) {
	cases := []struct {
		value, max any
		want       float64
	}{
		{0, 1000, 0},
		{500, 1000, 50},
		{1000, 1000, 100},
		{2000, 1000, 100},
		{-5, 1000, 0},
		{5, 0, 0},
	}
	for _, tc := range cases {
		if got := pct(tc.value, tc.max); got != tc.want {
			t.Errorf("pct(%v, %v) = %v, want %v", tc.value, tc.max, got, tc.want)
		}
	}
}

// Pagination has to survive 1000 steps without rendering 1000 links.
func TestPaginate(t *testing.T) {
	cases := []struct {
		name           string
		current, total int
		want           []int
	}{
		{"single page", 1, 1, []int{1}},
		{"no pages", 1, 0, nil},
		{"short run has no gaps", 3, 5, []int{1, 2, 3, 4, 5}},
		{"gap on the right", 2, 100, []int{1, 2, 3, 4, 0, 100}},
		{"gaps on both sides", 50, 100, []int{1, 0, 48, 49, 50, 51, 52, 0, 100}},
		{"gap on the left", 99, 100, []int{1, 0, 97, 98, 99, 100}},
		{"current below range", 0, 10, []int{1, 2, 3, 0, 10}},
		{"current above range", 99, 10, []int{1, 0, 8, 9, 10}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := paginate(tc.current, tc.total)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("paginate(%d, %d) = %v, want %v", tc.current, tc.total, got, tc.want)
			}
			if len(got) > 12 {
				t.Errorf("paginate rendered %d links; the window should keep it small", len(got))
			}
		})
	}
}

// Every component the gallery advertises must actually render. This is the
// check that catches a renamed or deleted component before a page 500s.
func TestGalleryRendersEveryComponent(t *testing.T) {
	body := renderGallery(t)

	for name := range componentCalls {
		if !strings.Contains(body, name) {
			t.Errorf("gallery does not mention %s", name)
		}
	}
	if strings.Contains(body, "%!s(") || strings.Contains(body, "<no value>") {
		t.Error("a component rendered a Go formatting error into the page")
	}
}

// The gallery exercises every component with real arguments, so it doubles as
// a smoke test that each one produces the markup its behaviour depends on.
func TestGalleryMarkupIsWired(t *testing.T) {
	body := renderGallery(t)

	for _, hook := range []string{
		"data-ui-editor", "data-ui-editor-input", "data-ui-editor-gutter",
		"data-ui-copy", "data-ui-run", "data-ui-share", "data-ui-output",
		"data-ui-tabs", "data-ui-modal", "data-ui-dismiss",
		`role="tablist"`, `role="progressbar"`, `aria-modal="true"`,
		"ui-toasts",
	} {
		if !strings.Contains(body, hook) {
			t.Errorf("gallery is missing %q, so static/ui.js cannot wire it", hook)
		}
	}
}

// Private pages are served locally but never published.
func TestPrivatePagesAreExcludedFromThePublishedSite(t *testing.T) {
	var private []string
	for _, p := range Site() {
		if p.Private {
			private = append(private, p.Path)
		}
	}
	if len(private) == 0 {
		t.Fatal("no page is marked private; the gallery should be")
	}

	for _, p := range PublicSite() {
		if p.Private {
			t.Errorf("PublicSite returned the private page %s", p.Path)
		}
	}
	if len(PublicSite())+len(private) != len(Site()) {
		t.Errorf("PublicSite has %d pages, Site has %d, %d private", len(PublicSite()), len(Site()), len(private))
	}

	// The gallery specifically, since that is the one that must not ship.
	var found bool
	for _, p := range Site() {
		if p.Path == "/components" && p.Private {
			found = true
		}
	}
	if !found {
		t.Error("/components is not marked private")
	}
}

// The nav link follows the page: present on the dev server, gone once built.
func TestGalleryNavLinkIsDevOnly(t *testing.T) {
	restore := Dev
	defer func() { Dev = restore }()

	Dev = true
	if !strings.Contains(renderHome(t), `href="/components"`) {
		t.Error("the dev server should link to the gallery")
	}

	Dev = false
	if strings.Contains(renderHome(t), `href="/components"`) {
		t.Error("a built page should not link to the gallery")
	}
}

func renderHome(t *testing.T) string {
	t.Helper()
	var buf strings.Builder
	if err := RenderTo(&buf, homePage()); err != nil {
		t.Fatalf("render: %v", err)
	}
	return buf.String()
}
