package toggle

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own toggle-rtl demo: an outline, sm Toggle pairing a
// bookmark icon with a translated label, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Toggle variant="outline" size="sm" aria-label="Toggle bookmark">
			<icon.Bookmark/>
			إشارة مرجعية
		</ui.Toggle>
	</div>
}
