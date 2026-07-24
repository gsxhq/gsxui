package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestInputOTPPinned(t *testing.T) {
	// Exact full-render pin: the container wraps the ONE real, visually-
	// hidden <input> (opacity-0, z-10, absolute inset-0 covering the whole
	// slots row — relative on the container is what anchors it) plus the
	// caller's slot markup. No maxlength/pattern/name here — see
	// TestInputOTPAttrsFallThrough for those, all via attrs.
	got := render(t, ui.InputOTP(gsx.Raw("x"), nil))
	want := `<div data-slot="input-otp" data-gsxui-input-otp class="relative flex items-center gap-2 has-[input:disabled]:opacity-50"><input data-gsxui-input-otp-input inputmode="numeric" autocomplete="one-time-code" class="absolute inset-0 z-10 h-full w-full cursor-text opacity-0 disabled:cursor-not-allowed">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputOTPAttrsFallThrough(t *testing.T) {
	// maxlength/pattern/name/value/disabled/aria-invalid all fall through
	// via attrs onto the real input, matching ui/input.gsx's own
	// convention (no explicit Go params for any of them).
	got := render(t, ui.InputOTP(gsx.Raw("x"), gsx.Attrs{
		{Key: "maxlength", Value: "6"},
		{Key: "pattern", Value: "[0-9]*"},
		{Key: "name", Value: "otp"},
		{Key: "data-gsxui-input-otp-pattern", Value: "[0-9]"},
	}))
	for _, want := range []string{
		`maxlength="6"`,
		`pattern="[0-9]*"`,
		`name="otp"`,
		`data-gsxui-input-otp-pattern="[0-9]"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputOTPCallerClassMerges(t *testing.T) {
	got := render(t, ui.InputOTP(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "max-w-sm"}}))
	if !strings.Contains(got, "max-w-sm") || !strings.Contains(got, "relative flex items-center") {
		t.Errorf("want caller class plus surviving structural classes\nin: %s", got)
	}
}

func TestInputOTPGroupPinned(t *testing.T) {
	// Nova adds a group-level has-aria-invalid ring block absent from
	// new-york-v4's InputOTPGroup entirely (which only ever has per-slot
	// aria-invalid styling) — adopted per the 2026-07-24 controls source
	// map's `## input-otp` nova deltas table.
	got := render(t, ui.InputOTPGroup(gsx.Raw("x"), nil))
	want := `<div data-slot="input-otp-group" class="flex items-center has-aria-invalid:ring-destructive/20 dark:has-aria-invalid:ring-destructive/40 has-aria-invalid:border-destructive rounded-lg has-aria-invalid:ring-3">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputOTPSlotPinned(t *testing.T) {
	// Renders EMPTY at server-render time (data-active="false", no char, no
	// caret markup) — gsx has no client React context to pre-populate
	// char/isActive from an initial value the way shadcn's SSR-capable React
	// version can; ui/input-otp.js populates every slot after mount. This is
	// a real first-paint divergence from shadcn, ledgered in
	// docs/jsx-parity.md. Nova deltas applied: size-8 (was h-9 w-9), no
	// shadow-xs, first:rounded-l-lg/last:rounded-r-lg (was -md), ring-3 (was
	// ring-[3px]) — data-[active=true]:z-10 is KEPT regardless of nova's
	// excerpt (functionally necessary, not a deliberate drop, see the map).
	got := render(t, ui.InputOTPSlot(nil))
	want := `<div data-slot="input-otp-slot" data-active="false" class="relative flex size-8 items-center justify-center border-y border-r border-input text-sm transition-all outline-none first:rounded-l-lg first:border-l last:rounded-r-lg aria-invalid:border-destructive data-[active=true]:z-10 data-[active=true]:border-ring data-[active=true]:ring-3 data-[active=true]:ring-ring/50 data-[active=true]:aria-invalid:border-destructive data-[active=true]:aria-invalid:ring-destructive/20 dark:bg-input/30 dark:data-[active=true]:aria-invalid:ring-destructive/40"></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputOTPSlotNoCaretServerSide(t *testing.T) {
	// The fake-caret overlay (h-4 w-px animate-caret-blink bg-foreground
	// duration-1000, inside a centered pointer-events-none absolute inset-0
	// wrapper) is JS-only — ui/input-otp.js renders it into the active
	// slot's empty char only after mount, on input/selectionchange/
	// focus/blur. It must never appear in the server-rendered slot.
	got := render(t, ui.InputOTPSlot(nil))
	for _, absent := range []string{"animate-caret-blink", "pointer-events-none absolute inset-0"} {
		if strings.Contains(got, absent) {
			t.Errorf("server-rendered slot must not contain caret markup %q\nin: %s", absent, got)
		}
	}
}

func TestInputOTPSlotNoIndexParam(t *testing.T) {
	// Deliberate API departure from shadcn's InputOTPSlot(index): binding
	// decision is DOM-order stamping (option b, command.js's
	// data-gsxui-index source-order-stamp precedent) — InputOTPSlot takes
	// NO index param at all; ui/input-otp.js stamps data-index positionally
	// at mount. This test exists to document the call signature — it
	// wouldn't even compile if a caller tried ui.InputOTPSlot(0, nil).
	got := render(t, ui.InputOTPSlot(nil))
	if strings.Contains(got, "data-index") {
		t.Errorf("server render must not stamp data-index — that's ui/input-otp.js's job at mount\nin: %s", got)
	}
}

func TestInputOTPSeparatorPinned(t *testing.T) {
	// icon.Minus, static, unchanged from shadcn — plus nova's
	// [&_svg:not([class*='size-'])]:size-4 safeguard (a likely no-op since
	// icon.Minus already defaults to size-4, carried regardless per the
	// map's nova deltas table).
	got := render(t, ui.InputOTPSeparator(nil))
	want := `<div data-slot="input-otp-separator" role="separator" class="[&amp;_svg:not([class*=&#39;size-&#39;])]:size-4"><svg data-slot="icon" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="size-4"><path d="M5 12h14"/></svg></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputOTPGroupCallerClassMerges(t *testing.T) {
	got := render(t, ui.InputOTPGroup(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "gap-1"}}))
	if !strings.Contains(got, "gap-1") || !strings.Contains(got, "flex items-center") {
		t.Errorf("want caller class plus surviving structural classes\nin: %s", got)
	}
}

func TestInputOTPSlotAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputOTPSlot(gsx.Attrs{{Key: "aria-invalid", Value: "true"}}))
	if !strings.Contains(got, `aria-invalid="true"`) {
		t.Errorf("missing aria-invalid fallthrough\nin: %s", got)
	}
}
