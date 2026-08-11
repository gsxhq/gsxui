package ui

import "github.com/gsxhq/gsx"

// Textarea is the shadcn/ui Textarea. HTML <textarea> takes its initial
// content as a text child, not a value attribute — shadcn's `...props`
// value pass-through has no gsx equivalent for that reason. value renders
// as the (escaped) text child instead (ADAPT, see docs/jsx-parity.md).
component Textarea(value string, attrs gsx.Attrs) {
	<textarea
		class={
			"border-transparent border-b-input bg-transparent focus-visible:border-b-ring aria-invalid:border-b-destructive dark:aria-invalid:border-b-destructive/50 resize-none rounded-none border px-0 py-3 text-base transition-[color,border-color] md:text-sm flex"
		}
		{ attrs... }
		data-gsxui-slot-textarea
	>{ value }</textarea>
}
