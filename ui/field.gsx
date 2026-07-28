package ui

import (
	"strconv"

	"github.com/gsxhq/gsx"
)

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
		{ attrs... }
		data-gsxui-slot-field-legend
	>
		{ children }
	</legend>
}

component FieldGroup(children gsx.Node, attrs gsx.Attrs) {
	<div
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
		{ attrs... }
		data-gsxui-slot-field
	>
		{ children }
	</div>
}

component FieldContent(children gsx.Node, attrs gsx.Attrs) {
	<div
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
		{ attrs... }
		data-gsxui-slot-field-title
	>
		{ children }
	</div>
}

component FieldDescription(children gsx.Node, attrs gsx.Attrs) {
	<p
		{ attrs... }
		data-gsxui-slot-field-description
	>
		{ children }
	</p>
}

// FieldSeparator composes ui.Separator with ordered tokens
// "separator field-separator". The wrapper has its own token because it
// owns layout while the nested separator owns the rule. data-content mirrors
// shadcn's `data-content={!!
// children}` boolean stamp, using strconv.FormatBool to retain its valued
// "true"/"false" contract. The optional label span only renders when children
// is present.
component FieldSeparator(children gsx.Node, attrs gsx.Attrs) {
	<div
		data-content={strconv.FormatBool(children != nil)}
		{ attrs... }
		data-gsxui-slot-field-separator-wrapper
	>
		<Separator data-gsxui-slot-field-separator/>
		{ if children != nil {
			<span
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
		<div role="alert" { attrs... } data-gsxui-slot-field-error>
			{ children }
		</div>
	} }
}
