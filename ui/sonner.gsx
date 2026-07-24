package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Toast is the server-rendered toast card — the single source of truth for
// the toast <li> markup. shadcn's own sonner.tsx renders nothing but a
// re-themed <Sonner> passthrough (the toast library owns 100% of the toast
// DOM from a non-Tailwind stylesheet), so there is no upstream markup to
// port; gsxui reconstructs the look as plain Tailwind classes here. This is
// the ONE place the card is authored: ui/sonner.js clones a pre-rendered
// Toast (one per type, shipped as inert <template>s by Toaster) rather than
// building the card from JS string DOM — the old "icon paths hand-copied
// into a JS module" maintenance seam is gone (docs/jsx-parity.md ## sonner).
//
// toastType is one of default/success/info/warning/error/loading (the Go
// keyword `type` forces the param name); empty is normalised to "default".
// The type drives the icon (via ui/icon), the data-type attribute (which the
// class list tints the icon from), and the aria-live level: an error toast
// announces assertively, every other type politely. description/action/
// cancel are optional — an empty string renders the part absent, matching
// the JS `toast(msg, { description, action, cancel })` option surface; the
// action/cancel buttons carry the data-action/data-cancel hooks ui/sonner.js
// wires clicks onto. A custom auto-dismiss is a data-duration attr passed
// through attrs (ui/sonner.js reads it on adoption; loading defaults to no
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
		data-slot="toast"
		data-gsxui-toast
		data-type={t}
		role="status"
		aria-live={ariaLive}
		aria-atomic="true"
		class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&>[data-icon]]:text-emerald-500 data-[type=info]:[&>[data-icon]]:text-sky-500 data-[type=warning]:[&>[data-icon]]:text-amber-500 data-[type=error]:[&>[data-icon]]:text-destructive"
		{ attrs... }
	>
		{ if t != "default" {
			<div data-icon class="mt-0.5 shrink-0 [&>svg]:size-4">
				{ switch t {
				case "success":
					<icon.CircleCheck/>
				case "info":
					<icon.Info/>
				case "warning":
					<icon.TriangleAlert/>
				case "error":
					<icon.OctagonX/>
				case "loading":
					<icon.LoaderCircle class="animate-spin"/>
				} }
			</div>
		} }
		<div data-content class="flex flex-1 flex-col gap-1">
			<div data-title class="font-medium text-foreground">{ title }</div>
			{ if description != "" {
				<div data-description class="text-muted-foreground">{ description }</div>
			} }
		</div>
		{ if action != "" {
			<button
				type="button"
				data-action
				class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline"
			>{ action }</button>
		} }
		{ if cancel != "" {
			<button
				type="button"
				data-cancel
				class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline"
			>{ cancel }</button>
		} }
		<button
			type="button"
			data-close-button
			aria-label="Close"
			class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"
		>
			<icon.X class="size-3"/>
		</button>
	</li>
}

// Toaster is the always-present, positioned toast region. Mount it ONCE per
// page (typically the root layout, same convention as shadcn's <Toaster/> in
// app/layout.tsx). v1 ships only the default bottom-right position — the
// other five sonner positions are a ledgered gap (docs/jsx-parity.md
// ## sonner).
//
// The <section> is the aria landmark ("Notifications"). The <ol> is the
// mount point ui/sonner.js observes: every toast <li> — whether inserted by
// the imperative toast() API, cloned from a template by the declarative
// trigger, or appended by the server (a full-page-load flash rendered inline,
// or an HTMX out-of-band swap `hx-swap-oob="beforeend:#gsxui-toaster"`) —
// lands here and is adopted by a MutationObserver into the same stacking /
// timer / dismiss lifecycle. It carries a stable id="gsxui-toaster" (caller-
// overridable via attrs) so server OOB/partial appends have a fixed target,
// and pointer-events-none so clicks fall through the empty gutter (each toast
// re-enables pointer-events on itself).
//
// After the <ol> come six inert <template>s, one per type — the same idiom as
// a server flash viewport's per-severity templates. ui/sonner.js clones the
// matching type's template on each toast() call and fills or removes the
// title/description/action/cancel parts, so the card markup lives in exactly
// one place (the Toast component above), never duplicated in JS. Their
// placeholder texts are always overwritten or removed on clone.
component Toaster(attrs gsx.Attrs) {
	<section aria-label="Notifications" tabindex="-1">
		<ol
			id="gsxui-toaster"
			data-slot="toaster"
			data-gsxui-toaster
			class="pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0"
			{ attrs... }
		></ol>
		<template data-gsxui-toast-template="default">
			<Toast toastType="default" title="Title" description="Description" action="Action" cancel="Cancel"/>
		</template>
		<template data-gsxui-toast-template="success">
			<Toast toastType="success" title="Title" description="Description" action="Action" cancel="Cancel"/>
		</template>
		<template data-gsxui-toast-template="info">
			<Toast toastType="info" title="Title" description="Description" action="Action" cancel="Cancel"/>
		</template>
		<template data-gsxui-toast-template="warning">
			<Toast toastType="warning" title="Title" description="Description" action="Action" cancel="Cancel"/>
		</template>
		<template data-gsxui-toast-template="error">
			<Toast toastType="error" title="Title" description="Description" action="Action" cancel="Cancel"/>
		</template>
		<template data-gsxui-toast-template="loading">
			<Toast toastType="loading" title="Title" description="Description" action="Action" cancel="Cancel"/>
		</template>
	</section>
}
