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
		mux.Handle("/static/", v.StaticHandler())
	}

	if _, err := structpages.Mount(mux, pages.Pages{}, "/", "gsxui",
		structpages.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			var se pages.ErrorWithStatus
			if errors.As(err, &se) {
				http.Error(w, se.Message, se.Status)
				return
			}
			log.Printf("error rendering %s: %v", r.URL.Path, err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}),
	); err != nil {
		log.Fatal(err)
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
