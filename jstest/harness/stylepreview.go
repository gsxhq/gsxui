package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/site/stylepreview/luma"
	"github.com/gsxhq/gsxui/site/stylepreview/lyra"
	"github.com/gsxhq/gsxui/site/stylepreview/maia"
	"github.com/gsxhq/gsxui/site/stylepreview/mira"
	"github.com/gsxhq/gsxui/site/stylepreview/nova"
	"github.com/gsxhq/gsxui/site/stylepreview/rhea"
	"github.com/gsxhq/gsxui/site/stylepreview/sera"
	"github.com/gsxhq/gsxui/site/stylepreview/vega"
)

// stylePreviewGalleries renders one style's Gallery composition standalone
// (not the twice-in-one-document theme editor preview), for
// jstest/specs/style-visual.spec.ts's per-style, per-card snapshot gate.
var stylePreviewGalleries = map[string]func(idp string) gsx.Node{
	"vega": vega.Gallery,
	"nova": nova.Gallery,
	"maia": maia.Gallery,
	"lyra": lyra.Gallery,
	"mira": mira.Gallery,
	"luma": luma.Gallery,
	"sera": sera.Gallery,
	"rhea": rhea.Gallery,
}

// registerStylePreviewRoutes serves GET /stylepreview/{style}, the harness
// route the visual gate navigates to for each style under test. It reuses
// the production preview entry (jstest/harness-preview.js, the browser-
// native equivalent of web/preview.js) so the gallery renders with real
// component behaviour — menus, dialogs, and tooltips resting closed — the
// same as the theme editor's own preview iframe.
func registerStylePreviewRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /stylepreview/{style}", func(w http.ResponseWriter, r *http.Request) {
		style := r.PathValue("style")
		gallery, ok := stylePreviewGalleries[style]
		if !ok {
			http.NotFound(w, r)
			return
		}
		var buf bytes.Buffer
		if err := gallery(style).Render(r.Context(), &buf); err != nil {
			http.Error(w, "render: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		documentTitle := "stylepreview/" + style
		if err := renderShell(w, documentTitle, stylesheetFor(r), "/static/jstest/harness-preview.js", template.HTML(buf.String())); err != nil {
			log.Printf("rendering stylepreview shell for %s: %v", style, err)
		}
	})
}
