package stylegen

import (
	"bytes"
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"sort"
	"strconv"
	"strings"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"
)

type ResolveReport struct {
	UsedTokens []string
}

type literalEdit struct {
	start int
	end   int
	value string
}

type resolver struct {
	filename string
	src      []byte
	fset     *token.FileSet
	recipes  *Recipes
	used     map[string]struct{}
	edits    []literalEdit
	err      error
}

func Resolve(filename string, src []byte, recipes Recipes) ([]byte, ResolveReport, error) {
	declared, err := declaredRecipeTokens(filename, recipes)
	if err != nil {
		return nil, ResolveReport{}, err
	}

	result, err := inspectRecipeTokens(filename, src, &recipes)
	if err != nil {
		return nil, ResolveReport{}, err
	}
	missing, unused := recipeSetDifference(result.tokens, declared)
	if len(missing) != 0 || len(unused) != 0 {
		return nil, ResolveReport{}, fmt.Errorf(
			"%s: recipe token set mismatch: missing=%q unused=%q",
			filename,
			missing,
			unused,
		)
	}

	resolved, err := applyLiteralEdits(src, result.edits)
	if err != nil {
		return nil, ResolveReport{}, fmt.Errorf("%s: apply recipe literals: %w", filename, err)
	}
	formatted, err := gen.Format(filename, resolved)
	if err != nil {
		return nil, ResolveReport{}, fmt.Errorf("%s: format resolved GSX: %w", filename, err)
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, formatted, 0); err != nil {
		return nil, ResolveReport{}, fmt.Errorf("%s: reparse resolved GSX: %w", filename, err)
	}
	remaining, err := RecipeTokens(filename, formatted)
	if err != nil {
		return nil, ResolveReport{}, fmt.Errorf("%s: validate resolved GSX recipe use: %w", filename, err)
	}
	if len(remaining) != 0 {
		return nil, ResolveReport{}, fmt.Errorf("%s: resolved GSX still contains recipe tokens %q", filename, remaining)
	}
	return formatted, ResolveReport{UsedTokens: result.tokens}, nil
}

func RecipeTokens(filename string, src []byte) ([]string, error) {
	result, err := inspectRecipeTokens(filename, src, nil)
	if err != nil {
		return nil, err
	}
	return result.tokens, nil
}

type inspectionResult struct {
	tokens []string
	edits  []literalEdit
}

func inspectRecipeTokens(filename string, src []byte, recipes *Recipes) (inspectionResult, error) {
	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return inspectionResult{}, fmt.Errorf("%s: %w", filename, err)
	}

	r := &resolver{
		filename: filename,
		src:      src,
		fset:     fset,
		recipes:  recipes,
		used:     make(map[string]struct{}),
	}
	gsxast.Inspect(file, func(node gsxast.Node) bool {
		if r.err != nil {
			return false
		}
		switch node := node.(type) {
		case *gsxast.ClassAttr:
			r.inspectClassAttr(node)
			return false
		case *gsxast.StaticAttr:
			r.inspectStaticAttr(node)
		case *gsxast.ExprAttr:
			r.inspectNonClassExpr(node.Expr, node.ExprPos, node.Name)
		case *gsxast.EmbeddedAttr:
			r.inspectEmbeddedAttr(node)
		case *gsxast.Interp:
			r.inspectNonClassExpr(node.Expr, node.ExprPos, "interpolation")
		}
		return r.err == nil
	})
	if r.err != nil {
		return inspectionResult{}, r.err
	}

	tokens := make([]string, 0, len(r.used))
	for token := range r.used {
		tokens = append(tokens, token)
	}
	sort.Strings(tokens)
	return inspectionResult{tokens: tokens, edits: r.edits}, nil
}

func (r *resolver) inspectClassAttr(attr *gsxast.ClassAttr) {
	for i := range attr.Parts {
		part := &attr.Parts[i]
		if part.Expr != "" {
			if attr.Name == "class" {
				r.inspectClassExpr(part.Expr, part.ExprPos)
			} else {
				r.inspectNonClassExpr(part.Expr, part.ExprPos, attr.Name)
			}
		}
		if part.Cond != "" {
			r.inspectNonClassExpr(part.Cond, part.CondPos, attr.Name+" condition")
		}
		r.inspectPipelineStages(part.Stages, attr.Name+" pipeline stage")
		if len(part.CSSSegments) != 0 {
			r.inspectEmbeddedSegments(part.CSSSegments, part.Pos(), attr.Name+" CSS segment")
		}
		if part.CF == nil {
			continue
		}
		if part.CF.If != nil {
			r.inspectValueIf(part.CF.If, attr.Name)
		}
		if part.CF.Switch != nil {
			r.inspectValueSwitch(part.CF.Switch, attr.Name)
		}
	}
}

func (r *resolver) inspectClassExpr(expr string, exprPos token.Pos) {
	if r.err != nil {
		return
	}
	parsed, err := goparser.ParseExpr(expr)
	if err != nil {
		r.err = r.positionedError(exprPos, "class recipe expression must be a Go string literal: %v", err)
		return
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
	resolved, changed, err := r.resolveClassLiteral(value)
	if err != nil {
		r.err = r.positionedError(exprPos, "%v", err)
		return
	}
	if !changed {
		return
	}

	exprOffset := r.fset.Position(exprPos).Offset
	if exprOffset < 0 {
		r.err = r.positionedError(exprPos, "class expression has no source position")
		return
	}
	start := exprOffset + int(literal.Pos()) - 1
	end := exprOffset + int(literal.End()) - 1
	if start < 0 || end < start || end > len(r.src) {
		r.err = r.positionedError(exprPos, "class string literal span is outside source")
		return
	}
	literal.Value = strconv.Quote(resolved)
	r.edits = append(r.edits, literalEdit{start: start, end: end, value: literal.Value})
}

func (r *resolver) resolveClassLiteral(value string) (string, bool, error) {
	var out strings.Builder
	out.Grow(len(value))
	changed := false

	for pos := 0; pos < len(value); {
		if isHTMLASCIIWhitespace(value[pos]) {
			out.WriteByte(value[pos])
			pos++
			continue
		}
		end := pos + 1
		for end < len(value) && !isHTMLASCIIWhitespace(value[end]) {
			end++
		}
		classToken := value[pos:end]
		prefix := strings.Index(classToken, RecipePrefix)
		switch {
		case prefix < 0:
			out.WriteString(classToken)
		case prefix > 0:
			return "", false, fmt.Errorf("recipe token must be a whole class token, found %q", classToken)
		default:
			r.used[classToken] = struct{}{}
			if r.recipes == nil {
				out.WriteString(classToken)
				break
			}
			recipe, ok := r.recipes.Lookup(classToken)
			if !ok {
				out.WriteString(classToken)
				break
			}
			out.WriteString(strings.Join(recipe.Utilities, " "))
			changed = true
		}
		pos = end
	}
	return out.String(), changed, nil
}

func (r *resolver) inspectStaticAttr(attr *gsxast.StaticAttr) {
	if !strings.Contains(attr.Value, RecipePrefix) {
		return
	}
	if attr.Name == "class" {
		r.err = r.positionedError(attr.Pos(), "recipe token %q is not allowed in a static class attribute", attr.Value)
		return
	}
	r.err = r.positionedError(attr.Pos(), "recipe token %q appears in non-class attribute %q", attr.Value, attr.Name)
}

func (r *resolver) inspectEmbeddedAttr(attr *gsxast.EmbeddedAttr) {
	r.inspectEmbeddedSegments(attr.Segments, attr.Pos(), attr.Name+" interpolation")
	r.inspectPipelineStages(attr.Stages, attr.Name+" pipeline stage")
}

func (r *resolver) inspectNonClassExpr(expr string, exprPos token.Pos, context string) {
	if r.err != nil || expr == "" {
		return
	}
	parsed, err := goparser.ParseExpr(expr)
	if err != nil {
		return
	}
	if recipe, ok := recipeLiteralInExpr(parsed); ok {
		r.err = r.positionedError(exprPos, "recipe token %q appears in non-class %s expression", recipe, context)
	}
}

func (r *resolver) inspectNonClassExprList(exprs string, exprPos token.Pos, context string) {
	if r.err != nil || exprs == "" {
		return
	}
	parsed, err := goparser.ParseExpr("scan(" + exprs + ")")
	if err != nil {
		return
	}
	if recipe, ok := recipeLiteralInExpr(parsed); ok {
		r.err = r.positionedError(exprPos, "recipe token %q appears in non-class %s expression", recipe, context)
	}
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
		r.inspectClassExpr(arm.Expr, arm.ExprPos)
	} else {
		r.inspectNonClassExpr(arm.Expr, arm.ExprPos, context)
	}
	r.inspectPipelineStages(arm.Stages, context+" pipeline stage")
}

func (r *resolver) inspectValueIf(valueIf *gsxast.ValueIf, attrName string) {
	for valueIf != nil {
		r.inspectNonClassExpr(valueIf.Cond, valueIf.CondPos, attrName+" value-if condition")
		r.inspectValueArm(valueIf.Then, false, attrName+" value-if arm")
		r.inspectValueArm(valueIf.Else, false, attrName+" value-if arm")
		valueIf = valueIf.ElseIf
	}
}

func (r *resolver) inspectValueSwitch(valueSwitch *gsxast.ValueSwitch, attrName string) {
	r.inspectNonClassExpr(valueSwitch.Tag, valueSwitch.TagPos, attrName+" switch tag")
	for _, classCase := range valueSwitch.Cases {
		if classCase.List != "" {
			r.inspectNonClassExprList(classCase.List, classCase.ListPos, attrName+" switch case")
		}
		r.inspectValueArm(classCase.Value, attrName == "class", attrName+" switch arm")
	}
}

func (r *resolver) inspectEmbeddedSegments(segments []gsxast.Markup, pos token.Pos, context string) {
	recipe, ok := recipeInEmbeddedSegments(segments)
	if !ok {
		return
	}
	r.err = r.positionedError(
		pos,
		"recipe token %q must be a whole class string literal; found in non-class %s",
		recipe,
		context,
	)
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
	var found string
	goast.Inspect(expr, func(node goast.Node) bool {
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
		if strings.Contains(value, RecipePrefix) {
			found = value
			return false
		}
		return true
	})
	if found == "" {
		if value, ok := constantRecipeAssembly(expr); ok {
			found = value
		}
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
			if strings.Contains(run.String(), RecipePrefix) {
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

func recipeInEmbeddedSegments(segments []gsxast.Markup) (string, bool) {
	var assembled strings.Builder
	for _, segment := range segments {
		switch segment := segment.(type) {
		case *gsxast.Text:
			assembled.WriteString(segment.Value)
		case *gsxast.Interp:
			parsed, err := goparser.ParseExpr(segment.Expr)
			if err != nil {
				assembled.Reset()
				continue
			}
			if recipe, ok := recipeLiteralInExpr(parsed); ok {
				return recipe, true
			}
			prefix, suffix, fullyConstant := constantBoundaryFragments(parsed)
			assembled.WriteString(prefix)
			if strings.Contains(assembled.String(), RecipePrefix) {
				return assembled.String(), true
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
		if strings.Contains(assembled.String(), RecipePrefix) {
			return assembled.String(), true
		}
	}
	return "", false
}

func declaredRecipeTokens(filename string, recipes Recipes) (map[string]struct{}, error) {
	declared := make(map[string]struct{}, len(recipes.ordered))
	for _, recipe := range recipes.ordered {
		if _, exists := declared[recipe.Token]; exists {
			return nil, fmt.Errorf("%s: duplicate recipe declaration %q", filename, recipe.Token)
		}
		declared[recipe.Token] = struct{}{}
	}
	return declared, nil
}

func recipeSetDifference(used []string, declared map[string]struct{}) (missing, unused []string) {
	usedSet := make(map[string]struct{}, len(used))
	for _, token := range used {
		usedSet[token] = struct{}{}
		if _, ok := declared[token]; !ok {
			missing = append(missing, token)
		}
	}
	for token := range declared {
		if _, ok := usedSet[token]; !ok {
			unused = append(unused, token)
		}
	}
	sort.Strings(missing)
	sort.Strings(unused)
	return missing, unused
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
