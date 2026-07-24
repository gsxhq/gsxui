package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

// TestToasterPinned pins the full server surface of the sonner port: the
// always-present bottom-right region (now with a stable id="gsxui-toaster"
// so server OOB/partial appends have a fixed target) followed by the six
// inert per-type <template>s ui/sonner.js clones each toast card from. The
// card markup lives in exactly one place — the ui.Toast component — so the
// templates and every ui.Toast call render one identical code path.
func TestToasterPinned(t *testing.T) {
	got := render(t, ui.Toaster(nil))
	want := `<section aria-label="Notifications" tabindex="-1"><ol id="gsxui-toaster" data-slot="toaster" data-gsxui-toaster class="pointer-events-none fixed z-100 flex flex-col gap-2 p-6 bottom-0 right-0"></ol><template data-gsxui-toast-template="default"><li data-slot="toast" data-gsxui-toast data-type="default" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Title</div><div data-description class="text-muted-foreground">Description</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Action</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li></template><template data-gsxui-toast-template="success"><li data-slot="toast" data-gsxui-toast data-type="success" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Title</div><div data-description class="text-muted-foreground">Description</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Action</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li></template><template data-gsxui-toast-template="info"><li data-slot="toast" data-gsxui-toast data-type="info" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Title</div><div data-description class="text-muted-foreground">Description</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Action</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li></template><template data-gsxui-toast-template="warning"><li data-slot="toast" data-gsxui-toast data-type="warning" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Title</div><div data-description class="text-muted-foreground">Description</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Action</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li></template><template data-gsxui-toast-template="error"><li data-slot="toast" data-gsxui-toast data-type="error" role="status" aria-live="assertive" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m15 9-6 6"/><path d="M2.586 16.726A2 2 0 0 1 2 15.312V8.688a2 2 0 0 1 .586-1.414l4.688-4.688A2 2 0 0 1 8.688 2h6.624a2 2 0 0 1 1.414.586l4.688 4.688A2 2 0 0 1 22 8.688v6.624a2 2 0 0 1-.586 1.414l-4.688 4.688a2 2 0 0 1-1.414.586H8.688a2 2 0 0 1-1.414-.586z"/><path d="m9 9 6 6"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Title</div><div data-description class="text-muted-foreground">Description</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Action</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li></template><template data-gsxui-toast-template="loading"><li data-slot="toast" data-gsxui-toast data-type="loading" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4 animate-spin"><path d="M21 12a9 9 0 1 1-6.219-8.56"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Title</div><div data-description class="text-muted-foreground">Description</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Action</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Cancel</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li></template></section>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// Attrs merge onto the <ol>; a caller id overrides the default gsxui-toaster.
func TestToasterAttrsMerge(t *testing.T) {
	got := render(t, ui.Toaster(gsx.Attrs{{Key: "id", Value: "my-toaster"}}))
	if !strings.Contains(got, `id="my-toaster"`) {
		t.Errorf("caller id missing\nin: %s", got)
	}
	if strings.Contains(got, `id="gsxui-toaster"`) {
		t.Errorf("default id should be overridden by caller\nin: %s", got)
	}
	if !strings.Contains(got, `data-gsxui-toaster`) {
		t.Errorf("data-gsxui-toaster hook missing\nin: %s", got)
	}
}

// TestToastPinned pins the toast card — the single source of truth for the
// <li> markup ui/sonner.js clones. This is a success card with a title, a
// description, an action and a cancel: every optional part present.
func TestToastPinnedFull(t *testing.T) {
	got := render(t, ui.Toast("error", "Failed", "Try again", "Retry", "Dismiss", nil))
	want := `<li data-slot="toast" data-gsxui-toast data-type="error" role="status" aria-live="assertive" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="m15 9-6 6"/><path d="M2.586 16.726A2 2 0 0 1 2 15.312V8.688a2 2 0 0 1 .586-1.414l4.688-4.688A2 2 0 0 1 8.688 2h6.624a2 2 0 0 1 1.414.586l4.688 4.688A2 2 0 0 1 22 8.688v6.624a2 2 0 0 1-.586 1.414l-4.688 4.688a2 2 0 0 1-1.414.586H8.688a2 2 0 0 1-1.414-.586z"/><path d="m9 9 6 6"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Failed</div><div data-description class="text-muted-foreground">Try again</div></div><button type="button" data-action class="shrink-0 self-center text-sm font-medium underline-offset-4 hover:underline">Retry</button><button type="button" data-cancel class="shrink-0 self-center text-sm text-muted-foreground underline-offset-4 hover:underline">Dismiss</button><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// A success card with title+description but no action/cancel — the common
// server-flash shape.
func TestToastPinnedSuccess(t *testing.T) {
	got := render(t, ui.Toast("success", "Changes saved", "Your settings are updated", "", "", nil))
	want := `<li data-slot="toast" data-gsxui-toast data-type="success" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-icon class="mt-0.5 shrink-0 [&amp;&gt;svg]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/></svg></div><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Changes saved</div><div data-description class="text-muted-foreground">Your settings are updated</div></div><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// A default toast (empty toastType normalises to "default") has no icon slot
// and no description/action/cancel — title only.
func TestToastPinnedDefault(t *testing.T) {
	got := render(t, ui.Toast("", "Hello", "", "", "", nil))
	want := `<li data-slot="toast" data-gsxui-toast data-type="default" role="status" aria-live="polite" aria-atomic="true" class="pointer-events-auto absolute bottom-6 right-6 flex w-[356px] items-start gap-3 rounded-2xl border border-border bg-popover p-4 text-sm text-popover-foreground shadow-lg origin-bottom transition-[transform,opacity] duration-300 ease-out data-[type=success]:[&amp;&gt;[data-icon]]:text-emerald-500 data-[type=info]:[&amp;&gt;[data-icon]]:text-sky-500 data-[type=warning]:[&amp;&gt;[data-icon]]:text-amber-500 data-[type=error]:[&amp;&gt;[data-icon]]:text-destructive"><div data-content class="flex flex-1 flex-col gap-1"><div data-title class="font-medium text-foreground">Hello</div></div><button type="button" data-close-button aria-label="Close" class="absolute -top-1.5 -right-1.5 flex size-5 items-center justify-center rounded-full border border-border bg-background text-foreground shadow-sm"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-3"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg></button></li>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

// Each type renders its own Lucide glyph (via ui/icon) and the correct
// aria-live level: error announces assertively, every other type politely.
// The default type has no icon slot at all.
func TestToastTypeIconsAndAria(t *testing.T) {
	cases := []struct {
		typ      string
		ariaLive string
		glyph    string // a distinctive path fragment of that type's icon
	}{
		{"default", "polite", ""},
		{"success", "polite", `<path d="m9 12 2 2 4-4"/>`},
		{"info", "polite", `<path d="M12 16v-4"/>`},
		{"warning", "polite", `d="m21.73 18-8-14`},
		{"error", "assertive", `<path d="m15 9-6 6"/>`},
		{"loading", "polite", `d="M21 12a9 9 0 1 1-6.219-8.56"`},
	}
	for _, c := range cases {
		got := render(t, ui.Toast(c.typ, "Title", "", "", "", nil))
		if !strings.Contains(got, `data-type="`+c.typ+`"`) {
			t.Errorf("%s: data-type missing\nin: %s", c.typ, got)
		}
		if !strings.Contains(got, `aria-live="`+c.ariaLive+`"`) {
			t.Errorf("%s: want aria-live %q\nin: %s", c.typ, c.ariaLive, got)
		}
		if c.typ == "default" {
			if strings.Contains(got, "<div data-icon") {
				t.Errorf("default toast should have no icon slot\nin: %s", got)
			}
			continue
		}
		if !strings.Contains(got, "<div data-icon") {
			t.Errorf("%s: icon slot missing\nin: %s", c.typ, got)
		}
		if !strings.Contains(got, c.glyph) {
			t.Errorf("%s: glyph %q missing\nin: %s", c.typ, c.glyph, got)
		}
	}
	// The loading glyph spins.
	loading := render(t, ui.Toast("loading", "Title", "", "", "", nil))
	if !strings.Contains(loading, "animate-spin") {
		t.Errorf("loading icon should carry animate-spin\nin: %s", loading)
	}
}

// Description is optional: an empty description renders no data-description.
func TestToastDescriptionOptional(t *testing.T) {
	with := render(t, ui.Toast("info", "Title", "Detail", "", "", nil))
	if !strings.Contains(with, `data-description`) {
		t.Errorf("description present should render data-description\nin: %s", with)
	}
	without := render(t, ui.Toast("info", "Title", "", "", "", nil))
	if strings.Contains(without, `data-description`) {
		t.Errorf("empty description should render no data-description\nin: %s", without)
	}
}

// Action/cancel are optional: present labels render the data-action/
// data-cancel buttons ui/sonner.js wires clicks onto; empty labels render
// neither.
func TestToastActionCancelPresence(t *testing.T) {
	with := render(t, ui.Toast("default", "Title", "", "Retry", "Dismiss", nil))
	if !strings.Contains(with, `data-action`) || !strings.Contains(with, ">Retry<") {
		t.Errorf("action button missing\nin: %s", with)
	}
	if !strings.Contains(with, `data-cancel`) || !strings.Contains(with, ">Dismiss<") {
		t.Errorf("cancel button missing\nin: %s", with)
	}
	without := render(t, ui.Toast("default", "Title", "", "", "", nil))
	if strings.Contains(without, `data-action`) {
		t.Errorf("empty action should render no data-action\nin: %s", without)
	}
	if strings.Contains(without, `data-cancel`) {
		t.Errorf("empty cancel should render no data-cancel\nin: %s", without)
	}
}

// A custom auto-dismiss is a data-duration attr passed through attrs, which
// ui/sonner.js reads on adoption.
func TestToastDurationPassthrough(t *testing.T) {
	got := render(t, ui.Toast("success", "Saved", "", "", "", gsx.Attrs{{Key: "data-duration", Value: "8000"}}))
	if !strings.Contains(got, `data-duration="8000"`) {
		t.Errorf("data-duration passthrough missing\nin: %s", got)
	}
}
