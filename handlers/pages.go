package handlers

// stepGoal is the size of the collection this project is working toward.
const stepGoal = 1000

// Step is one numbered exercise, as the index renders it.
type Step struct {
	Number int
	Slug   string
	Title  string
}

// Path is where the step is served.
func (s Step) Path() string { return "/steps/" + s.Slug }

// publishedSteps returns the steps that exist today, in order. Empty until the
// first exercise is written.
func publishedSteps() []Step { return nil }
