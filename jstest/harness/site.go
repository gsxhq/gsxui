package main

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing/fstest"

	"github.com/gsxhq/gsxui/site/pages"
	"github.com/gsxhq/vite"
	"github.com/jackielii/structpages"
)

// testSiteManifest makes the real site pages use the harness's compiled CSS
// and browser-native module entry. The manifest key stays web/main.js because
// siteHead is production code; only the resolved test assets differ.
var testSiteManifest fs.FS = fstest.MapFS{
	".vite/manifest.json": {
		Data: []byte(`{
			"web/main.js": {
				"file": "jstest/harness-site.js",
				"src": "web/main.js",
				"isEntry": true,
				"css": ["jstest/.tmp/site.css"]
			},
			"web/preview.js": {
				"file": "jstest/harness-preview.js",
				"src": "web/preview.js",
				"isEntry": true,
				"css": ["jstest/.tmp/site.css"]
			}
		}`),
	},
}

func registerSiteRoutes(mux *http.ServeMux) {
	pagesMux := http.NewServeMux()
	if _, err := structpages.Mount(
		pagesMux,
		pages.Pages{},
		"/",
		"gsxui",
		structpages.WithErrorHandler(func(w http.ResponseWriter, _ *http.Request, err error) {
			if statusError, ok := errors.AsType[pages.ErrorWithStatus](err); ok {
				http.Error(w, statusError.Message, statusError.Status)
				return
			}
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}),
	); err != nil {
		panic("mounting site pages in test harness: " + err.Error())
	}

	v, err := vite.New(vite.Config{Dist: testSiteManifest})
	if err != nil {
		panic("loading test site manifest: " + err.Error())
	}
	siteHandler := withBrowserImportMap(v.Middleware(pagesMux))
	mux.Handle("/site/", http.StripPrefix("/site", siteHandler))
	mux.Handle("/", siteHandler)
}

func withBrowserImportMap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := httptest.NewRecorder()
		next.ServeHTTP(recorder, r)

		result := recorder.Result()
		defer result.Body.Close()

		for key, values := range result.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		body := recorder.Body.Bytes()
		if bytes.HasPrefix([]byte(result.Header.Get("Content-Type")), []byte("text/html")) {
			// The test manifest swaps the production bundle for browser-native
			// modules. Keep their import map before the exact module tag emitted
			// by siteHead so the browser resolves bare package imports.
			moduleTag := []byte(`<script type="module"`)
			body = bytes.Replace(body, moduleTag, append([]byte(browserImportMap+"\n"), moduleTag...), 1)
			w.Header().Del("Content-Length")
		}

		w.WriteHeader(result.StatusCode)
		_, _ = w.Write(body)
	})
}
