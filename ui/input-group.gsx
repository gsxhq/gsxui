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
		data-slot="input-group"
		role="group"
		class={
			"group/input-group relative flex w-full items-center rounded-lg border border-input transition-[color,box-shadow] outline-none dark:bg-input/30",
			"h-8 min-w-0 has-[>textarea]:h-auto",
			"has-[>[data-align=inline-start]]:[&>input]:pl-1.5 has-[>[data-align=inline-end]]:[&>input]:pr-1.5 has-[>[data-align=block-start]]:h-auto has-[>[data-align=block-start]]:flex-col has-[>[data-align=block-start]]:[&>input]:pb-3 has-[>[data-align=block-end]]:h-auto has-[>[data-align=block-end]]:flex-col has-[>[data-align=block-end]]:[&>input]:pt-3",
			"has-[[data-slot=input-group-control]:focus-visible]:border-ring has-[[data-slot=input-group-control]:focus-visible]:ring-[3px] has-[[data-slot=input-group-control]:focus-visible]:ring-ring/50",
			"has-[[data-slot][aria-invalid=true]]:border-destructive has-[[data-slot][aria-invalid=true]]:ring-destructive/20 dark:has-[[data-slot][aria-invalid=true]]:ring-destructive/40"
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// InputGroupAddon's inputGroupAddonVariants cva map picks between four
// static class blocks by the JS-resolved align value — no
// data-[align=...]: selectors in the source to preserve for THIS class
// string (contrast InputGroup's own class above, which does key off
// data-align on the addon child) — so it ports as a switch inside class={},
// the same idiom as item/button-group. data-align is still stamped (WIN
// pattern, |> default) since InputGroup's own has-[>[data-align=...]] rules
// depend on it.
component InputGroupAddon(align string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-slot="input-group-addon"
		data-align={align |> default("inline-start")}
		class={
			"flex h-auto cursor-text items-center justify-center gap-2 py-1.5 text-sm font-medium text-muted-foreground select-none group-data-[disabled=true]/input-group:opacity-50 [&>kbd]:rounded-[calc(var(--radius)-5px)] [&>svg:not([class*='size-'])]:size-4",
			switch align {
			case "inline-end":
				"order-last pr-2 has-[>button]:mr-[-0.3rem] has-[>kbd]:mr-[-0.15rem]"
			case "block-start":
				"order-first w-full justify-start px-2.5 pt-2 group-has-[>input]/input-group:pt-2 [.border-b]:pb-2"
			case "block-end":
				"order-last w-full justify-start px-2.5 pb-2 group-has-[>input]/input-group:pb-2 [.border-t]:pt-2"
			default:
				"order-first pl-2 has-[>button]:ml-[-0.3rem] has-[>kbd]:ml-[-0.15rem]"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}

// InputGroupButton composes ui.Button — the input-group -> button
// dependency. shadcn's own version deliberately does NOT pass its `size`
// prop through to the underlying Button's own `size` (its type is
// `Omit<ComponentProps<typeof Button>, "size">`): Button renders with its
// own default size classes, and inputGroupButtonVariants({size})'s classes
// are merged on top of them by cn(), so tailwind-merge — not a size prop —
// is what actually resolves the visible height/padding/rounding. This port
// mirrors that exactly: `size` is never forwarded to Button's own `size`
// param, only used for the overlay switch and the data-size stamp. `variant`
// defaults to "ghost" (Button's own zero-value default is "default"/primary)
// and IS forwarded to Button's own `variant` param, matching shadcn's
// `variant = "ghost"` passthrough. data-size is set as an explicit
// non-parameter attribute on the `<Button>` call, so it lands in Button's
// own attrs bag and overrides Button's internal `data-size={size}` stamp
// (which would otherwise read Button's own, unset, size param) — the same
// competing-defaults override mechanism as ItemSeparator/ButtonGroupSeparator
// overriding Separator's data-slot (see ui/item.gsx, ui/button-group.gsx).
component InputGroupButton(variant string, size string, children gsx.Node, attrs gsx.Attrs) {
	<Button
		data-size={size |> default("xs")}
		variant={variant |> default("ghost")}
		class={
			"flex items-center gap-2 text-sm shadow-none",
			switch size {
			case "sm":
				"h-8 gap-1.5 rounded-md px-2.5 has-[>svg]:px-2.5"
			case "icon-xs":
				"size-6 rounded-[calc(var(--radius)-3px)] p-0 has-[>svg]:p-0"
			case "icon-sm":
				"size-8 p-0 has-[>svg]:p-0"
			default:
				"h-6 gap-1 rounded-[calc(var(--radius)-3px)] px-1.5 has-[>svg]:px-1.5 [&>svg:not([class*='size-'])]:size-3.5"
			}
		}
		{ attrs... }
	>
		{ children }
	</Button>
}

// InputGroupText carries no data-slot in shadcn's own source (unlike every
// other input-group part) — ported as-is, the same unmatched-data-slot call
// as ButtonGroupText (see ui/button-group.gsx).
component InputGroupText(children gsx.Node, attrs gsx.Attrs) {
	<span
		class="flex items-center gap-2 text-sm text-muted-foreground [&_svg]:pointer-events-none [&_svg:not([class*='size-'])]:size-4"
		{ attrs... }
	>
		{ children }
	</span>
}

// InputGroupInput composes ui.Input directly (flat package). data-slot is
// overridden from Input's own "input" to "input-group-control" as an
// explicit non-parameter attribute on the call — same override mechanism as
// InputGroupButton's data-size above — so InputGroup's own
// has-[[data-slot=input-group-control]:focus-visible] rules can key off it.
component InputGroupInput(attrs gsx.Attrs) {
	<Input
		data-slot="input-group-control"
		class="flex-1 rounded-none border-0 bg-transparent shadow-none focus-visible:ring-0 dark:bg-transparent"
		{ attrs... }
	/>
}

// InputGroupTextarea composes ui.Textarea directly (flat package), forwarding
// `value` into Textarea's own `value` param (Textarea's text-child ADAPT,
// see ui/textarea.gsx) the same way ItemSeparator forwards `orientation`
// into Separator's own param.
component InputGroupTextarea(value string, attrs gsx.Attrs) {
	<Textarea
		value={value}
		data-slot="input-group-control"
		class="flex-1 resize-none rounded-none border-0 bg-transparent py-2 shadow-none focus-visible:ring-0 dark:bg-transparent"
		{ attrs... }
	/>
}
