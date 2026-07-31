package ui_test

import (
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui"
)

func TestTableParts(t *testing.T) {
	cases := []struct {
		name string
		node gsx.Node
		want []string
	}{
		{"Table", ui.Table(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-container`, `data-gsxui-slot-table`}},
		{"TableHeader", ui.TableHeader(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-header`}},
		{"TableBody", ui.TableBody(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-body`}},
		{"TableFooter", ui.TableFooter(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-footer`}},
		{"TableRow", ui.TableRow(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-row`}},
		{"TableHead", ui.TableHead(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-head`}},
		{"TableCell", ui.TableCell(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-cell`}},
		{"TableCaption", ui.TableCaption(gsx.Raw("x"), nil), []string{`data-gsxui-slot-table-caption`}},
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

func TestTablePinned(t *testing.T) {
	// Exact full-render pin, verified token-by-token against shadcn's Table
	// (registry/new-york-v4/ui/table.tsx) and docs/jsx-parity.md — a straight
	// port, no divergences. Covers both the container div and the table.
	got := render(t, ui.Table(gsx.Raw("Content"), nil))
	want := `<div class="relative w-full overflow-x-auto" data-gsxui-slot-table-container><table class="w-full caption-bottom text-sm" data-gsxui-slot-table>Content</table></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTableComposition(t *testing.T) {
	// 2x2 table through all parts: header row + body row, two columns each.
	got := render(t, ui.Table(
		gsx.Fragment(
			ui.TableHeader(
				ui.TableRow(
					gsx.Fragment(
						ui.TableHead(gsx.Raw("Name"), nil),
						ui.TableHead(gsx.Raw("Age"), nil),
					),
					nil,
				),
				nil,
			),
			ui.TableBody(
				ui.TableRow(
					gsx.Fragment(
						ui.TableCell(gsx.Raw("Alice"), nil),
						ui.TableCell(gsx.Raw("30"), nil),
					),
					nil,
				),
				nil,
			),
		),
		nil,
	))
	for _, want := range []string{
		`data-gsxui-slot-table-container`,
		`data-gsxui-slot-table`,
		`data-gsxui-slot-table-header`,
		`data-gsxui-slot-table-body`,
		`data-gsxui-slot-table-row`,
		`data-gsxui-slot-table-head`,
		`data-gsxui-slot-table-cell`,
		">Name<", ">Age<", ">Alice<", ">30<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTableCaption(t *testing.T) {
	got := render(t, ui.Table(
		gsx.Fragment(
			ui.TableCaption(gsx.Raw("A list of results."), nil),
			ui.TableBody(ui.TableRow(ui.TableCell(gsx.Raw("x"), nil), nil), nil),
		),
		nil,
	))
	for _, want := range []string{`data-gsxui-slot-table-caption`, ">A list of results.<"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTableCallerClassIsForwardedOnce(t *testing.T) {
	got := render(t, ui.Table(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "caption-top"}}))
	if !strings.Contains(got, "caption-top") || strings.Count(got, "class=") != 2 {
		t.Errorf("caller class must merge onto the table exactly once\nin: %s", got)
	}
}

func TestTableAttrsFallthrough(t *testing.T) {
	// attrs land on the <table>, not the table-container div.
	got := render(t, ui.Table(gsx.Raw("x"), gsx.Attrs{
		{Key: "id", Value: "results"},
		{Key: "aria-describedby", Value: "results-caption"},
	}))
	for _, want := range []string{`id="results"`, `aria-describedby="results-caption"`} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
	if strings.Contains(got, `id="results"`) {
		// sanity: id must appear inside <table ...> not on the outer <div ...>
		div := got[:strings.Index(got, "<table")]
		if strings.Contains(div, `id="results"`) {
			t.Errorf("attrs must land on <table>, not the table-container div\nin: %s", got)
		}
	}
}

func TestTableRowAttrsFallthrough(t *testing.T) {
	got := render(t, ui.TableRow(gsx.Raw("x"), gsx.Attrs{{Key: "data-state", Value: "selected"}}))
	if !strings.Contains(got, `data-state="selected"`) {
		t.Errorf("missing data-state attr\nin: %s", got)
	}
}
