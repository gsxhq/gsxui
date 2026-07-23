// Command site serves the gsxui documentation/showcase site: the component
// registry rendered as browsable pages, backed by structpages routing and
// gsx templates, with a Vite-built frontend bundle.
package main

import (
	"cmp"
	"context"
	"embed"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gsxhq/gsxui/site/pages"
	"github.com/gsxhq/vite"
	"github.com/jackielii/structpages"
)

// dist holds the Vite production build (npm run build / `make site`). Empty
// in dev — vite.New only reads it in prod mode (DevURL == "").
//
//go:embed all:dist
var distFS embed.FS

// listenPort resolves the bind port. GO_PORT (dev): set via .env / matched
// by vite.config.ts's proxy target — wins so the dev loop stays coherent.
// PORT (containers): the Dockerfile sets ENV PORT=8080 and Cloud Run injects
// it. Final fallback is 7777 — the SAME default vite.config.ts and `gsx dev`
// assume — so a fresh checkout's `make site-dev` works with no .env at all;
// an 8080-style fallback here desynchronizes the proxy and shows the
// "backend unavailable" interstitial forever.
func listenPort(getenv func(string) string) string {
	return cmp.Or(getenv("GO_PORT"), getenv("PORT"), "7777")
}

// immutableAssets marks the Vite bundle under /static/assets/ cacheable
// forever: those filenames carry a content hash, so a changed file is a new
// URL and stale copies are unreachable by construction. The prefix check is
// the fingerprint guarantee — anything else under /static/ (e.g. the Vite
// manifest) is not hashed and must not get a forever header.
func immutableAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/static/assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		}
		next.ServeHTTP(w, r)
	})
}

// pageCache stamps rendered pages cacheable at the CDN edge (Cloudflare
// sits proxied in front of ui.gsxhq.dev): browsers always revalidate
// (max-age=0) so a deploy is visible immediately to direct visitors, while
// the edge may hold a copy for s-maxage to absorb traffic without waking a
// stopped Fly machine. The pages are fully static — no per-user content —
// so shared caching is unconditionally safe. errorHandler must clear this
// header before writing a failure (a cached 404/500 would outlive the
// fault).
func pageCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=0, s-maxage=300")
		next.ServeHTTP(w, r)
	})
}

// errorHandler is the structpages error handler. It clears the Cache-Control
// header pageCache stamped before routing ran — http.Error does not remove
// preset headers, and an edge-cached error page would keep serving after the
// underlying fault (or missing route) is fixed.
func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	w.Header().Del("Cache-Control")
	var se pages.ErrorWithStatus
	if errors.As(err, &se) {
		http.Error(w, se.Message, se.Status)
		return
	}
	log.Printf("error rendering %s: %v", r.URL.Path, err)
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func main() {
	devURL := os.Getenv("VITE_DEV_URL") // "" in prod
	v, err := vite.New(vite.Config{DevURL: devURL, DevBase: "/__vite/", Dist: distFS, DistDir: "dist"})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	if !v.Dev() {
		mux.Handle("/static/", immutableAssets(v.StaticHandler()))
	}

	// Pages mount on their own mux so the cache middleware wraps exactly the
	// rendered routes — never /healthz, and /static/ keeps its own policy.
	// Dev skips the wrapper: the dev loop should never see a caching header.
	pagesMux := http.NewServeMux()
	if _, err := structpages.Mount(pagesMux, pages.Pages{}, "/", "gsxui",
		structpages.WithErrorHandler(errorHandler),
	); err != nil {
		log.Fatal(err)
	}
	if v.Dev() {
		mux.Handle("/", pagesMux)
	} else {
		mux.Handle("/", pageCache(pagesMux))
	}

	// v.Middleware injects *vite.Vite into each request's context so components
	// read the asset bundle from ctx (no prop threading).
	//
	port := listenPort(os.Getenv)
	srv := &http.Server{Addr: ":" + port, Handler: v.Middleware(mux)}

	// Serve in the background so the main goroutine can wait for a shutdown
	// signal. gsx dev sends SIGTERM on each rebuild; shutting down gracefully
	// releases the port BEFORE exit, so the next build re-binds cleanly.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("listening on http://localhost:%s", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
