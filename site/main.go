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
	// GO_PORT (dev): set by .env / `gsx dev`, wired to vite.config.ts's proxy
	// target — takes precedence so the dev loop is unaffected. PORT (prod):
	// Cloud Run injects this and expects the container to bind it; 8080 is
	// Cloud Run's own default, used as the final fallback for a bare
	// `docker run` with neither var set.
	port := cmp.Or(os.Getenv("GO_PORT"), os.Getenv("PORT"), "8080")
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
