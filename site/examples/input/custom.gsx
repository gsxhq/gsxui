package input

import "github.com/gsxhq/gsxui/ui"

// Custom overrides Input's default height by passing class="h-12 text-base" at the
// call site. Input's own h-8 is part of the same class-merge group, so the
// caller's h-12 wins and h-8 is dropped entirely — not appended alongside
// it — rather than the two classes fighting in the cascade.
component Custom() {
	<ui.Input class="h-12 text-base" placeholder="Larger input"/>
}
