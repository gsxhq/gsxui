package ui

import "github.com/gsxhq/gsx"

// Table and its parts are the shadcn/ui Table compound set. Parts are plain
// sibling components — compose them in markup; no shared state, no context.
//
// Table renders a scroll-container <div data-gsxui-slot-table-container>
// wrapping the <table data-gsxui-slot-table>; attrs land on the <table>, not the
// container (see docs/jsx-parity.md).

component Table(children gsx.Node, attrs gsx.Attrs) {
	<div data-gsxui-slot-table-container>
		<table { attrs... } data-gsxui-slot-table>{ children }</table>
	</div>
}

component TableHeader(children gsx.Node, attrs gsx.Attrs) {
	<thead { attrs... } data-gsxui-slot-table-header>{ children }</thead>
}

component TableBody(children gsx.Node, attrs gsx.Attrs) {
	<tbody { attrs... } data-gsxui-slot-table-body>{ children }</tbody>
}

component TableFooter(children gsx.Node, attrs gsx.Attrs) {
	<tfoot { attrs... } data-gsxui-slot-table-footer>
		{ children }
	</tfoot>
}

component TableRow(children gsx.Node, attrs gsx.Attrs) {
	<tr { attrs... } data-gsxui-slot-table-row>
		{ children }
	</tr>
}

component TableHead(children gsx.Node, attrs gsx.Attrs) {
	<th { attrs... } data-gsxui-slot-table-head>
		{ children }
	</th>
}

component TableCell(children gsx.Node, attrs gsx.Attrs) {
	<td { attrs... } data-gsxui-slot-table-cell>
		{ children }
	</td>
}

component TableCaption(children gsx.Node, attrs gsx.Attrs) {
	<caption { attrs... } data-gsxui-slot-table-caption>{ children }</caption>
}
