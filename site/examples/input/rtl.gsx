package input

import "github.com/gsxhq/gsxui/ui"

// Rtl keeps Basic's single default Input, wrapped in dir="rtl". shadcn's
// own input-rtl demo wraps it in a Field/FieldLabel/FieldDescription API
// key form, which this directory's Basic doesn't compose; DEVIATED from
// rather than invented.
component Rtl() {
	<div dir="rtl" lang="ar">
		<ui.Input type="email" placeholder="انت@مثال.com"/>
	</div>
}
