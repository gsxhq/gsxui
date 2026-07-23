package textarea

import uitextarea "github.com/gsxhq/gsxui/ui/textarea"

// States renders Textarea with a pre-filled value, then disabled and
// invalid variants.
component States() {
	<div class="flex max-w-sm flex-col gap-3">
		<uitextarea.Textarea value="Filled in already." placeholder="Message"/>
		<uitextarea.Textarea value="" placeholder="Disabled" disabled/>
		<uitextarea.Textarea value="" placeholder="Invalid" aria-invalid="true"/>
	</div>
}
