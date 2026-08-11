package luma

import "github.com/gsxhq/gsx"

// Textarea is the shadcn/ui Textarea. HTML <textarea> takes its initial
// content as a text child, not a value attribute — shadcn's `...props`
// value pass-through has no gsx equivalent for that reason. value renders
// as the (escaped) text child instead (ADAPT, see docs/jsx-parity.md).
component Textarea(value string, attrs gsx.Attrs) {
	<textarea
		class={
			"field-sizing-content min-h-16 w-full outline-none placeholder:text-muted-foreground disabled:cursor-not-allowed disabled:opacity-50 bg-input/50 border-transparent focus-visible:border-ring focus-visible:ring-ring/30 aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive dark:aria-invalid:border-destructive/50 resize-none rounded-2xl border px-3 py-3 text-base transition-[color,box-shadow,background-color] focus-visible:ring-3 aria-invalid:ring-3 md:text-sm flex"
		}
		{ attrs... }
		data-gsxui-slot-textarea
	>{ value }</textarea>
}
