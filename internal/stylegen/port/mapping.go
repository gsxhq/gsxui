package port

import "strings"

// ignoredSections are upstream MARK sections with no gsxui component
// counterpart at all: seven blocks belong to an AI-chat/composer surface
// gsxui doesn't ship (Bubble, Attachment, Marker, Questionnaire, Message
// Scroller, Message), one is a chart surface with no component or shape yet
// (Chart), and one is a translucent-menu variant (Menu Translucent). "React
// Aria" is a ninth: its 30 rules are all `-aria`-suffixed duplicates of
// rules that already live in their own component's section — Ignored's
// suffix check would catch every one of them individually, but excluding
// the whole section here means a stray non-`-aria` class landing in it by
// accident can never slip through as a fabricated 55th component.
var ignoredSections = map[string]bool{
	"React Aria":       true,
	"Chart":            true,
	"Menu Translucent": true,
	"Bubble":           true,
	"Attachment":       true,
	"Marker":           true,
	"Questionnaire":    true,
	"Message Scroller": true,
	"Message":          true,
}

// sectionComponentOverrides are declared exceptions to the default
// lowercase-and-hyphenate MARK-name derivation. Verified against the full
// 54-component survey (dossier §6): "Radio Group" is the ONLY override
// across all 54 — default derivation would produce "radio-group", but
// registry/canonical/shapes/radio.go names the shape "radio" (a single
// native <input type="radio"> root slot; "group" describes the surrounding
// ui.RadioGroup composition, not this shape).
var sectionComponentOverrides = map[string]string{
	"Radio Group": "radio",
}

// SectionComponent maps an upstream MARK name to a gsxui component name.
// Returns ok=false for sections we deliberately ignore: the synthetic ""
// section (rules appearing before any MARK comment — never a component) and
// the ignoredSections table above.
func SectionComponent(mark string) (component string, ok bool) {
	if mark == "" || ignoredSections[mark] {
		return "", false
	}
	if override, has := sectionComponentOverrides[mark]; has {
		return override, true
	}
	return strings.ToLower(strings.ReplaceAll(mark, " ", "-")), true
}

// Ignored reports classes skipped by policy: React-Aria's shadow rules.
// These duplicate a component's primary rules under React Aria's own
// data-* attribute vocabulary (data-selected instead of data-checked, and
// so on) rather than Tailwind's native data-[attr=value] arbitrary variant,
// so they carry no information the primary rule doesn't already have. They
// appear both inline within a component's own MARK section (e.g. Select's
// cn-select-value-aria, cn-select-item-aria, cn-select-empty-aria) and in
// the trailing "React Aria" section (already excluded by SectionComponent
// above) — matching by suffix catches both regardless of which section a
// rule physically lives in.
func Ignored(class string) bool {
	return strings.HasSuffix(class, "-aria")
}

// styleInvariantComponents take the same recipe in every style: no upstream
// MARK section exists for them in any of the 8 style-<name>.css files, so
// they stay hand-authored residue outside the porter's reach entirely (spec
// §4.4). Confirmed against dossier §6: these are the only 4 of the 54
// canonical shapes with no upstream section at all.
var styleInvariantComponents = map[string]bool{
	"aspect-ratio": true,
	"collapsible":  true,
	"spinner":      true,
	"toaster":      true,
}

// StyleInvariant reports whether component takes the same recipe in every
// style.
func StyleInvariant(component string) bool {
	return styleInvariantComponents[component]
}

// slotOverride is a declared exception to the default class→slot
// derivation: either a fixed target slot (ok=true) or an explicit "this
// upstream class has no gsxui slot" (ok=false, reported by the caller as
// unmapped, never silently dropped).
type slotOverride struct {
	slot string
	ok   bool
}

// slotOverrides holds per-component exceptions to the default class→slot
// derivation, keyed by [component][upstream class]. "At minimum" per the
// plan: select's native-popover scroll buttons have no gsxui slot at all
// (the browser handles overflow scrolling on its own — no JS needed), and
// the four native-<dialog>-backed compound components fuse upstream's
// separate overlay rule onto their content slot, because gsxui's <dialog>
// paints the backdrop as ::backdrop, a pseudo-element of dialog-content
// itself, not a separate node upstream's .cn-*-overlay class could ever
// land on.
var slotOverrides = map[string]map[string]slotOverride{
	"select": {
		// Native popover: no JS scroll buttons to style.
		"cn-select-scroll-up-button":   {ok: false},
		"cn-select-scroll-down-button": {ok: false},
	},
	"dialog": {
		// Native <dialog>'s ::backdrop is part of dialog-content; upstream's
		// separate overlay rule's utilities land on the content rule.
		"cn-dialog-overlay": {slot: "content", ok: true},
	},
	"alert-dialog": {
		"cn-alert-dialog-overlay": {slot: "content", ok: true},
	},
	"drawer": {
		"cn-drawer-overlay": {slot: "content", ok: true},
	},
	"sheet": {
		"cn-sheet-overlay": {slot: "content", ok: true},
	},
}

// upstreamPrefix returns the upstream cn-* class prefix (without "cn-") a
// gsxui component's own rules appear under. Every component shares its name
// with the upstream prefix except radio, whose section is "Radio Group" (see
// sectionComponentOverrides): its rules are cn-radio-group, cn-radio-group-
// item, and so on, never cn-radio-*.
func upstreamPrefix(component string) string {
	if component == "radio" {
		return "radio-group"
	}
	return component
}

// SlotFor maps an upstream class within a component to a gsxui slot name.
// Returns ok=false when the class has no gsxui slot (reported, not
// dropped): either by a declared override (e.g. select's scroll buttons) or
// because the class doesn't belong to this component's upstream prefix at
// all (an unrelated stray).
func SlotFor(component, class string) (slot string, ok bool) {
	if overrides, has := slotOverrides[component]; has {
		if override, matched := overrides[class]; matched {
			return override.slot, override.ok
		}
	}

	prefix := "cn-" + upstreamPrefix(component)
	if class == prefix {
		return "", true
	}
	suffix, cut := strings.CutPrefix(class, prefix+"-")
	if !cut {
		return "", false
	}
	return suffix, true
}
