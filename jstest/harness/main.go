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
			var buf bytes.Buffer
			if err := ex.Node.Render(r.Context(), &buf); err != nil {
				http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			title := component + "/" + name
			if err := renderShell(w, title, "/ui/index.js", template.HTML(buf.String())); err != nil {
				log.Printf("rendering shell for %s: %v", title, err)
			}
			return
		}
		http.NotFound(w, r)
	})

	registerModuleRoutes(mux, root)

	return mux
}
