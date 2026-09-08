package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if err := LoadTemplates("../templates"); err != nil {
		panic(err)
	}
	m.Run()
}

// renderGallery returns the rendered /components page.
func renderGallery(t *testing.T) string {
	t.Helper()
	rec := httptest.NewRecorder()
	ComponentsHandler(rec, httptest.NewRequest(http.MethodGet, "/components", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("components page status = %d, want 200", rec.Code)
	}
	return rec.Body.String()
}

func TestPagesRender(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		handler http.HandlerFunc
	}{
		{"landing", "/", HomeHandler},
		{"steps", "/steps", StepsHandler},
		{"components", "/components", ComponentsHandler},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			tc.handler(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rec.Code)
			}
			if !strings.Contains(rec.Body.String(), "</html>") {
				t.Error("page did not render through the layout")
			}
		})
	}
}

// TestCountsAgreeAcrossPages stops the landing page and the index from ever
// claiming different numbers.
func TestCountsAgreeAcrossPages(t *testing.T) {
	claim := strconv.Itoa(len(publishedSteps())) + " published"

	for path, h := range map[string]http.HandlerFunc{"/": HomeHandler, "/steps": StepsHandler} {
		rec := httptest.NewRecorder()
		h(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if !strings.Contains(rec.Body.String(), claim) {
			t.Errorf("%s does not report %q", path, claim)
		}
	}
}

// TestStepsIndexLinksOnlyToRealSteps guards the empty collection: the index
// must not link to a step that does not exist.
func TestStepsIndexLinksOnlyToRealSteps(t *testing.T) {
	rec := httptest.NewRecorder()
	StepsHandler(rec, httptest.NewRequest(http.MethodGet, "/steps", nil))

	if len(publishedSteps()) == 0 && strings.Contains(rec.Body.String(), `href="/steps/`) {
		t.Error("empty index links to a step")
	}
}

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	HealthHandler(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.HasPrefix(rec.Body.String(), "ok:") {
		t.Errorf("body = %q, want an ok report", rec.Body.String())
	}
}

// TestLegacyLinksRedirect keeps every published-and-moved URL pointing
// somewhere real. These are also written as static stubs by cmd/build.
func TestLegacyLinksRedirect(t *testing.T) {
	for _, path := range []string{"/sessions", "/sessions/variables", "/resources", "/go/anything"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		LegacyRedirectHandler(rec, req)

		if rec.Code != http.StatusMovedPermanently {
			t.Errorf("%s: status = %d, want 301", path, rec.Code)
		}
		if got := rec.Header().Get("Location"); got != "/steps" {
			t.Errorf("%s: Location = %q, want /steps", path, got)
		}
	}
}

// An unknown path is a 404, not a redirect.
func TestUnknownPathIsNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	LegacyRedirectHandler(rec, httptest.NewRequest(http.MethodGet, "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "404") {
		t.Error("the 404 page did not render")
	}
}

// Every page the static build writes must also be reachable on the dev
// server, and vice versa: the two run off handlers.Site().
func TestSiteRendersEveryPage(t *testing.T) {
	seen := map[string]bool{}
	for _, p := range Site() {
		if seen[p.File] {
			t.Errorf("two pages both write to %s", p.File)
		}
		seen[p.File] = true

		var buf strings.Builder
		if err := RenderTo(&buf, p); err != nil {
			t.Errorf("%s: %v", p.Path, err)
			continue
		}
		if !strings.Contains(buf.String(), "</html>") {
			t.Errorf("%s did not render through the layout", p.Path)
		}
	}
	if len(seen) != len(Site()) {
		t.Errorf("Site() has %d pages but %d distinct files", len(Site()), len(seen))
	}
}
