package main

import "github.com/gsxhq/gsxui/ui"

component StyleContractFixture() {
	<div class="grid gap-4 p-6">
		<div><ui.Button>Default</ui.Button></div>
		<div><ui.Button class="h-12 rounded-none">Caller override</ui.Button></div>
		<div><ui.Input value="Ada"/></div>
		<div><ui.Card><ui.CardContent>Card</ui.CardContent></ui.Card></div>
		<div><ui.Badge>Badge</ui.Badge></div>
	</div>
}
