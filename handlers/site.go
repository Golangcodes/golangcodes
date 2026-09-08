package handlers

import (
	"io"
	"net/http"
)

// Page is one page of the site: where it lives, which template renders it,
// and the data that template needs.
//
// This list is the single source of truth. The dev server serves it over HTTP
// and cmd/build writes the same pages to disk, so the two cannot drift.
type Page struct {
	Path     string // the URL, e.g. "/steps"
	File     string // where it lands on disk, e.g. "steps/index.html"
	Template string
	Data     any

	// Private keeps a page off the published site. It is still served by the
	// dev server, so internal tooling stays one URL away while you work
	// without being deployed for everyone to find.
	Private bool
}

// Site returns every page, in order, private ones included.
func Site() []Page {
	return []Page{homePage(), stepsPage(), componentsPage(), notFoundPage()}
}

// PublicSite returns only the pages that get published. cmd/build uses this;
// the dev server uses Site.
func PublicSite() []Page {
	var out []Page
	for _, p := range Site() {
		if !p.Private {
			out = append(out, p)
		}
	}
	return out
}

// Dev is true when the dev server is running, rather than a static build. It
// gates anything that should exist locally but never ship — the component
// gallery and its nav link.
var Dev bool

// Redirect is an old URL kept alive by pointing at a new one.
type Redirect struct {
	From string
	To   string
}

// Redirects are the links this site has already published and should not
// break. The dev server answers them with a 301; the static build writes a
// meta-refresh page, because a static host cannot send one.
func Redirects() []Redirect {
	return []Redirect{
		{"/sessions", "/steps"},
		{"/resources", "/steps"},
		{"/go", "/steps"},
	}
}

func homePage() Page {
	steps := publishedSteps()
	return Page{
		Path:     "/",
		File:     "index.html",
		Template: "landing.html",
		Data: struct {
			page
			Goal      int
			Published int
		}{
			page:      newPage("Golang Codes", "1000 Go exercises.", "/", "home"),
			Goal:      stepGoal,
			Published: len(steps),
		},
	}
}

func stepsPage() Page {
	steps := publishedSteps()
	return Page{
		Path:     "/steps",
		File:     "steps/index.html",
		Template: "steps.html",
		Data: struct {
			page
			Goal      int
			Published int
			Steps     []Step
		}{
			page:      newPage("Steps — Golang Codes", "Every exercise, in order.", "/steps", "steps"),
			Goal:      stepGoal,
			Published: len(steps),
			Steps:     steps,
		},
	}
}

func componentsPage() Page {
	return Page{
		Path:     "/components",
		File:     "components/index.html",
		Template: "components.html",
		Private:  true,
		Data: struct {
			page
			Calls    map[string]string
			Sample   string
			Starter  string
			Goal     int
			Progress int
		}{
			page:     newPage("Components — Golang Codes", "Every UI component in this project, rendered live, with the call that produced it.", "/components", "components"),
			Calls:    componentCalls,
			Sample:   sampleGo,
			Starter:  sampleStarter,
			Goal:     stepGoal,
			Progress: 120,
		},
	}
}

func notFoundPage() Page {
	return Page{
		Path:     "/404",
		File:     "404.html",
		Template: "404.html",
		Data: struct{ page }{
			page: newPage("Not found — Golang Codes", "That page does not exist.", "/404", ""),
		},
	}
}

// RenderTo writes a page to w. The static build uses this directly; the HTTP
// handlers go through render, which adds the buffering and the status code.
func RenderTo(w io.Writer, p Page) error {
	t, ok := pages[p.Template]
	if !ok {
		return errUnknownTemplate(p.Template)
	}
	return t.ExecuteTemplate(w, "layout", p.Data)
}

type errUnknownTemplate string

func (e errUnknownTemplate) Error() string { return "unknown template " + string(e) }

// handlerFor serves one page over HTTP.
func handlerFor(build func() Page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := build()
		render(w, p.Template, p.Data)
	}
}

// The routes main.go wires up. Each rebuilds its page per request so an edit
// during `make dev` shows up without a restart.
var (
	HomeHandler       = handlerFor(homePage)
	StepsHandler      = handlerFor(stepsPage)
	ComponentsHandler = handlerFor(componentsPage)
)

// NotFoundHandler serves the 404 page with the right status.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	p := notFoundPage()
	w.WriteHeader(http.StatusNotFound)
	render(w, p.Template, p.Data)
}

// LegacyRedirectHandler answers every published-and-moved URL.
func LegacyRedirectHandler(w http.ResponseWriter, r *http.Request) {
	for _, rd := range Redirects() {
		if r.URL.Path == rd.From || len(r.URL.Path) > len(rd.From) && r.URL.Path[:len(rd.From)+1] == rd.From+"/" {
			http.Redirect(w, r, rd.To, http.StatusMovedPermanently)
			return
		}
	}
	NotFoundHandler(w, r)
}
