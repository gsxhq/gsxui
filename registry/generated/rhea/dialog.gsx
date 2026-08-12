package ui

import "github.com/gsxhq/gsx"

// Dialog uses the native <dialog> top layer. Trigger/content wiring is scoped
// by the dedicated root hook and implemented by ui/dialog.js.
component Dialog(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "contents" } { attrs... } data-gsxui-slot-dialog>{ children }</div>
}

component DialogTrigger(children gsx.Node, attrs gsx.Attrs) {
	<button
		type="button"
		aria-haspopup="dialog"
		aria-expanded="false"
		{ attrs... }
		data-gsxui-slot-dialog-trigger
	>
		{ children }
	</button>
}

component DialogContent(hideCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<dialog
		class={
			"transition-none backdrop:transition-none data-[state=open]:backdrop:animate-in data-[state=closed]:backdrop:animate-out data-[state=closed]:backdrop:fade-out-0 data-[state=open]:backdrop:fade-in-0 backdrop:bg-black/30 backdrop:duration-100 supports-backdrop-filter:backdrop:backdrop-blur-sm bg-popover text-popover-foreground data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 ring-foreground/5 dark:ring-foreground/10 open:grid max-w-[calc(100%-2rem)] gap-6 rounded-[min(var(--radius-4xl),24px)] p-6 text-sm shadow-xl ring-1 duration-100 sm:max-w-md fixed z-50 top-1/2 left-1/2 w-full -translate-x-1/2 -translate-y-1/2 outline-none"
		}
		data-state="closed"
		{ attrs... }
		data-gsxui-slot-dialog-content
	>
		{ children }
		{ if !hideCloseButton {
			<button
				type="button"
				class={ "absolute top-4 end-4 bg-secondary" }
				data-gsxui-dialog-close
				data-gsxui-slot-dialog-close-button
				data-gsxui-slot-dialog-close
			>
				<svg
					aria-hidden="true"
					xmlns="http://www.w3.org/2000/svg"
					width="24"
					height="24"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
					class={ "size-4 shrink-0 pointer-events-none" }
					data-gsxui-slot-dialog-close-icon
				>
					<path d="M18 6 6 18"/>
					<path d="m6 6 12 12"/>
				</svg>
				<span data-gsxui-slot-dialog-close-label>Close</span>
			</button>
		} }
	</dialog>
}

component DialogHeader(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-1.5 flex flex-col" } { attrs... } data-gsxui-slot-dialog-header>{ children }</div>
}

component DialogFooter(showCloseButton bool, children gsx.Node, attrs gsx.Attrs) {
	<div class={ "gap-2 flex flex-col-reverse sm:flex-row sm:justify-end" } { attrs... } data-gsxui-slot-dialog-footer>
		{ children }
		{ if showCloseButton {
			<Button
				variant="outline"
				data-gsxui-dialog-close
				data-gsxui-slot-dialog-footer-close
			>
				Close
			</Button>
		} }
	</div>
}

component DialogTitle(children gsx.Node, attrs gsx.Attrs) {
	<h2 class={ "text-base leading-none font-medium" } { attrs... } data-gsxui-slot-dialog-title>{ children }</h2>
}

component DialogDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		class={ "text-muted-foreground *:[a]:hover:text-foreground text-sm *:[a]:underline *:[a]:underline-offset-3" }
		{ attrs... }
		data-gsxui-slot-dialog-description
	>
		{ children }
	</p>
}

component DialogClose(children gsx.Node, attrs gsx.Attrs) {
	<button data-gsxui-dialog-close type="button" { attrs... } data-gsxui-slot-dialog-close>{ children }</button>
}
