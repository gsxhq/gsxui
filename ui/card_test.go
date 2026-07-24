package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestCardParts(t *testing.T) {
	cases := []struct {
		name string
		node gsx.Node
		want []string
	}{
		{"Card", ui.Card(gsx.Raw("x"), nil), []string{`data-slot="card"`, "flex flex-col gap-4 rounded-xl border bg-card"}},
		{"CardHeader", ui.CardHeader(gsx.Raw("x"), nil), []string{`data-slot="card-header"`, "@container/card-header"}},
		{"CardTitle", ui.CardTitle(gsx.Raw("x"), nil), []string{`data-slot="card-title"`, "leading-snug font-medium"}},
		{"CardDescription", ui.CardDescription(gsx.Raw("x"), nil), []string{`data-slot="card-description"`, "text-muted-foreground"}},
		{"CardAction", ui.CardAction(gsx.Raw("x"), nil), []string{`data-slot="card-action"`, "col-start-2 row-span-2"}},
		{"CardContent", ui.CardContent(gsx.Raw("x"), nil), []string{`data-slot="card-content"`, "px-4"}},
		{"CardFooter", ui.CardFooter(gsx.Raw("x"), nil), []string{`data-slot="card-footer"`, "flex items-center rounded-b-xl border-t p-4"}},
	}
	for _, tc := range cases {
		got := render(t, tc.node)
		for _, want := range tc.want {
			if !strings.Contains(got, want) {
				t.Errorf("%s: missing %q\nin: %s", tc.name, want, got)
			}
		}
	}
}

func TestCardPinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's Card
	// (registry/new-york-v4/ui/card.tsx) and docs/jsx-parity.md — a straight
	// port, no divergences.
	got := render(t, ui.Card(gsx.Raw("Content"), nil))
	want := `<div data-slot="card" class="flex flex-col gap-4 rounded-xl border bg-card py-4 text-sm text-card-foreground has-data-[slot=card-footer]:pb-0">Content</div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestCardComposition(t *testing.T) {
	got := render(t, ui.Card(
		gsx.Fragment(
			ui.CardHeader(ui.CardTitle(gsx.Raw("Title"), nil), nil),
			ui.CardContent(gsx.Raw("Body"), nil),
		),
		gsx.Attrs{{Key: "class", Value: "py-8"}},
	))
	for _, want := range []string{`data-slot="card-title"`, ">Title<", ">Body<", "py-8"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, "py-4") {
		t.Errorf("caller py-8 must drop default py-4\nin: %s", got)
	}
}
