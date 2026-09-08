package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRebasePrefixesRootRelativeLinksOnly(t *testing.T) {
	in := `<a href="/steps">s</a><script src="/static/ui.js"></script>` +
		`<a href="https://go.dev/play">go</a><a href="//cdn.example/x.js">cdn</a>` +
		`<a href="#anchor">a</a><form action="/search">`

	got := string(rebase([]byte(in), "/golangcodes"))

	for _, want := range []string{
		`href="/golangcodes/steps"`,
		`src="/golangcodes/static/ui.js"`,
		`href="https://go.dev/play"`, // absolute, untouched
		`href="//cdn.example/x.js"`,  // protocol-relative, untouched
		`href="#anchor"`,             // fragment, untouched
		`action="/golangcodes/search"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %s", want, got)
		}
	}
}

func TestRebaseIsANoOpWithoutABase(t *testing.T) {
	in := []byte(`<a href="/steps">s</a>`)
	if got := string(rebase(in, "")); got != string(in) {
		t.Errorf("rebase changed the html with no base: %s", got)
	}
}

// The output directory gets wiped, so pointing -out at the wrong place must
// fail rather than delete someone's files.
func TestPrepareRefusesADirectoryItDidNotCreate(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "important.txt"), []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := prepare(dir, false); err == nil {
		t.Fatal("prepare should refuse a directory with no build marker")
	}
	if _, err := os.Stat(filepath.Join(dir, "important.txt")); err != nil {
		t.Error("prepare deleted a file it should have left alone")
	}

	// With the marker, it is one of ours and may be replaced.
	if err := os.WriteFile(filepath.Join(dir, marker), nil, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := prepare(dir, false); err != nil {
		t.Fatalf("prepare refused its own output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "important.txt")); !os.IsNotExist(err) {
		t.Error("prepare should have emptied its own output directory")
	}
}

func TestPrepareAcceptsAnEmptyOrMissingDirectory(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nested", "out")
	if err := prepare(missing, false); err != nil {
		t.Fatalf("missing directory: %v", err)
	}
	if err := prepare(missing, false); err != nil {
		t.Fatalf("empty directory: %v", err)
	}
}

// The whole build, end to end, against the real templates.
func TestBuildWritesTheWholeSite(t *testing.T) {
	out := filepath.Join(t.TempDir(), "site")

	if err := run(out, "../../templates", "../../static", "", "", false); err != nil {
		t.Fatalf("run: %v", err)
	}

	for _, want := range []string{
		"index.html",
		"steps/index.html",
		"404.html",
		"sessions/index.html", // redirect stub
		"resources/index.html",
		"static/ui.js",
		"static/ui.css",
		marker,
	} {
		if _, err := os.Stat(filepath.Join(out, want)); err != nil {
			t.Errorf("missing %s", want)
		}
	}

	if _, err := os.Stat(filepath.Join(out, "CNAME")); !os.IsNotExist(err) {
		t.Error("CNAME written without -domain")
	}

	home, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(home), "1000 Go exercises") {
		t.Error("the landing page did not render its content")
	}

	stub, err := os.ReadFile(filepath.Join(out, "sessions", "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`http-equiv="refresh"`, `href="/steps"`, "noindex"} {
		if !strings.Contains(string(stub), want) {
			t.Errorf("redirect stub is missing %q", want)
		}
	}
}

func TestBuildWithBaseAndDomain(t *testing.T) {
	out := filepath.Join(t.TempDir(), "site")

	if err := run(out, "../../templates", "../../static", "golangcodes", "golang.codes", false); err != nil {
		t.Fatalf("run: %v", err)
	}

	cname, err := os.ReadFile(filepath.Join(out, "CNAME"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(cname)) != "golang.codes" {
		t.Errorf("CNAME = %q", cname)
	}

	home, err := os.ReadFile(filepath.Join(out, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(home), `href="/golangcodes/steps"`) {
		t.Error("links were not rebased onto the project path")
	}
	if !strings.Contains(string(home), `href="/golangcodes/static/ui.css"`) &&
		!strings.Contains(string(home), `src="/golangcodes/static/ui.js"`) {
		t.Error("assets were not rebased onto the project path")
	}
}

// The component gallery is internal tooling. It must never reach the published
// site: not as a page, not as a nav link, not as a stray reference.
func TestGalleryIsNeverPublished(t *testing.T) {
	out := filepath.Join(t.TempDir(), "site")

	if err := run(out, "../../templates", "../../static", "", "", false); err != nil {
		t.Fatalf("run: %v", err)
	}

	if _, err := os.Stat(filepath.Join(out, "components")); !os.IsNotExist(err) {
		t.Error("the component gallery was written to the published site")
	}

	err := filepath.WalkDir(out, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return err
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.Contains(string(body), `href="/components"`) {
			t.Errorf("%s links to the gallery", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
