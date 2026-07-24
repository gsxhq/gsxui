// Package scrollarea holds the site's example gsx components for
// ui/scroll-area. "scroll-area" can't be a directory/Go package name
// (hyphen), so the directory drops it — same selectbox/switchctl-style
// workaround as select/switch. The registered example key stays the
// hyphenated "scroll-area" (see scrollarea.go).
package scrollarea

import "github.com/gsxhq/gsxui/ui"

// tags mirrors shadcn's own scroll-area-demo.tsx tag list shape
// (Array.from({length: 50}).map((_, i) => `v1.2.0-beta.${50 - i}`)),
// trimmed to 15 rows — enough to force real overflow in a h-72 box without
// padding the example out unnecessarily.
var tags = []string{
	"v1.2.0-beta.15", "v1.2.0-beta.14", "v1.2.0-beta.13", "v1.2.0-beta.12",
	"v1.2.0-beta.11", "v1.2.0-beta.10", "v1.2.0-beta.9", "v1.2.0-beta.8",
	"v1.2.0-beta.7", "v1.2.0-beta.6", "v1.2.0-beta.5", "v1.2.0-beta.4",
	"v1.2.0-beta.3", "v1.2.0-beta.2", "v1.2.0-beta.1",
}

// Basic mirrors shadcn's own scroll-area-demo.tsx: a h-72 w-48 bordered
// box, vertical (default orientation), a title row, then every tag row
// separated by a Separator — a good test of the CSS-first thumb showing up
// against a border-radius-clipped container (rounded-[inherit] on the
// collapsed div matters here, see ui/scroll-area.gsx's own doc comment).
component Basic() {
	<ui.ScrollArea class="h-72 w-48 rounded-md border">
		<div class="p-4">
			<h4 class="mb-4 text-sm leading-none font-medium">Tags</h4>
			{ for _, tag := range tags {
				<div>
					<div class="text-sm">{ tag }</div>
					<ui.Separator class="my-2"/>
				</div>
			} }
		</div>
	</ui.ScrollArea>
}
