package input

import uiinput "github.com/gsxhq/gsxui/ui/input"

// Custom overrides Input's default height by passing class="h-12" at the
// call site. Input's own h-9 is part of the same class-merge group, so the
// caller's h-12 wins and h-9 is dropped entirely — not appended alongside
// it — rather than the two classes fighting in the cascade.
component Custom() {
	<uiinput.Input class="h-12 text-base" placeholder="Larger input"/>
}
