package main

import (
	"github.com/gsxhq/gsxui/ui"
	"github.com/gsxhq/gsxui/ui/icon"
)

component DirectionalDrawerFixture(name string, direction string) {
	<ui.Drawer>
		<ui.DrawerTrigger>{ "Open drawer " + name }</ui.DrawerTrigger>
		<ui.DrawerContent direction={direction} data-style-contract={"drawer-" + name}>
			<ui.DrawerHeader>
				<ui.DrawerTitle>{ "Drawer " + name }</ui.DrawerTitle>
				<ui.DrawerDescription>{ "Drawer " + name + " description" }</ui.DrawerDescription>
			</ui.DrawerHeader>
			<ui.DrawerClose>{ "Close drawer " + name }</ui.DrawerClose>
		</ui.DrawerContent>
	</ui.Drawer>
}

component DirectionalSheetFixture(name string, side string) {
	<ui.Sheet>
		<ui.SheetTrigger>{ "Open sheet " + name }</ui.SheetTrigger>
		<ui.SheetContent side={side} hideCloseButton={true} data-style-contract={"sheet-" + name}>
			<ui.SheetHeader>
				<ui.SheetTitle>{ "Sheet " + name }</ui.SheetTitle>
				<ui.SheetDescription>{ "Sheet " + name + " description" }</ui.SheetDescription>
			</ui.SheetHeader>
			<ui.SheetClose>{ "Close sheet " + name }</ui.SheetClose>
		</ui.SheetContent>
	</ui.Sheet>
}

component StyleContractFixture() {
	<div class="grid gap-4 p-6">
		<div>
			<ui.Button>Default</ui.Button>
		</div>
		<div>
			<ui.Button class="h-12 rounded-none">Caller override</ui.Button>
		</div>
		<div>
			<ui.Input value="Ada"/>
		</div>
		<div>
			<ui.Card>
				<ui.CardContent>Card</ui.CardContent>
			</ui.Card>
		</div>
		<div>
			<ui.CardHeader class="border-b pb-10" data-style-contract="card-header-caller">
				<ui.CardTitle>Card header caller override</ui.CardTitle>
			</ui.CardHeader>
		</div>
		<div data-disabled="true">
			<ui.Label class="opacity-100" data-style-contract="label-disabled-caller">Disabled label caller override</ui.Label>
		</div>
		<div>
			<ui.Badge>Badge</ui.Badge>
		</div>
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
			<div
				class="border border-destructive ring-3 ring-destructive/20"
				data-style-contract-reference="otp-invalid-light"
			></div>
			<div
				class="border border-destructive ring-3 ring-destructive/40"
				data-style-contract-reference="otp-invalid-dark"
			></div>
			<div class="bg-input/50" data-style-contract-reference="dark-outline-hover"></div>
			<div class="bg-accent/50" data-style-contract-reference="dark-ghost-hover"></div>
			<div class="bg-overlay" data-style-contract-reference="overlay"></div>
		</div>
		<div class="flex gap-2">
			<ui.PaginationPrevious href="#" data-style-contract="pagination-previous"/>
			<ui.PaginationNext href="#" data-style-contract="pagination-next"/>
			<ui.PaginationPrevious href="#" class="pl-12" data-style-contract="pagination-previous-caller"/>
		</div>
		<div class="flex gap-2">
			<ui.Input class="text-lg" data-style-contract="input-caller-text-size"/>
			<ui.Textarea class="text-lg" data-style-contract="textarea-caller-text-size"/>
		</div>
		<ui.AlertDialog>
			<ui.AlertDialogContent class="max-w-md" data-style-contract="alert-dialog-content-caller">
				<ui.AlertDialogTitle>Caller max-w override</ui.AlertDialogTitle>
			</ui.AlertDialogContent>
		</ui.AlertDialog>
		<ui.ButtonGroup data-style-contract="button-group-tail">
			<ui.Button data-style-contract="button-group-earlier">Earlier</ui.Button>
			<ui.Button data-style-contract="button-group-visible-tail">Visible tail</ui.Button>
			<select aria-hidden="true">
				<option>Hidden native tail</option>
			</select>
		</ui.ButtonGroup>
		<div>
			<ui.InputOTPGroup>
				<ui.InputOTPSlot data-active="true" aria-invalid="true" data-style-contract="otp-active-invalid"/>
			</ui.InputOTPGroup>
		</div>
		<div>
			<ui.ToggleGroup groupType="multiple" variant="outline" size="sm" spacing="0" aria-label="Joined formatting">
				<ui.ToggleGroupItem
					groupType="multiple"
					variant="outline"
					size="sm"
					spacing="0"
					value="text"
					data-style-contract="toggle-group-sm-first"
				>
					Text
				</ui.ToggleGroupItem>
				<ui.ToggleGroupItem
					groupType="multiple"
					variant="outline"
					size="sm"
					spacing="0"
					value="icon"
					data-style-contract="toggle-group-sm-icon"
					aria-label="Icon"
				>
					<icon.Bold/>
				</ui.ToggleGroupItem>
				<ui.ToggleGroupItem
					groupType="multiple"
					variant="outline"
					size="sm"
					spacing="0"
					value="last"
					data-style-contract="toggle-group-sm-last"
				>
					Last
				</ui.ToggleGroupItem>
			</ui.ToggleGroup>
		</div>
		<div class="flex gap-2">
			<ui.ToggleGroup groupType="multiple" size="default" spacing="2">
				<ui.ToggleGroupItem
					groupType="multiple"
					size="default"
					spacing="2"
					value="default"
					data-style-contract="toggle-group-default"
				>
					Default
				</ui.ToggleGroupItem>
			</ui.ToggleGroup>
			<ui.ToggleGroup groupType="multiple" size="lg" spacing="2">
				<ui.ToggleGroupItem
					groupType="multiple"
					size="lg"
					spacing="2"
					value="large"
					data-style-contract="toggle-group-large"
				>
					Large
				</ui.ToggleGroupItem>
			</ui.ToggleGroup>
			<ui.ToggleGroup groupType="multiple" size="sm" spacing="2">
				<ui.ToggleGroupItem
					groupType="multiple"
					size="sm"
					spacing="2"
					value="caller"
					class="px-8"
					data-style-contract="toggle-group-caller"
				>
					Caller
				</ui.ToggleGroupItem>
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
		<div class="flex flex-wrap gap-2">
			<DirectionalDrawerFixture name="default"/>
			<DirectionalDrawerFixture name="top" direction="top"/>
			<DirectionalDrawerFixture name="left" direction="left"/>
			<DirectionalDrawerFixture name="right" direction="right"/>
			<DirectionalSheetFixture name="default"/>
			<DirectionalSheetFixture name="left" side="left"/>
			<DirectionalSheetFixture name="top" side="top"/>
			<DirectionalSheetFixture name="bottom" side="bottom"/>
		</div>
		<ui.Tooltip>
			<ui.TooltipTrigger>Show contract tooltip</ui.TooltipTrigger>
			<ui.TooltipContent data-style-contract="tooltip-kbd">
				Press <ui.Kbd>K</ui.Kbd>
			</ui.TooltipContent>
		</ui.Tooltip>
		<ui.AccordionItem name="contract-accordion">
			<ui.AccordionTrigger>Contract accordion</ui.AccordionTrigger>
			<ui.AccordionContent id="accordion-caller-content" class="pb-8" data-style-contract="accordion-caller-content">
				Contract accordion body
			</ui.AccordionContent>
		</ui.AccordionItem>
		<ui.CommandDialog
			title="Contract command"
			description="Command dialog style contract"
			data-style-contract="command-dialog"
		>
			<ui.CommandInput placeholder="Search contract commands"/>
			<ui.CommandList>
				<ui.CommandGroup heading="Suggestions">
					<ui.CommandItem value="alpha">Alpha</ui.CommandItem>
				</ui.CommandGroup>
			</ui.CommandList>
		</ui.CommandDialog>
		<div>
			<ui.DropdownMenuItem variant="destructive" tabindex="0" data-style-contract="dropdown-destructive">
				Delete dropdown item
			</ui.DropdownMenuItem>
			<ui.ContextMenuItem variant="destructive" tabindex="0" data-style-contract="context-destructive">
				Delete context item
			</ui.ContextMenuItem>
			<ui.MenubarItem variant="destructive" tabindex="0" data-style-contract="menubar-destructive">
				Delete menubar item
			</ui.MenubarItem>
			<div class="bg-destructive/10 text-destructive" data-style-contract-reference="menu-destructive-focus"></div>
		</div>
		<div>
			{/* §10b caller-override pins: data-inset/data-disabled are stamped by
			    the CALLER through attrs (not by the component), and their
			    pl-8/opacity-50 rules deliberately stayed in
			    assets/css/styles/default/menu.css's own @layer components block
			    — not migrated onto the recipe — so a caller's plain utility still
			    wins. pointer-events-none for [data-disabled] WAS migrated (baked
			    onto the recipe on purpose), so only opacity is exercised here. */}
			<ui.DropdownMenuItem data-inset="true" data-disabled="true" class="pl-2 opacity-100" tabindex="0" data-style-contract="dropdown-menu-item-caller">
				Inset + disabled dropdown item, caller-overridden
			</ui.DropdownMenuItem>
			<ui.ContextMenuItem data-inset="true" data-disabled="true" class="pl-2 opacity-100" tabindex="0" data-style-contract="context-menu-item-caller">
				Inset + disabled context item, caller-overridden
			</ui.ContextMenuItem>
			<ui.MenubarItem data-inset="true" data-disabled="true" class="pl-2 opacity-100" tabindex="0" data-style-contract="menubar-item-caller">
				Inset + disabled menubar item, caller-overridden
			</ui.MenubarItem>
			{/* Same §10b shape for Command/Combobox: data-disabled is stamped by
			    the CALLER through attrs on both (CommandItem/ComboboxItem's own
			    doc comments), so their opacity-50 rules stayed behind in
			    assets/css/styles/default/{command,combobox}.css rather than
			    moving onto the recipe. pointer-events-none WAS migrated for
			    both, on purpose. */}
			<ui.CommandItem data-disabled="true" class="opacity-100" data-style-contract="command-item-caller">
				Disabled command item, caller-overridden
			</ui.CommandItem>
		</div>
		<div>
			<ui.DropdownMenuCheckboxItem checked={true} value="dropdown-on" data-style-contract="dropdown-checked">
				Dropdown on
			</ui.DropdownMenuCheckboxItem>
			<ui.DropdownMenuCheckboxItem checked={false} value="dropdown-off" data-style-contract="dropdown-unchecked">
				Dropdown off
			</ui.DropdownMenuCheckboxItem>
			<ui.ContextMenuRadioItem checked={true} value="context-on" data-style-contract="context-checked">
				Context on
			</ui.ContextMenuRadioItem>
			<ui.ContextMenuRadioItem checked={false} value="context-off" data-style-contract="context-unchecked">
				Context off
			</ui.ContextMenuRadioItem>
			<ui.MenubarCheckboxItem checked={true} value="menubar-on" data-style-contract="menubar-checked">
				Menubar on
			</ui.MenubarCheckboxItem>
			<ui.MenubarCheckboxItem checked={false} value="menubar-off" data-style-contract="menubar-unchecked">
				Menubar off
			</ui.MenubarCheckboxItem>
			<ui.ComboboxItem value="combobox-on" selected={true} data-style-contract="combobox-checked">
				Combobox on
			</ui.ComboboxItem>
			<ui.ComboboxItem value="combobox-off" selected={false} data-style-contract="combobox-unchecked">
				Combobox off
			</ui.ComboboxItem>
		</div>
		<div>
			<ui.DropdownMenuSubContent data-side="right" data-style-contract="submenu-right">
				Right submenu
			</ui.DropdownMenuSubContent>
		</div>
		<ui.NavigationMenu data-style-contract="navigation-root">
			<ui.NavigationMenuList>
				<ui.NavigationMenuItem>
					<ui.NavigationMenuTrigger data-style-contract="navigation-trigger">
						Contract products
					</ui.NavigationMenuTrigger>
					<ui.NavigationMenuContent data-style-contract="navigation-content">
						<ui.NavigationMenuLink href="#">Contract product</ui.NavigationMenuLink>
					</ui.NavigationMenuContent>
				</ui.NavigationMenuItem>
				<ui.NavigationMenuIndicator data-style-contract="navigation-indicator"/>
			</ui.NavigationMenuList>
		</ui.NavigationMenu>
	</div>
}
