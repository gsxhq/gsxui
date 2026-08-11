package port

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// titleOverrides corrects componentTitle's naive per-word capitalization for
// the handful of components whose upstream MARK name isn't just Title Case
// (an acronym, mainly) — cosmetic only: comments are invisible to
// recipe.ParseStyle (dossier §1.1), so this affects provenance-header
// readability, never validation.
var titleOverrides = map[string]string{
	"input-otp": "Input OTP",
}

// componentTitle renders a gsxui component name for the provenance header,
// e.g. "alert-dialog" -> "Alert Dialog". It is derived from Ported.Component
// alone (Render's signature carries no separate display name), so it will
// not always reproduce the exact upstream MARK name a mapping-table override
// changed (radio's section was "Radio Group"; its title here is "Radio") —
// that is a header-prose nicety, not a provenance fact Render can lose: the
// SHA/path/line-range triple is the part that matters, and it survives
// intact.
func componentTitle(component string) string {
	if title, ok := titleOverrides[component]; ok {
		return title
	}
	words := strings.Split(component, "-")
	for i, w := range words {
		if w == "" {
			continue
		}
		words[i] = strings.ToUpper(w[:1]) + w[1:]
	}
	return strings.Join(words, " ")
}

// wrapWords greedily wraps s into lines no longer than width (best-effort —
// a single word longer than width still gets its own line), preserving word
// order and single-space separation.
func wrapWords(s string, width int) []string {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return nil
	}
	lines := []string{fields[0]}
	for _, word := range fields[1:] {
		last := lines[len(lines)-1]
		if len(last)+1+len(word) > width {
			lines = append(lines, word)
			continue
		}
		lines[len(lines)-1] = last + " " + word
	}
	return lines
}

// renderHeader builds the leading `/* Source: ... */` provenance comment per
// dossier §1.1's grammar: a fixed "Source: <repo>@<sha>" first line, then the
// upstream path, component title and exact line range the section came from,
// then a terse sentence naming the porter (the dossier's own guidance: the
// porter has no hand-authoring judgment to narrate, so its headers stay
// closer to button.css's terse two-line form than accordion/select's prose).
//
// SrcStart == SrcEnd == 0 marks a StyleInvariant component (mapping.go's
// StyleInvariant): no upstream MARK section exists for it at all, so Transform
// ran against a synthetic empty Section and every rule is 100% carried from
// nova (Ported.Carried, marked per-rule below) — the header says so instead
// of citing a nonsensical "lines 0-0".
func renderHeader(p Ported, upstreamSHA, upstreamPath string) []byte {
	var sentence string
	if p.SrcStart == 0 && p.SrcEnd == 0 {
		sentence = fmt.Sprintf(
			"%s has no %s section — this component is style-invariant. "+
				"Ported by internal/stylegen/port; every rule is carried from nova "+
				"(see the per-rule comments below).",
			upstreamPath, componentTitle(p.Component))
	} else {
		sentence = fmt.Sprintf(
			"%s, %s lines %d-%d. Ported by internal/stylegen/port; utilities are "+
				"upstream Tailwind classes, re-slotted onto gsxui's shape-declared "+
				"classes.",
			upstreamPath, componentTitle(p.Component), p.SrcStart, p.SrcEnd)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "/* Source: shadcn-ui/ui@%s\n", upstreamSHA)
	lines := wrapWords(sentence, 76)
	for i, line := range lines {
		buf.WriteString("   ")
		buf.WriteString(line)
		if i == len(lines)-1 {
			buf.WriteString(" */")
		}
		buf.WriteByte('\n')
	}
	return buf.Bytes()
}

// Render produces the recipe file bytes (header + @layer components block)
// for one ported component. Classes are composed only via
// shape.BaseClass/ValueClass (never string concatenation); rules are emitted
// in shape-declaration order (each Base:true slot, then each of its
// dimensions' values in declared order — the same nesting Conform and
// CheckConflicts walk); a carried rule (Ported.Carried) is preceded by a
// `/* carried: no upstream counterpart */` comment line.
func Render(p Ported, style, upstreamSHA, upstreamPath string) ([]byte, error) {
	shape, ok := shapes.All()[p.Component]
	if !ok {
		return nil, fmt.Errorf("render %s: no declared shape for style %s", p.Component, style)
	}

	var buf bytes.Buffer
	buf.Write(renderHeader(p, upstreamSHA, upstreamPath))
	buf.WriteString("@layer components {\n")

	first := true
	writeRule := func(class string, utilities []string) error {
		if len(utilities) == 0 {
			return fmt.Errorf("render %s: class %s has no utilities", p.Component, class)
		}
		if !first {
			buf.WriteByte('\n')
		}
		first = false
		if p.Carried[class] {
			buf.WriteString("  /* carried: no upstream counterpart */\n")
		}
		fmt.Fprintf(&buf, "  .%s {\n", class)
		fmt.Fprintf(&buf, "    @apply %s;\n", strings.Join(utilities, " "))
		buf.WriteString("  }\n")
		return nil
	}

	for _, slot := range shape.Slots {
		if slot.Base {
			if err := writeRule(shape.BaseClass(slot.Name), p.Base[slot.Name]); err != nil {
				return nil, err
			}
		}
		for _, dimension := range slot.Dimensions {
			for _, value := range dimension.Values {
				class := shape.ValueClass(slot.Name, dimension.Name, value)
				if err := writeRule(class, p.Values[slot.Name][dimension.Name][value]); err != nil {
					return nil, err
				}
			}
		}
	}

	buf.WriteString("}\n")
	return buf.Bytes(), nil
}

// Validate runs the repo's own three gates on rendered recipe bytes — the
// same trio internal/stylegen/generate.go's resolveAll runs before writing
// any generated artifact (dossier §3, §4.2): parse the recipe CSS, check it
// conforms exactly to shape (bidirectional coverage, no extra or missing
// classes), then check the Tailwind-merged utility lists carry no conflict
// or repetition.
func Validate(path string, shape recipe.Shape, src []byte) error {
	parsed, err := recipe.ParseStyle(path, src)
	if err != nil {
		return err
	}
	resolved, err := recipe.Conform(path, shape, parsed)
	if err != nil {
		return err
	}
	return recipe.CheckConflicts(path, resolved, merge.Merge)
}
