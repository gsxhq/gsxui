package examples

import (
	"fmt"
	"net/url"

	"github.com/gsxhq/gsx"
)

// PreviewQueryKey selects one exact named preview document for an isolated
// example with multiple cases. The leading underscore reserves it for the
// gallery rather than an example's own request state.
const PreviewQueryKey = "_preview"

// Preview is one independently rendered document within a registered
// example. It lets a single shared source block demonstrate multiple
// viewport-owning cases without mounting those fixed layouts together.
type Preview struct {
	Name  string
	Title string
	Node  gsx.Node
}

// Example is one live demo on a component page. Name/SourcePath key the
// embedded source file (SourcePath is "{component}/{Name}.gsx", relative
// to this package — the same path Source reads), Title is the display
// heading, and Node is the already-constructed example ready to render.
// Examples that take no params (the common case) just call their
// zero-arg gsx component constructor directly, since gsx components
// render lazily — Node is safe to store before render time.
type Example struct {
	Name       string
	Title      string
	Node       gsx.Node
	SourcePath string
	// Isolated renders the live example in its own same-origin document.
	// Use it for application-shell components that intentionally own the
	// viewport; source display and registration otherwise stay unchanged.
	Isolated bool
	// ViewportWidth gives an isolated example a stable inner canvas in CSS
	// pixels. The parent preview surface remains fluid and clips this canvas;
	// zero keeps the iframe fluid with the surface.
	ViewportWidth int
	// Previews, when non-empty, renders one isolated document per named case
	// while retaining this example's single title and source block.
	Previews []Preview
	// PreviewRTL marks an isolated example whose document is right-to-left.
	// The parent preview surface clips a fixed-width canvas from its inline
	// end; setting this flips the surface to dir="rtl" so the canvas's RIGHT
	// edge — where an RTL layout anchors its chrome — stays visible instead
	// of the left one.
	PreviewRTL bool

	// Query, when non-nil, re-renders the example per-request from the
	// request's raw query parameters instead of using the static Node —
	// for an example whose content depends on request state (the calendar
	// examples' own ?month=, read by jstest/harness's /x/{component}/
	// {example} handler). This mechanism knows nothing about "month" or
	// "calendar": it is a generic url.Values -> gsx.Node hook, and the
	// interpretation of any given query key lives entirely in the example
	// package that sets this field (see site/examples/calendar.go).
	// Node must still be set for every example, Query or not — it backs
	// the static/default render (the site's own component pages, the
	// manifest, the drift and render-smoke tests in examples_test.go, and
	// any request with no matching query parameter).
	Query func(url.Values) gsx.Node
}

var (
	registry = map[string][]Example{}
	order    []string
)

// Register appends ex to component's example list, in call order. Each
// component's examples are registered together from that component's own
// init() (see button.go), so component pages iterate a stable,
// source-defined order — Task 3 batches append one file each, same pattern.
func Register(component string, ex Example) {
	if ex.ViewportWidth < 0 {
		panic(fmt.Sprintf(
			"examples: component %q example %q has negative viewport width %d",
			component,
			ex.Name,
			ex.ViewportWidth,
		))
	}
	if ex.ViewportWidth > 0 && !ex.Isolated {
		panic(fmt.Sprintf(
			"examples: component %q example %q has viewport width %d but is not isolated",
			component,
			ex.Name,
			ex.ViewportWidth,
		))
	}
	for _, registered := range registry[component] {
		if registered.Name == ex.Name {
			panic(fmt.Sprintf(
				"examples: duplicate registration for component %q example %q",
				component,
				ex.Name,
			))
		}
	}
	if _, ok := registry[component]; !ok {
		order = append(order, component)
	}
	registry[component] = append(registry[component], ex)
}

// For returns component's registered examples in registration order, or
// nil if component has none registered.
func For(component string) []Example {
	return registry[component]
}

// Find resolves one exact registered component/example pair and, when the
// reserved selector is present, exactly one named preview. Query-backed
// examples render from the supplied values; all others return their static
// node. Unknown, empty, or ambiguous keys are rejected rather than normalized
// or reconstructed.
func Find(component string, name string, query url.Values) (string, gsx.Node, bool) {
	for _, example := range registry[component] {
		if example.Name != name {
			continue
		}
		if previewNames, selected := query[PreviewQueryKey]; selected {
			if len(previewNames) != 1 || previewNames[0] == "" {
				return "", nil, false
			}
			for _, preview := range example.Previews {
				if preview.Name == previewNames[0] {
					return example.Title + " · " + preview.Title, preview.Node, true
				}
			}
			return "", nil, false
		}
		node := example.Node
		if example.Query != nil {
			node = example.Query(query)
		}
		return example.Title, node, true
	}
	return "", nil, false
}

// Components returns the names of components with at least one registered
// example, in registration order. This is a strict subset of
// internal/registry.Components() (all shipped ui/ components) until every
// component has examples registered (Task 3).
func Components() []string {
	return append([]string(nil), order...)
}
