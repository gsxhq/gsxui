package main

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

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
			<div class="border border-destructive ring-3 ring-destructive/20" data-style-contract-reference="otp-invalid-light"></div>
			<div class="border border-destructive ring-3 ring-destructive/40" data-style-contract-reference="otp-invalid-dark"></div>
			<div class="bg-input/50" data-style-contract-reference="dark-outline-hover"></div>
			<div class="bg-accent/50" data-style-contract-reference="dark-ghost-hover"></div>
			<div class="bg-overlay" data-style-contract-reference="overlay"></div>
		</div>
		<div class="flex gap-2">
			<ui.PaginationPrevious href="#" data-style-contract="pagination-previous"/>
			<ui.PaginationNext href="#" data-style-contract="pagination-next"/>
			<ui.PaginationPrevious href="#" class="pl-12" data-style-contract="pagination-previous-caller"/>
		</div>
		<div>
			<ui.InputOTPGroup>
				<ui.InputOTPSlot data-active="true" aria-invalid="true" data-style-contract="otp-active-invalid"/>
			</ui.InputOTPGroup>
		</div>
		<div>
			<ui.ToggleGroup groupType="multiple" variant="outline" size="sm" spacing="0" aria-label="Joined formatting">
				<ui.ToggleGroupItem groupType="multiple" variant="outline" size="sm" spacing="0" value="text" data-style-contract="toggle-group-sm-first">Text</ui.ToggleGroupItem>
				<ui.ToggleGroupItem groupType="multiple" variant="outline" size="sm" spacing="0" value="icon" data-style-contract="toggle-group-sm-icon" aria-label="Icon"><icon.Bold/></ui.ToggleGroupItem>
				<ui.ToggleGroupItem groupType="multiple" variant="outline" size="sm" spacing="0" value="last" data-style-contract="toggle-group-sm-last">Last</ui.ToggleGroupItem>
			</ui.ToggleGroup>
		</div>
		<div class="flex gap-2">
			<ui.ToggleGroup groupType="multiple" size="default" spacing="2">
				<ui.ToggleGroupItem groupType="multiple" size="default" spacing="2" value="default" data-style-contract="toggle-group-default">Default</ui.ToggleGroupItem>
			</ui.ToggleGroup>
			<ui.ToggleGroup groupType="multiple" size="lg" spacing="2">
				<ui.ToggleGroupItem groupType="multiple" size="lg" spacing="2" value="large" data-style-contract="toggle-group-large">Large</ui.ToggleGroupItem>
			</ui.ToggleGroup>
			<ui.ToggleGroup groupType="multiple" size="sm" spacing="2">
				<ui.ToggleGroupItem groupType="multiple" size="sm" spacing="2" value="caller" class="px-8" data-style-contract="toggle-group-caller">Caller</ui.ToggleGroupItem>
			</ui.ToggleGroup>
		</div>
		<ui.FieldGroup data-variant="outline">
			<ui.FieldSeparator data-style-contract="field-outline-separator"/>
		</ui.FieldGroup>
		<ui.Dialog>
			<ui.DialogTrigger>Open contract dialog</ui.DialogTrigger>
			<ui.DialogContent class="rounded-none" data-style-contract="dialog-caller">
				<ui.DialogTitle>Contract dialog title</ui.DialogTitle>
				<ui.DialogDescription>Contract dialog description</ui.DialogDescription>
			</ui.DialogContent>
		</ui.Dialog>
		<ui.Drawer>
			<ui.DrawerTrigger>Open contract drawer</ui.DrawerTrigger>
			<ui.DrawerContent data-style-contract="drawer-bottom">
				<ui.DrawerHeader data-style-contract="drawer-header">
					<ui.DrawerTitle>Contract drawer title</ui.DrawerTitle>
					<ui.DrawerDescription>Contract drawer description</ui.DrawerDescription>
				</ui.DrawerHeader>
			</ui.DrawerContent>
		</ui.Drawer>
		<ui.Tooltip>
			<ui.TooltipTrigger>Show contract tooltip</ui.TooltipTrigger>
			<ui.TooltipContent data-style-contract="tooltip-kbd">Press <ui.Kbd>K</ui.Kbd></ui.TooltipContent>
		</ui.Tooltip>
		<ui.AccordionItem name="contract-accordion">
			<ui.AccordionTrigger>Contract accordion</ui.AccordionTrigger>
			<ui.AccordionContent>Contract accordion body</ui.AccordionContent>
		</ui.AccordionItem>
	</div>
}
