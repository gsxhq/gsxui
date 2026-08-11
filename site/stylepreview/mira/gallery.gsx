// Package stylepreview is a placeholder clause: internal/stylegen reads this
// file and rewrites the package to each authored style (nova, maia, ...),
// emitting site/stylepreview/<style>/gallery.gsx. Author HERE, never in the
// generated copies — go run ./cmd/stylegen refreshes them and --check fails
// on drift.
//
// The gallery is the theme editor's preview document: one composed showcase
// (in the spirit of ui.shadcn.com/themes) that renders every component in the
// catalogue as realistic product UI. site/stylepreview/gallery_test.go
// asserts every component in registry/generated/recipes.json appears at
// least once, so a future component cannot silently miss the gallery.
//
// The preview document loads web/preview.js only — no component JS — so
// every popup (menus, dialogs, popovers, tooltips) renders in its resting
// closed state. What the gallery showcases is trigger and surface chrome
// theming, not behavior.
package mira

import (
	"time"

	"github.com/gsxhq/gsxui/ui/icon"
)

// galleryMonth is fixed so the preview never depends on the day it renders.
var galleryMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// galleryAvatarSVG is a stand-in portrait tile, the same data-URI idiom as
// site/examples/avatar/basic.gsx — no image asset, no network fetch.
var galleryAvatarSVG = []byte("<svg xmlns='http://www.w3.org/2000/svg' width='64' height='64'><rect width='64' height='64' fill='#6d28d9'/><text x='32' y='34' text-anchor='middle' dominant-baseline='central' font-family='sans-serif' font-weight='600' font-size='26' fill='#fff'>AL</text></svg>")

// Gallery renders the full component showcase. idp prefixes every id and
// form-control name: the theme preview document renders the gallery twice
// (once per style section), so unprefixed ids would collide across sections.
component Gallery(idp string) {
	<div data-theme-preview-gallery class="mx-auto grid w-full max-w-7xl items-start gap-4 p-4 sm:p-6 md:grid-cols-2 xl:grid-cols-3">
		<galleryButtonsCard/>
		<galleryLoginCard idp={idp}/>
		<gallerySettingsCard idp={idp}/>
		<galleryCalendarCard idp={idp}/>
		<galleryMenusCard/>
		<galleryFeedbackCard/>
		<galleryTeamCard/>
		<galleryTableCard/>
		<galleryChartCard/>
		<galleryTabsCard idp={idp}/>
		<galleryNavigationCard/>
		<galleryControlsCard idp={idp}/>
		<galleryOverlaysCard idp={idp}/>
		<galleryEmptyCard/>
		<galleryMediaCard/>
		<gallerySidebarCard idp={idp}/>
	</div>
}

// galleryButtonsCard keeps the old Matrix's coverage: every variant, every
// size, and the states the theme editor keys on (data-theme-preview-case is
// web/theme-preview.js's focus hook and pages_test.go's marker set).
component galleryButtonsCard() {
	<Card>
		<CardHeader>
			<CardTitle>Buttons</CardTitle>
			<CardDescription>Every variant, size, and state.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<div class="flex flex-wrap items-center gap-2">
				<Button data-theme-preview-case="text">Default</Button>
				<Button variant="secondary" data-theme-preview-case="text">Secondary</Button>
				<Button variant="destructive" data-theme-preview-case="text">Destructive</Button>
				<Button variant="outline" data-theme-preview-case="text">Outline</Button>
				<Button variant="ghost" data-theme-preview-case="text">Ghost</Button>
				<Button variant="link" data-theme-preview-case="text">Link</Button>
			</div>
			<div class="flex flex-wrap items-center gap-2">
				<Button size="xs" data-theme-preview-case="text">Extra small</Button>
				<Button size="sm" data-theme-preview-case="text">Small</Button>
				<Button data-theme-preview-case="text">Default</Button>
				<Button size="lg" data-theme-preview-case="text">Large</Button>
				<Button size="icon-xs" aria-label="Icon extra small" data-theme-preview-case="icon"><icon.Plus/></Button>
				<Button size="icon-sm" aria-label="Icon small" data-theme-preview-case="icon"><icon.Plus/></Button>
				<Button size="icon" aria-label="Icon" data-theme-preview-case="icon"><icon.Plus/></Button>
				<Button size="icon-lg" aria-label="Icon large" data-theme-preview-case="icon"><icon.Plus/></Button>
			</div>
			<div class="flex flex-wrap items-center gap-2">
				<Button disabled data-theme-preview-case="disabled">Disabled</Button>
				<Button variant="outline" href="/docs/getting-started" data-theme-preview-case="link">Link</Button>
				<Button autofocus data-theme-preview-case="focus-visible">Focus visible</Button>
				<Button aria-invalid="true" data-theme-preview-case="invalid">Invalid</Button>
				{/* Proves caller utilities beat the recipe (radius and padding here)
				   without reading as a defect: a pill override is a thing a real
				   caller would ship. */}
				<Button
					variant="secondary"
					class="rounded-full px-5"
					data-caller-marker
					data-theme-preview-case="caller-composition"
				>
					Caller override
				</Button>
			</div>
		</CardContent>
	</Card>
}

component galleryLoginCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle>Create an account</CardTitle>
			<CardDescription>Enter your email below to create your account.</CardDescription>
			<CardAction>
				<Badge variant="secondary">Beta</Badge>
			</CardAction>
		</CardHeader>
		<CardContent>
			<form>
				<FieldGroup>
					<Field>
						<FieldLabel for={idp + "-login-email"}>Email</FieldLabel>
						<Input id={idp + "-login-email"} type="email" placeholder="m@example.com" autocomplete="email"/>
						<FieldDescription>We never share your email with anyone.</FieldDescription>
					</Field>
					<Field>
						<FieldLabel for={idp + "-login-password"}>Password</FieldLabel>
						<Input id={idp + "-login-password"} type="password" value="" autocomplete="current-password"/>
					</Field>
					<FieldSeparator>or</FieldSeparator>
					<div class="flex items-center gap-2">
						<Checkbox id={idp + "-login-remember"} checked/>
						<Label for={idp + "-login-remember"}>Remember me for 30 days</Label>
					</div>
				</FieldGroup>
			</form>
		</CardContent>
		<CardFooter class="flex-col gap-2 pb-4">
			<Button class="w-full">Create account</Button>
			<Button variant="outline" class="w-full">Sign in with Google</Button>
		</CardFooter>
	</Card>
}

component gallerySettingsCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle>Workspace settings</CardTitle>
			<CardDescription>Configure notifications and appearance.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-5">
			<div class="flex items-center justify-between gap-4">
				<Label for={idp + "-settings-emails"}>Email notifications</Label>
				<Switch id={idp + "-settings-emails"} checked/>
			</div>
			<div class="flex items-center justify-between gap-4">
				<Label for={idp + "-settings-digest"}>Weekly digest</Label>
				<Switch id={idp + "-settings-digest"}/>
			</div>
			<Separator/>
			<div class="flex flex-col gap-2">
				<Label for={idp + "-settings-density"}>Density</Label>
				<Select name={idp + "-settings-density"}>
					{/* The label's for= must point at a labelable element; Select's
					   own attrs land on a display:contents wrapper, so the id goes
					   on the trigger button. */}
					<SelectTrigger id={idp + "-settings-density"} class="w-full">
						<SelectValue placeholder="Select density"/>
					</SelectTrigger>
					<SelectContent>
						<SelectGroup>
							<SelectLabel>Density</SelectLabel>
							<SelectItem value="comfortable" selected>Comfortable</SelectItem>
							<SelectItem value="compact">Compact</SelectItem>
							<SelectItem value="spacious" disabled>Spacious</SelectItem>
						</SelectGroup>
					</SelectContent>
				</Select>
			</div>
			<div class="flex flex-col gap-2">
				<Label for={idp + "-settings-timezone"}>Timezone</Label>
				<NativeSelect id={idp + "-settings-timezone"} name={idp + "-settings-timezone"}>
					<NativeSelectOption value="utc" selected>UTC</NativeSelectOption>
					<NativeSelectOption value="cet">Central European Time</NativeSelectOption>
					<NativeSelectOption value="pt">Pacific Time</NativeSelectOption>
				</NativeSelect>
			</div>
			<div class="flex flex-col gap-3">
				<Label>Sidebar position</Label>
				<div class="flex items-center gap-2">
					<Radio id={idp + "-settings-left"} name={idp + "-settings-side"} checked/>
					<Label for={idp + "-settings-left"}>Left</Label>
				</div>
				<div class="flex items-center gap-2">
					<Radio id={idp + "-settings-right"} name={idp + "-settings-side"}/>
					<Label for={idp + "-settings-right"}>Right</Label>
				</div>
			</div>
			<div class="flex flex-col gap-3">
				<Label for={idp + "-settings-volume"}>Notification volume</Label>
				<Slider id={idp + "-settings-volume"} value={35} min={0} max={100} step={1} aria-label="Notification volume"/>
			</div>
			<div class="flex flex-col gap-2">
				<Label for={idp + "-settings-note"}>Status message</Label>
				<Textarea id={idp + "-settings-note"} value="" placeholder="Out of office until Monday."/>
			</div>
		</CardContent>
	</Card>
}

component galleryCalendarCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle>Schedule a review</CardTitle>
			<CardDescription>Pick a date for the quarterly review.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col items-start gap-4">
			<Calendar
				mode="single"
				month={galleryMonth}
				selected={[]time.Time{time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)}}
				weekStartsOn={time.Sunday}
				showOutsideDays={true}
				captionLayout="label"
			/>
			<Popover>
				<Button
					variant="outline"
					data-gsxui-slot-popover-trigger
					aria-expanded="false"
					class="w-[220px] justify-start text-left font-normal text-muted-foreground"
				>
					<icon.Calendar/>
					Pick a date
				</Button>
				<PopoverContent class="w-auto p-0">
					<Calendar
						id={idp + "-datepicker-calendar"}
						mode="single"
						month={galleryMonth}
						weekStartsOn={time.Sunday}
						showOutsideDays={true}
						captionLayout="label"
					/>
				</PopoverContent>
			</Popover>
		</CardContent>
	</Card>
}

component galleryMenusCard() {
	<Card>
		<CardHeader>
			<CardTitle>Menus</CardTitle>
			<CardDescription>Menubar, dropdown, and context menu chrome.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<Menubar>
				<MenubarMenu>
					<MenubarTrigger>File</MenubarTrigger>
					<MenubarContent>
						<MenubarItem>
							New Tab
							<MenubarShortcut>⌘T</MenubarShortcut>
						</MenubarItem>
						<MenubarItem>
							New Window
							<MenubarShortcut>⌘N</MenubarShortcut>
						</MenubarItem>
						<MenubarSeparator/>
						<MenubarItem>
							Print...
							<MenubarShortcut>⌘P</MenubarShortcut>
						</MenubarItem>
					</MenubarContent>
				</MenubarMenu>
				<MenubarMenu>
					<MenubarTrigger>Edit</MenubarTrigger>
					<MenubarContent>
						<MenubarItem>Undo</MenubarItem>
						<MenubarItem>Redo</MenubarItem>
					</MenubarContent>
				</MenubarMenu>
			</Menubar>
			<div class="flex flex-wrap items-center gap-2">
				<DropdownMenu>
					<Button
						variant="outline"
						data-gsxui-slot-dropdown-menu-trigger
						aria-haspopup="menu"
						aria-expanded="false"
					>
						Options
					</Button>
					<DropdownMenuContent>
						<DropdownMenuGroup>
							<DropdownMenuLabel>My Account</DropdownMenuLabel>
							<DropdownMenuSeparator/>
							<DropdownMenuItem>Profile</DropdownMenuItem>
							<DropdownMenuItem>Billing</DropdownMenuItem>
							<DropdownMenuItem>
								Settings <DropdownMenuShortcut>⌘,</DropdownMenuShortcut>
							</DropdownMenuItem>
						</DropdownMenuGroup>
					</DropdownMenuContent>
				</DropdownMenu>
			</div>
			<ContextMenu>
				<ContextMenuTrigger class="flex h-24 w-full items-center justify-center rounded-md border border-dashed text-sm">
					Right click here
				</ContextMenuTrigger>
				<ContextMenuContent class="w-52">
					<ContextMenuItem>
						Back
						<ContextMenuShortcut>⌘[</ContextMenuShortcut>
					</ContextMenuItem>
					<ContextMenuItem>
						Reload
						<ContextMenuShortcut>⌘R</ContextMenuShortcut>
					</ContextMenuItem>
					<ContextMenuSeparator/>
					<ContextMenuItem variant="destructive">Delete</ContextMenuItem>
				</ContextMenuContent>
			</ContextMenu>
		</CardContent>
	</Card>
}

component galleryFeedbackCard() {
	<Card>
		<CardHeader>
			<CardTitle>Sync status</CardTitle>
			<CardDescription>Alerts, badges, progress, and toasts.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<Alert>
				<icon.CircleCheck/>
				<AlertTitle>Backup complete</AlertTitle>
				<AlertDescription>Your workspace was backed up 5 minutes ago.</AlertDescription>
			</Alert>
			<Alert variant="destructive">
				<icon.CircleAlert/>
				<AlertTitle>Payment failed</AlertTitle>
				<AlertDescription>Please verify your billing information and try again.</AlertDescription>
			</Alert>
			<div class="flex flex-wrap items-center gap-2">
				<Badge>Active</Badge>
				<Badge variant="secondary">Draft</Badge>
				<Badge variant="destructive">Overdue</Badge>
				<Badge variant="outline">Archived</Badge>
			</div>
			<div class="flex items-center gap-3">
				<Progress value={66} class="flex-1"/>
				<span class="text-xs text-muted-foreground">66%</span>
			</div>
			<div class="flex items-center gap-4">
				<Skeleton class="size-10 rounded-full"/>
				<div class="grid flex-1 gap-2">
					<Skeleton class="h-3 w-3/4"/>
					<Skeleton class="h-3 w-1/2"/>
				</div>
				<Spinner/>
			</div>
			<div class="flex items-center gap-3 text-sm text-muted-foreground">
				Press
				<KbdGroup>
					<Kbd>⌘</Kbd>
					<Kbd>K</Kbd>
				</KbdGroup>
				to open the command palette
			</div>
			<ul class="list-none">
				<Toast
					toastType="success"
					title="Profile updated"
					description="Your changes have been saved."
					action="Undo"
					class="relative w-full"
				/>
			</ul>
		</CardContent>
	</Card>
}

component galleryTeamCard() {
	<Card>
		<CardHeader>
			<CardTitle>Team members</CardTitle>
			<CardDescription>Invite and manage your team.</CardDescription>
			<CardAction>
				<div class="flex -space-x-2">
					<Avatar class="ring-2 ring-background">
						<AvatarImage src={galleryAvatarSVG |> dataURL("image/svg+xml")} alt="Ada Lovelace"/>
						<AvatarFallback>AL</AvatarFallback>
					</Avatar>
					<Avatar class="ring-2 ring-background">
						<AvatarFallback>GH</AvatarFallback>
					</Avatar>
					<Avatar class="ring-2 ring-background">
						<AvatarFallback>AT</AvatarFallback>
					</Avatar>
				</div>
			</CardAction>
		</CardHeader>
		<CardContent>
			<ItemGroup>
				<Item>
					<ItemMedia variant="icon">
						<icon.User/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>Jamie Lee</ItemTitle>
						<ItemDescription>jamie@example.com</ItemDescription>
					</ItemContent>
					<ItemActions>
						<Button variant="ghost" size="sm">Remove</Button>
					</ItemActions>
				</Item>
				<ItemSeparator/>
				<Item>
					<ItemMedia variant="icon">
						<icon.User/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>Alex Chen</ItemTitle>
						<ItemDescription>alex@example.com</ItemDescription>
					</ItemContent>
					<ItemActions>
						<Button variant="ghost" size="sm">Remove</Button>
					</ItemActions>
				</Item>
			</ItemGroup>
			<HoverCard>
				<HoverCardTrigger>
					<Button variant="link">@gsxui</Button>
				</HoverCardTrigger>
				<HoverCardContent class="w-72">
					<div class="flex justify-between gap-4">
						<Avatar>
							<AvatarFallback>GX</AvatarFallback>
						</Avatar>
						<div class="space-y-1">
							<h4 class="text-sm font-semibold">@gsxui</h4>
							<p class="text-sm">Server-rendered shadcn/ui components for Go.</p>
							<div class="text-xs text-muted-foreground">Joined July 2026</div>
						</div>
					</div>
				</HoverCardContent>
			</HoverCard>
		</CardContent>
	</Card>
}

component galleryTableCard() {
	<Card>
		<CardHeader>
			<CardTitle>Recent invoices</CardTitle>
			<CardDescription>Your billing history for this quarter.</CardDescription>
		</CardHeader>
		<CardContent>
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Invoice</TableHead>
						<TableHead>Status</TableHead>
						<TableHead>Amount</TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow>
						<TableCell>INV-2041</TableCell>
						<TableCell><Badge variant="outline">Paid</Badge></TableCell>
						<TableCell>$250.00</TableCell>
					</TableRow>
					<TableRow data-state="selected">
						<TableCell>INV-2042</TableCell>
						<TableCell><Badge variant="secondary">Pending</Badge></TableCell>
						<TableCell>$125.00</TableCell>
					</TableRow>
					<TableRow>
						<TableCell>INV-2043</TableCell>
						<TableCell><Badge variant="destructive">Overdue</Badge></TableCell>
						<TableCell>$75.00</TableCell>
					</TableRow>
				</TableBody>
				<TableFooter>
					<TableRow>
						<TableCell>Total</TableCell>
						<TableCell></TableCell>
						<TableCell>$450.00</TableCell>
					</TableRow>
				</TableFooter>
			</Table>
		</CardContent>
	</Card>
}

// galleryChartCard is a static color swatch shaped like a chart, painted
// with var(--chart-1)..var(--chart-5) — a preview for the chart-color
// axis, deliberately not a real charting component (no data binding, no
// tooltips; the theme preview loads no component JS to begin with, per
// the doc comment at the top of this file).
component galleryChartCard() {
	<Card>
		<CardHeader>
			<CardTitle>Chart colors</CardTitle>
			<CardDescription>chart-1 through chart-5.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			{/* Percentage heights need a parent with a DEFINITE height, so the
			   bars are direct children of the h-32 row (not nested inside a
			   flex-1 column, which has no height of its own to resolve
			   against) — the number labels live in a separate row underneath. */}
			<div class="flex h-32 items-end gap-3">
				<div class="flex-1 rounded-t-sm" style="height: 45%; background-color: var(--chart-1)"></div>
				<div class="flex-1 rounded-t-sm" style="height: 70%; background-color: var(--chart-2)"></div>
				<div class="flex-1 rounded-t-sm" style="height: 55%; background-color: var(--chart-3)"></div>
				<div class="flex-1 rounded-t-sm" style="height: 90%; background-color: var(--chart-4)"></div>
				<div class="flex-1 rounded-t-sm" style="height: 65%; background-color: var(--chart-5)"></div>
			</div>
			<div class="flex gap-3 text-center text-xs text-muted-foreground">
				<span class="flex-1">1</span>
				<span class="flex-1">2</span>
				<span class="flex-1">3</span>
				<span class="flex-1">4</span>
				<span class="flex-1">5</span>
			</div>
			<div class="flex flex-wrap gap-x-4 gap-y-1.5 text-xs text-muted-foreground">
				<span class="flex items-center gap-1.5"><span class="size-2.5 rounded-full" style="background-color: var(--chart-1)"></span>chart-1</span>
				<span class="flex items-center gap-1.5"><span class="size-2.5 rounded-full" style="background-color: var(--chart-2)"></span>chart-2</span>
				<span class="flex items-center gap-1.5"><span class="size-2.5 rounded-full" style="background-color: var(--chart-3)"></span>chart-3</span>
				<span class="flex items-center gap-1.5"><span class="size-2.5 rounded-full" style="background-color: var(--chart-4)"></span>chart-4</span>
				<span class="flex items-center gap-1.5"><span class="size-2.5 rounded-full" style="background-color: var(--chart-5)"></span>chart-5</span>
			</div>
		</CardContent>
	</Card>
}

component galleryTabsCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle>Project overview</CardTitle>
			<CardDescription>Tabs and frequently asked questions.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<Tabs value="overview">
				<TabsList>
					<TabsTrigger value="overview" selected>Overview</TabsTrigger>
					<TabsTrigger value="analytics">Analytics</TabsTrigger>
					<TabsTrigger value="reports">Reports</TabsTrigger>
				</TabsList>
				<TabsContent value="overview" selected>
					<p class="text-sm text-muted-foreground">Deployments are healthy. Last release shipped 2 hours ago.</p>
				</TabsContent>
				<TabsContent value="analytics">
					<p class="text-sm text-muted-foreground">Traffic is up 12% week over week.</p>
				</TabsContent>
				<TabsContent value="reports">
					<p class="text-sm text-muted-foreground">Monthly reports are generated on the 1st.</p>
				</TabsContent>
			</Tabs>
			<Accordion name={idp + "-faq"}>
				<AccordionItem name={idp + "-faq"} open>
					<AccordionTrigger>Is it accessible?</AccordionTrigger>
					<AccordionContent>Yes, it follows the WAI-ARIA design pattern.</AccordionContent>
				</AccordionItem>
				<AccordionItem name={idp + "-faq"}>
					<AccordionTrigger>Is it themed?</AccordionTrigger>
					<AccordionContent>Yes, every surface reads the semantic theme tokens.</AccordionContent>
				</AccordionItem>
			</Accordion>
		</CardContent>
	</Card>
}

component galleryNavigationCard() {
	<Card>
		<CardHeader>
			<CardTitle>Navigation</CardTitle>
			<CardDescription>Breadcrumbs, menus, and pagination.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<Breadcrumb>
				<BreadcrumbList>
					<BreadcrumbItem>
						<BreadcrumbLink href="#">Home</BreadcrumbLink>
					</BreadcrumbItem>
					<BreadcrumbSeparator/>
					<BreadcrumbItem>
						<BreadcrumbEllipsis/>
					</BreadcrumbItem>
					<BreadcrumbSeparator/>
					<BreadcrumbItem>
						<BreadcrumbLink href="#">Components</BreadcrumbLink>
					</BreadcrumbItem>
					<BreadcrumbSeparator/>
					<BreadcrumbItem>
						<BreadcrumbPage>Gallery</BreadcrumbPage>
					</BreadcrumbItem>
				</BreadcrumbList>
			</Breadcrumb>
			<NavigationMenu>
				<NavigationMenuList>
					<NavigationMenuItem>
						<NavigationMenuLink variant="trigger" active={true} href="#">Home</NavigationMenuLink>
					</NavigationMenuItem>
					<NavigationMenuItem>
						<NavigationMenuTrigger>Components</NavigationMenuTrigger>
						<NavigationMenuContent>
							<div class="grid w-64 gap-2">
								<NavigationMenuLink href="#">
									<div class="text-sm font-medium">Dialog</div>
									<div class="text-muted-foreground">A window overlaid on the page.</div>
								</NavigationMenuLink>
								<NavigationMenuLink href="#">
									<div class="text-sm font-medium">Tooltip</div>
									<div class="text-muted-foreground">A popup shown on hover.</div>
								</NavigationMenuLink>
							</div>
						</NavigationMenuContent>
					</NavigationMenuItem>
					<NavigationMenuItem>
						<NavigationMenuLink variant="trigger" href="#">Docs</NavigationMenuLink>
					</NavigationMenuItem>
					<NavigationMenuIndicator/>
				</NavigationMenuList>
			</NavigationMenu>
			<Pagination>
				<PaginationContent>
					<PaginationItem>
						<PaginationPrevious href="#"/>
					</PaginationItem>
					<PaginationItem>
						<PaginationLink href="#">1</PaginationLink>
					</PaginationItem>
					<PaginationItem>
						<PaginationLink href="#" isActive>2</PaginationLink>
					</PaginationItem>
					<PaginationItem>
						<PaginationLink href="#">3</PaginationLink>
					</PaginationItem>
					<PaginationItem>
						<PaginationEllipsis/>
					</PaginationItem>
					<PaginationItem>
						<PaginationNext href="#"/>
					</PaginationItem>
				</PaginationContent>
			</Pagination>
		</CardContent>
	</Card>
}

component galleryControlsCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle>Editor controls</CardTitle>
			<CardDescription>Toggles, groups, and composed inputs.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<div class="flex flex-wrap items-center gap-3">
				<Toggle pressed={true} aria-label="Toggle bold">
					<icon.Bold/>
				</Toggle>
				<Toggle variant="outline" aria-label="Toggle italic">
					<icon.Italic/>
				</Toggle>
				<ToggleGroup groupType="multiple" variant="outline" aria-label="Text formatting">
					<ToggleGroupItem groupType="multiple" variant="outline" value="bold" pressed={true} aria-label="Toggle bold">
						<icon.Bold/>
					</ToggleGroupItem>
					<ToggleGroupItem groupType="multiple" variant="outline" value="italic" aria-label="Toggle italic">
						<icon.Italic/>
					</ToggleGroupItem>
					<ToggleGroupItem groupType="multiple" variant="outline" value="underline" aria-label="Toggle underline">
						<icon.Underline/>
					</ToggleGroupItem>
				</ToggleGroup>
			</div>
			<div class="flex flex-wrap items-start gap-3">
				<ButtonGroup>
					<Button variant="outline">Archive</Button>
					<Button variant="outline">Report</Button>
				</ButtonGroup>
				<ButtonGroup aria-label="Quantity">
					<Button variant="outline" size="icon" aria-label="Decrease quantity">
						<icon.Minus/>
					</Button>
					<ButtonGroupText>42</ButtonGroupText>
					<Button variant="outline" size="icon" aria-label="Increase quantity">
						<icon.Plus/>
					</Button>
				</ButtonGroup>
				<ButtonGroup>
					<Button variant="secondary">Copy</Button>
					<ButtonGroupSeparator/>
					<Button variant="secondary">Paste</Button>
				</ButtonGroup>
			</div>
			<InputGroup>
				<InputGroupAddon>
					<icon.Search class="size-4"/>
				</InputGroupAddon>
				<InputGroupInput placeholder="Search projects..."/>
				<InputGroupAddon align="inline-end">
					<InputGroupButton aria-label="Send">
						<icon.Send/>
					</InputGroupButton>
				</InputGroupAddon>
			</InputGroup>
			<div class="flex flex-col gap-2">
				<Label>Verification code</Label>
				<InputOTP maxlength="6">
					<InputOTPGroup>
						<InputOTPSlot/>
						<InputOTPSlot/>
						<InputOTPSlot/>
					</InputOTPGroup>
					<InputOTPSeparator/>
					<InputOTPGroup>
						<InputOTPSlot/>
						<InputOTPSlot/>
						<InputOTPSlot/>
					</InputOTPGroup>
				</InputOTP>
			</div>
			<div class="flex flex-col gap-2">
				<Label for={idp + "-controls-framework"}>Framework</Label>
				<Combobox name={idp + "-controls-framework"} value="">
					<ComboboxInput id={idp + "-controls-framework"} placeholder="Search framework..." showTrigger/>
					<ComboboxContent>
						<ComboboxList>
							<ComboboxEmpty>No framework found.</ComboboxEmpty>
							<ComboboxItem value="next.js" selected={false}>Next.js</ComboboxItem>
							<ComboboxItem value="sveltekit" selected={false}>SvelteKit</ComboboxItem>
							<ComboboxItem value="astro" selected={false}>Astro</ComboboxItem>
						</ComboboxList>
					</ComboboxContent>
				</Combobox>
			</div>
		</CardContent>
	</Card>
}

component galleryOverlaysCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle>Overlays</CardTitle>
			<CardDescription>Dialogs, sheets, popovers, and the command palette.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<div class="flex flex-wrap items-center gap-2">
				<Dialog>
					<Button
						variant="outline"
						data-gsxui-slot-dialog-trigger
						aria-haspopup="dialog"
						aria-expanded="false"
					>
						Edit profile
					</Button>
					<DialogContent>
						<DialogHeader>
							<DialogTitle>Edit profile</DialogTitle>
							<DialogDescription>Make changes to your profile here.</DialogDescription>
						</DialogHeader>
						<DialogFooter>
							<Button variant="outline" data-gsxui-dialog-close>Cancel</Button>
							<Button>Save changes</Button>
						</DialogFooter>
					</DialogContent>
				</Dialog>
				<AlertDialog>
					<Button
						variant="outline"
						data-gsxui-slot-alert-dialog-trigger
						aria-haspopup="dialog"
						aria-expanded="false"
					>
						Delete account
					</Button>
					<AlertDialogContent>
						<AlertDialogHeader>
							<AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
							<AlertDialogDescription>This action cannot be undone.</AlertDialogDescription>
						</AlertDialogHeader>
						<AlertDialogFooter>
							<AlertDialogCancel>Cancel</AlertDialogCancel>
							<AlertDialogAction>Continue</AlertDialogAction>
						</AlertDialogFooter>
					</AlertDialogContent>
				</AlertDialog>
				<Sheet>
					<Button
						variant="outline"
						data-gsxui-slot-sheet-trigger
						aria-haspopup="dialog"
						aria-expanded="false"
					>
						Open sheet
					</Button>
					<SheetContent side="" hideCloseButton={false}>
						<SheetHeader>
							<SheetTitle>Edit profile</SheetTitle>
							<SheetDescription>Make changes to your profile here.</SheetDescription>
						</SheetHeader>
						<SheetFooter>
							<Button data-gsxui-dialog-close>Save changes</Button>
						</SheetFooter>
					</SheetContent>
				</Sheet>
				<Drawer>
					<Button
						variant="outline"
						data-gsxui-slot-drawer-trigger
						aria-haspopup="dialog"
						aria-expanded="false"
					>
						Open drawer
					</Button>
					<DrawerContent direction="">
						<DrawerHeader>
							<DrawerTitle>Move goal</DrawerTitle>
							<DrawerDescription>Set your daily activity goal.</DrawerDescription>
						</DrawerHeader>
						<DrawerFooter>
							<Button>Submit</Button>
							<Button variant="outline" data-gsxui-dialog-close>Cancel</Button>
						</DrawerFooter>
					</DrawerContent>
				</Drawer>
				<Popover>
					<Button
						variant="outline"
						data-gsxui-slot-popover-trigger
						aria-expanded="false"
					>
						Dimensions
					</Button>
					<PopoverContent class="w-72">
						<div class="grid gap-3">
							<p class="text-sm text-muted-foreground">Set the dimensions for the layer.</p>
							<div class="grid grid-cols-3 items-center gap-3">
								<Label for={idp + "-overlay-width"}>Width</Label>
								<Input id={idp + "-overlay-width"} value="100%" class="col-span-2"/>
							</div>
						</div>
					</PopoverContent>
				</Popover>
				<Tooltip>
					<Button variant="outline" data-gsxui-slot-tooltip-trigger>Hover me</Button>
					<TooltipContent>Add to library</TooltipContent>
				</Tooltip>
			</div>
			<Command class="rounded-lg border shadow-md">
				<CommandInput placeholder="Type a command or search..."/>
				<CommandList>
					<CommandEmpty>No results found.</CommandEmpty>
					<CommandGroup heading="Suggestions">
						<CommandItem>
							<icon.Calendar/>
							<span>Calendar</span>
						</CommandItem>
						<CommandItem>
							<icon.Smile/>
							<span>Search Emoji</span>
						</CommandItem>
					</CommandGroup>
					<CommandSeparator/>
					<CommandGroup heading="Settings">
						<CommandItem>
							<icon.User/>
							<span>Profile</span>
							<CommandShortcut>⌘P</CommandShortcut>
						</CommandItem>
					</CommandGroup>
				</CommandList>
			</Command>
		</CardContent>
	</Card>
}

component galleryEmptyCard() {
	<Card>
		<CardContent class="flex flex-col gap-4">
			<Empty>
				<EmptyHeader>
					<EmptyMedia variant="icon">
						<icon.Inbox/>
					</EmptyMedia>
					<EmptyTitle>No messages</EmptyTitle>
					<EmptyDescription>You're all caught up. New messages will appear here.</EmptyDescription>
				</EmptyHeader>
				<EmptyContent>
					<Button>Compose message</Button>
				</EmptyContent>
			</Empty>
			<Collapsible open class="flex flex-col gap-2">
				<CollapsibleTrigger class="flex cursor-default items-center justify-between gap-4 px-1 text-sm font-semibold">
					3 starred repositories
				</CollapsibleTrigger>
				<div class="rounded-md border px-4 py-2 font-mono text-sm">gsxhq/gsx</div>
				<CollapsibleContent class="flex flex-col gap-2">
					<div class="rounded-md border px-4 py-2 font-mono text-sm">gsxhq/gsxui</div>
					<div class="rounded-md border px-4 py-2 font-mono text-sm">gsxhq/vite</div>
				</CollapsibleContent>
			</Collapsible>
		</CardContent>
	</Card>
}

component galleryMediaCard() {
	<Card>
		<CardHeader>
			<CardTitle>Media and layout</CardTitle>
			<CardDescription>Ratio boxes, scroll areas, carousels, and panes.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<AspectRatio
				ratio="16 / 9"
				class="flex items-center justify-center rounded-lg border bg-muted text-sm text-muted-foreground"
			>
				16 / 9
			</AspectRatio>
			<ScrollArea class="h-40 w-full rounded-md border">
				<div class="p-4">
					<h4 class="mb-3 text-sm leading-none font-medium">Releases</h4>
					{ for _, tag := range []string{
						"v1.2.0-beta.5", "v1.2.0-beta.4", "v1.2.0-beta.3",
						"v1.2.0-beta.2", "v1.2.0-beta.1", "v1.1.9", "v1.1.8", "v1.1.7",
					} {
						<div>
							<div class="text-sm">{ tag }</div>
							<Separator class="my-2"/>
						</div>
					} }
				</div>
			</ScrollArea>
			<Carousel class="mx-auto w-full max-w-[240px]">
				<CarouselContent>
					{ for _, n := range []int{1, 2, 3} {
						<CarouselItem>
							<div class="p-1">
								<Card>
									<CardContent class="flex aspect-square items-center justify-center p-6">
										<span class="text-4xl font-semibold">{ n }</span>
									</CardContent>
								</Card>
							</div>
						</CarouselItem>
					} }
				</CarouselContent>
				<CarouselPrevious/>
				<CarouselNext/>
			</Carousel>
			<ResizablePanelGroup orientation="horizontal" class="rounded-lg border">
				<ResizablePanel defaultSize="50%">
					<div class="flex h-24 items-center justify-center p-4">
						<span class="text-sm font-semibold">One</span>
					</div>
				</ResizablePanel>
				<ResizableHandle orientation="horizontal" withHandle={true}/>
				<ResizablePanel defaultSize="50%">
					<div class="flex h-24 items-center justify-center p-4">
						<span class="text-sm font-semibold">Two</span>
					</div>
				</ResizablePanel>
			</ResizablePanelGroup>
		</CardContent>
	</Card>
}

component gallerySidebarCard(idp string) {
	<Card class="md:col-span-2 xl:col-span-3">
		<CardHeader>
			<CardTitle>Application shell</CardTitle>
			<CardDescription>Sidebar navigation with inset content.</CardDescription>
		</CardHeader>
		<CardContent>
			{/* The sidebar panel is position:fixed by design; the transform on
			    this wrapper makes it the panel's containing block, so the shell
			    stays inside its card instead of overlaying the page. */}
			<div class="relative isolate transform-gpu overflow-hidden rounded-lg border">
				<SidebarProvider open={true} class="min-h-[24rem]">
				<Sidebar open={true}>
					<SidebarHeader>
						<div class="px-2 py-1 text-sm font-semibold">Acme Inc</div>
						<SidebarInput placeholder="Search navigation"/>
					</SidebarHeader>
					<SidebarSeparator/>
					<SidebarContent>
						<SidebarGroup>
							<SidebarGroupLabel>Application</SidebarGroupLabel>
							<SidebarGroupContent>
								<SidebarMenu>
									<SidebarMenuItem>
										<SidebarMenuButton isActive={true}>
											<icon.House/>
											<span>Home</span>
										</SidebarMenuButton>
									</SidebarMenuItem>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<icon.Inbox/>
											<span>Inbox</span>
										</SidebarMenuButton>
										<SidebarMenuBadge>24</SidebarMenuBadge>
									</SidebarMenuItem>
									<SidebarMenuItem>
										<SidebarMenuButton>
											<icon.Settings/>
											<span>Settings</span>
										</SidebarMenuButton>
									</SidebarMenuItem>
								</SidebarMenu>
							</SidebarGroupContent>
						</SidebarGroup>
					</SidebarContent>
					<SidebarFooter>
						<SidebarMenu>
							<SidebarMenuItem>
								<SidebarMenuButton>
									<icon.User/>
									<span>Account</span>
								</SidebarMenuButton>
							</SidebarMenuItem>
						</SidebarMenu>
					</SidebarFooter>
				</Sidebar>
				<SidebarInset>
					<header class="flex h-12 items-center gap-2 border-b px-4">
						<SidebarTrigger/>
						<span class="text-sm text-muted-foreground">Dashboard</span>
					</header>
					<div class="p-4 text-sm text-muted-foreground">
						Overview of { idp } workspace activity.
					</div>
				</SidebarInset>
				</SidebarProvider>
			</div>
		</CardContent>
	</Card>
}
