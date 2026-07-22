package table_test

import (
	"context"
	"strings"
	"testing"

	gsx "github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/table"
)

func render(t *testing.T, n gsx.Node) string {
	t.Helper()
	var sb strings.Builder
	if err := n.Render(context.Background(), &sb); err != nil {
		t.Fatal(err)
	}
	return sb.String()
}

func TestTableParts(t *testing.T) {
	cases := []struct {
		name string
		node gsx.Node
		want []string
	}{
		{"Table", table.Table(gsx.Raw("x"), nil), []string{`data-slot="table-container"`, `data-slot="table"`, "relative w-full overflow-x-auto", "w-full caption-bottom text-sm"}},
		{"TableHeader", table.TableHeader(gsx.Raw("x"), nil), []string{`data-slot="table-header"`, "[&amp;_tr]:border-b"}},
		{"TableBody", table.TableBody(gsx.Raw("x"), nil), []string{`data-slot="table-body"`, "[&amp;_tr:last-child]:border-0"}},
		{"TableFooter", table.TableFooter(gsx.Raw("x"), nil), []string{`data-slot="table-footer"`, "border-t bg-muted/50 font-medium"}},
		{"TableRow", table.TableRow(gsx.Raw("x"), nil), []string{`data-slot="table-row"`, "border-b transition-colors hover:bg-muted/50"}},
		{"TableHead", table.TableHead(gsx.Raw("x"), nil), []string{`data-slot="table-head"`, "h-10 px-2 text-left align-middle font-medium"}},
		{"TableCell", table.TableCell(gsx.Raw("x"), nil), []string{`data-slot="table-cell"`, "p-2 align-middle whitespace-nowrap"}},
		{"TableCaption", table.TableCaption(gsx.Raw("x"), nil), []string{`data-slot="table-caption"`, "mt-4 text-sm text-muted-foreground"}},
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
	got := render(t, table.Table(gsx.Raw("Content"), nil))
	want := `<div data-slot="table-container" class="relative w-full overflow-x-auto"><table data-slot="table" class="w-full caption-bottom text-sm">Content</table></div>`
	if got != want {
		t.Errorf("pinned render mismatch\n got: %s\nwant: %s", got, want)
	}
}

func TestTableComposition(t *testing.T) {
	// 2x2 table through all parts: header row + body row, two columns each.
	got := render(t, table.Table(
		gsx.Fragment(
			table.TableHeader(
				table.TableRow(
					gsx.Fragment(
						table.TableHead(gsx.Raw("Name"), nil),
						table.TableHead(gsx.Raw("Age"), nil),
					),
					nil,
				),
				nil,
			),
			table.TableBody(
				table.TableRow(
					gsx.Fragment(
						table.TableCell(gsx.Raw("Alice"), nil),
						table.TableCell(gsx.Raw("30"), nil),
					),
					nil,
				),
				nil,
			),
		),
		nil,
	))
	for _, want := range []string{
		`data-slot="table-container"`,
		`data-slot="table"`,
		`data-slot="table-header"`,
		`data-slot="table-body"`,
		`data-slot="table-row"`,
		`data-slot="table-head"`,
		`data-slot="table-cell"`,
		">Name<", ">Age<", ">Alice<", ">30<",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTableCaption(t *testing.T) {
	got := render(t, table.Table(
		gsx.Fragment(
			table.TableCaption(gsx.Raw("A list of results."), nil),
			table.TableBody(table.TableRow(table.TableCell(gsx.Raw("x"), nil), nil), nil),
		),
		nil,
	))
	for _, want := range []string{`data-slot="table-caption"`, ">A list of results.<"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q\nin: %s", want, got)
		}
	}
}

func TestTableClassMerge(t *testing.T) {
	// Caller class must merge, and a conflicting utility must drop the base
	// token (w-full vs a caller override).
	got := render(t, table.Table(gsx.Raw("x"), gsx.Attrs{{Key: "class", Value: "caption-top"}}))
	if !strings.Contains(got, "caption-top") {
		t.Errorf("caller class must merge\nin: %s", got)
	}
	if strings.Contains(got, "caption-bottom") {
		t.Errorf("caller caption-top must drop default caption-bottom\nin: %s", got)
	}
}

func TestTableAttrsFallthrough(t *testing.T) {
	// attrs land on the <table>, not the table-container div.
	got := render(t, table.Table(gsx.Raw("x"), gsx.Attrs{
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
	got := render(t, table.TableRow(gsx.Raw("x"), gsx.Attrs{{Key: "data-state", Value: "selected"}}))
	if !strings.Contains(got, `data-state="selected"`) {
		t.Errorf("missing data-state attr\nin: %s", got)
	}
}
