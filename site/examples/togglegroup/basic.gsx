// Package togglegroup holds the site's example gsx components for
// ui/toggle-group. "toggle-group" can't be a Go package name (hyphen), so
// the directory drops it — same selectbox/switchctl-style workaround as
// select/switch, not a Go-keyword issue this time. The registered example
// key stays the hyphenated "toggle-group" (see togglegroup.go).
package togglegroup

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Basic mirrors shadcn's own toggle-group-demo.tsx: type="multiple",
// variant="outline", 3 icon-only items, none pressed at first paint —
// every ToggleGroupItem repeats the group's groupType/variant/size/spacing
// explicitly (see ui/toggle-group.gsx's own GAP doc comment on why: no
// context, the caller resolves and passes down). groupType, not type — see
// ToggleGroup's own doc comment on the Go-keyword param rename.
component Basic() {
	<ui.ToggleGroup groupType="multiple" variant="outline" aria-label="Text formatting">
		<ui.ToggleGroupItem groupType="multiple" variant="outline" value="bold" aria-label="Toggle bold">
			<icon.Bold/>
		</ui.ToggleGroupItem>
		<ui.ToggleGroupItem groupType="multiple" variant="outline" value="italic" aria-label="Toggle italic">
			<icon.Italic/>
		</ui.ToggleGroupItem>
		<ui.ToggleGroupItem groupType="multiple" variant="outline" value="underline" aria-label="Toggle underline">
			<icon.Underline/>
		</ui.ToggleGroupItem>
	</ui.ToggleGroup>
}
