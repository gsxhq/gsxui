package sera

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
			"group/input-group",
			"border-transparent border-b-input bg-transparent has-[[data-gsxui-slot-input-group-control]:focus-visible]:border-b-ring has-[[aria-invalid=true]]:border-b-destructive dark:has-[[aria-invalid=true]]:border-b-destructive/50 h-10 rounded-none border transition-[color,border-color] [[data-gsxui-slot-combobox-content]_&]:focus-within:border-inherit [[data-gsxui-slot-combobox-content]_&]:focus-within:ring-0 has-data-[align=block-end]:rounded-none has-data-[align=block-start]:rounded-none has-[textarea]:rounded-none has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>[data-align=block-end]]:[&>input]:pt-3 has-[>[data-align=block-start]]:[&>input]:pb-3 flex"
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
			"text-muted-foreground [&_[data-gsxui-slot-kbd]]:bg-muted-foreground/10 h-auto gap-2 py-2 text-sm font-medium group-data-[disabled=true]/input-group:opacity-50 [&_[data-gsxui-slot-kbd]]:rounded-none [&_[data-gsxui-slot-kbd]]:px-1.5 [&>svg:not([class*='size-'])]:size-3.5 flex",
			switch align {
			case "inline-end":
				"order-last pe-2"
			case "block-start":
				"pt-3 group-has-[>input]/input-group:pt-3.5 [.border-b]:pb-3.5"
			case "block-end":
				"pb-3 group-has-[>input]/input-group:pb-3.5 [.border-t]:pt-3.5"
			default:
				"order-first ps-2"
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
		class={ "text-muted-foreground gap-2 text-sm [&_svg:not([class*='size-'])]:size-3.5 flex" }
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
