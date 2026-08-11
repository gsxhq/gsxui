package ui

import "github.com/gsxhq/gsx"

// Textarea is the shadcn/ui Textarea. HTML <textarea> takes its initial
// content as a text child, not a value attribute — shadcn's `...props`
// value pass-through has no gsx equivalent for that reason. value renders
// as the (escaped) text child instead (ADAPT, see docs/jsx-parity.md).
component Textarea(value string, attrs gsx.Attrs) {
	<textarea
		class={
			"field-sizing-content min-h-16 w-full outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 border-input dark:bg-input/30 focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 disabled:bg-input/50 dark:disabled:bg-input/80 rounded-lg border bg-transparent px-2.5 py-2 text-base transition-colors focus-visible:ring-3 aria-invalid:ring-3 md:text-sm flex"
		}
		{ attrs... }
		data-gsxui-slot-textarea
	>{ value }</textarea>
}
