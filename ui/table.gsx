package ui

import "github.com/gsxhq/gsx"

// Table and its parts are the shadcn/ui Table compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.
//
// Table renders a scroll-container <div data-gsxui-slot="table-container">
// wrapping the <table data-gsxui-slot="table">; attrs land on the <table>, not the
// container (see docs/jsx-parity.md).

component Table(children gsx.Node, attrs gsx.Attrs) {
	<div { withSlot("table-container", nil)... }>
		<table { withSlot("table", attrs)... }>{ children }</table>
	</div>
}

component TableHeader(children gsx.Node, attrs gsx.Attrs) {
	<thead { withSlot("table-header", attrs)... }>{ children }</thead>
}

component TableBody(children gsx.Node, attrs gsx.Attrs) {
	<tbody { withSlot("table-body", attrs)... }>{ children }</tbody>
}

component TableFooter(children gsx.Node, attrs gsx.Attrs) {
	<tfoot { withSlot("table-footer", attrs)... }>
		{ children }
	</tfoot>
}

component TableRow(children gsx.Node, attrs gsx.Attrs) {
	<tr { withSlot("table-row", attrs)... }>
		{ children }
	</tr>
}

component TableHead(children gsx.Node, attrs gsx.Attrs) {
	<th { withSlot("table-head", attrs)... }>
		{ children }
	</th>
}

component TableCell(children gsx.Node, attrs gsx.Attrs) {
	<td { withSlot("table-cell", attrs)... }>
		{ children }
	</td>
}

component TableCaption(children gsx.Node, attrs gsx.Attrs) {
	<caption { withSlot("table-caption", attrs)... }>{ children }</caption>
}
