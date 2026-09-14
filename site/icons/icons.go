// Package icons holds the site favicon set (from the gsxui brand design) and
// serves it at stable root URLs. Layout's <head> links /favicon.svg,
// /favicon-32.png and /apple-touch-icon.png, so every server that renders the
// site pages — the site binary and the jstest harness — mounts this package
// and the links resolve identically in dev, prod and tests.
package icons

import (
	"embed"
	"net/http"
)

//go:embed favicon.svg favicon-32.png apple-touch-icon.png
var files embed.FS

// Register mounts each icon at its root URL, straight from the binary rather
// than through the Vite pipeline: favicons need unhashed, root-level URLs. A
// day of caching is plenty — the files are tiny and change roughly never, but
// unlike /static/assets/ they are not content-hashed, so no immutable header.
func Register(mux *http.ServeMux) {
	for _, name := range []string{"favicon.svg", "favicon-32.png", "apple-touch-icon.png"} {
		mux.HandleFunc("/"+name, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=86400")
			http.ServeFileFS(w, r, files, name)
		})
	}
}
