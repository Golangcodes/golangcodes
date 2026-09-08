package handlers

// componentCalls is the copy-paste line for each component on /components.
//
// These live in Go rather than in the gallery template because the gallery is
// itself a template: a literal {{template ...}} written in the HTML would be
// executed instead of shown.
var componentCalls = map[string]string{
	"ui/button": `{{template "ui/button" dict "Label" "Run" "Variant" "primary" "Size" "md"}}

Variants: primary secondary ghost danger success
Sizes:    sm md lg
Also:     Href Type Class Attrs Disabled Loading Full`,

	"ui/badge": `{{template "ui/badge" dict "Label" "passed" "Variant" "pass"}}

Variants: default go pass fail warn info. Dot adds a leading dot.`,

	"ui/card": `{{template "ui/card" dict
    "Title" "Slices share memory"
    "Sub"   "Appending can mutate what someone else is holding."
    "Body"  (safe "<p class=\"text-sm text-slate-500\">Any HTML.</p>")}}

Href turns the whole card into a link.`,

	"ui/progress": `{{template "ui/progress" dict
    "Value" .Published "Max" .Goal
    "Label" "12 published" "Sub" "1000 planned"}}`,

	"ui/stat": `{{template "ui/stat" dict "Value" "0012" "Label" "published" "Sub" "of 1000"}}`,

	"ui/empty": `{{template "ui/empty" dict
    "Glyph" "0000" "Title" "No steps published yet"
    "Body"  "Step 1 goes up first."}}`,

	"ui/spinner":  `{{template "ui/spinner" dict "Size" "md"}}`,
	"ui/skeleton": `{{template "ui/skeleton" dict "Class" "h-4 w-48"}}`,

	"ui/section": `{{template "ui/section" dict "Title" "The steps" "Sub" "In order."}}`,

	"ui/divider": `{{template "ui/divider" dict "Label" "or"}}`,

	"ui/input": `{{template "ui/input" dict
    "Name" "answer" "Label" "Your answer" "Type" "text"
    "Placeholder" "len(\"héllo\")" "Hint" "Bytes, not runes."}}

Add "Error" to switch the field into its invalid state.`,

	"ui/textarea": `{{template "ui/textarea" dict
    "Name" "notes" "Label" "In your own words" "Rows" 4}}`,

	"ui/select": `{{template "ui/select" dict
    "Name" "topic" "Label" "Topic" "Value" "slices"
    "Options" (slice
        (dict "Value" "slices" "Label" "Slices")
        (dict "Value" "maps"   "Label" "Maps"))}}`,

	"ui/checkbox": `{{template "ui/checkbox" dict
    "Name" "done" "Label" "I ran the tests" "Checked" true}}`,

	"ui/radio-group": `{{template "ui/radio-group" dict
    "Name" "guess" "Legend" "What does this print?" "Value" "b"
    "Options" (slice
        (dict "Value" "a" "Label" "6")
        (dict "Value" "b" "Label" "5" "Hint" "Runes, not bytes"))}}`,

	"ui/toggle": `{{template "ui/toggle" dict
    "Name" "hints" "Label" "Show hints" "Checked" true}}`,

	"ui/search": `{{template "ui/search" dict "Placeholder" "Search 1000 steps"}}`,

	"ui/code": `{{template "ui/code" dict "File" "main.go" "Code" .Source}}

The copy button is wired automatically.`,

	"ui/editor": `{{template "ui/editor" dict
    "Name" "step-0001" "File" "main.go" "Code" .Starter "Height" "18rem"}}

Line numbers, Tab to indent, Cmd/Ctrl+Enter to run.
NoRun and NoShare drop those buttons.`,

	"ui/share": `{{template "ui/share" dict "Code" .Source "Label" "Share this snippet"}}

Posts to /api/share and copies a go.dev/play/p link.`,

	"ui/diff": `{{template "ui/diff" dict "Want" "5" "Got" "6"}}`,

	"ui/alert": `{{template "ui/alert" dict
    "Variant" "warn" "Title" "Careful"
    "Body" "len counts bytes." "Dismissible" true}}

Variants: info success warn error`,

	"ui/note": `{{template "ui/note" dict "Body" "The spec settles this in one line."}}`,

	"UI.toast": `UI.toast("Copied", {variant: "success"})

Variants: info success warn error. timeout: 0 keeps it up.`,

	"ui/modal": `{{template "ui/modal" dict
    "Id" "confirm" "Title" "Reset your answer?"
    "Body" (safe "<p>This clears the buffer.</p>")
    "Footer" (safe "...buttons...")}}

Open with UI.modal.open("confirm"), or data-ui-modal-open="confirm".
Never use window.confirm: it freezes the page.`,

	"ui/tabs": `{{template "ui/tabs" dict "Id" "step" "Tabs" (slice
    (dict "Id" "prompt" "Label" "Prompt" "Body" (safe "..."))
    (dict "Id" "tests"  "Label" "Tests"  "Body" (safe "...")))}}

Arrow keys, Home and End move between tabs.`,

	"ui/breadcrumb": `{{template "ui/breadcrumb" dict "Items" (slice
    (dict "Label" "Steps" "Href" "/steps")
    (dict "Label" "0001"))}}`,

	"ui/pagination": `{{template "ui/pagination" dict
    "Current" 42 "Total" 100 "Base" "/steps?page="}}

Windowed, so 1000 pages stay a handful of links.`,

	"ui/prevnext": `{{template "ui/prevnext" dict
    "Prev" (dict "Href" "/steps/1" "Label" "Hello, world")
    "Next" (dict "Href" "/steps/3" "Label" "Zero values")}}`,

	"ui/back": `{{template "ui/back" dict "Href" "/steps" "Label" "All steps"}}`,

	"ui/step-card": `{{template "ui/step-card" dict
    "Number" 1 "Title" "Hello, world" "Href" "/steps/1"
    "Topic" "basics" "Status" "passed"}}

Status: "" | "attempted" | "passed"`,

	"ui/prompt": `{{template "ui/prompt" dict
    "Number" 1 "Title" "Make the test pass" "Topic" "basics"
    "Body" (safe "<p>Return the greeting.</p>")}}`,

	"ui/hint": `{{template "ui/hint" dict
    "Label" "Hint" "Body" (safe "<p>len counts bytes.</p>")}}

Collapsed until asked for, so it is never spoiled.`,

	"ui/test-result": `{{template "ui/test-result" dict
    "Passed" false "Failed" 1 "Total" 3
    "Output" "--- FAIL: TestGreet"}}`,

	"ui/solution": `{{template "ui/solution" dict "Code" .Solution "Locked" false}}

Locked true hides it until the learner's tests pass.`,

	"ui/checklist": `{{template "ui/checklist" dict "Items" (slice
    (dict "Label" "Read the prompt" "Done" true)
    (dict "Label" "Make the test pass"))}}`,
}

const sampleGo = `package main

import "fmt"

func main() {
	// len counts bytes, not runes.
	fmt.Println(len("héllo"))
}`

const sampleStarter = `package main

import "fmt"

// Greet returns a greeting for name.
func Greet(name string) string {
	return "" // TODO
}

func main() {
	fmt.Println(Greet("world"))
}`
