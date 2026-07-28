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
	if attr.Name != "class" {
		for i := range attr.Parts {
			part := &attr.Parts[i]
			if part.Expr != "" {
				r.inspectNonClassExpr(part.Expr, part.ExprPos, attr.Name)
			}
			r.inspectNonClassValueCF(part.CF, attr.Name)
		}
		return
	}

	for i := range attr.Parts {
		part := &attr.Parts[i]
		if part.Expr != "" {
			r.inspectClassExpr(part.Expr, part.ExprPos)
		}
		if part.Cond != "" {
			r.inspectNonClassExpr(part.Cond, part.CondPos, "class condition")
		}
		if part.CF == nil {
			continue
		}
		if part.CF.Switch != nil {
			for _, classCase := range part.CF.Switch.Cases {
				if classCase.Value != nil {
					r.inspectClassExpr(classCase.Value.Expr, classCase.Value.ExprPos)
				}
			}
		}
		if part.CF.If != nil {
			r.inspectNonClassValueIf(part.CF.If)
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
	recipe, ok := recipeInEmbeddedSegments(attr.Segments)
	if !ok {
		return
	}
	if attr.Name == "class" {
		r.err = r.positionedError(
			attr.Pos(),
			"recipe token %q must be one whole string literal; interpolation is not allowed",
			recipe,
		)
		return
	}
	r.err = r.positionedError(attr.Pos(), "recipe token %q appears in non-class attribute %q", recipe, attr.Name)
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

func (r *resolver) inspectNonClassValueCF(cf *gsxast.ValueCF, context string) {
	if cf == nil {
		return
	}
	if cf.Switch != nil {
		for _, classCase := range cf.Switch.Cases {
			if classCase.Value != nil {
				r.inspectNonClassExpr(classCase.Value.Expr, classCase.Value.ExprPos, context)
			}
		}
	}
	if cf.If != nil {
		r.inspectNonClassValueIf(cf.If)
	}
}

func (r *resolver) inspectNonClassValueIf(valueIf *gsxast.ValueIf) {
	for valueIf != nil {
		if valueIf.Then != nil {
			r.inspectNonClassExpr(valueIf.Then.Expr, valueIf.Then.ExprPos, "class value-form if")
		}
		if valueIf.Else != nil {
			r.inspectNonClassExpr(valueIf.Else.Expr, valueIf.Else.ExprPos, "class value-form if")
		}
		valueIf = valueIf.ElseIf
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
		if value, ok := constantStringValue(expr); ok && strings.Contains(value, RecipePrefix) {
			found = value
		}
	}
	return found, found != ""
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
	for _, segment := range segments {
		switch segment := segment.(type) {
		case *gsxast.Text:
			if strings.Contains(segment.Value, RecipePrefix) {
				return segment.Value, true
			}
		case *gsxast.Interp:
			parsed, err := goparser.ParseExpr(segment.Expr)
			if err == nil {
				if recipe, ok := recipeLiteralInExpr(parsed); ok {
					return recipe, true
				}
			}
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
