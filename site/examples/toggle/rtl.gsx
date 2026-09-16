package toggle

import (
	"github.com/gsxhq/gsxui/site/uirtl"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Rtl mirrors shadcn's own toggle-rtl demo: an outline, sm Toggle pairing a
// bookmark icon with a translated label, wrapped in dir="rtl".
component Rtl() {
	<div dir="rtl" lang="ar">
		<uirtl.Toggle variant="outline" size="sm" aria-label="Toggle bookmark">
			<icon.Bookmark/>
			إشارة مرجعية
		</uirtl.Toggle>
	</div>
}
