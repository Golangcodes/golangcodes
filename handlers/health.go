package handlers

import (
	"fmt"
	"net/http"
)

// HealthHandler answers container and platform health checks.
//
// Templates are the liveness signal: they are parsed once at startup, so a
// process serving none of them is up without being able to answer anything.
// The step count is reported but deliberately not part of the verdict — zero
// published steps is the honest state of a collection being written by hand.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")

	if len(pages) == 0 {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprint(w, "unhealthy: no templates parsed\n")
		return
	}

	fmt.Fprintf(w, "ok: %d templates, %d/%d steps published\n", len(pages), len(publishedSteps()), stepGoal)
}
