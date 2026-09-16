package rtl

import (
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"strconv"
	"strings"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/srcedit"
)

// quoting is how the source spelled a class list, and so how a rewritten
// one must be spelled back into the same span.
type quoting uint8

const (
	goString    quoting = iota // "…", a Go interpreted string literal
	goRawString                // `…`, a Go raw string literal
	attrText                   // bare text between a static attribute's quotes
)

// classList is one class-carrying string found in a gsx file: its decoded
// value, the byte span in src that holds it (the text only, never the
// surrounding quotes or backticks), how that span is quoted, and whether it
// sits in a side-keyed arm.
type classList struct {
	Value     string
	Start     int
	End       int
	Quoting   quoting
	SideKeyed bool
}

// source spells value the way the span it replaces was spelled. A static
// attribute's span is bare text inside the HTML quotes, so it is written
// verbatim rather than through strconv.Quote, whose Go escaping would
// corrupt any backslash a Tailwind arbitrary value carries.
func (l classList) source(value string) (string, error) {
	switch l.Quoting {
	case goRawString:
		if strings.ContainsRune(value, '`') {
			return "", fmt.Errorf("class list %q cannot be written as a raw string literal", value)
		}
		return "`" + value + "`", nil
	case attrText:
		if strings.ContainsRune(value, '"') {
			return "", fmt.Errorf("class list %q cannot be written as a static attribute value", value)
		}
		return value, nil
	default:
		return strconv.Quote(value), nil
	}
}

// GSX rewrites every class attribute literal in one gsx source file and
// returns the gsx-formatted result. A file whose class lists are already
// logical and single-spaced comes back byte-identical; a list that only
// needed its whitespace normalised still counts as an edit and the file is
// reformatted. Only string literals that are class-list elements or pair
// keys change; a literal that is a condition operand, another attribute, or
// text is untouched.
func GSX(filename string, src []byte) ([]byte, error) {
	lists, err := classLists(filename, src)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	var edits []srcedit.Edit
	for _, l := range lists {
		mapped := classes(l.Value, l.SideKeyed)
		if mapped == l.Value {
			continue
		}
		spelled, err := l.source(mapped)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", filename, err)
		}
		edits = append(edits, srcedit.Edit{Start: l.Start, End: l.End, Value: spelled})
	}
	if len(edits) == 0 {
		return src, nil
	}
	edited, err := srcedit.Apply(src, edits)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", filename, err)
	}
	formatted, err := gen.Format(filename, edited)
	if err != nil {
		return nil, fmt.Errorf("%s: format after rtl rewrite: %w", filename, err)
	}
	if _, err := gsxparser.ParseFile(token.NewFileSet(), filename, formatted, 0); err != nil {
		return nil, fmt.Errorf("%s: reparse after rtl rewrite: %w", filename, err)
	}
	return formatted, nil
}

// classLists finds every class literal: class="…" static attributes, and
// inside class={ … } every element literal, pair key, and value-form
// if/switch arm literal. Arms keyed on a placement parameter — `switch side`
// / `if side …`, and Drawer's `switch direction` — are side-keyed.
func classLists(filename string, src []byte) ([]classList, error) {
	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		return nil, err
	}
	var lists []classList
	var walkErr error
	gsxast.Inspect(file, func(n gsxast.Node) bool {
		if walkErr != nil {
			return false
		}
		switch a := n.(type) {
		case *gsxast.StaticAttr:
			if a.Name != "class" {
				return true
			}
			end := fset.Position(a.End()).Offset - 1 // closing quote
			start := end - len(a.Value)
			if start < 0 || string(src[start:end]) != a.Value {
				walkErr = fmt.Errorf("%s: static class span mismatch at %v", filename, fset.Position(a.Pos()))
				return false
			}
			lists = append(lists, classList{Value: a.Value, Start: start, End: end, Quoting: attrText})
		case *gsxast.ComposedAttr:
			if a.Name != "class" {
				return true
			}
			for i := range a.Parts {
				part := &a.Parts[i]
				if part.Expr != "" {
					if l, ok, err := literalIn(fset, src, part.Expr, part.ExprPos, false); err != nil {
						walkErr = err
						return false
					} else if ok {
						lists = append(lists, l)
					}
				}
				if part.CF == nil {
					continue
				}
				if part.CF.If != nil {
					sideKeyed := mentionsPlacement(part.CF.If.Cond)
					for vi := part.CF.If; vi != nil; vi = vi.ElseIf {
						sideKeyed = sideKeyed || mentionsPlacement(vi.Cond)
						if vi.Then != nil {
							if l, ok, err := armLiteral(fset, src, vi.Then, sideKeyed); err != nil {
								walkErr = err
								return false
							} else if ok {
								lists = append(lists, l)
							}
						}
						if vi.ElseIf == nil && vi.Else != nil {
							if l, ok, err := armLiteral(fset, src, vi.Else, sideKeyed); err != nil {
								walkErr = err
								return false
							} else if ok {
								lists = append(lists, l)
							}
						}
					}
				}
				if part.CF.Switch != nil {
					sideKeyed := mentionsPlacement(part.CF.Switch.Tag)
					for _, c := range part.CF.Switch.Cases {
						if c.Value == nil {
							continue
						}
						if l, ok, err := armLiteral(fset, src, c.Value, sideKeyed); err != nil {
							walkErr = err
							return false
						} else if ok {
							lists = append(lists, l)
						}
					}
				}
			}
		}
		return true
	})
	return lists, walkErr
}

func armLiteral(fset *token.FileSet, src []byte, arm *gsxast.ValueArm, sideKeyed bool) (classList, bool, error) {
	if arm.Segments != nil || arm.Expr == "" {
		return classList{}, false, nil // f`…` literal arms carry no class table entries in shipped source
	}
	return literalIn(fset, src, arm.Expr, arm.ExprPos, sideKeyed)
}

// literalIn parses one Go expression from a class list and, when it is a
// plain string literal (parentheses allowed), returns its span in src.
// go/parser.ParseExpr positions start at 1, so offset = exprOffset + pos - 1.
func literalIn(fset *token.FileSet, src []byte, expr string, exprPos token.Pos, sideKeyed bool) (classList, bool, error) {
	parsed, err := goparser.ParseExpr(expr)
	if err != nil {
		return classList{}, false, nil // not a literal (a helper call, a variable): nothing to rewrite
	}
	for {
		p, ok := parsed.(*goast.ParenExpr)
		if !ok {
			break
		}
		parsed = p.X
	}
	lit, ok := parsed.(*goast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return classList{}, false, nil
	}
	base := fset.Position(exprPos).Offset
	start := base + int(lit.Pos()) - 1
	end := base + int(lit.End()) - 1
	if start < 0 || end > len(src) || string(src[start:end]) != lit.Value {
		return classList{}, false, fmt.Errorf("class literal span mismatch at %v", fset.Position(exprPos))
	}
	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return classList{}, false, fmt.Errorf("class literal %s: %w", lit.Value, err)
	}
	quoting := goString
	if lit.Value[0] == '`' {
		quoting = goRawString
	}
	return classList{Value: value, Start: start, End: end, Quoting: quoting, SideKeyed: sideKeyed}, true, nil
}

// mentionsPlacement reports whether a switch tag or if condition is keyed on
// the component's physical-placement parameter, whose arms stay physical.
// Sheet and Sidebar name it `side`; Drawer names the same prop `direction`,
// after vaul's own API. Both identifiers count, and no other class-list
// switch or condition in ui/*.gsx uses either name for anything else.
func mentionsPlacement(goExpr string) bool {
	parsed, err := goparser.ParseExpr(goExpr)
	if err != nil {
		return false
	}
	found := false
	goast.Inspect(parsed, func(n goast.Node) bool {
		if id, ok := n.(*goast.Ident); ok && (id.Name == "side" || id.Name == "direction") {
			found = true
		}
		return !found
	})
	return found
}
