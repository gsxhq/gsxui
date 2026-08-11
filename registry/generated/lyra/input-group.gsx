package ui

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
			"border-input dark:bg-input/30 has-[[data-gsxui-slot-input-group-control]:focus-visible]:border-ring has-[[data-gsxui-slot-input-group-control]:focus-visible]:ring-ring/50 has-[[aria-invalid=true]]:ring-destructive/20 has-[[aria-invalid=true]]:border-destructive dark:has-[[aria-invalid=true]]:ring-destructive/40 has-disabled:bg-input/50 dark:has-disabled:bg-input/80 h-8 rounded-none border transition-colors [[data-gsxui-slot-combobox-content]_&]:focus-within:border-inherit [[data-gsxui-slot-combobox-content]_&]:focus-within:ring-0 has-disabled:opacity-50 has-[[data-gsxui-slot-input-group-control]:focus-visible]:ring-1 has-[[aria-invalid=true]]:ring-1 has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>[data-align=block-end]]:[&>input]:pt-3 has-[>[data-align=block-start]]:[&>input]:pb-3 has-[>[data-align=inline-end]]:[&>input]:pr-1.5 has-[>[data-align=inline-start]]:[&>input]:pl-1.5 flex"
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
			"text-muted-foreground h-auto gap-2 py-1.5 text-xs font-medium group-data-[disabled=true]/input-group:opacity-50 [&>kbd]:rounded-none [&>svg:not([class*='size-'])]:size-4 flex",
			switch align {
			case "inline-end":
				"pr-2 has-[>button]:mr-[-0.3rem] has-[>kbd]:mr-[-0.15rem]"
			case "block-start":
				"px-2.5 pt-2 group-has-[>input]/input-group:pt-2 [.border-b]:pb-2"
			case "block-end":
				"px-2.5 pb-2 group-has-[>input]/input-group:pb-2 [.border-t]:pt-2"
			default:
				"pl-2 has-[>button]:ml-[-0.3rem] has-[>kbd]:ml-[-0.15rem]"
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
		class={ "text-muted-foreground gap-2 text-xs [&_svg:not([class*='size-'])]:size-4 flex" }
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
