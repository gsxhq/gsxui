package ui

import "github.com/gsxhq/gsx"

// Textarea is the shadcn/ui Textarea. HTML <textarea> takes its initial
// content as a text child, not a value attribute — shadcn's `...props`
// value pass-through has no gsx equivalent for that reason. value renders
// as the (escaped) text child instead (ADAPT, see docs/jsx-parity.md).
component Textarea(value string, attrs gsx.Attrs) {
	<textarea
		class={
			"border-input bg-input/20 dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 resize-none rounded-md border px-2 py-2 text-sm transition-colors focus-visible:ring-2 aria-invalid:ring-2 md:text-xs/relaxed flex"
		}
		{ attrs... }
		data-gsxui-slot-textarea
	>{ value }</textarea>
}
