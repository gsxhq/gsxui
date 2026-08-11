package stylegen

import (
	"bytes"
	"encoding/json"
	"fmt"
	goast "go/ast"
	goparser "go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"

	gsxast "github.com/gsxhq/gsx/ast"
	"github.com/gsxhq/gsx/gen"
	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/recipe"
)

var buttonRecipeTokens = []string{
	"gsxui-recipe-button",
	"gsxui-recipe-button-size-default",
	"gsxui-recipe-button-size-icon",
	"gsxui-recipe-button-size-icon-lg",
	"gsxui-recipe-button-size-icon-sm",
	"gsxui-recipe-button-size-icon-xs",
	"gsxui-recipe-button-size-lg",
	"gsxui-recipe-button-size-sm",
	"gsxui-recipe-button-size-xs",
	"gsxui-recipe-button-variant-default",
	"gsxui-recipe-button-variant-destructive",
	"gsxui-recipe-button-variant-ghost",
	"gsxui-recipe-button-variant-link",
	"gsxui-recipe-button-variant-outline",
	"gsxui-recipe-button-variant-secondary",
}

func TestCanonicalButtonKeepsAxesMarkerAndCallerPrecedence(t *testing.T) {
	t.Parallel()

	filename, src := canonicalButtonSource(t)
	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile() error = %v", err)
	}

	found := make(map[string]bool)
	gsxast.Inspect(file, func(node gsxast.Node) bool {
		element, ok := node.(*gsxast.Element)
		if !ok || (element.Tag != "a" && element.Tag != "button") {
			return true
		}
		found[element.Tag] = true

		classIndex := -1
		spreadIndex := -1
		var variant, size *gsxast.ExprAttr
		var marker *gsxast.BoolAttr
		for i, attr := range element.Attrs {
			switch attr := attr.(type) {
			case *gsxast.ComposedAttr:
				if attr.Name == "class" {
					classIndex = i
				}
			case *gsxast.SpreadAttr:
				if attr.Expr == "attrs" {
					spreadIndex = i
				}
			case *gsxast.ExprAttr:
				switch attr.Name {
				case "data-variant":
					variant = attr
				case "data-size":
					size = attr
				}
			case *gsxast.BoolAttr:
				if attr.Name == "data-gsxui-slot-button" {
					marker = attr
				}
			}
		}

		if classIndex < 0 {
			t.Errorf("<%s> has no class expression", element.Tag)
		}
		if spreadIndex != classIndex+1 {
			t.Errorf("<%s> caller attrs index = %d, want immediately after class index %d", element.Tag, spreadIndex, classIndex)
		}
		if marker == nil || marker.Pos() < element.Attrs[spreadIndex].End() {
			t.Errorf("<%s> slot marker must be bare and authored after caller attrs", element.Tag)
		}
		assertDefaultedAxis(t, element.Tag, variant, "variant")
		assertDefaultedAxis(t, element.Tag, size, "size")
		return true
	})
	for _, tag := range []string{"a", "button"} {
		if !found[tag] {
			t.Errorf("canonical Button has no <%s> branch", tag)
		}
	}
}

func TestButtonRecipesDeclareExactCanonicalRoleSet(t *testing.T) {
	t.Parallel()

	for _, style := range []string{"nova", "maia"} {
		t.Run(style, func(t *testing.T) {
			t.Parallel()

			filename := filepath.Join("..", "..", "registry", "styles", style, "button.css")
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			recipes, err := recipe.ParseStyle(filename, src)
			if err != nil {
				t.Fatalf("ParseStyle() error = %v", err)
			}
			got := recipes.Classes()
			sort.Strings(got)
			if !reflect.DeepEqual(got, buttonRecipeTokens) {
				t.Errorf("%s recipe tokens = %q, want exact canonical role set %q", style, got, buttonRecipeTokens)
			}
		})
	}
}

func TestButtonRecipesCarryRecognizableConcretePresentationForEveryRole(t *testing.T) {
	t.Parallel()

	tests := map[string]map[string][]string{
		"nova": {
			"gsxui-recipe-button":                     {"inline-flex", "rounded-lg", "border-transparent", "disabled:opacity-50"},
			"gsxui-recipe-button-variant-default":     {"bg-primary", "text-primary-foreground", "hover:bg-primary/80"},
			"gsxui-recipe-button-variant-destructive": {"bg-destructive/10", "text-destructive", "dark:bg-destructive/20"},
			"gsxui-recipe-button-variant-outline":     {"border-border", "bg-background", "dark:bg-input/30"},
			"gsxui-recipe-button-variant-secondary":   {"bg-secondary", "hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)]"},
			"gsxui-recipe-button-variant-ghost":       {"hover:bg-muted", "dark:hover:bg-muted/50"},
			"gsxui-recipe-button-variant-link":        {"text-primary", "hover:underline"},
			"gsxui-recipe-button-size-default":        {"h-8", "px-2.5", "has-data-[icon=inline-start]:pl-2"},
			"gsxui-recipe-button-size-xs":             {"h-6", "rounded-[min(var(--radius-md),10px)]", "text-xs"},
			"gsxui-recipe-button-size-sm":             {"h-7", "rounded-[min(var(--radius-md),12px)]", "text-[0.8rem]"},
			"gsxui-recipe-button-size-lg":             {"h-9", "px-2.5"},
			"gsxui-recipe-button-size-icon":           {"size-8"},
			"gsxui-recipe-button-size-icon-xs":        {"size-6", "rounded-[min(var(--radius-md),10px)]"},
			"gsxui-recipe-button-size-icon-sm":        {"size-7", "rounded-[min(var(--radius-md),12px)]"},
			"gsxui-recipe-button-size-icon-lg":        {"size-9"},
		},
		"maia": {
			"gsxui-recipe-button":                     {"inline-flex", "rounded-4xl", "border-transparent", "active:not-aria-[haspopup]:translate-y-px"},
			"gsxui-recipe-button-variant-default":     {"bg-primary", "text-primary-foreground", "hover:bg-primary/80"},
			"gsxui-recipe-button-variant-destructive": {"bg-destructive/10", "text-destructive", "dark:hover:bg-destructive/30"},
			"gsxui-recipe-button-variant-outline":     {"border-border", "bg-input/30", "aria-expanded:bg-muted"},
			"gsxui-recipe-button-variant-secondary":   {"bg-secondary", "hover:bg-[color-mix(in_oklch,var(--secondary),var(--foreground)_5%)]"},
			"gsxui-recipe-button-variant-ghost":       {"hover:bg-muted", "dark:hover:bg-muted/50"},
			"gsxui-recipe-button-variant-link":        {"text-primary", "hover:underline"},
			"gsxui-recipe-button-size-default":        {"h-9", "px-3", "has-data-[icon=inline-start]:pl-2.5"},
			"gsxui-recipe-button-size-xs":             {"h-6", "px-2.5", "text-xs"},
			"gsxui-recipe-button-size-sm":             {"h-8", "px-3"},
			"gsxui-recipe-button-size-lg":             {"h-10", "px-4"},
			"gsxui-recipe-button-size-icon":           {"size-9"},
			"gsxui-recipe-button-size-icon-xs":        {"size-6", "[&_svg:not([class*='size-'])]:size-3"},
			"gsxui-recipe-button-size-icon-sm":        {"size-8"},
			"gsxui-recipe-button-size-icon-lg":        {"size-10"},
		},
	}

	for style, roles := range tests {
		t.Run(style, func(t *testing.T) {
			t.Parallel()

			filename := filepath.Join("..", "..", "registry", "styles", style, "button.css")
			src, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			recipes, err := recipe.ParseStyle(filename, src)
			if err != nil {
				t.Fatalf("ParseStyle() error = %v", err)
			}
			for role, wantUtilities := range roles {
				rule, ok := recipes.Lookup(role)
				if !ok {
					t.Errorf("recipe %q is missing", role)
					continue
				}
				for _, utility := range wantUtilities {
					if !containsString(rule.Utilities, utility) {
						t.Errorf("recipe %q utilities = %q, want recognizable %q", role, rule.Utilities, utility)
					}
				}
			}
		})
	}
}

func TestGeneratedPreviewsMatchRegistryAST(t *testing.T) {
	t.Parallel()

	root := repoRoot(t)
	for _, style := range []string{"nova", "maia"} {
		registryPath := filepath.Join(root, "registry", "generated", style, "button.gsx")
		previewPath := filepath.Join(root, "site", "stylepreview", style, "button.gsx")

		registryAST := comparableButtonAST(t, registryPath, "preview")
		previewAST := comparableButtonAST(t, previewPath, "preview")
		if !bytes.Equal(previewAST, registryAST) {
			t.Errorf("%s preview fixture differs from registry artifact after package normalization\n--- preview AST ---\n%s\n--- registry AST ---\n%s", style, previewAST, registryAST)
		}
	}
}

func TestGenerateAllWritesDeterministicallyAndCheckNeverWrites(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	copyRepoFixture(t, root)
	if err := GenerateAll(root, false); err != nil {
		t.Fatalf("GenerateAll(write) error = %v", err)
	}
	first := readGeneratedButtons(t, root)

	if err := GenerateAll(root, false); err != nil {
		t.Fatalf("GenerateAll(second write) error = %v", err)
	}
	second := readGeneratedButtons(t, root)
	if !reflect.DeepEqual(second, first) {
		t.Error("repeated generation changed Button artifacts")
	}
	if err := GenerateAll(root, true); err != nil {
		t.Fatalf("GenerateAll(clean check) error = %v", err)
	}

	maiaPath := filepath.Join(root, "registry", "generated", "maia", "button.gsx")
	if err := os.WriteFile(maiaPath, []byte("locally modified\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeCheck := readGeneratedButtons(t, root)
	err := GenerateAll(root, true)
	if err == nil {
		t.Fatal("GenerateAll(drift check) error = nil, want drift")
	}
	if !strings.Contains(err.Error(), filepath.ToSlash("registry/generated/maia/button.gsx")) {
		t.Errorf("GenerateAll(drift check) error %q does not identify Maia artifact", err)
	}
	afterCheck := readGeneratedButtons(t, root)
	if !reflect.DeepEqual(afterCheck, beforeCheck) {
		t.Error("check mode mutated generated artifacts")
	}

	if err := GenerateAll(root, false); err != nil {
		t.Fatalf("GenerateAll(restore) error = %v", err)
	}
	previewPath := filepath.Join(root, "site", "stylepreview", "nova", "button.gsx")
	previewDrift := []byte("locally modified preview\n")
	if err := os.WriteFile(previewPath, previewDrift, 0o644); err != nil {
		t.Fatal(err)
	}
	err = GenerateAll(root, true)
	if err == nil {
		t.Fatal("GenerateAll(preview drift check) error = nil, want drift")
	}
	if !strings.Contains(err.Error(), filepath.ToSlash("site/stylepreview/nova/button.gsx")) {
		t.Errorf("GenerateAll(preview drift check) error %q does not identify Nova preview artifact", err)
	}
	afterPreviewCheck, readErr := os.ReadFile(previewPath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if !bytes.Equal(afterPreviewCheck, previewDrift) {
		t.Error("check mode mutated generated preview artifact")
	}
}

func canonicalButtonSource(t *testing.T) (string, []byte) {
	t.Helper()

	filename := filepath.Join("..", "..", "registry", "canonical", "button.gsx")
	src, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return filename, src
}

func readGeneratedButtons(t *testing.T, root string) map[string][]byte {
	t.Helper()

	result := make(map[string][]byte)
	for _, style := range []string{"nova", "maia"} {
		path := filepath.Join(root, "registry", "generated", style, "button.gsx")
		src, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		result[style] = src
	}
	return result
}

func assertGeneratedButtonSource(t *testing.T, filename string, src []byte) {
	t.Helper()

	if bytes.Contains(src, []byte(recipe.Prefix)) {
		t.Errorf("%s contains recipe tokens", filename)
	}
	for _, helper := range []string{"variantClass", "sizeClass", "buttonClass"} {
		if bytes.Contains(src, []byte(helper)) {
			t.Errorf("%s contains private styling helper %q", filename, helper)
		}
	}

	formatted, err := gen.Format(filename, src)
	if err != nil {
		t.Fatalf("gen.Format(%s) error = %v", filename, err)
	}
	if !bytes.Equal(formatted, src) {
		t.Errorf("%s is not idempotently formatted", filename)
	}
	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile(%s) error = %v", filename, err)
	}
	var button *gsxast.Component
	gsxast.Inspect(file, func(node gsxast.Node) bool {
		component, ok := node.(*gsxast.Component)
		if ok && component.Name == "Button" {
			button = component
		}
		return true
	})
	if button == nil {
		t.Errorf("%s has no Button component", filename)
		return
	}
	const wantParams = "variant string, size string, href string, disabled bool, children gsx.Node, attrs gsx.Attrs"
	if button.Params != wantParams {
		t.Errorf("%s Button params = %q, want %q", filename, button.Params, wantParams)
	}

	classes := buttonComposedAttrs(t, filename, src)
	for _, tag := range []string{"a", "button"} {
		class := classes[tag]
		if class == nil {
			t.Errorf("%s has no <%s> class expression", filename, tag)
			continue
		}
		if len(class.Parts) == 0 {
			t.Errorf("%s <%s> class expression is empty", filename, tag)
			continue
		}
		assertButtonStructuralLiteral(t, filename+" <"+tag+"> structural marker", class.Parts[0].Expr)
	}
}

func comparableButtonAST(t *testing.T, filename, packageName string) []byte {
	t.Helper()

	src, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	file, err := gsxparser.ParseFile(token.NewFileSet(), filename, src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile(%s) error = %v", filename, err)
	}
	file.Package = packageName
	encoded, err := json.MarshalIndent(comparableASTValue(reflect.ValueOf(file)), "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent(%s AST) error = %v", filename, err)
	}
	return encoded
}

var tokenPositionType = reflect.TypeOf(token.Pos(0))

func comparableASTValue(value reflect.Value) any {
	if !value.IsValid() {
		return nil
	}
	if value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		return comparableASTValue(value.Elem())
	}
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		return comparableASTValue(value.Elem())
	}
	if value.Type() == tokenPositionType {
		return nil
	}

	switch value.Kind() {
	case reflect.Struct:
		result := map[string]any{"$type": value.Type().String()}
		for i := range value.NumField() {
			field := value.Type().Field(i)
			if field.PkgPath != "" || field.Type == tokenPositionType {
				continue
			}
			result[field.Name] = comparableASTValue(value.Field(i))
		}
		return result
	case reflect.Slice, reflect.Array:
		result := make([]any, value.Len())
		for i := range value.Len() {
			result[i] = comparableASTValue(value.Index(i))
		}
		return result
	case reflect.Map:
		result := make(map[string]any, value.Len())
		iter := value.MapRange()
		for iter.Next() {
			result[fmt.Sprint(iter.Key().Interface())] = comparableASTValue(iter.Value())
		}
		return result
	default:
		return value.Interface()
	}
}

func assertButtonStylesDifferVisibly(t *testing.T, nova, maia []byte) {
	t.Helper()

	checks := []struct {
		name string
		nova string
		maia string
	}{
		{name: "default geometry", nova: "h-8", maia: "h-9"},
		{name: "large geometry", nova: "h-9", maia: "h-10"},
		{name: "radius", nova: "rounded-lg", maia: "rounded-4xl"},
		{name: "outline background", nova: "bg-background", maia: "bg-input/30"},
		{name: "dark outline border", nova: "dark:border-input", maia: "dark:aria-invalid:border-destructive/50"},
	}
	for _, check := range checks {
		if !bytes.Contains(nova, []byte(check.nova)) {
			t.Errorf("Nova output lacks %s marker %q", check.name, check.nova)
		}
		if !bytes.Contains(maia, []byte(check.maia)) {
			t.Errorf("Maia output lacks %s marker %q", check.name, check.maia)
		}
	}
	if bytes.Equal(nova, maia) {
		t.Error("Nova and Maia generated Button sources are byte-identical")
	}
}

func buttonComposedAttrs(t *testing.T, filename string, src []byte) map[string]*gsxast.ComposedAttr {
	t.Helper()

	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile() error = %v", err)
	}
	classes := make(map[string]*gsxast.ComposedAttr)
	gsxast.Inspect(file, func(node gsxast.Node) bool {
		element, ok := node.(*gsxast.Element)
		if !ok || (element.Tag != "a" && element.Tag != "button") {
			return true
		}
		for _, attr := range element.Attrs {
			class, ok := attr.(*gsxast.ComposedAttr)
			if ok && class.Name == "class" {
				classes[element.Tag] = class
			}
		}
		return true
	})
	return classes
}

func assertButtonStructuralLiteral(t *testing.T, name, expr string) {
	t.Helper()

	parsed, err := goparser.ParseExpr(expr)
	if err != nil {
		t.Errorf("%s expression %q does not parse: %v", name, expr, err)
		return
	}
	literal, ok := parsed.(*goast.BasicLit)
	if !ok || literal.Kind != token.STRING {
		t.Errorf("%s expression = %q, want one string literal", name, expr)
		return
	}
	got, err := strconv.Unquote(literal.Value)
	if err != nil {
		t.Errorf("%s expression %q does not unquote: %v", name, expr, err)
		return
	}
	if got != "group/button" {
		t.Errorf("%s = %q, want sole invariant structural marker %q", name, got, "group/button")
	}
}

func assertDefaultedAxis(t *testing.T, tag string, attr *gsxast.ExprAttr, axis string) {
	t.Helper()

	if attr == nil {
		t.Errorf("<%s> has no persistent data-%s expression", tag, axis)
		return
	}
	if attr.Expr != axis {
		t.Errorf("<%s> data-%s expression = %q, want %q", tag, axis, attr.Expr, axis)
	}
	if len(attr.Stages) != 1 || attr.Stages[0].Name != "default" || attr.Stages[0].Args != `"default"` {
		t.Errorf("<%s> data-%s stages = %+v, want default(\"default\")", tag, axis, attr.Stages)
	}
}

func containsString(values []string, want string) bool {
	return slices.Contains(values, want)
}
