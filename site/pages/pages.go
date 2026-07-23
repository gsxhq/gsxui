// Package pages is the gsxui site's structpages route tree: the landing
// page today, and (later tasks) the component showcase, docs, and theme
// editor as sibling route groups.
package pages

import "net/http"

// Pages is the route tree root, mounted at "/". Later tasks add Components,
// Docs, and Theme fields here.
type Pages struct {
	Home `route:"/{$} Home"`
}

// ServeHTTP is the fallback for any path under "/" that no child route
// matches (e.g. /components/{name} before Task 2 registers that tree).
func (Pages) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "page not found", http.StatusNotFound)
}
