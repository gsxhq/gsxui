package maia

import "github.com/gsxhq/gsx"

// InputGroup and its parts are the shadcn/ui InputGroup
// (registry/new-york-v4/ui/input-group.tsx) — no Radix primitive underneath;
// it's a plain styled `<div role="group">` wrapping an Input/Textarea plus
// leading/trailing "addon" content (icons, buttons, text). InputGroupInput
// and InputGroupTextarea compose ui.Input/ui.Textarea directly (flat
// package, no re-implementation) and InputGroupButton composes ui.Button —
// the input-group -> [button input textarea] dependency internal/registry
// derives from those calls and registry_test.go pins.
//
// GAP: InputGroupAddon's onClick handler (focuses the group's own <input> on
// a click that doesn't land on a nested <button>) is client JS with no
// equivalent here — zero client JS for this component, per the Tier 1 plan's
// Tech Stack constraint. Dropped; the addon still renders and styles
// identically, it just isn't click-to-focus.
//
// Retargeted to nova density (2026-07-24 nova density map, `## input-group`).
// InputGroupButton's xs/icon-xs radius arithmetic moves from
// rounded-[calc(var(--radius)-5px)] to rounded-[calc(var(--radius)-3px)],
// matching nova's smaller inset from the control radius; its sm size and
// InputGroupAddon's own [&>kbd]:rounded-[calc(var(--radius)-5px)] have no
// nova counterpart and are left unchanged (nova ships only xs/icon-xs/
// icon-sm for this button; the kbd radius is a separate, untouched token).
component InputGroup(children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		class={
			"relative flex h-8 min-w-0 w-full items-center rounded-lg border border-input transition-[color,box-shadow] outline-none dark:bg-input/30 has-[>[data-gsxui-slot-textarea][data-gsxui-slot-input-group-control]]:h-auto has-[>[data-gsxui-slot-input-group-addon][data-align=block-start]]:h-auto has-[>[data-gsxui-slot-input-group-addon][data-align=block-start]]:flex-col has-[>[data-gsxui-slot-input-group-addon][data-align=block-end]]:h-auto has-[>[data-gsxui-slot-input-group-addon][data-align=block-end]]:flex-col has-[[data-gsxui-slot-input-group-control]:focus-visible]:border-ring has-[[data-gsxui-slot-input-group-control]:focus-visible]:ring-[3px] has-[[data-gsxui-slot-input-group-control]:focus-visible]:ring-ring/50 has-[[aria-invalid=true]]:border-destructive has-[[aria-invalid=true]]:ring-destructive/20 dark:has-[[aria-invalid=true]]:ring-destructive/40"
		}
		{ attrs... }
		data-gsxui-slot-input-group
	>
		{ children }
	</div>
}

// data-align is the public CSS axis for the four addon placements and is
// also consumed by InputGroup's relational layout rules.
component InputGroupAddon(align string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-align={align |> default("inline-start")}
		class={
			"flex h-auto cursor-text items-center justify-center gap-2 py-1.5 text-sm font-medium text-muted-foreground select-none [&>kbd]:rounded-[calc(var(--radius)-5px)] [&>svg:not([class*='size-'])]:size-4",
			switch align {
			case "inline-end":
				"order-last pr-2"
			case "block-start":
				"order-first w-full justify-start px-2.5 pt-2 [[data-gsxui-slot-input-group]:has(>[data-gsxui-slot-input])>&]:pt-2 [&.border-b]:pb-2"
			case "block-end":
				"order-last w-full justify-start px-2.5 pb-2 [[data-gsxui-slot-input-group]:has(>[data-gsxui-slot-input])>&]:pb-2 [&.border-t]:pt-2"
			default:
				"order-first pl-2"
			}
		}
		{ attrs... }
		data-gsxui-slot-input-group-addon
	>
		{ children }
	</div>
}

// InputGroupButton composes ui.Button — the input-group -> button
// dependency. It forwards both public styling axes into Button, defaulting
// size to xs and variant to ghost, so Button's rendered data axes and concrete
// generated variant/size classes always describe the same state.
component InputGroupButton(variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	<Button
		size={size |> default("xs")}
		variant={variant |> default("ghost")}
		{ attrs... }
		data-gsxui-slot-input-group-button
	>
		{ children }
	</Button>
}

// InputGroupText has its own theme token.
component InputGroupText(children gsx.Node, attrs gsx.Attrs) {
	<span
		class={
			"flex items-center gap-2 text-sm text-muted-foreground [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
		}
		{ attrs... }
		data-gsxui-slot-input-group-text
	>
		{ children }
	</span>
}

// InputGroupInput composes ordered tokens "input input-group-control";
// InputGroup keys focus and invalid relations off the latter.
component InputGroupInput(attrs gsx.Attrs) {
	<Input
		class={ "flex-1 bg-transparent shadow-none focus-visible:ring-0 dark:bg-transparent" }
		{ attrs... }
		data-gsxui-slot-input-group-control
	/>
}

// InputGroupTextarea composes ui.Textarea directly (flat package), forwarding
// `value` into Textarea's own `value` param (Textarea's text-child ADAPT,
// see ui/textarea.gsx) the same way ItemSeparator forwards `orientation`
// into Separator's own param.
component InputGroupTextarea(value string, attrs gsx.Attrs) {
	<Textarea
		value={value}
		class={ "flex-1 bg-transparent shadow-none focus-visible:ring-0 dark:bg-transparent" }
		{ attrs... }
		data-gsxui-slot-input-group-control
	/>
}
