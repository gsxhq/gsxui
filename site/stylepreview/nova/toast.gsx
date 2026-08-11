package nova

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Toast is the server-rendered toast card — the single source of truth for
// the toast <li> markup, and a standalone component: it is usable on its own
// (a server-rendered flash row, a static markup showcase) without Toaster.
// That is why it is registered separately rather than shipped as part of
// ui/toaster.gsx, matching upstream shadcn's split between the toast card and
// the toaster region.
//
// shadcn's own sonner.tsx renders nothing but a re-themed <Sonner>
// passthrough (the toast library owns 100% of the toast DOM from a
// non-Tailwind stylesheet), so there is no upstream markup to port. This is
// the ONE place the card is authored: ui/toaster.js clones a pre-rendered
// Toast (one per type, shipped as inert <template>s by Toaster) rather than
// building the card from JS string DOM — the old "icon paths hand-copied
// into a JS module" maintenance seam is gone (docs/jsx-parity.md ## sonner).
//
// toastType is one of default/success/info/warning/error/loading (the Go
// keyword `type` forces the param name); empty is normalised to "default".
// The type drives the icon (via ui/icon), the data-type styling axis, and
// the aria-live level: an error toast
// announces assertively, every other type politely. description/action/
// cancel are optional — an empty string renders the part absent, matching
// the JS `toast(msg, { description, action, cancel })` option surface; the
// action/cancel buttons carry dedicated behavior hooks ui/toaster.js wires
// clicks onto. A custom auto-dismiss is a data-duration attr passed
// through attrs (ui/toaster.js reads it on adoption; loading defaults to no
// auto-dismiss).
component Toast(toastType string, title string, description string, action string, cancel string, attrs gsx.Attrs) {
	{{
		t := toastType
		if t == "" {
			t = "default"
		}
		ariaLive := "polite"
		if t == "error" {
			ariaLive = "assertive"
		}
	}}
	<li
		data-type={t}
		role="status"
		aria-live={ariaLive}
		aria-atomic="true"
		class={
			"rounded-2xl flex w-[356px] items-start gap-3 border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg transition-[transform,opacity] duration-300 ease-out"
		}
		{ attrs... }
		data-gsxui-slot-toast
	>
		{ if t != "default" {
			{ switch t {
			case "success":
				<icon.CircleCheck
					class={
						"mt-0.5 size-4 shrink-0 [[data-gsxui-slot-toast][data-type='success']>&]:text-success [[data-gsxui-slot-toast][data-type='info']>&]:text-info [[data-gsxui-slot-toast][data-type='warning']>&]:text-warning [[data-gsxui-slot-toast][data-type='error']>&]:text-destructive [[data-gsxui-slot-toast][data-type='loading']>&]:animate-spin"
					}
					data-gsxui-slot-toast-icon
				/>
			case "info":
				<icon.Info
					class={
						"mt-0.5 size-4 shrink-0 [[data-gsxui-slot-toast][data-type='success']>&]:text-success [[data-gsxui-slot-toast][data-type='info']>&]:text-info [[data-gsxui-slot-toast][data-type='warning']>&]:text-warning [[data-gsxui-slot-toast][data-type='error']>&]:text-destructive [[data-gsxui-slot-toast][data-type='loading']>&]:animate-spin"
					}
					data-gsxui-slot-toast-icon
				/>
			case "warning":
				<icon.TriangleAlert
					class={
						"mt-0.5 size-4 shrink-0 [[data-gsxui-slot-toast][data-type='success']>&]:text-success [[data-gsxui-slot-toast][data-type='info']>&]:text-info [[data-gsxui-slot-toast][data-type='warning']>&]:text-warning [[data-gsxui-slot-toast][data-type='error']>&]:text-destructive [[data-gsxui-slot-toast][data-type='loading']>&]:animate-spin"
					}
					data-gsxui-slot-toast-icon
				/>
			case "error":
				<icon.OctagonX
					class={
						"mt-0.5 size-4 shrink-0 [[data-gsxui-slot-toast][data-type='success']>&]:text-success [[data-gsxui-slot-toast][data-type='info']>&]:text-info [[data-gsxui-slot-toast][data-type='warning']>&]:text-warning [[data-gsxui-slot-toast][data-type='error']>&]:text-destructive [[data-gsxui-slot-toast][data-type='loading']>&]:animate-spin"
					}
					data-gsxui-slot-toast-icon
				/>
			case "loading":
				<icon.LoaderCircle
					class={
						"mt-0.5 size-4 shrink-0 [[data-gsxui-slot-toast][data-type='success']>&]:text-success [[data-gsxui-slot-toast][data-type='info']>&]:text-info [[data-gsxui-slot-toast][data-type='warning']>&]:text-warning [[data-gsxui-slot-toast][data-type='error']>&]:text-destructive [[data-gsxui-slot-toast][data-type='loading']>&]:animate-spin"
					}
					data-gsxui-slot-toast-icon
				/>
			} }
		} }
		<div class={ "flex flex-1 flex-col gap-1" } data-gsxui-slot-toast-content>
			<div class={ "font-medium text-foreground" } data-gsxui-slot-toast-title>{ title }</div>
			{ if description != "" {
				<div class={ "text-muted-foreground" } data-gsxui-slot-toast-description>{ description }</div>
			} }
		</div>
		{ if action != "" {
			<button
				type="button"
				class={ "shrink-0 self-center text-sm underline-offset-4 font-medium hover:underline" }
				data-gsxui-slot-toast-action
			>
				{ action }
			</button>
		} }
		{ if cancel != "" {
			<button
				type="button"
				class={ "shrink-0 self-center text-sm underline-offset-4 text-muted-foreground hover:underline" }
				data-gsxui-slot-toast-cancel
			>
				{ cancel }
			</button>
		} }
		<button
			type="button"
			class={
				"absolute -top-1.5 -end-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"
			}
			aria-label="Close"
			data-gsxui-slot-toast-close
		>
			<icon.X class={ "size-3" } data-gsxui-slot-toast-close-icon/>
		</button>
	</li>
}
