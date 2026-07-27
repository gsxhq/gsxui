// Command harness serves gsxui's component examples one per page, for the
// Playwright suite in jstest/. It is a test tool: it binds loopback only and
// is never built into the shipped site.
//
//	go run ./jstest/harness -addr 127.0.0.1:7799 -root .
//	go run ./jstest/harness -manifest jstest/.tmp/examples.json
package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gsxhq/gsxui/site/examples"
	"github.com/gsxhq/gsxui/site/pages"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:7799", "loopback address to listen on")
	root := flag.String("root", ".", "repo root, served read-only under /static/")
	manifest := flag.String("manifest", "", "write the example manifest to this path and exit")
	flag.Parse()

	if *manifest != "" {
		if err := writeManifest(*manifest); err != nil {
			log.Fatalf("writing manifest: %v", err)
		}
		return
	}

	if !strings.HasPrefix(*addr, "127.0.0.1:") && !strings.HasPrefix(*addr, "localhost:") {
		log.Fatalf("-addr %q is not loopback; the harness serves the repo tree and must not be reachable off-host", *addr)
	}

	fmt.Fprintf(os.Stderr, "harness listening on http://%s\n", *addr)
	if err := http.ListenAndServe(*addr, newMux(*root)); err != nil {
		log.Fatal(err)
	}
}

func newMux(root string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// /f/ contains fixtures that exercise real GSX components without adding
	// synthetic entries to the production example registry.
	mux.HandleFunc("GET /f/style-contract", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := StyleContractFixture().Render(r.Context(), &buf); err != nil {
			http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := renderShell(w, "style-contract", stylesheetFor(r), "/ui/index.js", template.HTML(buf.String())); err != nil {
			log.Printf("rendering style-contract fixture: %v", err)
		}
	})

	mux.HandleFunc("GET /f/sidebar-contract", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := SidebarContractFixture(r.URL.Query().Get("case")).Render(r.Context(), &buf); err != nil {
			http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := renderShell(w, "sidebar-contract", stylesheetFor(r), "/ui/index.js", template.HTML(buf.String())); err != nil {
			log.Printf("rendering sidebar-contract fixture: %v", err)
		}
	})

	mux.HandleFunc("GET /theme", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		if err := pages.ThemeEditor().Render(r.Context(), &buf); err != nil {
			http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := renderShell(w, "theme", stylesheetFor(r), "/web/theme.js", template.HTML(buf.String())); err != nil {
			log.Printf("rendering theme editor: %v", err)
		}
	})

	// /static/ serves the repo tree read-only. The compiled stylesheet lands
	// at jstest/.tmp/site.css, and Tailwind's bundled @fontsource imports
	// carry url() references relative to that output file — serving from the
	// repo root is what lets those resolve instead of 404ing into the
	// clean-load invariant.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir(root))))

	mux.HandleFunc("GET /x/{component}/{example}", func(w http.ResponseWriter, r *http.Request) {
		component := r.PathValue("component")
		name := r.PathValue("example")

		for _, ex := range examples.For(component) {
			if ex.Name != name {
				continue
			}
			// ex.Query, when the example sets one, re-renders from the
			// request's own query parameters (site/examples/registry.go's
			// Example.Query doc comment) — generic across every component;
			// this handler forwards the raw query values without ever
			// inspecting component or name itself. An example with no Query
			// hook (the common case) always falls through to its static Node.
			node := ex.Node
			if ex.Query != nil {
				node = ex.Query(r.URL.Query())
			}
			var buf bytes.Buffer
			if err := node.Render(r.Context(), &buf); err != nil {
				http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			title := component + "/" + name
			if err := renderShell(w, title, stylesheetFor(r), "/ui/index.js", template.HTML(buf.String())); err != nil {
				log.Printf("rendering shell for %s: %v", title, err)
			}
			return
		}
		http.NotFound(w, r)
	})

	registerModuleRoutes(mux, root)

	return mux
}

func stylesheetFor(r *http.Request) string {
	if r.URL.Query().Get("css") == "foundation" {
		return "/static/jstest/.tmp/foundation.css"
	}
	return "/static/jstest/.tmp/site.css"
}
