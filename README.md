# Golang Codes

1000 Go exercises, written one at a time.

![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

**0 of 1000 published.**

## Run it

Needs [Go 1.24+](https://go.dev/dl/).

```bash
make tools   # installs air, the hot-reload runner
make dev     # http://localhost:8080
```

Or in Docker, with no Go toolchain: `make up-dev`.

`make help` lists the rest. `make test`, `make fmt`, `make vet` before pushing.

## Deploy

The site ships as static files on GitHub Pages. `cmd/build` renders every page
to disk from `handlers.Site()` — the same list the dev server serves, so the
two cannot drift.

```bash
make site      # render into dist/
make preview   # render, then serve dist/ exactly as Pages will
```

Pushing to `main` runs `.github/workflows/pages.yml`, which tests, builds and
deploys. It picks the right shape on its own:

| `SITE_DOMAIN` repo variable | Result |
|---|---|
| set, e.g. `golang.codes` | Built for the domain root, `CNAME` written |
| unset | Built for `<user>.github.io/<repo>`, every root-relative link prefixed |

Set it under **Settings → Secrets and variables → Actions → Variables**, and
point the DNS at GitHub before you do — a `CNAME` for a domain that does not
resolve takes the site offline at both addresses.

There is no server in production. Run and Share call the Go playground straight
from the browser (it sends `access-control-allow-origin: *`), so nothing has to
be hosted, nothing sleeps, and a playground link outlives this site.

`Dockerfile`, `compose.yaml` and `/healthz` are now only for running the dev
server in a container. Nothing deploys with them.

## Components

Every UI component is a Go template in `templates/components/`, parsed into
every page. There is nothing to import — call one and go:

```gotemplate
{{template "ui/button" dict "Label" "Run" "Variant" "primary"}}
{{template "ui/editor" dict "Name" "step-0001" "Code" .Starter}}
{{template "ui/test-result" dict "Passed" false "Failed" 1 "Total" 3}}
```

**Run the site and open [/components](http://localhost:8080/components)** — it
renders every component live with the exact call to copy.

The gallery is **local only**. Its page is marked `Private` in
`handlers/site.go`, so `cmd/build` leaves it out and `layout.html` hides its nav
link unless `handlers.Dev` is set, which only `main.go` does. Tests assert it
never reaches a build. Mark any other internal page the same way.

| File | Components |
|------|-----------|
| `primitives.html` | button, badge, card, progress, stat, empty, spinner, skeleton, section, divider |
| `forms.html` | input, textarea, select, checkbox, radio-group, toggle, search |
| `code.html` | code, editor, output, share, diff |
| `feedback.html` | alert, toast-host, modal, note |
| `navigation.html` | tabs, breadcrumb, pagination, prevnext, back |
| `exercise.html` | step-card, prompt, hint, test-result, solution, checklist |

Behaviour is wired from `data-` attributes by `static/ui.js`, so a component
dropped into any template works with no extra script. Its public surface:

```js
UI.toast("Copied", {variant: "success"})   // never window.alert
UI.copy(text)                              // clipboard, with an http:// fallback
UI.modal.open("id") / UI.modal.close()
UI.run("editor-id")                        // compiles on the Go playground
UI.share("editor-id")                      // copies a go.dev/play/p link
```

## Layout

```
main.go                 the dev server
cmd/build/              renders the site to static files
handlers/site.go        every page, once, for both of the above
handlers/templates.go   the template funcs components are built on
templates/components/   the component library
templates/              layout.html is the shell; the pages fill it in
static/ui.js            component behaviour
static/ui.css           editor gutter, toast motion
```

Exercises will live in `exercises/` as real Go packages. That directory does not
exist yet, so `publishedSteps` in `handlers/pages.go` returns nothing and the
site reports zero.

## License

MIT — see [LICENSE](LICENSE).
