// Package scrollarea holds the site's example gsx components for
// ui/scroll-area. "scroll-area" can't be a directory/Go package name
// (hyphen), so the directory drops it — same selectbox/switchctl-style
// workaround as select/switch. The registered example key stays the
// hyphenated "scroll-area" (see scrollarea.go).
package scrollarea

import "github.com/gsxhq/gsxui/ui"

// tagsRtl mirrors this dir's own basic.gsx tags list (shadcn's own
// scroll-area-demo.tsx tag shape), trimmed the same way.
var tagsRtl = []string{
	"v1.2.0-beta.15", "v1.2.0-beta.14", "v1.2.0-beta.13", "v1.2.0-beta.12",
	"v1.2.0-beta.11", "v1.2.0-beta.10", "v1.2.0-beta.9", "v1.2.0-beta.8",
	"v1.2.0-beta.7", "v1.2.0-beta.6", "v1.2.0-beta.5", "v1.2.0-beta.4",
	"v1.2.0-beta.3", "v1.2.0-beta.2", "v1.2.0-beta.1",
}

// Rtl mirrors shadcn's own scroll-area-rtl demo: the same bordered h-72
// w-48 box of tag rows as this dir's own Basic, title translated to
// Arabic and wrapped in dir="rtl".
component Rtl() {
	<ui.ScrollArea dir="rtl" lang="ar" class="h-72 w-48 rounded-md border">
		<div class="p-4">
			<h4 class="mb-4 text-sm leading-none font-medium">العلامات</h4>
			{ for _, tag := range tagsRtl {
				<div>
					<div class="text-sm">{ tag }</div>
					<ui.Separator class="my-2"/>
				</div>
			} }
		</div>
	</ui.ScrollArea>
}
