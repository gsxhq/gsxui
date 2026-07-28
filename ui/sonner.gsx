package ui

import (
	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// Toast is the server-rendered toast card — the single source of truth for
// the toast <li> markup. shadcn's own sonner.tsx renders nothing but a
// re-themed <Sonner> passthrough (the toast library owns 100% of the toast
// DOM from a non-Tailwind stylesheet), so there is no upstream markup to
// port. This is the ONE place the card is authored: ui/sonner.js clones a pre-rendered
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
// action/cancel buttons carry dedicated behavior hooks ui/sonner.js wires
// clicks onto. A custom auto-dismiss is a data-duration attr passed
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
		data-gsxui-toast
		data-type={t}
		role="status"
		aria-live={ariaLive}
		aria-atomic="true"
		{ attrs... }
		data-gsxui-slot-toast
	>
		{ if t != "default" {
			{ switch t {
			case "success":
				<icon.CircleCheck data-gsxui-toast-icon data-gsxui-slot-toast-icon/>
			case "info":
				<icon.Info data-gsxui-toast-icon data-gsxui-slot-toast-icon/>
			case "warning":
				<icon.TriangleAlert data-gsxui-toast-icon data-gsxui-slot-toast-icon/>
			case "error":
				<icon.OctagonX data-gsxui-toast-icon data-gsxui-slot-toast-icon/>
			case "loading":
				<icon.LoaderCircle data-gsxui-toast-icon data-gsxui-slot-toast-icon/>
			} }
		} }
		<div data-gsxui-slot-toast-content>
			<div data-gsxui-toast-title data-gsxui-slot-toast-title>{ title }</div>
			{ if description != "" {
				<div data-gsxui-toast-description data-gsxui-slot-toast-description>{ description }</div>
			} }
		</div>
		{ if action != "" {
			<button
				type="button"
				data-gsxui-toast-action
				data-gsxui-slot-toast-action
			>
				{ action }
			</button>
		} }
		{ if cancel != "" {
			<button
				type="button"
				data-gsxui-toast-cancel
				data-gsxui-slot-toast-cancel
			>
				{ cancel }
			</button>
		} }
		<button
			type="button"
			data-gsxui-toast-close
			aria-label="Close"
			data-gsxui-slot-toast-close
		>
			<icon.X data-gsxui-slot-toast-close-icon/>
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
// and pointer events pass through the empty gutter (each toast re-enables
// pointer events on itself through live lifecycle state).
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
			data-gsxui-toaster
			{ attrs... }
			data-gsxui-slot-toaster
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
