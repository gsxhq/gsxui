package ui

import "github.com/gsxhq/gsx"

// FieldSet, FieldLegend, FieldGroup, Field, FieldContent, FieldLabel,
// FieldTitle, FieldDescription, FieldSeparator, and FieldError are the
// shadcn/ui Field family (registry/new-york-v4/ui/field.tsx) — no Radix
// primitive underneath; every part is already a plain styled element.
// FieldLabel composes ui.Label and FieldSeparator composes ui.Separator
// directly (flat package, no re-implementation) — the field -> [label
// separator] dependency internal/registry derives from those calls and
// registry_test.go pins.
//
// ADAPT: FieldError's react-hook-form `errors` prop (an
// `Array<{message?: string} | undefined>`, deduplicated and rendered as a
// single message or a `<ul>` of messages via a useMemo) is not ported —
// there is no react-hook-form in a server-rendered gsx tree to produce that
// shape. FieldError keeps only shadcn's other content path: plain
// `children`. A caller with more than one message renders its own `<ul>`
// child (the same markup shadcn's own multi-error branch would have
// produced) — no functionality is lost, only the automatic
// errors-array-to-list plumbing.
component FieldSet(children gsx.Node, attrs gsx.Attrs) {
	<fieldset
		data-slot="field-set"
		class="flex flex-col gap-6 has-[>[data-slot=checkbox-group]]:gap-3 has-[>[data-slot=radio-group]]:gap-3"
		{ attrs... }
	>
		{ children }
	</fieldset>
}

// FieldLegend's variant cva only ever types `variant` — both
// `data-[variant=legend]:text-base` and `data-[variant=label]:text-sm` are
// present unconditionally in shadcn's own class string; data-variant plus
// Tailwind's data-[variant=...] selectors are what actually pick one, the
// same "no switch needed, single verbatim class string dispatches on the
// stamped attribute" shape as Separator's data-orientation (see
// ui/separator.gsx), not a static-block switch like item/button-group's cva
// maps.
component FieldLegend(variant string, children gsx.Node, attrs gsx.Attrs) {
	<legend
		data-slot="field-legend"
		data-variant={variant |> default("legend")}
		class="mb-3 font-medium data-[variant=legend]:text-base data-[variant=label]:text-sm"
		{ attrs... }
	>
		{ children }
	</legend>
}

component FieldGroup(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="field-group"
		class="group/field-group @container/field-group flex w-full flex-col gap-7 data-[slot=checkbox-group]:gap-3 [&>[data-slot=field-group]]:gap-4"
		{ attrs... }
	>
		{ children }
	</div>
}

// Field's fieldVariants cva map picks between three static class blocks by
// the JS-resolved orientation value (WIN: switch inside class={}, same idiom
// as item/button-group/empty) — but UNLIKE button-group's orientation
// (nothing downstream ever reads button-group's data-orientation), Field's
// own data-orientation IS read by a sibling: FieldDescription's
// `group-has-[[data-orientation=horizontal]]/field:text-balance` selector
// keys directly off it. So both apply here: the switch picks Field's own
// class blocks, AND data-orientation is stamped (not merely for
// stamp-everything consistency, but because it's structurally load-bearing
// for FieldDescription).
component Field(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-slot="field"
		data-orientation={orientation |> default("vertical")}
		class={
			"group/field flex w-full gap-3 data-[invalid=true]:text-destructive",
			switch orientation {
			case "horizontal":
				"flex-row items-center [&>[data-slot=field-label]]:flex-auto has-[>[data-slot=field-content]]:items-start has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px"
			case "responsive":
				"flex-col @md/field-group:flex-row @md/field-group:items-center [&>*]:w-full @md/field-group:[&>*]:w-auto [&>.sr-only]:w-auto @md/field-group:[&>[data-slot=field-label]]:flex-auto @md/field-group:has-[>[data-slot=field-content]]:items-start @md/field-group:has-[>[data-slot=field-content]]:[&>[role=checkbox],[role=radio]]:mt-px"
			default:
				"flex-col [&>*]:w-full [&>.sr-only]:w-auto"
			}
		}
		{ attrs... }
	>
		{ children }
	</div>
}

component FieldContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="field-content"
		class="group/field-content flex flex-1 flex-col gap-1.5 leading-snug"
		{ attrs... }
	>
		{ children }
	</div>
}

// FieldLabel composes ui.Label directly (flat package) — the field -> label
// dependency. data-slot is overridden from Label's own "label" to
// "field-label" as an explicit non-parameter attribute on the call, the
// same override mechanism as ItemSeparator/ButtonGroupSeparator overriding
// Separator's data-slot (see ui/item.gsx, ui/button-group.gsx).
component FieldLabel(children gsx.Node, attrs gsx.Attrs) {
	<Label
		data-slot="field-label"
		class="group/field-label peer/field-label flex w-fit gap-2 leading-snug group-data-[disabled=true]/field:opacity-50 has-[>[data-slot=field]]:w-full has-[>[data-slot=field]]:flex-col has-[>[data-slot=field]]:rounded-md has-[>[data-slot=field]]:border [&>*]:data-[slot=field]:p-4 has-data-[state=checked]:border-primary has-data-[state=checked]:bg-primary/5 dark:has-data-[state=checked]:bg-primary/10"
		{ attrs... }
	>
		{ children }
	</Label>
}

// FieldTitle renders a <div>, matching shadcn's own actual returned element
// (its TypeScript prop type reads React.ComponentProps<"div">, consistent
// with the tag). Its data-slot is "field-label" — the SAME value FieldLabel
// itself stamps — in shadcn's own source; ported as-is (token-for-token),
// the same unmatched/shared-data-slot call as EmptyMedia's "empty-icon" (see
// ## empty) and ButtonGroupText's missing data-slot (see ## button-group).
component FieldTitle(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="field-label"
		class="flex w-fit items-center gap-2 text-sm leading-snug font-medium group-data-[disabled=true]/field:opacity-50"
		{ attrs... }
	>
		{ children }
	</div>
}

component FieldDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		data-slot="field-description"
		class="text-sm leading-normal font-normal text-muted-foreground group-has-[[data-orientation=horizontal]]/field:text-balance last:mt-0 nth-last-2:-mt-1 [[data-variant=legend]+&]:-mt-1.5 [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary"
		{ attrs... }
	>
		{ children }
	</p>
}

// FieldSeparator composes ui.Separator directly (flat package) — the field
// -> separator dependency. data-content mirrors shadcn's `data-content={!!
// children}` boolean stamp (gsx renders a bool expression as "true"/"false"
// text directly, the same mechanism as pagination.gsx's data-active — see
// ui/pagination.gsx); the optional label span only renders when children is
// present.
component FieldSeparator(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-slot="field-separator"
		data-content={children != nil}
		class="relative -my-2 h-5 text-sm group-data-[variant=outline]/field-group:-mb-2"
		{ attrs... }
	>
		<Separator class="absolute inset-0 top-1/2"/>
		{ if children != nil {
			<span
				class="relative mx-auto block w-fit bg-background px-2 text-muted-foreground"
				data-slot="field-separator-content"
			>
				{ children }
			</span>
		} }
	</div>
}

// FieldError renders nothing when children is nil — the gsx equivalent of
// shadcn's `if (!content) return null`, now driven by children alone (see
// the file-level ADAPT comment above for the dropped errors prop).
component FieldError(children gsx.Node, attrs gsx.Attrs) {
	{ if children != nil {
		<div role="alert" data-slot="field-error" class="text-sm font-normal text-destructive" { attrs... }>
			{ children }
		</div>
	} }
}
