package main

import (
	"errors"
	"io/fs"
	"net/http"
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
	mux.Handle("/", v.Middleware(pagesMux))
}
