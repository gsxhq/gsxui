package pages

type docTOCItem struct {
	ID    string
	Title string
	Depth int
}

component docHeading(item docTOCItem) {
	{ if item.Depth == 3 {
		<h3
			id={item.ID}
			data-doc-heading
			tabindex="-1"
			class="scroll-mt-20 text-lg font-semibold tracking-tight"
		>
			{ item.Title }
		</h3>
	} else {
		<h2
			id={item.ID}
			data-doc-heading
			tabindex="-1"
			class="scroll-mt-20 text-xl font-semibold tracking-tight"
		>
			{ item.Title }
		</h2>
	} }
}

component docTableOfContents(items []docTOCItem) {
	<nav aria-label="On this page" class="flex flex-col gap-3 text-sm">
		<div class="font-medium">On this page</div>
		<div class="flex flex-col gap-2">
			{ for _, item := range items {
				<a
					data-site-toc-link
					href={"#" + item.ID}
					class={
						"text-muted-foreground transition-colors hover:text-foreground data-[active]:font-medium data-[active]:text-foreground",
						"pl-3": item.Depth == 3
					}
				>
					{ item.Title }
				</a>
			} }
		</div>
	</nav>
}
