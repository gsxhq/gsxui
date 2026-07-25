package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

//go:embed shim.js
var shimSource []byte

// registerModuleRoutes wires the two module trees:
//
//	/ui/*.js    the real ui/ modules, byte-for-byte, as native ES modules
//	/shim/*.js  the same, except gsxui.js is shim.js
//
// Serving the real source with no bundler is deliberate: nothing sits
// between a test and the code under test, so there is no transform to keep
// in sync and no bundle cache to invalidate.
func registerModuleRoutes(mux *http.ServeMux, root string) {
	uiDir := filepath.Join(root, "ui")

	serve := func(w http.ResponseWriter, r *http.Request, shim bool) {
		name := r.PathValue("file")
		// Only bare .js filenames. filepath.Base collapses any traversal
		// attempt, and the equality check rejects anything that changed.
		if filepath.Base(name) != name || !strings.HasSuffix(name, ".js") {
			http.NotFound(w, r)
			return
		}
		if shim && name == "gsxui.js" {
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
			w.Write(shimSource)
			return
		}
		b, err := os.ReadFile(filepath.Join(uiDir, name))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		w.Write(b)
	}

	mux.HandleFunc("GET /ui/{file}", func(w http.ResponseWriter, r *http.Request) {
		serve(w, r, false)
	})
	mux.HandleFunc("GET /shim/{file}", func(w http.ResponseWriter, r *http.Request) {
		serve(w, r, true)
	})

	// A blank page whose only job is to import every behavior module through
	// the shim, so window.__gsxuiRegistrations holds the full registry.
	mux.HandleFunc("GET /registrations", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := renderShell(w, "registrations", "/shim/index.js", template.HTML("")); err != nil {
			log.Printf("rendering shell for registrations: %v", err)
		}
	})
}
