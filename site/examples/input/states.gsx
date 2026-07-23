package input

import "github.com/gsxhq/gsxui/ui"

// States renders Input in its enabled, disabled, and invalid states —
// disabled is a bare boolean attribute, aria-invalid is a plain string
// attribute; both fall through Input's { attrs... } spread untouched.
component States() {
	<div class="flex max-w-sm flex-col gap-3">
		<ui.Input placeholder="Enabled"/>
		<ui.Input placeholder="Disabled" disabled/>
		<ui.Input placeholder="Invalid" aria-invalid="true" value="not-an-email"/>
	</div>
}
