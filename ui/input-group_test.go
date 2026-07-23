package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestInputGroupPinned(t *testing.T) {
	got := render(t, ui.InputGroup(gsx.Raw("x"), nil))
	want := `<div data-slot="input-group" role="group" class="group/input-group relative flex w-full items-center rounded-md border border-input shadow-xs transition-[color,box-shadow] outline-none dark:bg-input/30 h-9 min-w-0 has-[&gt;textarea]:h-auto has-[&gt;[data-align=inline-start]]:[&amp;&gt;input]:pl-2 has-[&gt;[data-align=inline-end]]:[&amp;&gt;input]:pr-2 has-[&gt;[data-align=block-start]]:h-auto has-[&gt;[data-align=block-start]]:flex-col has-[&gt;[data-align=block-start]]:[&amp;&gt;input]:pb-3 has-[&gt;[data-align=block-end]]:h-auto has-[&gt;[data-align=block-end]]:flex-col has-[&gt;[data-align=block-end]]:[&amp;&gt;input]:pt-3 has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-[3px] has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50 has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-destructive/20 dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroup(nil, gsx.Attrs{{Key: "id", Value: "ig1"}}))
	if !strings.Contains(got, `id="ig1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestInputGroupCallerClassMerges(t *testing.T) {
	got := render(t, ui.InputGroup(nil, gsx.Attrs{{Key: "class", Value: "max-w-sm"}}))
	if !strings.Contains(got, "max-w-sm") {
		t.Errorf("missing caller class max-w-sm\nin: %s", got)
	}
}

// TestInputGroupAddonDefaultPinned pins the zero-value ("inline-start") align.
func TestInputGroupAddonDefaultPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("", gsx.Raw("x"), nil))
	want := `<div role="group" data-slot="input-group-addon" data-align="inline-start" class="flex h-auto cursor-text items-center justify-center gap-2 py-1.5 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [&amp;&gt;kbd]:rounded-[calc(var(--radius)-5px)] [&amp;&gt;svg:not([class*=&#39;size-&#39;])]:size-4 order-first pl-3 has-[&gt;button]:ml-[-0.45rem] has-[&gt;kbd]:ml-[-0.35rem]">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupAddonInlineEndPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("inline-end", gsx.Raw("x"), nil))
	for _, want := range []string{`data-align="inline-end"`, "order-last pr-3", "has-[&gt;button]:mr-[-0.45rem]"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupAddonBlockStartPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("block-start", gsx.Raw("x"), nil))
	want := `<div role="group" data-slot="input-group-addon" data-align="block-start" class="flex h-auto cursor-text items-center gap-2 py-1.5 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [&amp;&gt;kbd]:rounded-[calc(var(--radius)-5px)] [&amp;&gt;svg:not([class*=&#39;size-&#39;])]:size-4 order-first w-full justify-start px-3 pt-3 group-has-[&gt;input]/input-group:pt-2.5 [.border-b]:pb-3">x</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "justify-center") {
		t.Errorf("justify-center should be dropped by justify-start\nin: %s", got)
	}
}

func TestInputGroupAddonBlockEndPinned(t *testing.T) {
	got := render(t, ui.InputGroupAddon("block-end", gsx.Raw("x"), nil))
	for _, want := range []string{`data-align="block-end"`, "order-last w-full justify-start px-3 pb-3", "group-has-[&gt;input]/input-group:pb-2.5", "[.border-t]:pt-3"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupAddonAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupAddon("", nil, gsx.Attrs{{Key: "id", Value: "a1"}}))
	if !strings.Contains(got, `id="a1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestInputGroupButtonDefaultPinned proves InputGroupButton actually
// composes ui.Button (data-slot="button" and Button's own base/variant
// classes all come through) and that the "xs" overlay classes win their
// tailwind-merge conflicts against Button's own default size classes: h-9 ->
// h-6, rounded-md -> rounded-[calc(var(--radius)-5px)], px-4 -> px-2,
// has-[>svg]:px-3 -> has-[>svg]:px-2. py-2 (Button's own default size class,
// which the xs overlay never mentions) survives untouched, matching
// shadcn's own cn() merge exactly (see ui/input-group.gsx's own comment).
func TestInputGroupButtonDefaultPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", gsx.Raw("x"), nil))
	want := `<button data-slot="button" data-variant="ghost" type="button" class="shrink-0 justify-center font-medium whitespace-nowrap transition-all outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 [&amp;_svg]:pointer-events-none [&amp;_svg]:shrink-0 [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4 hover:bg-accent hover:text-accent-foreground dark:hover:bg-accent/50 py-2 flex items-center text-sm shadow-none h-6 gap-1 rounded-[calc(var(--radius)-5px)] px-2 has-[&gt;svg]:px-2 [&amp;&gt;svg:not([class*=&#39;size-&#39;])]:size-3.5" data-size="xs">x</button>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "h-9") || strings.Contains(got, "px-4") || strings.Contains(got, "rounded-md\"") {
		t.Errorf("Button's own default size classes should be overridden by the xs overlay\nin: %s", got)
	}
}

func TestInputGroupButtonSmPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="sm"`, "h-8 gap-1.5 rounded-md px-2.5 has-[&gt;svg]:px-2.5"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupButtonIconXsPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "icon-xs", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="icon-xs"`, "size-6 rounded-[calc(var(--radius)-5px)] p-0 has-[&gt;svg]:p-0"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestInputGroupButtonIconSmPinned(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "icon-sm", gsx.Raw("x"), nil))
	for _, want := range []string{`data-size="icon-sm"`, "size-8 p-0 has-[&gt;svg]:p-0"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

// TestInputGroupButtonVariantOverride proves variant is forwarded to
// Button's own variant param (default "ghost" per shadcn's own passthrough).
func TestInputGroupButtonVariantOverride(t *testing.T) {
	got := render(t, ui.InputGroupButton("outline", "", gsx.Raw("x"), nil))
	if !strings.Contains(got, `data-variant="outline"`) {
		t.Errorf("missing data-variant=outline override\nin: %s", got)
	}
	if !strings.Contains(got, "dark:border-input dark:bg-input/30") {
		t.Errorf("missing outline variant classes\nin: %s", got)
	}
}

func TestInputGroupButtonAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", nil, gsx.Attrs{{Key: "id", Value: "b1"}}))
	if !strings.Contains(got, `id="b1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

func TestInputGroupButtonCallerClassMerges(t *testing.T) {
	got := render(t, ui.InputGroupButton("", "", nil, gsx.Attrs{{Key: "class", Value: "w-full"}}))
	if !strings.Contains(got, "w-full") {
		t.Errorf("missing caller class w-full\nin: %s", got)
	}
}

func TestInputGroupTextPinned(t *testing.T) {
	got := render(t, ui.InputGroupText(gsx.Raw("x"), nil))
	want := `<span class="flex items-center gap-2 text-sm text-muted-foreground [&amp;_svg]:pointer-events-none [&amp;_svg:not([class*=&#39;size-&#39;])]:size-4">x</span>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, "data-slot") {
		t.Errorf("InputGroupText must not carry data-slot, matching shadcn's own source\nin: %s", got)
	}
}

func TestInputGroupTextAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupText(nil, gsx.Attrs{{Key: "id", Value: "t1"}}))
	if !strings.Contains(got, `id="t1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestInputGroupInputPinned proves InputGroupInput composes ui.Input
// (Input's own base classes come through) and that data-slot is overridden
// from Input's own "input" to "input-group-control", the attribute
// InputGroup's own has-[[data-slot=input-group-control]...] selectors key
// off. rounded-md/border/shadow-xs/dark:bg-input%2F30/ring-[3px] are all
// overridden by the overlay's rounded-none/border-0/shadow-none/
// dark:bg-transparent/ring-0 via the ordinary tailwind-merge conflict
// mechanism (see ui/input-group.gsx's own comment).
func TestInputGroupInputPinned(t *testing.T) {
	got := render(t, ui.InputGroupInput(nil))
	want := `<input type="text" class="h-9 w-full min-w-0 border-input px-3 py-1 text-base transition-[color,box-shadow] outline-none selection:bg-primary selection:text-primary-foreground file:inline-flex file:h-7 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-50 md:text-sm focus-visible:border-ring focus-visible:ring-ring/50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 flex-1 rounded-none border-0 bg-transparent shadow-none focus-visible:ring-0 dark:bg-transparent" data-slot="input-group-control"/>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
	if strings.Contains(got, `data-slot="input"`) {
		t.Errorf("data-slot should be overridden to input-group-control\nin: %s", got)
	}
}

func TestInputGroupInputAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupInput(gsx.Attrs{{Key: "placeholder", Value: "Search..."}}))
	if !strings.Contains(got, `placeholder="Search..."`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// TestInputGroupTextareaPinned proves InputGroupTextarea composes
// ui.Textarea (value renders as the text child, Textarea's own ADAPT — see
// ui/textarea.gsx) with the same data-slot override and class-merge shape as
// InputGroupInput above.
func TestInputGroupTextareaPinned(t *testing.T) {
	got := render(t, ui.InputGroupTextarea("hi", nil))
	want := `<textarea class="flex field-sizing-content min-h-16 w-full border-input px-3 text-base transition-[color,box-shadow] outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-destructive/20 md:text-sm dark:aria-invalid:ring-destructive/40 flex-1 resize-none rounded-none border-0 bg-transparent py-3 shadow-none focus-visible:ring-0 dark:bg-transparent" data-slot="input-group-control">hi</textarea>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestInputGroupTextareaAttrsFallThrough(t *testing.T) {
	got := render(t, ui.InputGroupTextarea("", gsx.Attrs{{Key: "id", Value: "ta1"}}))
	if !strings.Contains(got, `id="ta1"`) {
		t.Errorf("missing fallthrough attr\nin: %s", got)
	}
}

// Realistic composition: a search-style InputGroup (leading icon addon,
// InputGroupInput, trailing InputGroupButton) — the site example's own
// shape.
func TestInputGroupSearchComposition(t *testing.T) {
	got := render(t, ui.InputGroup(
		gsx.Fragment(
			ui.InputGroupAddon("", gsx.Raw("<svg/>"), nil),
			ui.InputGroupInput(gsx.Attrs{{Key: "placeholder", Value: "Search..."}}),
			ui.InputGroupAddon("inline-end", ui.InputGroupButton("", "icon-xs", gsx.Raw("<svg/>"), gsx.Attrs{{Key: "aria-label", Value: "Send"}}), nil),
		),
		nil,
	))
	for _, want := range []string{
		`data-slot="input-group"`,
		`data-align="inline-start"`,
		`data-slot="input-group-control"`,
		`placeholder="Search..."`,
		`data-align="inline-end"`,
		`data-slot="button"`,
		`data-size="icon-xs"`,
		`aria-label="Send"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}
