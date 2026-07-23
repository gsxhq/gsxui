package ui

import "github.com/gsxhq/gsx"

// Table and its parts are the shadcn/ui Table compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.
//
// Table renders a scroll-container <div data-slot="table-container"> wrapping
// the <table data-slot="table">; attrs land on the <table>, not the
// container (see docs/jsx-parity.md).

component Table(children gsx.Node, attrs gsx.Attrs) {
	<div data-slot="table-container" class="relative w-full overflow-x-auto">
		<table data-slot="table" class="w-full caption-bottom text-sm" { attrs... }>{ children }</table>
	</div>
}

component TableHeader(children gsx.Node, attrs gsx.Attrs) {
	<thead data-slot="table-header" class="[&_tr]:border-b" { attrs... }>{ children }</thead>
}

component TableBody(children gsx.Node, attrs gsx.Attrs) {
	<tbody data-slot="table-body" class="[&_tr:last-child]:border-0" { attrs... }>{ children }</tbody>
}

component TableFooter(children gsx.Node, attrs gsx.Attrs) {
	<tfoot data-slot="table-footer" class="border-t bg-muted/50 font-medium [&>tr]:last:border-b-0" { attrs... }>
		{ children }
	</tfoot>
}

component TableRow(children gsx.Node, attrs gsx.Attrs) {
	<tr
		data-slot="table-row"
		class="border-b transition-colors hover:bg-muted/50 has-aria-expanded:bg-muted/50 data-[state=selected]:bg-muted"
		{ attrs... }
	>
		{ children }
	</tr>
}

component TableHead(children gsx.Node, attrs gsx.Attrs) {
	<th
		data-slot="table-head"
		class="h-10 px-2 text-left align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]"
		{ attrs... }
	>
		{ children }
	</th>
}

component TableCell(children gsx.Node, attrs gsx.Attrs) {
	<td
		data-slot="table-cell"
		class="p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pr-0 [&>[role=checkbox]]:translate-y-[2px]"
		{ attrs... }
	>
		{ children }
	</td>
}

component TableCaption(children gsx.Node, attrs gsx.Attrs) {
	<caption data-slot="table-caption" class="mt-4 text-sm text-muted-foreground" { attrs... }>{ children }</caption>
}
