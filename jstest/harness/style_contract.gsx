package main

import "github.com/gsxhq/gsxui/ui"

component StyleContractFixture() {
	<div class="grid gap-4 p-6">
		<div><ui.Button>Default</ui.Button></div>
		<div><ui.Button class="h-12 rounded-none">Caller override</ui.Button></div>
		<div><ui.Input value="Ada"/></div>
		<div><ui.Card><ui.CardContent>Card</ui.CardContent></ui.Card></div>
		<div><ui.Badge>Badge</ui.Badge></div>
		<div class="flex flex-wrap gap-2">
			<ui.Button variant="destructive" data-style-contract="dark-button-destructive">Destructive</ui.Button>
			<ui.Button aria-invalid="true" data-style-contract="dark-button-invalid">Invalid</ui.Button>
			<ui.Button variant="outline" data-style-contract="dark-button-outline">Outline</ui.Button>
			<ui.Button variant="ghost" data-style-contract="dark-button-ghost">Ghost</ui.Button>
			<ui.Badge variant="destructive" data-style-contract="dark-badge-destructive">Destructive badge</ui.Badge>
			<ui.Badge aria-invalid="true" tabindex="0" data-style-contract="dark-badge-invalid">Invalid badge</ui.Badge>
		</div>
		<div class="hidden">
			<div class="bg-destructive/60" data-style-contract-reference="dark-destructive"></div>
			<div class="bg-destructive/90" data-style-contract-reference="destructive-hover"></div>
			<div class="ring-[3px] ring-destructive/40" data-style-contract-reference="dark-invalid-ring"></div>
			<div class="bg-input/50" data-style-contract-reference="dark-outline-hover"></div>
			<div class="bg-accent/50" data-style-contract-reference="dark-ghost-hover"></div>
		</div>
		<div class="flex gap-2">
			<ui.PaginationPrevious href="#" data-style-contract="pagination-previous"/>
			<ui.PaginationNext href="#" data-style-contract="pagination-next"/>
			<ui.PaginationPrevious href="#" class="pl-12" data-style-contract="pagination-previous-caller"/>
		</div>
	</div>
}
