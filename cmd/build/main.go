// Command build renders the site to static files.
//
// Production is GitHub Pages, which serves files and runs nothing, so every
// page is written to disk here. The pages come from handlers.Site(), the same
// list the dev server serves, so the two cannot drift.
//
//	go run ./cmd/build                      # into docs/
//	go run ./cmd/build -base /golangcodes   # for a project page
//	go run ./cmd/build -domain golang.codes # writes a CNAME
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.codes/handlers"
)

// marker is written into every build. Its presence is what makes wiping the
// output directory safe: without it, the directory is something else.
const marker = ".nojekyll"

func main() {
	var (
		out       = flag.String("out", "docs", "directory to write the site into")
		templates = flag.String("templates", "templates", "templates directory")
		static    = flag.String("static", "static", "static assets directory")
		base      = flag.String("base", "", "path prefix, for a project page (e.g. /golangcodes)")
		domain    = flag.String("domain", "", "custom domain to write into CNAME")
		force     = flag.Bool("force", false, "wipe the output directory even if this build did not create it")
	)
	flag.Parse()

	if err := run(*out, *templates, *static, *base, *domain, *force); err != nil {
		log.Fatalf("build: %v", err)
	}
}

func run(out, templates, static, base, domain string, force bool) error {
	if base != "" {
		base = "/" + strings.Trim(base, "/")
	}

	if err := handlers.LoadTemplates(templates); err != nil {
		return err
	}
	if err := prepare(out, force); err != nil {
		return err
	}

	var written int

	for _, p := range handlers.PublicSite() {
		var buf bytes.Buffer
		if err := handlers.RenderTo(&buf, p); err != nil {
			return fmt.Errorf("rendering %s: %w", p.Path, err)
		}
		if err := writeFile(filepath.Join(out, p.File), rebase(buf.Bytes(), base)); err != nil {
			return err
		}
		written++
	}

	// A static host cannot send a 301, so a published URL that moved gets a
	// page that redirects instead. The canonical link keeps search engines
	// pointed at the destination rather than the stub.
	for _, r := range handlers.Redirects() {
		file := filepath.Join(out, strings.Trim(r.From, "/"), "index.html")
		if err := writeFile(file, rebase(redirectHTML(r.To), base)); err != nil {
			return err
		}
		written++
	}

	n, err := copyTree(static, filepath.Join(out, "static"))
	if err != nil {
		return err
	}

	if err := writeFile(filepath.Join(out, marker), nil); err != nil {
		return err
	}
	if domain != "" {
		if err := writeFile(filepath.Join(out, "CNAME"), []byte(domain+"\n")); err != nil {
			return err
		}
	}

	log.Printf("wrote %d pages and %d assets to %s/", written, n, out)
	return nil
}

// prepare empties the output directory, refusing one this build did not make.
func prepare(out string, force bool) error {
	info, err := os.Stat(out)
	switch {
	case os.IsNotExist(err):
		return os.MkdirAll(out, 0o755)
	case err != nil:
		return err
	case !info.IsDir():
		return fmt.Errorf("%s is not a directory", out)
	}

	entries, err := os.ReadDir(out)
	if err != nil {
		return err
	}
	if len(entries) > 0 && !force {
		if _, err := os.Stat(filepath.Join(out, marker)); err != nil {
			return fmt.Errorf("%s is not empty and was not created by this build; pass -force to overwrite it", out)
		}
	}

	if err := os.RemoveAll(out); err != nil {
		return err
	}
	return os.MkdirAll(out, 0o755)
}

// attrPath matches a root-relative URL in an HTML attribute.
var attrPath = regexp.MustCompile(`\b(href|src|action)="(/[^"]*)"`)

// rebase prefixes root-relative links, so the site works when it is served
// from a subdirectory rather than a domain root. Protocol-relative URLs
// (//example.com) are left alone.
func rebase(html []byte, base string) []byte {
	if base == "" {
		return html
	}
	return attrPath.ReplaceAllFunc(html, func(m []byte) []byte {
		g := attrPath.FindSubmatch(m)
		if bytes.HasPrefix(g[2], []byte("//")) {
			return m
		}
		return []byte(string(g[1]) + `="` + base + string(g[2]) + `"`)
	})
}

func redirectHTML(to string) []byte {
	return []byte(`<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8" />
<title>Moved</title>
<link rel="canonical" href="` + to + `" />
<meta name="robots" content="noindex" />
<meta http-equiv="refresh" content="0; url=` + to + `" />
</head>
<body>
<p>This page moved to <a href="` + to + `">` + to + `</a>.</p>
</body>
</html>
`)
}

func writeFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

// copyTree copies src into dst, returning how many files it wrote.
func copyTree(src, dst string) (int, error) {
	var n int
	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		outFile, err := os.Create(target)
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, in); err != nil {
			return err
		}
		n++
		return nil
	})
	return n, err
}
