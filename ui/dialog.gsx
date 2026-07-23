package ui

import "github.com/gsxhq/gsx"

// Dialog is the shadcn/ui Dialog on the native <dialog> element: the top
// layer replaces Radix's Portal, ::backdrop replaces Overlay, and Esc-to-
// close is browser-native. Trigger and content are wired by proximity —
// DialogTrigger opens the <dialog> inside the same Dialog root, no ids.
// Requires the dialog behavior module (ui/dialog/dialog.js).

component Dialog(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dialog" data-gsxui-dialog class="contents" { attrs... }>{ children }</div>
}

component DialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="dialog-trigger" data-gsxui-dialog-trigger type="button" aria-haspopup="dialog" aria-expanded="false" { attrs... }>{ children }</button>
}

// DialogContent renders the native <dialog>. hideCloseButton omits the
// injected top-right X button (shadcn's showCloseButton default-true,
// inverted so the gsx zero value keeps the shadcn default).
component DialogContent(hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		data-slot="dialog-content"
		data-gsxui-dialog-content
		data-state="closed"
		class={
			"fixed top-[50%] left-[50%] z-50 open:grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border bg-background p-6 text-foreground shadow-lg duration-200 outline-none sm:max-w-lg",
			"data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[state=open]:animate-in data-[state=open]:fade-in-0 data-[state=open]:zoom-in-95",
			"backdrop:bg-black/50",
		}
		{ attrs... }
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				data-slot="dialog-close"
				data-gsxui-dialog-close
				aria-label="Close"
				class="absolute top-4 right-4 rounded-xs opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:ring-2 focus:ring-ring focus:ring-offset-2 focus:outline-hidden disabled:pointer-events-none [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>
			</button>
		} }
	</dialog>
}

component DialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dialog-header" class="flex flex-col gap-2 text-center sm:text-left" { attrs... }>{ children }</div>
}

// DialogFooter is the shadcn/ui DialogFooter. showCloseButton (zero value
// false, matching shadcn's default) appends an outline Close button — the
// data-attribute idiom standing in for shadcn's <DialogClose asChild>.
component DialogFooter(showCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="dialog-footer" class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end" { attrs... }>
		{ children }
		{ if showCloseButton {
			<Button variant="outline" data-gsxui-dialog-close>Close</Button>
		} }
	</div>
}

component DialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 data-slot="dialog-title" class="text-lg leading-none font-semibold" { attrs... }>{ children }</h2>
}

component DialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p data-slot="dialog-description" class="text-sm text-muted-foreground" { attrs... }>{ children }</p>
}

component DialogClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-slot="dialog-close" data-gsxui-dialog-close type="button" { attrs... }>{ children }</button>
}
