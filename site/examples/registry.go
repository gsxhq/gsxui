package examples

import "github.com/gsxhq/gsx"

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

// Components returns the names of components with at least one registered
// example, in registration order. This is a strict subset of
// internal/registry.Components() (all shipped ui/ components) until every
// component has examples registered (Task 3).
func Components() []string {
	return append([]string(nil), order...)
}
