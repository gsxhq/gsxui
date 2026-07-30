package stylegen

import (
	"fmt"
	"go/token"
	"go/types"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gsxhq/gsxui/internal/recipe"
)

// componentIdentifier derives the Go identifier a component's generated code
// and its .gsx accessor receiver use, from the component's registry name.
//
// The registry name is the CSS identity — it appears in every
// data-gsxui-slot-* marker, every recipe class, every file name and the
// parity docs — and eleven of the fifty-two components spell it with a hyphen.
// A hyphen is not a Go identifier: "input-group" + "Recipe" is the type
// `input-groupRecipe`, which format.Source rejects, and `input-group.Root()`
// inside a class expression parses as SUBTRACTION rather than a method call.
// So the registry name stays put and the Go side is DERIVED from it.
//
// The derivation is exactly the one accessorName already applies to slots
// ("menu-button" -> "MenuButton"), lowered at the initial because both the
// generated type and the package-level receiver variable are unexported:
//
//	registry / CSS identity   input-group
//	Go identifier             inputGroup   (type inputGroupRecipe)
//
// "button" has no hyphen, so it derives "button" and nothing about the
// already-migrated component moves.
//
// The same "registry stays put, Go side is derived" answer applies when the
// derived name is a Go keyword ("switch", "select") or a predeclared
// identifier ("len", "new"): rather than reject it, a single trailing
// underscore is appended ("switch" -> "switch_"), Go's own convention for
// exactly this collision. The registry/CSS identity is untouched — the
// generated CSS class is still gsxui-recipe-switch, only the Go receiver
// and type name pick up the suffix.
//
// Derivation is neither total nor injective, so every remaining failure mode
// (empty, digit-initial, not-an-identifier after suffixing) is an error
// naming the component rather than a panic or a silent fallback. Injectivity
// is checkComponentIdentifiers' job; legality is this function's.
func componentIdentifier(component string) (string, error) {
	if component == "" {
		return "", fmt.Errorf("component name is empty: nothing to derive a Go identifier from")
	}
	derived := lowerInitial(accessorName(component))
	if derived == "" {
		return "", fmt.Errorf(
			"component %q derives an empty Go identifier; it must contain at least one identifier character",
			component)
	}
	if first, _ := utf8.DecodeRuneInString(derived); unicode.IsDigit(first) {
		return "", fmt.Errorf(
			"component %q derives the Go identifier %q, which starts with a digit; "+
				"a component name's first segment must not start with a digit",
			component, derived)
	}
	// A keyword or predeclared identifier is not illegal to derive — the
	// registry/CSS identity ("switch", "select") is public surface (every
	// data-gsxui-slot-* marker, every recipe class, every file name) and
	// stays put exactly the way a hyphenated name does; only the GO side is
	// derived. So this collision gets the same answer hyphenation already
	// got: extend the derivation instead of renaming the component. A single
	// trailing underscore is Go's own convention for exactly this situation
	// (protoc-gen-go does the same for a field named e.g. "type"), and it's
	// deterministic — no per-component override map to keep in sync.
	if token.IsKeyword(derived) || types.Universe.Lookup(derived) != nil {
		derived += "_"
	}
	if !token.IsIdentifier(derived) {
		return "", fmt.Errorf(
			"component %q derives %q, which is not a valid Go identifier; "+
				"a component name must be made of letters, digits and hyphens",
			component, derived)
	}
	return derived, nil
}

// componentRecipeTypeName is the generated accessor type for one component.
func componentRecipeTypeName(component string) (string, error) {
	identifier, err := componentIdentifier(component)
	if err != nil {
		return "", err
	}
	return identifier + "Recipe", nil
}

// lowerInitial lowers the first rune. accessorName exports its result because
// slots become method names; a component becomes an unexported type and an
// unexported package-level variable instead.
func lowerInitial(name string) string {
	if name == "" {
		return ""
	}
	first, width := utf8.DecodeRuneInString(name)
	return string(unicode.ToLower(first)) + name[width:]
}

// identifierFold is the form two registry names must not share. Derivation
// erases the separators it title-cases over and the case distinction it
// creates, so "input-group", "input_group" and "inputGroup" are three spellings
// that reach the same Go identifier — or near enough that having both in the
// registry is a mistake worth failing on rather than a naming choice. Folding
// on the erased characters catches the pairs a plain string comparison of the
// derived identifiers misses.
func identifierFold(component string) string {
	return strings.ToLower(strings.NewReplacer("-", "", "_", "").Replace(component))
}

// checkComponentIdentifiers requires every component in the set to derive a
// legal Go identifier, and requires no two of them to derive the same one.
//
// A collision is not a cosmetic problem: two components whose recipe types are
// both `inputGroupRecipe` collide in package canonical and fail the build after
// the files are written, and — the quieter half — an accessor receiver written
// in either component's .gsx binds against whichever shape is being resolved,
// so a call meant for one component desugars against the other's utilities
// without anything complaining. This runs before either can happen.
func checkComponentIdentifiers(components []string) error {
	byFold := map[string]string{}
	sorted := append([]string(nil), components...)
	sort.Strings(sorted)
	for _, component := range sorted {
		if _, err := componentIdentifier(component); err != nil {
			return err
		}
		fold := identifierFold(component)
		if previous, taken := byFold[fold]; taken {
			previousID, _ := componentIdentifier(previous)
			currentID, _ := componentIdentifier(component)
			return fmt.Errorf(
				"components %q and %q both normalize to the Go identifier %q/%q — "+
					"their generated recipe types would collide in package canonical and an accessor "+
					"receiver could not say which one it means; rename one of them",
				previous, component, previousID, currentID)
		}
		byFold[fold] = component
	}
	return nil
}

// checkComponentNaming runs every naming gate a shape must pass before its
// accessors are generated or any accessor call is resolved against it: the
// component's own derived identifier, its uniqueness across the registry, and
// the shape's per-slot accessor names.
//
// It runs at all three entry points — GenerateAccessors, Resolve and
// HelperCalls — for the reason checkAccessorNames already does: resolution
// binds accessor calls to slots, so a name that codegen would have rejected
// must not be able to bind first.
func checkComponentNaming(shape recipe.Shape) error {
	if err := checkComponentIdentifiers(knownComponents(shape)); err != nil {
		return err
	}
	return checkAccessorNames(shape)
}
