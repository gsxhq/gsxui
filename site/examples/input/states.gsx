package input

import uiinput "github.com/gsxhq/gsxui/ui/input"

// States renders Input in its enabled, disabled, and invalid states —
// disabled is a bare boolean attribute, aria-invalid is a plain string
// attribute; both fall through Input's { attrs... } spread untouched.
component States() {
	<div class="flex max-w-sm flex-col gap-3">
		<uiinput.Input placeholder="Enabled"/>
		<uiinput.Input placeholder="Disabled" disabled/>
		<uiinput.Input placeholder="Invalid" aria-invalid="true" value="not-an-email"/>
	</div>
}
