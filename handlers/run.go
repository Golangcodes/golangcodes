package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	maxOutputBytes = 1 * 1024 * 1024
	execTimeout    = 10 * time.Second
	playgroundURL  = "https://play.golang.org/compile"
)

type playgroundResponse struct {
	Errors string `json:"Errors"`
	Events []struct {
		Message string `json:"Message"`
		Kind    string `json:"Kind"` // "stdout" or "stderr"
	} `json:"Events"`
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		sendResult(w, "Error: No code provided", true)
		return
	}

	output, isError := runInPlayground(code)
	sendResult(w, output, isError)
}

func runInPlayground(code string) (string, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()

	form := url.Values{
		"body":    {code},
		"version": {"2"},
		"withVet": {"true"},
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		playgroundURL,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return fmt.Sprintf("Internal Error: %v", err), true
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Sprintf("Execution timed out (%s limit)", execTimeout), true
		}
		return fmt.Sprintf("Failed to reach execution service: %v", err), true
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxOutputBytes))
	if err != nil {
		return fmt.Sprintf("Failed to read response: %v", err), true
	}

	var result playgroundResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Sprintf("Failed to parse response: %v", err), true
	}

	// Compilation or vet errors come back in the Errors field.
	if result.Errors != "" {
		return result.Errors, true
	}

	// Events are ordered stdout/stderr messages from the running program.
	var sb strings.Builder
	hasStderr := false
	for _, event := range result.Events {
		sb.WriteString(event.Message)
		if event.Kind == "stderr" {
			hasStderr = true
		}
	}

	return sb.String(), hasStderr
}

func sendResult(w http.ResponseWriter, output string, isError bool) {
	w.Header().Set("Content-Type", "text/html")

	colorClass := "text-green-400"
	if isError {
		colorClass = "text-red-400"
	}

	if output == "" {
		output = "Program exited successfully with no output."
		colorClass = "text-gray-400"
	}

	html := fmt.Sprintf(`<pre class="whitespace-pre-wrap font-mono text-sm %s">%s</pre>`, colorClass, templateEscape(output))
	w.Write([]byte(html))
}

func templateEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
