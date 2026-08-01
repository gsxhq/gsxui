package ui

import "github.com/gsxhq/gsx"

// Table and its parts are the shadcn/ui Table compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.
//
// Table renders a scroll-container <div data-gsxui-slot-table-container>
// wrapping the <table data-gsxui-slot-table>; attrs land on the <table>, not the
// container (see docs/jsx-parity.md).

component Table(children gsx.Node, attrs gsx.Attrs) {
	<div class={ "relative w-full overflow-x-auto" } data-gsxui-slot-table-container>
		<table class={ "w-full caption-bottom text-sm" } { attrs... } data-gsxui-slot-table>{ children }</table>
	</div>
}

component TableHeader(children gsx.Node, attrs gsx.Attrs) {
	<thead { attrs... } data-gsxui-slot-table-header>{ children }</thead>
}

component TableBody(children gsx.Node, attrs gsx.Attrs) {
	<tbody { attrs... } data-gsxui-slot-table-body>{ children }</tbody>
}

component TableFooter(children gsx.Node, attrs gsx.Attrs) {
	<tfoot class={ "border-t bg-muted/50 font-medium" } { attrs... } data-gsxui-slot-table-footer>
		{ children }
	</tfoot>
}

component TableRow(children gsx.Node, attrs gsx.Attrs) {
	<tr
		class={
			"border-b transition-colors hover:bg-muted/50 has-[[aria-expanded=true]]:bg-muted/50 data-[state=selected]:bg-muted [[data-gsxui-slot-table-body]_&]:last:border-b-0 [[data-gsxui-slot-table-footer]_&]:last:border-b-0"
		}
		{ attrs... }
		data-gsxui-slot-table-row
	>
		{ children }
	</tr>
}

component TableHead(children gsx.Node, attrs gsx.Attrs) {
	<th
		class={
			"h-10 px-2 text-start align-middle font-medium whitespace-nowrap text-foreground has-[[role=checkbox]]:pe-0 [&>[role=checkbox]]:translate-y-[2px]"
		}
		{ attrs... }
		data-gsxui-slot-table-head
	>
		{ children }
	</th>
}

component TableCell(children gsx.Node, attrs gsx.Attrs) {
	<td
		class={ "p-2 align-middle whitespace-nowrap has-[[role=checkbox]]:pe-0 [&>[role=checkbox]]:translate-y-[2px]" }
		{ attrs... }
		data-gsxui-slot-table-cell
	>
		{ children }
	</td>
}

component TableCaption(children gsx.Node, attrs gsx.Attrs) {
	<caption class={ "mt-4 text-sm text-muted-foreground" } { attrs... } data-gsxui-slot-table-caption>
		{ children }
	</caption>
}
