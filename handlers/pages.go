package handlers

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

func slugify(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "-")
}

type Resource struct {
	Title   string
	Path    string
	Content template.HTML
}

type PageData struct {
	Resources      []Resource
	ActiveResource *Resource
	PrevResource   *Resource
	NextResource   *Resource
}

func getResources() []Resource {
	var resources []Resource
	dir := "Go"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return resources
	}

	// Define the linear order of topics (from basic to advanced)
	linearOrder := []string{
		"Introduction",
		"Variables in Go",
		"Data Structures",
		"Conditionals in Go",
		"Loops",
		"Functions",
		"Slice",
		"Maps",
		"Pointers",
		"Structs",
		"Packages and Modules",
		"Interfaces in Go",
		"Errors[1]",
		"Go Symbols",
		"Go standard library",
		"Channels in Go",
		"Mutexes in Go",
		"Context Package",
		"Generics in Go",
		"Enums in Go",
		"Understanding Allocations in Go",
		"Design Patterns in Go",
		"Go best lessions",
		"Community Resources",
	}

	// Create a map to hold the found entries for quick lookup
	entryMap := make(map[string]os.DirEntry)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			title := strings.TrimSuffix(entry.Name(), ".md")
			entryMap[title] = entry
		}
	}

	// Add resources in the defined linear order
	for _, expectedTitle := range linearOrder {
		if _, exists := entryMap[expectedTitle]; exists {
			resources = append(resources, Resource{
				Title: expectedTitle,
				Path:  "/go/" + slugify(expectedTitle),
			})
			// Remove so we know what's left
			delete(entryMap, expectedTitle)
		}
	}

	// Append any remaining items (just in case they weren't in the list)
	for title := range entryMap {
		resources = append(resources, Resource{
			Title: title,
			Path:  "/go/" + slugify(title),
		})
	}

	return resources
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/landing.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := PageData{
		Resources: getResources(),
	}

	tmpl.Execute(w, data)
}

func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/go/")

	// Normalize filename by removing any incoming extension early for our router checks
	slug = strings.TrimSuffix(slug, ".html")
	slug = strings.TrimSuffix(slug, ".md")

	// Make lookups case-insensitive
	slug = strings.ToLower(slug)

	if slug == "" || slug == "/" || slug == "index" {
		http.Redirect(w, r, "/go/introduction", http.StatusMovedPermanently)
		return
	}
	if strings.Contains(slug, "..") {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}

	resources := getResources()
	var actualTitle string
	var activeResource Resource
	var prevResource *Resource
	var nextResource *Resource

	normalizedPath := "/go/" + slug
	for i, res := range resources {
		if res.Path == normalizedPath {
			activeResource = res
			actualTitle = res.Title

			if i > 0 {
				prevResource = &resources[i-1]
			}
			if i < len(resources)-1 {
				nextResource = &resources[i+1]
			}
			break
		}
	}

	if actualTitle == "" {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	path := filepath.Join("Go", actualTitle+".md")

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Convert Markdown to HTML
	var htmlBuf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Linkify,
		),
	)
	if err := md.Convert(contentBytes, &htmlBuf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	content := htmlBuf.String()

	// Replace youtube links with iframe embeds (keep existing logic but adjusted for markdown compiled html)
	// Matches href="https://youtu.be/ID" or https://www.youtube.com/watch?v=ID
	youtubeRegex := regexp.MustCompile(`<a[^>]*href="https://(?:www\.)?youtu(?:be\.com/watch\?v=|\.be/)([^"&]+)[^"]*"[^>]*>[^<]*</a>`)
	content = youtubeRegex.ReplaceAllString(content, `<iframe class="w-full aspect-video rounded-lg shadow-md my-6" src="https://www.youtube.com/embed/$1" frameborder="0" allowfullscreen></iframe>`)

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	activeResource.Content = template.HTML(content) // Dangerous in prod without sanitization, but okay for local files

	data := PageData{
		Resources:      resources,
		ActiveResource: &activeResource,
		PrevResource:   prevResource,
		NextResource:   nextResource,
	}

	tmpl.Execute(w, data)
}
