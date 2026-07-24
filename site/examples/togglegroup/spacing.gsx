// Package togglegroup — see basic.gsx's own package doc comment.
package togglegroup

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Spacing mirrors shadcn's toggle-group-spacing.tsx: type="multiple",
// variant="outline", spacing="2", size="sm" — the spacing/rounding/size
// stress case (spacing != "0" means the data-[spacing=0]: squared-off/
// joined-pill selectors on ToggleGroupItem never fire, so items render as
// separate rounded buttons with a real gap between them, sized by the
// root's style="--gap: 2").
component Spacing() {
	<ui.ToggleGroup groupType="multiple" variant="outline" spacing="2" size="sm" aria-label="Text formatting">
		<ui.ToggleGroupItem groupType="multiple" variant="outline" spacing="2" size="sm" value="bold" aria-label="Toggle bold">
			<icon.Bold/>
		</ui.ToggleGroupItem>
		<ui.ToggleGroupItem groupType="multiple" variant="outline" spacing="2" size="sm" value="italic" aria-label="Toggle italic">
			<icon.Italic/>
		</ui.ToggleGroupItem>
		<ui.ToggleGroupItem groupType="multiple" variant="outline" spacing="2" size="sm" value="underline" aria-label="Toggle underline">
			<icon.Underline/>
		</ui.ToggleGroupItem>
	</ui.ToggleGroup>
}
