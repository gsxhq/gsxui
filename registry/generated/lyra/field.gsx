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
		class={ "gap-4 flex flex-col" }
		{ attrs... }
		data-gsxui-slot-field-set
	>
		{ children }
	</fieldset>
}

// data-variant is the public CSS axis for legend and label metrics.
component FieldLegend(variant string, children gsx.Node, attrs gsx.Attrs) {
	<legend
		data-variant={variant |> default("legend")}
		class={ "mb-2.5 font-medium", switch variant { case "label": "text-xs" default: "text-sm" } }
		{ attrs... }
		data-gsxui-slot-field-legend
	>
		{ children }
	</legend>
}

component FieldGroup(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={
			"group/field-group", "gap-5 [&>[data-gsxui-slot-field-group]]:gap-4 flex flex-col w-full @container/field-group"
		}
		{ attrs... }
		data-gsxui-slot-field-group
	>
		{ children }
	</div>
}

// data-orientation is the public CSS axis for layout and is also read by
// FieldDescription's relational text-balance rule.
component Field(orientation string, children gsx.Node, attrs gsx.Attrs) {
	<div
		role="group"
		data-orientation={orientation |> default("vertical")}
		class={
			"group/field",
			"data-[invalid=true]:text-destructive gap-2 flex w-full",
			switch orientation {
			case "horizontal":
				"flex-row items-center [&>[data-gsxui-slot-field-label]]:flex-auto has-[>[data-gsxui-slot-field-content]]:items-start has-[>[data-gsxui-slot-field-content]]:[&>:is([role=checkbox],[role=radio])]:mt-px"
			case "responsive":
				"flex-col [&>*]:w-full [&>[data-gsxui-slot-select-bridge]]:w-auto @min-[28rem]/field-group:flex-row @min-[28rem]/field-group:items-center @min-[28rem]/field-group:[&>*]:w-auto @min-[28rem]/field-group:[&>[data-gsxui-slot-field-label]]:flex-auto @min-[28rem]/field-group:has-[>[data-gsxui-slot-field-content]]:items-start @min-[28rem]/field-group:has-[>[data-gsxui-slot-field-content]]:[&>:is([role=checkbox],[role=radio])]:mt-px"
			default:
				"flex-col [&>*]:w-full [&>[data-gsxui-slot-select-bridge]]:w-auto"
			}
		}
		{ attrs... }
		data-gsxui-slot-field
	>
		{ children }
	</div>
}

component FieldContent(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-0.5 flex flex-col flex-1 leading-snug" }
		{ attrs... }
		data-gsxui-slot-field-content
	>
		{ children }
	</div>
}

// FieldLabel composes ui.Label directly, preserving ordered styling tokens
// "label field-label".
component FieldLabel(children gsx.Node, attrs gsx.Attrs) {
	<Label
		class={
			"has-data-checked:bg-primary/5 has-data-checked:border-primary/30 dark:has-data-checked:border-primary/20 dark:has-data-checked:bg-primary/10 gap-2 leading-snug has-[>[data-gsxui-slot-field]]:rounded-none has-[>[data-gsxui-slot-field]]:border flex w-fit has-[>[data-gsxui-slot-field]]:w-full has-[>[data-gsxui-slot-field]]:flex-col [&>[data-gsxui-slot-field]]:p-2"
		}
		{ attrs... }
		data-gsxui-slot-field-label
	>
		{ children }
	</Label>
}

// FieldTitle renders a <div> with a distinct token so themes can address it
// independently from the composed FieldLabel.
component FieldTitle(children gsx.Node, attrs gsx.Attrs) {
	<div
		class={ "gap-2 text-xs/relaxed flex w-fit items-center" }
		{ attrs... }
		data-gsxui-slot-field-title
	>
		{ children }
	</div>
}

component FieldDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		class={
			"text-muted-foreground text-left text-xs/relaxed [[data-variant=legend]+&]:-mt-1.5 leading-normal font-normal [&>a]:underline [&>a]:underline-offset-4 [&>a:hover]:text-primary"
		}
		{ attrs... }
		data-gsxui-slot-field-description
	>
		{ children }
	</p>
}

// FieldSeparator composes ui.Separator with ordered tokens
// "separator field-separator". The wrapper has its own token because it
// owns layout while the nested separator owns the rule. data-content is a
// presence marker: children emits it bare and no children omits it. The
// optional label span only renders when children is present.
component FieldSeparator(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-content={children != nil}
		class={ "-my-2 h-5 text-xs group-data-[variant=outline]/field-group:-mb-2 relative" }
		{ attrs... }
		data-gsxui-slot-field-separator-wrapper
	>
		<Separator class={ "absolute inset-0 top-1/2" } data-gsxui-slot-field-separator/>
		{ if children != nil {
			<span
				class={ "text-muted-foreground px-2 relative mx-auto block w-fit bg-background" }
				data-gsxui-slot-field-separator-content
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
		<div role="alert" class={ "text-destructive text-xs" } { attrs... } data-gsxui-slot-field-error>
			{ children }
		</div>
	} }
}
