// Package pages is the gsxui site's structpages route tree: the landing
// page today, and (later tasks) the component showcase, docs, and theme
// editor as sibling route groups.
package pages

import "net/http"

// Pages is the route tree root, mounted at "/". Later tasks add Docs here.
type Pages struct {
	Home      `route:"/{$} Home"`
	Component `route:"/components/{name} Component"`
	Theme     `route:"/theme Theme"`
}

// ServeHTTP is the fallback for any path under "/" that no child route
// matches.
func (Pages) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "page not found", http.StatusNotFound)
}

// ErrorWithStatus is a typed error a Props method can return to make the
// global structpages error handler (see site/main.go's WithErrorHandler)
// respond with a specific status instead of the default 500 — Props runs
// against a buffered writer, so it must return the error rather than
// calling http.Error itself (see the structpages skill's "Error handling
// in handlers" section).
type ErrorWithStatus struct {
	Status  int
	Message string
}

func (e ErrorWithStatus) Error() string { return e.Message }
