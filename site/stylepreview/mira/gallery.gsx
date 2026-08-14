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
// The preview document's entry (web/preview.js) already imports the full
// behaviour barrel (ui/index.js, see PR #15 "theme preview: load the
// behaviour barrel in the preview entries") — every popup (menus, dialogs,
// popovers, tooltips) is live, resting closed only because nothing in the
// gallery's own markup triggers it open, and galleryChartCard's charts
// render a real client-side SVG the same barrel already covers (ui/chart.js
// is one of the barrel's modules). What the gallery showcases is trigger
// and surface chrome theming, plus (galleryChartCard) live chart geometry —
// not scripted end-to-end behavior.
//
// Gallery renders TWO pages (theme-creator-parity Task 4c), mirroring the
// per-style section toggle site/pages/theme_preview.gsx and
// web/theme-preview.js already use: both page groups are always in the
// DOM, `hidden` decides which one is visible, and web/theme-preview.js
// flips that attribute the same way it already flips which style section
// is visible. "components" (galleryComponentsPage) is the original
// broad-coverage showcase — every component, every state, in a compact
// form. "product" (galleryProductPage) is newer: realistic product
// surfaces (an invoice, a chat thread, a pricing table, ...) built from
// the SAME 54 canonical components with plausible copy, closer to what a
// real screen built with gsxui looks like. Both pages count toward
// gallery_test.go's per-component coverage floor — the test scans this
// whole file's text, not just one page's markup.
//
// Every card's <CardTitle> carries an inline `font-family: var(--font-
// heading)` — this is the font pair picker's (theme-creator-parity Task 3)
// only visible manifestation in this preview: --font-sans already applies
// to the whole document for free (it's the CSS default font-family chain
// Tailwind's Preflight resolves from, so a runtime override on the root
// element cascades everywhere without any markup change), but nothing
// upstream of this file previously read --font-heading at all. This
// mirrors upstream's own approach exactly — a hand-authored .cn-font-
// heading marker on specific preview-only elements (dossier §3/§5(b)), not
// a change to the shared Card primitive (ui/card.gsx is untouched; a real
// project's own Card usage is unaffected by the font axis existing).
package mira

import (
	"time"

	"github.com/gsxhq/gsx"
	"github.com/gsxhq/gsxui/ui/icon"
)

// galleryMonth is fixed so the preview never depends on the day it renders.
var galleryMonth = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// galleryAvatarSVG is a stand-in portrait tile, the same data-URI idiom as
// site/examples/avatar/basic.gsx — no image asset, no network fetch.
var galleryAvatarSVG = []byte("<svg xmlns='http://www.w3.org/2000/svg' width='64' height='64'><rect width='64' height='64' fill='#6d28d9'/><text x='32' y='34' text-anchor='middle' dominant-baseline='central' font-family='sans-serif' font-weight='600' font-size='26' fill='#fff'>AL</text></svg>")

// Gallery renders both preview pages. idp prefixes every id and
// form-control name: the theme preview document renders the gallery twice
// per page (once per style section), so unprefixed ids would collide
// across sections.
component Gallery(idp string) {
	<galleryComponentsPage idp={idp}/>
	<galleryProductPage idp={idp}/>
}

// galleryComponentsPage is the original 15-card broad-coverage showcase:
// every component, every state, in a compact form. Visible by default —
// web/theme-state.js's createThemeState defaults state.page to
// "components".
component galleryComponentsPage(idp string) {
	<div
		data-theme-preview-page="components"
		data-theme-preview-gallery
		class="mx-auto grid w-full max-w-7xl items-start gap-4 p-4 sm:p-6 md:grid-cols-2 xl:grid-cols-3"
	>
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

// galleryProductPage is the realistic-product-surface showcase
// (theme-creator-parity Task 4b): fifteen cards built from the same 54
// canonical components as galleryComponentsPage, composed as plausible
// product screens (an invoice, a chat thread, a pricing table, ...) rather
// than component demonstrations. `hidden` by default; web/theme-preview.js
// reveals it when web/theme.js's page tab commits state.page = "product".
component galleryProductPage(idp string) {
	<div
		data-theme-preview-page="product"
		data-theme-preview-gallery
		hidden
		class="mx-auto grid w-full max-w-7xl items-start gap-4 p-4 sm:p-6 md:grid-cols-2 xl:grid-cols-3"
	>
		<galleryPricingCard/>
		<galleryStatsCard/>
		<galleryChatCard idp={idp}/>
		<galleryRolesCard/>
		<galleryBookingCard idp={idp}/>
		<galleryActivityCard/>
		<galleryUploadsCard/>
		<galleryNotificationsCard/>
		<gallerySearchCard idp={idp}/>
		<galleryProfileCard/>
		<galleryOnboardingCard/>
		<galleryProjectsEmptyCard/>
		<galleryVerifyCard/>
		<galleryNotificationPrefsCard idp={idp}/>
		<galleryBillingCard/>
	</div>
}

// galleryButtonsCard keeps the old Matrix's coverage: every variant, every
// size, and the states the theme editor keys on (data-theme-preview-case is
// web/theme-preview.js's focus hook and pages_test.go's marker set).
component galleryButtonsCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Buttons</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Create an account</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Workspace settings</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Schedule a review</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Menus</CardTitle>
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
						{/* data-inset mirrors upstream's menubar demo and gives every
						    style's per-family inset value a rendered element on the
						    preview document (menus stay closed there, but computed
						    style still resolves — theme-editor.spec.ts pins it). */}
						<MenubarItem data-inset="true">
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
			<CardTitle style="font-family: var(--font-heading)">Sync status</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Team members</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Recent invoices</CardTitle>
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

// galleryChartCard renders two real charts, both themed entirely off
// var(--chart-1)..var(--chart-5): a five-series BarChart (one series per
// chart-N color — the direct pendant of the old static swatch this card
// used to show) and a compact single-series AreaChart trending over time.
//
// Both are ui.Chart trees (Task 9, docs/superpowers/plans/
// 2026-08-14-chart-component.md), which need no extra wiring in
// web/preview.js: the entry already imports the full behaviour barrel
// (ui/index.js, see this file's own header comment), and ui/chart.js — one
// of that barrel's modules — lazily pulls in ui/chart.render.js the first
// time it finds a [data-gsxui-slot-chart] and draws the SVG client-side.
// This is also what makes the theme editor's chart-color picker demo real:
// picking a new chart-1 changes the BarChart's first (organic) bar's
// computed SVG fill live, through the CSS cascade alone
// (jstest/specs/theme-editor.spec.ts).
//
// Both Chart calls below carry a caller class of min-h-0: Chart's own div
// is "flex" and, as this card's CardContent flex column's own flex ITEM,
// picks up the flexbox default min-height:auto — a content-based floor
// that measured TALLER than the entrance animation's very first frame (the
// client renderer sizes its SVG off panel.clientHeight every animation
// frame, ui/chart.render.js's renderCartesian), so each frame's
// slightly-taller draw pushed clientHeight up again the next frame, a
// runaway growth for the whole 400ms bar entrance that only site/examples/
// chart's own wide, unconstrained layout happens to never trigger (its
// aspect-video height is already generous enough that content never
// exceeds it). min-h-0 is the standard fix for exactly this flex-item
// content-driven-growth gotcha, caller-overriding Chart's own class the
// same way class="aspect-[5/1]" below already does (attrs merge, caller
// wins per property) — no change to ui/chart.render.js or ui/chart.gsx.
component galleryChartCard() {
	{{
		channelConfig := ChartConfig{
			{Key: "organic", Label: "Organic", Color: "var(--chart-1)"},
			{Key: "paid", Label: "Paid", Color: "var(--chart-2)"},
			{Key: "referral", Label: "Referral", Color: "var(--chart-3)"},
			{Key: "email", Label: "Email", Color: "var(--chart-4)"},
			{Key: "social", Label: "Social", Color: "var(--chart-5)"},
		}
		channelData := []ChartDatum{
			{"quarter": "Q1", "organic": 186.0, "paid": 80.0, "referral": 45.0, "email": 60.0, "social": 30.0},
			{"quarter": "Q2", "organic": 210.0, "paid": 95.0, "referral": 52.0, "email": 58.0, "social": 41.0},
			{"quarter": "Q3", "organic": 237.0, "paid": 120.0, "referral": 61.0, "email": 65.0, "social": 55.0},
		}
		trendConfig := ChartConfig{
			{Key: "visits", Label: "Visits", Color: "var(--chart-1)"},
		}
		trendData := []ChartDatum{
			{"week": "W1", "visits": 186.0},
			{"week": "W2", "visits": 240.0},
			{"week": "W3", "visits": 214.0},
			{"week": "W4", "visits": 290.0},
			{"week": "W5", "visits": 305.0},
		}
	}}
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Revenue by channel</CardTitle>
			<CardDescription>chart-1 through chart-5, live-rendered.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<Chart config={channelConfig} class="min-h-0">
				{ BarChart(channelData, nil, gsx.Fragment(
					ChartCartesianGrid(&ChartCartesianGridOptions{Vertical: ChartBool(false)}),
					ChartXAxis("quarter", nil),
					ChartTooltip(nil),
					ChartBar("organic", nil),
					ChartBar("paid", nil),
					ChartBar("referral", nil),
					ChartBar("email", nil),
					ChartBar("social", nil),
					ChartLegend(nil),
				), nil) }
			</Chart>
			<div class="flex flex-col gap-1.5">
				<span class="text-xs text-muted-foreground">Weekly visits</span>
				<Chart config={trendConfig} class="aspect-[5/1] min-h-0">
					{ AreaChart(trendData, nil, gsx.Fragment(
						ChartXAxis("week", &ChartXAxisOptions{Hide: true}),
						ChartTooltip(&ChartTooltipOptions{Cursor: ChartBool(false)}),
						ChartArea("visits", &ChartAreaOptions{Type: ChartCurveNatural, FillOpacity: 0.3}),
					), nil) }
				</Chart>
			</div>
		</CardContent>
	</Card>
}

component galleryTabsCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Project overview</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Navigation</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Editor controls</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Overlays</CardTitle>
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
			<Command class="border shadow-md">
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
			<CardTitle style="font-family: var(--font-heading)">Media and layout</CardTitle>
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
			<CardTitle style="font-family: var(--font-heading)">Application shell</CardTitle>
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

// --- Product page (theme-creator-parity Task 4b) -------------------------
//
// Realistic product surfaces, not component showcases: plausible copy,
// real-sounding names, dates fixed to January/February 2026 (galleryMonth's
// own era) so nothing here depends on the day it renders.

// galleryPricingCard is a three-tier pricing table, the kind of page a
// theme actually needs to look right on: nested Cards (same idiom
// galleryMediaCard's Carousel already uses for its items) for each tier,
// with the middle tier highlighted as the recommended plan. Spans the full
// grid width (the same md:col-span-2 xl:col-span-3 idiom
// gallerySidebarCard already uses on the components page) — three tiers
// need real width, and a pricing table is realistically a page's widest
// element anyway.
component galleryPricingCard() {
	<Card class="md:col-span-2 xl:col-span-3">
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Pricing</CardTitle>
			<CardDescription>Simple plans that scale with your team.</CardDescription>
		</CardHeader>
		<CardContent class="grid gap-3 sm:grid-cols-3">
			<Card class="gap-3 shadow-none">
				<CardHeader>
					<CardTitle class="text-base">Starter</CardTitle>
					<CardDescription><span class="text-lg font-semibold text-foreground">$0</span>/mo</CardDescription>
				</CardHeader>
				<CardContent class="flex flex-1 flex-col gap-2 text-sm text-muted-foreground">
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0"/>5 projects</span>
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0"/>Community support</span>
					<div class="mt-auto pt-2"><Button variant="outline" class="w-full">Get started</Button></div>
				</CardContent>
			</Card>
			<Card class="gap-3 border-primary shadow-none">
				<CardHeader>
					<CardAction>
						<Badge>Popular</Badge>
					</CardAction>
					<CardTitle class="text-base">Team</CardTitle>
					<CardDescription><span class="text-lg font-semibold text-foreground">$24</span>/user/mo</CardDescription>
				</CardHeader>
				<CardContent class="flex flex-1 flex-col gap-2 text-sm">
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0 text-primary"/>Unlimited projects</span>
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0 text-primary"/>Priority support</span>
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0 text-primary"/>Single sign-on</span>
					<div class="mt-auto pt-2"><Button class="w-full">Start trial</Button></div>
				</CardContent>
			</Card>
			<Card class="gap-3 shadow-none">
				<CardHeader>
					<CardTitle class="text-base">Enterprise</CardTitle>
					<CardDescription>Custom pricing</CardDescription>
				</CardHeader>
				<CardContent class="flex flex-1 flex-col gap-2 text-sm text-muted-foreground">
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0"/>Dedicated support</span>
					<span class="flex items-start gap-1.5"><icon.Check class="mt-0.5 size-3.5 shrink-0"/>Audit logs</span>
					<div class="mt-auto pt-2"><Button variant="outline" class="w-full">Contact sales</Button></div>
				</CardContent>
			</Card>
		</CardContent>
	</Card>
}

// galleryStatsCard is a KPI/stat row — a dashboard's first screen, the
// kind of overview a theme's colors have to carry at a glance.
component galleryStatsCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Overview</CardTitle>
			<CardDescription>Key metrics for the last 30 days.</CardDescription>
		</CardHeader>
		<CardContent class="grid grid-cols-2 gap-4 sm:grid-cols-4">
			<div class="flex flex-col gap-1">
				<span class="text-xs text-muted-foreground">MRR</span>
				<span class="text-xl font-semibold">$48,230</span>
				<span class="flex items-center gap-1 text-xs text-primary"><icon.TrendingUp class="size-3.5"/>12.4%</span>
			</div>
			<div class="flex flex-col gap-1">
				<span class="text-xs text-muted-foreground">Active users</span>
				<span class="text-xl font-semibold">8,412</span>
				<span class="flex items-center gap-1 text-xs text-primary"><icon.TrendingUp class="size-3.5"/>3.1%</span>
			</div>
			<div class="flex flex-col gap-1">
				<span class="text-xs text-muted-foreground">Churn</span>
				<span class="text-xl font-semibold">1.8%</span>
				<span class="flex items-center gap-1 text-xs text-destructive"><icon.TrendingDown class="size-3.5"/>0.4%</span>
			</div>
			<div class="flex flex-col gap-1">
				<span class="text-xs text-muted-foreground">NPS</span>
				<span class="text-xl font-semibold">62</span>
				<span class="flex items-center gap-1 text-xs text-primary"><icon.TrendingUp class="size-3.5"/>5</span>
			</div>
		</CardContent>
	</Card>
}

// galleryChatCard is a support-inbox thread: alternating incoming/outgoing
// bubbles plus a composer row, exercising primary/muted surfaces the way a
// real messaging UI would.
component galleryChatCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Support inbox</CardTitle>
			<CardDescription>Conversation with a customer.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-3">
			<div class="flex flex-col items-start gap-1">
				<div class="max-w-[80%] rounded-lg bg-muted px-3 py-2 text-sm">I'm having trouble exporting my report to CSV.</div>
				<span class="text-xs text-muted-foreground">Priya · 9:41 AM</span>
			</div>
			<div class="flex flex-col items-end gap-1">
				<div class="max-w-[80%] rounded-lg bg-primary px-3 py-2 text-sm text-primary-foreground">Try Settings → Export → CSV. Let me know if that works!</div>
				<span class="text-xs text-muted-foreground">You · 9:44 AM</span>
			</div>
			<div class="flex flex-col items-start gap-1">
				<div class="max-w-[80%] rounded-lg bg-muted px-3 py-2 text-sm">That fixed it, thank you!</div>
				<span class="text-xs text-muted-foreground">Priya · 9:46 AM</span>
			</div>
		</CardContent>
		<CardFooter class="gap-2">
			<Input id={idp + "-chat-message"} placeholder="Type a message..." class="flex-1"/>
			<Button size="icon" aria-label="Send message">
				<icon.Send/>
			</Button>
		</CardFooter>
	</Card>
}

// galleryRolesCard is a team-roles table with per-row action menus —
// distinct from galleryTeamCard's invite-focused avatar list on the
// components page: this is the "manage access" screen, not the "who's on
// the team" screen.
component galleryRolesCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Team roles</CardTitle>
			<CardDescription>Manage access for your workspace.</CardDescription>
		</CardHeader>
		<CardContent>
			<Table>
				<TableHeader>
					<TableRow>
						<TableHead>Name</TableHead>
						<TableHead>Role</TableHead>
						<TableHead></TableHead>
					</TableRow>
				</TableHeader>
				<TableBody>
					<TableRow>
						<TableCell>Priya Sharma</TableCell>
						<TableCell><Badge variant="secondary">Owner</Badge></TableCell>
						<TableCell class="text-right">
							<DropdownMenu>
								<Button
									variant="ghost"
									size="icon-sm"
									data-gsxui-slot-dropdown-menu-trigger
									aria-haspopup="menu"
									aria-expanded="false"
									aria-label="Actions for Priya Sharma"
								>
									<icon.Ellipsis/>
								</Button>
								<DropdownMenuContent>
									<DropdownMenuItem>Change role</DropdownMenuItem>
									<DropdownMenuItem variant="destructive">Remove</DropdownMenuItem>
								</DropdownMenuContent>
							</DropdownMenu>
						</TableCell>
					</TableRow>
					<TableRow>
						<TableCell>Marcus Webb</TableCell>
						<TableCell><Badge variant="outline">Admin</Badge></TableCell>
						<TableCell class="text-right">
							<DropdownMenu>
								<Button
									variant="ghost"
									size="icon-sm"
									data-gsxui-slot-dropdown-menu-trigger
									aria-haspopup="menu"
									aria-expanded="false"
									aria-label="Actions for Marcus Webb"
								>
									<icon.Ellipsis/>
								</Button>
								<DropdownMenuContent>
									<DropdownMenuItem>Change role</DropdownMenuItem>
									<DropdownMenuItem variant="destructive">Remove</DropdownMenuItem>
								</DropdownMenuContent>
							</DropdownMenu>
						</TableCell>
					</TableRow>
					<TableRow>
						<TableCell>Elena Ruiz</TableCell>
						<TableCell><Badge variant="outline">Member</Badge></TableCell>
						<TableCell class="text-right">
							<DropdownMenu>
								<Button
									variant="ghost"
									size="icon-sm"
									data-gsxui-slot-dropdown-menu-trigger
									aria-haspopup="menu"
									aria-expanded="false"
									aria-label="Actions for Elena Ruiz"
								>
									<icon.Ellipsis/>
								</Button>
								<DropdownMenuContent>
									<DropdownMenuItem>Change role</DropdownMenuItem>
									<DropdownMenuItem variant="destructive">Remove</DropdownMenuItem>
								</DropdownMenuContent>
							</DropdownMenu>
						</TableCell>
					</TableRow>
				</TableBody>
			</Table>
		</CardContent>
	</Card>
}

// galleryBookingCard is a "book a demo" form — distinct from
// galleryCalendarCard's date-picker widget on the components page: this
// is fields-and-submit, not a calendar surface.
component galleryBookingCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Book a demo</CardTitle>
			<CardDescription>Schedule a 30-minute call with our team.</CardDescription>
		</CardHeader>
		<CardContent>
			<form>
				<FieldGroup>
					<Field>
						<FieldLabel for={idp + "-booking-topic"}>Topic</FieldLabel>
						<Select name={idp + "-booking-topic"}>
							<SelectTrigger id={idp + "-booking-topic"} class="w-full">
								<SelectValue placeholder="Select a topic"/>
							</SelectTrigger>
							<SelectContent>
								<SelectGroup>
									<SelectItem value="tour" selected>Product tour</SelectItem>
									<SelectItem value="pricing">Pricing questions</SelectItem>
									<SelectItem value="migration">Migration planning</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
					</Field>
					<Field>
						<FieldLabel for={idp + "-booking-time"}>Time</FieldLabel>
						<NativeSelect id={idp + "-booking-time"} name={idp + "-booking-time"}>
							<NativeSelectOption value="10:00" selected>Tue, Feb 3 · 10:00 AM</NativeSelectOption>
							<NativeSelectOption value="13:30">Tue, Feb 3 · 1:30 PM</NativeSelectOption>
							<NativeSelectOption value="15:00">Wed, Feb 4 · 3:00 PM</NativeSelectOption>
						</NativeSelect>
					</Field>
					<Field>
						<FieldLabel for={idp + "-booking-notes"}>Notes</FieldLabel>
						<Textarea id={idp + "-booking-notes"} value="" placeholder="Anything specific you'd like us to cover?"/>
					</Field>
				</FieldGroup>
			</form>
		</CardContent>
		<CardFooter>
			<Button class="w-full">Request time</Button>
		</CardFooter>
	</Card>
}

// galleryActivityCard is a project activity feed — icon-led events with
// relative timestamps, a common dashboard fixture.
component galleryActivityCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Activity</CardTitle>
			<CardDescription>Recent changes to this project.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-3 text-sm">
			<div class="flex items-start gap-2.5">
				<icon.GitCommitHorizontal class="mt-0.5 size-4 shrink-0 text-muted-foreground"/>
				<div class="flex-1">
					<div><span class="font-medium">Priya</span> pushed 3 commits to main</div>
					<div class="text-xs text-muted-foreground">2 hours ago</div>
				</div>
			</div>
			<Separator/>
			<div class="flex items-start gap-2.5">
				<icon.MessageSquare class="mt-0.5 size-4 shrink-0 text-muted-foreground"/>
				<div class="flex-1">
					<div><span class="font-medium">Marcus</span> commented on #482</div>
					<div class="text-xs text-muted-foreground">5 hours ago</div>
				</div>
			</div>
			<Separator/>
			<div class="flex items-start gap-2.5">
				<icon.CircleCheck class="mt-0.5 size-4 shrink-0 text-muted-foreground"/>
				<div class="flex-1">
					<div><span class="font-medium">Elena</span> closed issue #471</div>
					<div class="text-xs text-muted-foreground">Yesterday</div>
				</div>
			</div>
		</CardContent>
	</Card>
}

// galleryUploadsCard is a file list with one item mid-upload — Progress
// used inline inside an Item rather than standalone, unlike
// galleryFeedbackCard's bare progress bar.
component galleryUploadsCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Files</CardTitle>
			<CardDescription>Shared with your team.</CardDescription>
		</CardHeader>
		<CardContent>
			<ItemGroup>
				<Item>
					<ItemMedia variant="icon">
						<icon.FileText/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>Q4-report.pdf</ItemTitle>
						<ItemDescription>2.4 MB · Uploaded</ItemDescription>
					</ItemContent>
					<ItemActions>
						<Button variant="ghost" size="icon-sm" aria-label="Download Q4-report.pdf">
							<icon.Download/>
						</Button>
					</ItemActions>
				</Item>
				<ItemSeparator/>
				<Item>
					<ItemMedia variant="icon">
						<icon.Archive/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>brand-assets.zip</ItemTitle>
						<ItemDescription>18.1 MB · Uploaded</ItemDescription>
					</ItemContent>
					<ItemActions>
						<Button variant="ghost" size="icon-sm" aria-label="Download brand-assets.zip">
							<icon.Download/>
						</Button>
					</ItemActions>
				</Item>
				<ItemSeparator/>
				<Item>
					<ItemMedia variant="icon">
						<icon.FileImage/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>onboarding-deck.pptx</ItemTitle>
						<ItemDescription>Uploading — 64%</ItemDescription>
						<Progress value={64} class="mt-1.5 h-1.5"/>
					</ItemContent>
				</Item>
			</ItemGroup>
		</CardContent>
	</Card>
}

// galleryNotificationsCard is a notification centre: grouped, unread-dot
// rows and a bulk action — the dropdown/popover CONTENT shadcn's own
// preview-02 ships (e.g. Notification Settings), rendered here as a static
// resting card rather than inside a Popover, since the point is the row
// chrome, not the trigger.
component galleryNotificationsCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Notifications</CardTitle>
			<CardDescription>You have 3 unread.</CardDescription>
			<CardAction>
				<Button variant="ghost" size="sm">Mark all read</Button>
			</CardAction>
		</CardHeader>
		<CardContent class="flex flex-col gap-3">
			<div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">New</div>
			<div class="flex items-start gap-2.5">
				<span class="mt-1.5 size-2 shrink-0 rounded-full bg-primary" aria-hidden="true"></span>
				<div class="flex-1 text-sm">
					<div>Marcus Webb requested access to <span class="font-medium">Design System</span></div>
					<div class="text-xs text-muted-foreground">10 minutes ago</div>
				</div>
			</div>
			<div class="flex items-start gap-2.5">
				<span class="mt-1.5 size-2 shrink-0 rounded-full bg-primary" aria-hidden="true"></span>
				<div class="flex-1 text-sm">
					<div>Your export finished — <span class="font-medium">annual-report.csv</span></div>
					<div class="text-xs text-muted-foreground">1 hour ago</div>
				</div>
			</div>
			<Separator/>
			<div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">Earlier</div>
			<div class="flex items-start gap-2.5">
				<span class="mt-1.5 size-2 shrink-0 rounded-full bg-border" aria-hidden="true"></span>
				<div class="flex-1 text-sm text-muted-foreground">
					<div>Elena Ruiz left a comment on #471</div>
					<div class="text-xs">Yesterday</div>
				</div>
			</div>
		</CardContent>
	</Card>
}

// gallerySearchCard is a search-results surface: a query input over a list
// of typed results, an Item size="sm" variant not otherwise exercised in
// the gallery.
component gallerySearchCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Search</CardTitle>
			<CardDescription>Results for "deployment".</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-3">
			<InputGroup>
				<InputGroupAddon>
					<icon.Search class="size-4"/>
				</InputGroupAddon>
				<InputGroupInput id={idp + "-search-query"} value="deployment"/>
			</InputGroup>
			<ItemGroup>
				<Item size="sm">
					<ItemMedia variant="icon">
						<icon.FileText/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>Deployment checklist</ItemTitle>
						<ItemDescription>Docs / Runbooks</ItemDescription>
					</ItemContent>
				</Item>
				<ItemSeparator/>
				<Item size="sm">
					<ItemMedia variant="icon">
						<icon.MessageSquare/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>#482 deployment failing on staging</ItemTitle>
						<ItemDescription>Issues / Open</ItemDescription>
					</ItemContent>
				</Item>
				<ItemSeparator/>
				<Item size="sm">
					<ItemMedia variant="icon">
						<icon.Users/>
					</ItemMedia>
					<ItemContent>
						<ItemTitle>Deployment on-call rotation</ItemTitle>
						<ItemDescription>Team / Schedules</ItemDescription>
					</ItemContent>
				</Item>
			</ItemGroup>
		</CardContent>
	</Card>
}

// galleryProfileCard is a centered profile summary — deliberately skips
// CardHeader (like galleryEmptyCard on the components page), since a
// photo-first profile card is a different shape than the header/content
// pattern the rest of the gallery uses.
component galleryProfileCard() {
	<Card>
		<CardContent class="flex flex-col items-center gap-3 pt-6 text-center">
			<Avatar class="size-16">
				<AvatarFallback>SA</AvatarFallback>
			</Avatar>
			<div>
				<div class="font-semibold" style="font-family: var(--font-heading)">Sofia Alvarez</div>
				<div class="text-sm text-muted-foreground">Product Designer at Acme</div>
			</div>
			<p class="text-sm text-muted-foreground">Building calm, considered interfaces. Previously at Northwind.</p>
			<div class="flex items-center gap-6 text-sm">
				<div class="flex flex-col items-center">
					<span class="font-semibold">128</span>
					<span class="text-xs text-muted-foreground">Posts</span>
				</div>
				<div class="flex flex-col items-center">
					<span class="font-semibold">4,302</span>
					<span class="text-xs text-muted-foreground">Followers</span>
				</div>
				<div class="flex flex-col items-center">
					<span class="font-semibold">218</span>
					<span class="text-xs text-muted-foreground">Following</span>
				</div>
			</div>
			<Button class="w-full">Follow</Button>
		</CardContent>
	</Card>
}

// galleryOnboardingCard is a setup checklist: three done steps, two
// pending, plus a determinate progress bar in the same card (unlike
// galleryFeedbackCard, where the progress bar and the skeleton row are
// unrelated to each other).
component galleryOnboardingCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Get started</CardTitle>
			<CardDescription>3 of 5 steps complete.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-3">
			<Progress value={60}/>
			<ul class="flex flex-col gap-2 text-sm">
				<li class="flex items-center gap-2">
					<icon.CircleCheck class="size-4 text-primary"/>
					<span>Create your workspace</span>
				</li>
				<li class="flex items-center gap-2">
					<icon.CircleCheck class="size-4 text-primary"/>
					<span>Invite your team</span>
				</li>
				<li class="flex items-center gap-2">
					<icon.CircleCheck class="size-4 text-primary"/>
					<span>Connect a data source</span>
				</li>
				<li class="flex items-center gap-2 text-muted-foreground">
					<icon.Circle class="size-4"/>
					<span>Configure notifications</span>
				</li>
				<li class="flex items-center gap-2 text-muted-foreground">
					<icon.Circle class="size-4"/>
					<span>Invite a guest reviewer</span>
				</li>
			</ul>
		</CardContent>
		<CardFooter>
			<Button variant="outline" class="w-full">Continue setup</Button>
		</CardFooter>
	</Card>
}

// galleryProjectsEmptyCard is a second Empty composition (the components
// page already has one, for an inbox) — a distinct scenario and icon, both
// deliberate rather than a near-duplicate: an inbox empty state reads
// differently from a "nothing created yet" empty state in a real product,
// and this card is what proves the theme's Empty treatment holds up for
// both.
component galleryProjectsEmptyCard() {
	<Card>
		<CardContent>
			<Empty>
				<EmptyHeader>
					<EmptyMedia variant="icon">
						<icon.FolderOpen/>
					</EmptyMedia>
					<EmptyTitle>No projects yet</EmptyTitle>
					<EmptyDescription>Create your first project to start tracking work.</EmptyDescription>
				</EmptyHeader>
				<EmptyContent>
					<Button>New project</Button>
				</EmptyContent>
			</Empty>
		</CardContent>
	</Card>
}

// galleryVerifyCard is a two-factor verification card: InputOTP in its
// real product context (a security step), distinct from
// galleryControlsCard's bare "Verification code" demo on the components
// page.
component galleryVerifyCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Verify your identity</CardTitle>
			<CardDescription>Enter the 6-digit code we sent to +1 (415) •••-4471.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col items-start gap-3">
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
			<Button variant="link" class="h-auto p-0">Resend code</Button>
		</CardContent>
		<CardFooter>
			<Button class="w-full">Verify</Button>
		</CardFooter>
	</Card>
}

// galleryNotificationPrefsCard is a dense, all-Switch settings screen —
// distinct from gallerySettingsCard's mixed-control workspace settings on
// the components page: this is what a dedicated "notification
// preferences" screen looks like, grouped under plain section labels.
component galleryNotificationPrefsCard(idp string) {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Notification preferences</CardTitle>
			<CardDescription>Choose what you want to hear about.</CardDescription>
		</CardHeader>
		<CardContent class="flex flex-col gap-4">
			<div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">Email</div>
			<div class="flex items-center justify-between gap-4">
				<Label for={idp + "-notifprefs-comments"}>Comments on your work</Label>
				<Switch id={idp + "-notifprefs-comments"} checked/>
			</div>
			<div class="flex items-center justify-between gap-4">
				<Label for={idp + "-notifprefs-digest"}>Weekly digest</Label>
				<Switch id={idp + "-notifprefs-digest"} checked/>
			</div>
			<Separator/>
			<div class="text-xs font-medium text-muted-foreground uppercase tracking-wide">Push</div>
			<div class="flex items-center justify-between gap-4">
				<Label for={idp + "-notifprefs-mentions"}>Direct mentions</Label>
				<Switch id={idp + "-notifprefs-mentions"} checked/>
			</div>
			<div class="flex items-center justify-between gap-4">
				<Label for={idp + "-notifprefs-marketing"}>Product announcements</Label>
				<Switch id={idp + "-notifprefs-marketing"}/>
			</div>
		</CardContent>
		<CardFooter>
			<Button variant="outline" class="w-full">Save preferences</Button>
		</CardFooter>
	</Card>
}

// galleryBillingCard is an account billing summary — distinct from
// galleryTableCard's invoice line-item history on the components page:
// this is the "manage my subscription" screen, not a billing ledger.
component galleryBillingCard() {
	<Card>
		<CardHeader>
			<CardTitle style="font-family: var(--font-heading)">Billing</CardTitle>
			<CardDescription>Your current plan and payment method.</CardDescription>
			<CardAction>
				<Badge>Pro plan</Badge>
			</CardAction>
		</CardHeader>
		<CardContent class="flex flex-col gap-3">
			<div>
				<span class="text-xl font-semibold">$24</span>
				<span class="text-sm text-muted-foreground">/user/mo</span>
			</div>
			<div class="text-xs text-muted-foreground">Renews Feb 12, 2026</div>
			<Separator/>
			<div class="flex items-center justify-between gap-4">
				<div class="flex items-center gap-2 text-sm">
					<icon.CreditCard class="size-4 text-muted-foreground"/>
					<span>•••• 4242</span>
				</div>
				<Button variant="outline" size="sm">Update</Button>
			</div>
		</CardContent>
		<CardFooter>
			<Button variant="outline" class="w-full">Manage subscription</Button>
		</CardFooter>
	</Card>
}
