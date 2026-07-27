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
		{ withSlot("input-group", attrs)... }
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
		{ withSlot("input-group-addon", attrs)... }
	>
		{ children }
	</div>
}

// InputGroupButton composes ui.Button — the input-group -> button
// dependency. Its styling token composes after Button's token. `size` is
// intentionally an InputGroupButton axis rather than Button's size param;
// `variant` still forwards to Button and defaults to ghost.
component InputGroupButton(variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	<Button
		data-size={size |> default("xs")}
		variant={variant |> default("ghost")}
		{ withSlot("input-group-button", attrs)... }
	>
		{ children }
	</Button>
}

// InputGroupText has its own theme token.
component InputGroupText(children gsx.Node, attrs gsx.Attrs) {
	<span
		{ withSlot("input-group-text", attrs)... }
	>
		{ children }
	</span>
}

// InputGroupInput composes ordered tokens "input input-group-control";
// InputGroup keys focus and invalid relations off the latter.
component InputGroupInput(attrs gsx.Attrs) {
	<Input
		{ withSlot("input-group-control", attrs)... }
	/>
}

// InputGroupTextarea composes ui.Textarea directly (flat package), forwarding
// `value` into Textarea's own `value` param (Textarea's text-child ADAPT,
// see ui/textarea.gsx) the same way ItemSeparator forwards `orientation`
// into Separator's own param.
component InputGroupTextarea(value string, attrs gsx.Attrs) {
	<Textarea
		value={value}
		{ withSlot("input-group-control", attrs)... }
	/>
}
