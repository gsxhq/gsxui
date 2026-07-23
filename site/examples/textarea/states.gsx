package textarea

import "github.com/gsxhq/gsxui/ui"

// States renders Textarea with a pre-filled value, then disabled and
// invalid variants.
component States() {
	<div class="flex max-w-sm flex-col gap-3">
		<ui.Textarea value="Filled in already." placeholder="Message"/>
		<ui.Textarea value="" placeholder="Disabled" disabled/>
		<ui.Textarea value="" placeholder="Invalid" aria-invalid="true"/>
	</div>
}
