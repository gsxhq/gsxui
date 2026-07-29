package stylegen

import (
	"bytes"
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// Call is one recipe accessor call found in a canonical class expression.
// Slot is "" for the component's root element, and Dimension is empty for a
// slot's base accessor.
type Call struct {
	Component string
	Slot      string
	Dimension string
}

type literalEdit struct {
	start int
	end   int
	value string
}

type inspectionMode uint8

const (
	canonicalClassInspection inspectionMode = iota
	validationOnlyInspection
)

// classMode selects how a canonical class expression is treated. The two
// helper modes read the structural `<component>.<Method>(...)` form; the token
// mode only scans string literals for recipe tokens, so that RecipeTokens can
// report and reject stray token use without resolving anything.
type classMode uint8

const (
	tokenScanMode classMode = iota
	helperDesugarMode
	helperScanMode
)

type resolver struct {
	filename  string
	src       []byte
	fset      *token.FileSet
	classMode classMode
	resolved  recipe.Resolved
	shape     recipe.Shape
	used      map[string]struct{}
	calls     []Call
	edits     []literalEdit
	err       error
}

// Resolve rewrites a canonical component source into consumer source for one
// style: every generated slot accessor call in a class expression — a slot's
// base accessor or one of its dimension accessors — is replaced by the concrete
// utilities the resolved recipe supplies, so the output contains no recipe
// construct at all.
func Resolve(filename string, src []byte, resolved recipe.Resolved) ([]byte, error) {
	r, err := inspectSource(filename, src, func(r *resolver) {
		r.classMode = helperDesugarMode
		r.resolved = resolved
		r.shape = resolved.Shape
	})
	if err != nil {
		return nil, err
	}

	edited, err := applyLiteralEdits(src, r.edits)
	if err != nil {
		return nil, fmt.Errorf("%s: apply recipe edits: %w", filename, err)
	}
	formatted, err := gen.Format(filename, edited)
	if err != nil {
		return nil, fmt.Errorf("%s: format resolved GSX: %w", filename, err)
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, formatted, 0); err != nil {
		return nil, fmt.Errorf("%s: reparse resolved GSX: %w", filename, err)
	}
	remaining, err := RecipeTokens(filename, formatted)
	if err != nil {
		return nil, fmt.Errorf("%s: validate resolved GSX recipe use: %w", filename, err)
	}
	if len(remaining) != 0 {
		return nil, fmt.Errorf("%s: resolved GSX still contains recipe tokens %q", filename, remaining)
	}
	if bytes.Contains(formatted, []byte(recipe.Prefix)) {
		return nil, fmt.Errorf("%s: resolved GSX still contains recipe prefix %q", filename, recipe.Prefix)
	}
	if residue, ok := residualAccessorCall(formatted, knownComponents(resolved.Shape)); ok {
		return nil, fmt.Errorf("%s: resolved GSX still contains a recipe accessor call %q", filename, residue)
	}
	return formatted, nil
}

// knownComponents names every identifier that may legally receive a recipe
// accessor call, sorted. Deriving the residue check from this instead of a
// hardcoded ".Role(", ".Variant(", ".Size(" list is what keeps adding a slot or
// a dimension from being three uncoupled edits.
func knownComponents(extra ...recipe.Shape) []string {
	declared := shapes.All()
	names := make([]string, 0, len(declared)+len(extra))
	for name := range declared {
		names = append(names, name)
	}
	for _, shape := range extra {
		if shape.Component != "" && !slices.Contains(names, shape.Component) {
			names = append(names, shape.Component)
		}
	}
	sort.Strings(names)
	return names
}

// residualAccessorCall reports any surviving `<component>.<Method>(` whose
// receiver names a known component. Every such call should have been replaced
// by concrete utilities, so one reaching here means the desugaring missed it.
func residualAccessorCall(src []byte, components []string) (string, bool) {
	for _, component := range components {
		pattern := regexp.MustCompile(`(^|[^\w.])` + regexp.QuoteMeta(component) + `\.([A-Z][A-Za-z0-9_]*)\(`)
		if match := pattern.FindSubmatch(src); match != nil {
			return component + "." + string(match[2]) + "(", true
		}
	}
	return "", false
}

// HelperCalls reports every recipe accessor call a canonical source makes, in
// source order, without resolving any of them. Generation uses it to check a
// canonical against its declared shape before it touches any style. The shape
// is a parameter because an accessor's method name can only be split back into
// (slot, dimension) against the shape it was generated from.
func HelperCalls(filename string, src []byte, shape recipe.Shape) ([]Call, error) {
	r, err := inspectSource(filename, src, func(r *resolver) {
		r.classMode = helperScanMode
		r.shape = shape
	})
	if err != nil {
		return nil, err
	}
	return r.calls, nil
}

// RecipeTokens reports every recipe token a source names in a class value, and
// rejects every structurally invalid recipe use it finds along the way.
func RecipeTokens(filename string, src []byte) ([]string, error) {
	r, err := inspectSource(filename, src, func(r *resolver) {
		r.classMode = tokenScanMode
	})
	if err != nil {
		return nil, err
	}

	tokens := make([]string, 0, len(r.used))
	for token := range r.used {
		tokens = append(tokens, token)
	}
	sort.Strings(tokens)
	return tokens, nil
}

// inspectSource runs the one shared traversal. configure selects the class
// mode and supplies whatever that mode resolves against.
func inspectSource(filename string, src []byte, configure func(*resolver)) (*resolver, error) {
	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}

	r := &resolver{
		filename: filename,
		src:      src,
		fset:     fset,
		used:     make(map[string]struct{}),
	}
	if configure != nil {
		configure(r)
	}
	gsxast.Inspect(file, r.inspectVisitor(canonicalClassInspection))
	if r.err != nil {
		return nil, r.err
	}
	return r, nil
}

func (r *resolver) inspectVisitor(mode inspectionMode) func(gsxast.Node) bool {
	return func(node gsxast.Node) bool {
		return r.inspectNode(node, mode)
	}
}

func (r *resolver) inspectNode(node gsxast.Node, mode inspectionMode) bool {
	if r.err != nil {
		return false
	}
	switch node := node.(type) {
	case *gsxast.ClassAttr:
		r.inspectClassAttr(node, mode)
		return false
	case *gsxast.StaticAttr:
		r.inspectStaticAttr(node)
	case *gsxast.ExprAttr:
		r.inspectNonClassExpr(node.Expr, node.ExprPos, node.Name)
	case *gsxast.EmbeddedAttr:
		r.inspectEmbeddedAttr(node)
		return false
	case *gsxast.EmbeddedInterp:
		r.inspectEmbeddedSegments(node.Segments, node.Pos(), "embedded interpolation")
		r.inspectPipelineStages(node.Stages, "embedded interpolation pipeline stage")
		return false
	case *gsxast.Interp:
		r.inspectNonClassExpr(node.Expr, node.ExprPos, "interpolation")
		r.inspectPipelineStages(node.Stages, "interpolation pipeline stage")
	case *gsxast.Text:
		if mode == validationOnlyInspection {
			r.inspectNestedText(node)
		}
	}
	return r.err == nil
}

func (r *resolver) inspectClassAttr(attr *gsxast.ClassAttr, mode inspectionMode) {
	canonicalClass := mode == canonicalClassInspection && attr.Name == "class"
	context := attr.Name
	if mode == validationOnlyInspection {
		context = "nested " + context
	}
	for i := range attr.Parts {
		part := &attr.Parts[i]
		if part.Expr != "" {
			if canonicalClass {
				r.inspectClassExpr(part.Expr, part.ExprPos, part)
			} else {
				r.inspectNonClassExpr(part.Expr, part.ExprPos, context)
			}
		}
		if part.Cond != "" {
			r.inspectNonClassExpr(part.Cond, part.CondPos, context+" condition")
		}
		r.inspectPipelineStages(part.Stages, context+" pipeline stage")
		if len(part.CSSSegments) != 0 {
			r.inspectEmbeddedSegments(part.CSSSegments, part.Pos(), context+" CSS segment")
		}
		if part.CF == nil {
			continue
		}
		if part.CF.If != nil {
			r.inspectValueIf(part.CF.If, context)
		}
		if part.CF.Switch != nil {
			r.inspectValueSwitch(part.CF.Switch, context, canonicalClass)
		}
	}
}

func (r *resolver) inspectClassExpr(expr string, exprPos token.Pos, part *gsxast.ClassPart) {
	if r.err != nil {
		return
	}
	parsed, err := goparser.ParseExpr(expr)
	if err != nil {
		r.err = r.positionedError(exprPos, "class recipe expression must be a Go string literal: %v", err)
		return
	}

	if r.classMode != tokenScanMode {
		if call, ok := parsed.(*goast.CallExpr); ok {
			if selector, ok := call.Fun.(*goast.SelectorExpr); ok {
				r.inspectHelperCall(call, selector, exprPos, part)
				return
			}
		}
	}

	literal := unwrappedStringLiteral(parsed)
	if literal == nil {
		if recipe, ok := recipeLiteralInExpr(parsed); ok {
			if expressionUsesConcatenation(parsed) {
				r.err = r.positionedError(
					exprPos,
					"recipe token %q must be one whole string literal; concatenation is not allowed",
					recipe,
				)
			} else {
				r.err = r.positionedError(
					exprPos,
					"recipe token %q appears in a non-string class expression; use a whole string literal",
					recipe,
				)
			}
		}
		return
	}

	value, err := strconv.Unquote(literal.Value)
	if err != nil {
		r.err = r.positionedError(exprPos, "invalid class string literal: %v", err)
		return
	}
	if err := r.scanClassLiteral(value); err != nil {
		r.err = r.positionedError(exprPos, "%v", err)
	}
}

// inspectHelperCall matches `<component>.<Method>(...)` in a canonical class
// expression. In desugar mode it records an edit replacing the whole class part
// with the concrete style source; in scan mode it only records the call.
func (r *resolver) inspectHelperCall(
	call *goast.CallExpr,
	selector *goast.SelectorExpr,
	exprPos token.Pos,
	part *gsxast.ClassPart,
) {
	// The part shape is validated in both helper modes, so a scan reports a
	// accessor call that Resolve would later refuse rather than passing it as
	// well formed.
	if !r.validateHelperPart(exprPos, part) {
		return
	}
	receiver, ok := selector.X.(*goast.Ident)
	if !ok {
		r.err = r.positionedError(exprPos, "recipe accessor call receiver must be a bare component identifier")
		return
	}
	component := receiver.Name
	if component != r.shape.Component {
		r.err = r.positionedError(
			exprPos,
			"recipe accessor call on component %q, but this source declares component %q",
			component,
			r.shape.Component,
		)
		return
	}

	method := selector.Sel.Name
	slotName, dimensionName, ok := accessorTarget(r.shape, method)
	if !ok {
		r.err = r.positionedError(
			exprPos,
			"component %q declares no accessor %s.%s",
			r.shape.Component, component, method,
		)
		return
	}

	if dimensionName == "" {
		if len(call.Args) != 0 {
			r.err = r.positionedError(exprPos, "%s.%s takes no arguments, got %d", component, method, len(call.Args))
			return
		}
		r.calls = append(r.calls, Call{Component: component, Slot: slotName})
		if r.classMode != helperDesugarMode {
			return
		}
		r.recordPartEdit(exprPos, part,
			strconv.Quote(strings.Join(r.resolved.BaseUtilities(slotName), " ")))
		return
	}

	if len(call.Args) != 1 {
		r.err = r.positionedError(
			exprPos,
			"%s.%s takes exactly one argument, got %d",
			component, method, len(call.Args),
		)
		return
	}
	argument, ok := call.Args[0].(*goast.Ident)
	if !ok {
		r.err = r.positionedError(
			exprPos,
			"%s.%s argument must be a bare identifier naming a component parameter",
			component, method,
		)
		return
	}
	r.calls = append(r.calls, Call{Component: component, Slot: slotName, Dimension: dimensionName})
	if r.classMode != helperDesugarMode {
		return
	}
	slot, _ := r.shape.Slot(slotName)
	dimension, _ := slot.Dimension(dimensionName)
	r.recordPartEdit(exprPos, part, dimensionSwitch(r.resolved, slotName, dimension, argument.Name))
}

// accessorTarget resolves a generated accessor's method name back to the
// (slot, dimension) it was generated from. Name-splitting cannot do this:
// "MenuButtonSize" is slot "menu-button" + dimension "size", and nothing in
// the string says where the boundary is.
func accessorTarget(shape recipe.Shape, method string) (slot, dimension string, ok bool) {
	for _, s := range shape.Slots {
		if s.Base && accessorName(s.Name) == method {
			return s.Name, "", true
		}
		for _, d := range s.Dimensions {
			if dimensionAccessorName(s.Name, d.Name) == method {
				return s.Name, d.Name, true
			}
		}
	}
	return "", "", false
}

// validateHelperPart requires an accessor call to occupy a whole plain class part.
// The replacement span runs to part.End(), which also covers a part's condition,
// pipeline stages, control flow and CSS segments — none of which the generated
// replacement reproduces. Substituting anyway would silently delete them, so
// anything but a plain part is rejected. Runs in both helper modes.
func (r *resolver) validateHelperPart(exprPos token.Pos, part *gsxast.ClassPart) bool {
	if part == nil {
		r.err = r.positionedError(exprPos, "a recipe accessor call must be a whole class part")
		return false
	}
	if part.Cond != "" || len(part.Stages) != 0 || part.CF != nil || len(part.CSSSegments) != 0 {
		r.err = r.positionedError(exprPos, "a recipe accessor call must be a plain class part with no condition, pipeline or control flow")
		return false
	}
	return true
}

// recordPartEdit replaces a whole class part. Per the resolved spike a
// ClassPart's Pos includes the leading newline and indentation while End
// excludes the trailing comma, so the expression's own position is the correct
// start and the part's End is the correct end. gen.Format re-indents
// afterwards, so value carries no indentation. The caller has already run
// validateHelperPart, so part is a non-nil plain part here.
func (r *resolver) recordPartEdit(exprPos token.Pos, part *gsxast.ClassPart, value string) {
	start := r.fset.Position(exprPos).Offset
	end := r.fset.Position(part.End()).Offset
	if start < 0 || end < start || end > len(r.src) {
		r.err = r.positionedError(exprPos, "recipe accessor call span is outside source")
		return
	}
	r.edits = append(r.edits, literalEdit{start: start, end: end, value: value})
}

// dimensionSwitch generates the class switch for one dimension. Every declared
// value except the default gets a case; the declared default's utilities go in
// the default arm, so an unrecognized value renders the default rather than
// nothing. That mirrors recipe.Component's runtime behavior, which is what
// makes a behavior test written against the canonical true of every style.
func dimensionSwitch(resolved recipe.Resolved, slot string, dimension recipe.Dimension, argument string) string {
	var out strings.Builder
	fmt.Fprintf(&out, "switch %s {\n", argument)
	for _, value := range dimension.Values {
		if value == dimension.Default {
			continue
		}
		fmt.Fprintf(&out, "case %s:\n%s\n",
			strconv.Quote(value),
			strconv.Quote(strings.Join(resolved.Utilities(slot, dimension.Name, value), " ")))
	}
	fmt.Fprintf(&out, "default:\n%s\n}",
		strconv.Quote(strings.Join(resolved.Utilities(slot, dimension.Name, dimension.Default), " ")))
	return out.String()
}

// scanClassLiteral records every recipe token a class string literal names and
// rejects any token that is not a whole class token. It never rewrites: nothing
// authored is in token form any more, so a token reaching here is either a
// report for RecipeTokens or an error.
func (r *resolver) scanClassLiteral(value string) error {
	for pos := 0; pos < len(value); {
		if isHTMLASCIIWhitespace(value[pos]) {
			pos++
			continue
		}
		end := pos + 1
		for end < len(value) && !isHTMLASCIIWhitespace(value[end]) {
			end++
		}
		classToken := value[pos:end]
		switch prefix := strings.Index(classToken, recipe.Prefix); {
		case prefix < 0:
		case prefix > 0:
			return fmt.Errorf("recipe token must be a whole class token, found %q", classToken)
		default:
			r.used[classToken] = struct{}{}
		}
		pos = end
	}
	return nil
}

func (r *resolver) inspectStaticAttr(attr *gsxast.StaticAttr) {
	if !strings.Contains(attr.Value, recipe.Prefix) {
		return
	}
	if attr.Name == "class" {
		r.err = r.positionedError(attr.Pos(), "recipe token %q is not allowed in a static class attribute", attr.Value)
		return
	}
	r.err = r.positionedError(attr.Pos(), "recipe token %q appears in non-class attribute %q", attr.Value, attr.Name)
}

func (r *resolver) inspectNestedText(text *gsxast.Text) {
	if !strings.Contains(text.Value, recipe.Prefix) {
		return
	}
	r.err = r.positionedError(text.Pos(), "recipe token %q appears in nested element text", text.Value)
}

func (r *resolver) inspectEmbeddedAttr(attr *gsxast.EmbeddedAttr) {
	r.inspectEmbeddedSegments(attr.Segments, attr.Pos(), attr.Name+" interpolation")
	r.inspectPipelineStages(attr.Stages, attr.Name+" pipeline stage")
}

func (r *resolver) inspectNonClassExpr(expr string, exprPos token.Pos, context string) {
	if r.err != nil || expr == "" {
		return
	}
	r.inspectNonClassGo(expr, exprPos, context, false)
}

func (r *resolver) inspectNonClassExprList(exprs string, exprPos token.Pos, context string) {
	if r.err != nil || exprs == "" {
		return
	}
	r.inspectNonClassGo(exprs, exprPos, context, true)
}

func (r *resolver) inspectNonClassGo(src string, srcPos token.Pos, context string, expressionList bool) goast.Expr {
	parse := func(value string) (goast.Expr, error) {
		if expressionList {
			value = "scan(" + value + ")"
		}
		return goparser.ParseExpr(value)
	}

	parsed, err := parse(src)
	if err == nil {
		r.rejectRecipeInNonClassExpr(parsed, srcPos, context)
		return parsed
	}

	parts, splitErrs := gsxparser.SplitGoExprElements(r.fset, src, srcPos, nil)
	if len(splitErrs) != 0 {
		first := splitErrs[0]
		r.err = r.positionedError(first.Pos, "cannot parse non-class %s expression: %s", context, first.Msg)
		return nil
	}
	if len(parts) == 0 {
		// The containing file already passed the public GSX parser. Some raw GSX
		// metadata is intentionally not a standalone Go expression, so a failed
		// go/parser.ParseExpr is not itself grounds for rejecting valid GSX.
		return nil
	}

	normalized := normalizeGoParts(parts)
	parsed, err = parse(normalized)
	if err == nil {
		r.rejectRecipeInNonClassExpr(parsed, srcPos, context)
	}
	for _, part := range parts {
		if _, ok := part.(gsxast.GoText); ok {
			continue
		}
		gsxast.Inspect(part, r.inspectVisitor(validationOnlyInspection))
		if r.err != nil {
			return parsed
		}
	}
	return parsed
}

func (r *resolver) rejectRecipeInNonClassExpr(expr goast.Expr, exprPos token.Pos, context string) {
	if r.err != nil || expr == nil {
		return
	}
	if recipe, ok := recipeLiteralInExpr(expr); ok {
		r.err = r.positionedError(exprPos, "recipe token %q appears in non-class %s expression", recipe, context)
	}
}

func normalizeGoParts(parts []gsxast.GoPart) string {
	var normalized strings.Builder
	for _, part := range parts {
		if text, ok := part.(gsxast.GoText); ok {
			normalized.WriteString(text.Src)
			continue
		}
		normalized.WriteString("gsxDialectValue")
	}
	return normalized.String()
}

func (r *resolver) inspectPipelineStages(stages []gsxast.PipeStage, context string) {
	for _, stage := range stages {
		if !stage.HasArgs || stage.Args == "" {
			continue
		}
		r.inspectNonClassExprList(stage.Args, stage.ArgsPos, context)
	}
}

func (r *resolver) inspectValueArm(arm *gsxast.ValueArm, transform bool, context string) {
	if arm == nil {
		return
	}
	if transform {
		r.inspectClassExpr(arm.Expr, arm.ExprPos, nil)
	} else {
		r.inspectNonClassExpr(arm.Expr, arm.ExprPos, context)
	}
	r.inspectPipelineStages(arm.Stages, context+" pipeline stage")
}

func (r *resolver) inspectValueIf(valueIf *gsxast.ValueIf, attrName string) {
	for valueIf != nil {
		r.inspectControlClause("if", valueIf.Cond, valueIf.CondPos, attrName+" value-if condition")
		r.inspectValueArm(valueIf.Then, false, attrName+" value-if arm")
		r.inspectValueArm(valueIf.Else, false, attrName+" value-if arm")
		valueIf = valueIf.ElseIf
	}
}

func (r *resolver) inspectValueSwitch(valueSwitch *gsxast.ValueSwitch, attrName string, transform bool) {
	if valueSwitch.Tag != "" {
		r.inspectControlClause("switch", valueSwitch.Tag, valueSwitch.TagPos, attrName+" switch tag")
	}
	for _, classCase := range valueSwitch.Cases {
		if classCase.List != "" {
			r.inspectNonClassExprList(classCase.List, classCase.ListPos, attrName+" switch case")
		}
		r.inspectValueArm(classCase.Value, transform, attrName+" switch arm")
	}
}

func (r *resolver) inspectControlClause(keyword, clause string, clausePos token.Pos, context string) {
	if r.err != nil {
		return
	}
	src := "package p\nfunc _() { " + keyword + " " + clause + " {} }\n"
	file, err := goparser.ParseFile(token.NewFileSet(), "", src, 0)
	if err != nil {
		r.err = r.positionedError(clausePos, "cannot parse %s metadata: %v", context, err)
		return
	}
	if recipe, ok := recipeLiteralInNode(file); ok {
		r.err = r.positionedError(clausePos, "recipe token %q appears in non-class %s", recipe, context)
	}
}

func (r *resolver) inspectEmbeddedSegments(segments []gsxast.Markup, pos token.Pos, context string) {
	var assembled strings.Builder
	for _, markup := range segments {
		switch segment := markup.(type) {
		case *gsxast.Text:
			assembled.WriteString(segment.Value)
		case *gsxast.Interp:
			parsed := r.inspectNonClassGo(segment.Expr, segment.ExprPos, context, false)
			if r.err != nil {
				return
			}
			r.inspectPipelineStages(segment.Stages, context+" pipeline stage")
			if r.err != nil {
				return
			}
			if parsed == nil {
				assembled.Reset()
				continue
			}
			prefix, suffix, fullyConstant := constantBoundaryFragments(parsed)
			assembled.WriteString(prefix)
			if strings.Contains(assembled.String(), recipe.Prefix) {
				r.err = r.positionedError(
					pos,
					"recipe token %q must be a whole class string literal; found in non-class %s",
					assembled.String(),
					context,
				)
				return
			}
			if fullyConstant {
				continue
			}
			assembled.Reset()
			assembled.WriteString(suffix)
		default:
			assembled.Reset()
			continue
		}
		if strings.Contains(assembled.String(), recipe.Prefix) {
			r.err = r.positionedError(
				pos,
				"recipe token %q must be a whole class string literal; found in non-class %s",
				assembled.String(),
				context,
			)
			return
		}
	}
}

func (r *resolver) positionedError(pos token.Pos, format string, args ...any) error {
	position := r.fset.Position(pos)
	message := fmt.Sprintf(format, args...)
	if position.IsValid() {
		return fmt.Errorf("%s:%d:%d: %s", r.filename, position.Line, position.Column, message)
	}
	return fmt.Errorf("%s: %s", r.filename, message)
}

func unwrappedStringLiteral(expr goast.Expr) *goast.BasicLit {
	for {
		parenthesized, ok := expr.(*goast.ParenExpr)
		if !ok {
			break
		}
		expr = parenthesized.X
	}
	literal, ok := expr.(*goast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		return nil
	}
	return literal
}

func recipeLiteralInExpr(expr goast.Expr) (string, bool) {
	return recipeLiteralInNode(expr)
}

func recipeLiteralInNode(root goast.Node) (string, bool) {
	var found string
	goast.Inspect(root, func(node goast.Node) bool {
		if found != "" {
			return false
		}
		literal, ok := node.(*goast.BasicLit)
		if !ok || literal.Kind != token.STRING {
			return true
		}
		value, err := strconv.Unquote(literal.Value)
		if err != nil {
			return true
		}
		if strings.Contains(value, recipe.Prefix) {
			found = value
			return false
		}
		return true
	})
	if found == "" {
		goast.Inspect(root, func(node goast.Node) bool {
			if found != "" {
				return false
			}
			expr, ok := node.(goast.Expr)
			if !ok {
				return true
			}
			if value, ok := constantRecipeAssembly(expr); ok {
				found = value
				return false
			}
			return true
		})
	}
	return found, found != ""
}

func constantRecipeAssembly(expr goast.Expr) (string, bool) {
	var found string
	goast.Inspect(expr, func(node goast.Node) bool {
		if found != "" {
			return false
		}
		binary, ok := node.(*goast.BinaryExpr)
		if !ok || binary.Op != token.ADD {
			return true
		}
		var operands []goast.Expr
		flattenAddition(binary, &operands)
		var run strings.Builder
		for _, operand := range operands {
			value, ok := constantStringValue(operand)
			if !ok {
				run.Reset()
				continue
			}
			run.WriteString(value)
			if strings.Contains(run.String(), recipe.Prefix) {
				found = run.String()
				return false
			}
		}
		return true
	})
	return found, found != ""
}

func flattenAddition(expr goast.Expr, operands *[]goast.Expr) {
	if parenthesized, ok := expr.(*goast.ParenExpr); ok {
		flattenAddition(parenthesized.X, operands)
		return
	}
	binary, ok := expr.(*goast.BinaryExpr)
	if !ok || binary.Op != token.ADD {
		*operands = append(*operands, expr)
		return
	}
	flattenAddition(binary.X, operands)
	flattenAddition(binary.Y, operands)
}

func constantBoundaryFragments(expr goast.Expr) (prefix, suffix string, fullyConstant bool) {
	if value, ok := constantStringValue(expr); ok {
		return value, value, true
	}

	var operands []goast.Expr
	flattenAddition(expr, &operands)
	var leading strings.Builder
	for _, operand := range operands {
		value, ok := constantStringValue(operand)
		if !ok {
			break
		}
		leading.WriteString(value)
	}
	var trailing []string
	for i := len(operands) - 1; i >= 0; i-- {
		value, ok := constantStringValue(operands[i])
		if !ok {
			break
		}
		trailing = append(trailing, value)
	}
	var trailingText strings.Builder
	for i := len(trailing) - 1; i >= 0; i-- {
		trailingText.WriteString(trailing[i])
	}
	return leading.String(), trailingText.String(), false
}

func constantStringValue(expr goast.Expr) (string, bool) {
	switch expr := expr.(type) {
	case *goast.BasicLit:
		if expr.Kind != token.STRING {
			return "", false
		}
		value, err := strconv.Unquote(expr.Value)
		return value, err == nil
	case *goast.ParenExpr:
		return constantStringValue(expr.X)
	case *goast.BinaryExpr:
		if expr.Op != token.ADD {
			return "", false
		}
		left, leftOK := constantStringValue(expr.X)
		right, rightOK := constantStringValue(expr.Y)
		if !leftOK || !rightOK {
			return "", false
		}
		return left + right, true
	default:
		return "", false
	}
}

func expressionUsesConcatenation(expr goast.Expr) bool {
	found := false
	goast.Inspect(expr, func(node goast.Node) bool {
		if binary, ok := node.(*goast.BinaryExpr); ok && binary.Op == token.ADD {
			found = true
			return false
		}
		return !found
	})
	return found
}

func applyLiteralEdits(src []byte, edits []literalEdit) ([]byte, error) {
	sort.Slice(edits, func(i, j int) bool {
		return edits[i].start > edits[j].start
	})
	resolved := append([]byte(nil), src...)
	lastStart := len(src)
	for _, edit := range edits {
		if edit.start < 0 || edit.end < edit.start || edit.end > len(src) {
			return nil, fmt.Errorf("invalid literal span [%d:%d]", edit.start, edit.end)
		}
		if edit.end > lastStart {
			return nil, fmt.Errorf("overlapping literal span [%d:%d]", edit.start, edit.end)
		}
		var next bytes.Buffer
		next.Grow(len(resolved) - (edit.end - edit.start) + len(edit.value))
		next.Write(resolved[:edit.start])
		next.WriteString(edit.value)
		next.Write(resolved[edit.end:])
		resolved = next.Bytes()
		lastStart = edit.start
	}
	return resolved, nil
}

func isHTMLASCIIWhitespace(char byte) bool {
	switch char {
	case '\t', '\n', '\f', '\r', ' ':
		return true
	default:
		return false
	}
}
