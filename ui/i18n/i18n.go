// Package i18n holds the one type every gsxui component writes its own
// English strings through. It is a leaf package on purpose: a consumer's
// translator must name T, and the component code gsx generates must import
// the translator's package to call it, so T cannot live in package ui
// without an import cycle.
package i18n

// T is a message a component writes itself: an sr-only label, an
// aria-label, a visible "Close". Its value is the English text and doubles
// as the message id. It renders as written until the consuming module
// registers a renderer for it in gsx.toml:
//
//	[renderers]
//	"example.com/app/ui/i18n.T" = "example.com/app/ui/i18n.Translate"
//
// The translator lives beside this file, in a file you own (for example
// translate.go); gsxui add --overwrite rewrites only the files it vendored,
// so it survives re-vendoring:
//
//	func Translate(ctx context.Context, m T) string
//
// Renderers bind when gsx generate runs in your module, so only vendored
// components translate; importing github.com/gsxhq/gsxui/ui directly
// renders English.
//
// Text the caller supplies (children, props, attrs) never passes through T;
// only text the caller cannot reach does.
type T string
