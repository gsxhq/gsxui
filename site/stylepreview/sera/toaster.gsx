package sera

import "github.com/gsxhq/gsx"

// Toaster is the always-present, positioned toast region. Mount it ONCE per
// page (typically the root layout, same convention as shadcn's <Toaster/> in
// app/layout.tsx). v1 ships only the default bottom-right position — the
// other five sonner positions are a ledgered gap (docs/jsx-parity.md
// ## sonner).
//
// The <section> is the aria landmark ("Notifications"). The <ol> is the
// mount point ui/toaster.js observes: every toast <li> — whether inserted by
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
// a server flash viewport's per-severity templates. ui/toaster.js clones the
// matching type's template on each toast() call and fills or removes the
// title/description/action/cancel parts, so the card markup lives in exactly
// one place (the separately registered Toast component, ui/toast.gsx), never
// duplicated in JS. Their placeholder texts are always overwritten or removed
// on clone.
component Toaster(attrs gsx.Attrs) {
	<section aria-label="Notifications" tabindex="-1">
		<ol
			id="gsxui-toaster"
			class={ "[--gsxui-toast-offset:1.5rem] flex flex-col gap-2 p-6" }
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
