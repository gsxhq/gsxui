// Package hl serves syntax-highlighted code blocks for the docs site.
//
// The highlighting itself is NOT done here. tree-sitter is C, and the site
// binary is built CGO_ENABLED=0 into a distroless/static image (see
// site/Dockerfile), so the highlighter cannot be linked in. Instead
// site/hl/gen — a separate Go module, kept out of this one precisely so its
// cgo dependencies never reach the server build — renders every block ahead
// of time into blocks.gen.go, which is committed the same way the .x.go
// files are. This package is then pure data: a map lookup and a Raw node.
//
// Regenerate with `make highlight` after changing any example or snippet.
// TestBlocksCoverSources fails when the committed output goes stale.
package hl

import (
	"html"

	"github.com/gsxhq/gsx"
)

// Node returns the highlighted HTML for key as a renderable node, ready to
// drop inside a <pre><code>.
//
// Keys are source-relative paths: "button/basic.gsx" for a component example
// (an Example's SourcePath), "snippets/install.sh" for a doc snippet.
//
// An unknown key renders its own name as escaped text rather than failing:
// a doc page that references a block nobody generated should look obviously
// wrong on the page, not take the whole site down at render time. The
// generator's coverage test is what actually prevents that from shipping.
func Node(key string) gsx.Node {
	if h, ok := blocks[key]; ok {
		return gsx.Raw(h)
	}
	return gsx.Raw("<span class=\"ts-comment\">missing highlighted block: " + html.EscapeString(key) + "</span>")
}

// Has reports whether a highlighted block exists for key.
func Has(key string) bool {
	_, ok := blocks[key]
	return ok
}

// Keys returns every generated block key, for tests.
func Keys() []string {
	out := make([]string, 0, len(blocks))
	for k := range blocks {
		out = append(out, k)
	}
	return out
}
