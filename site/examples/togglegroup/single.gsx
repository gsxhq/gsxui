// Package togglegroup — see basic.gsx's own package doc comment.
package togglegroup

import "github.com/gsxhq/gsxui/ui"

// Single mirrors shadcn's toggle-group-single.tsx: groupType="single",
// default variant — the role="radio"/aria-checked branch (see
// ui/toggle-group.gsx's own doc comment) left untested by Basic's
// type="multiple" demo. "left" is server-rendered already pressed
// (pressed, aria-checked="true", data-state="on"); clicking "center" or
// "right" replaces it via toggle-group.js's single-type
// replace-on-activate mechanic.
component Single() {
	<ui.ToggleGroup groupType="single" aria-label="Text alignment">
		<ui.ToggleGroupItem groupType="single" value="left" pressed aria-label="Align left">
			Left
		</ui.ToggleGroupItem>
		<ui.ToggleGroupItem groupType="single" value="center" aria-label="Align center">
			Center
		</ui.ToggleGroupItem>
		<ui.ToggleGroupItem groupType="single" value="right" aria-label="Align right">
			Right
		</ui.ToggleGroupItem>
	</ui.ToggleGroup>
}
