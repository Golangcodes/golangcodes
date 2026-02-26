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
)

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
		"Erros[1]",
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
				Path:  "/go/" + expectedTitle,
			})
			// Remove so we know what's left
			delete(entryMap, expectedTitle)
		}
	}

	// Append any remaining items (just in case they weren't in the list)
	for title := range entryMap {
		resources = append(resources, Resource{
			Title: title,
			Path:  "/go/" + title,
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
	filename := strings.TrimPrefix(r.URL.Path, "/go/")

	// Normalize filename by removing any incoming extension early for our router checks
	filename = strings.TrimSuffix(filename, ".html")
	filename = strings.TrimSuffix(filename, ".md")

	if filename == "" || filename == "/" || strings.ToLower(filename) == "index" {
		http.Redirect(w, r, "/go/Introduction", http.StatusMovedPermanently)
		return
	}
	if strings.Contains(filename, "..") {
		http.Error(w, "Invalid resource", http.StatusBadRequest)
		return
	}

	path := filepath.Join("Go", filename+".md")

	contentBytes, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	}

	// Convert Markdown to HTML
	var htmlBuf bytes.Buffer
	if err := goldmark.Convert(contentBytes, &htmlBuf); err != nil {
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

	resources := getResources()
	var activeResource Resource
	var prevResource *Resource
	var nextResource *Resource

	normalizedPath := "/go/" + filename
	for i, res := range resources {
		if res.Path == normalizedPath {
			activeResource = res
			activeResource.Content = template.HTML(content) // Dangerous in prod without sanitization, but okay for local files

			if i > 0 {
				prevResource = &resources[i-1]
			}
			if i < len(resources)-1 {
				nextResource = &resources[i+1]
			}
			break
		}
	}

	data := PageData{
		Resources:      resources,
		ActiveResource: &activeResource,
		PrevResource:   prevResource,
		NextResource:   nextResource,
	}

	tmpl.Execute(w, data)
}
