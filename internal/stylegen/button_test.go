package stylegen

import (
	"bytes"
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

var buttonVariantRoles = []buttonRoleCase{
	{list: `"", "default"`, token: "gsxui-recipe-button-variant-default"},
	{list: `"destructive"`, token: "gsxui-recipe-button-variant-destructive"},
	{list: `"outline"`, token: "gsxui-recipe-button-variant-outline"},
	{list: `"secondary"`, token: "gsxui-recipe-button-variant-secondary"},
	{list: `"ghost"`, token: "gsxui-recipe-button-variant-ghost"},
	{list: `"link"`, token: "gsxui-recipe-button-variant-link"},
}

var buttonSizeRoles = []buttonRoleCase{
	{list: `"", "default"`, token: "gsxui-recipe-button-size-default"},
	{list: `"xs"`, token: "gsxui-recipe-button-size-xs"},
	{list: `"sm"`, token: "gsxui-recipe-button-size-sm"},
	{list: `"lg"`, token: "gsxui-recipe-button-size-lg"},
	{list: `"icon"`, token: "gsxui-recipe-button-size-icon"},
	{list: `"icon-xs"`, token: "gsxui-recipe-button-size-icon-xs"},
	{list: `"icon-sm"`, token: "gsxui-recipe-button-size-icon-sm"},
	{list: `"icon-lg"`, token: "gsxui-recipe-button-size-icon-lg"},
}

type buttonRoleCase struct {
	list  string
	token string
}

func TestCanonicalButtonRecipeTokenContract(t *testing.T) {
	t.Parallel()

	filename, src := canonicalButtonSource(t)
	got, err := RecipeTokens(filename, src)
	if err != nil {
		t.Fatalf("RecipeTokens() error = %v", err)
	}
	if !reflect.DeepEqual(got, buttonRecipeTokens) {
		t.Errorf("RecipeTokens() = %q, want exact Button role set %q", got, buttonRecipeTokens)
	}
}

func TestCanonicalButtonBranchesUseIdenticalInlineRoleSwitches(t *testing.T) {
	t.Parallel()

	filename, src := canonicalButtonSource(t)
	anchorClass, buttonClass := resolvedElementClassExpressions(t, filename, src)
	if !bytes.Equal(anchorClass, buttonClass) {
		t.Errorf("canonical anchor/button class expressions differ\n--- anchor ---\n%s\n--- button ---\n%s", anchorClass, buttonClass)
	}

	classes := buttonClassAttrs(t, filename, src)
	for _, tag := range []string{"a", "button"} {
		class := classes[tag]
		if class == nil {
			t.Errorf("<%s> has no class expression", tag)
			continue
		}
		if len(class.Parts) != 4 {
			t.Errorf("<%s> class parts = %d, want structural marker + base + variant switch + size switch", tag, len(class.Parts))
			continue
		}
		assertButtonStructuralLiteral(t, tag+" structural marker", class.Parts[0].Expr)
		assertButtonRoleLiteral(t, tag+" base", class.Parts[1].Expr, "gsxui-recipe-button")
		assertButtonRoleSwitch(t, tag+" variant", class.Parts[2].CF, "variant", buttonVariantRoles)
		assertButtonRoleSwitch(t, tag+" size", class.Parts[3].CF, "size", buttonSizeRoles)
	}
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
			case *gsxast.ClassAttr:
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
			recipes, err := ParseRecipes(filename, src)
			if err != nil {
				t.Fatalf("ParseRecipes() error = %v", err)
			}
			got := recipes.Tokens()
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
			"gsxui-recipe-button-variant-default":     {"bg-primary", "text-primary-foreground", "hover:bg-primary/90"},
			"gsxui-recipe-button-variant-destructive": {"bg-destructive", "text-contrast", "dark:bg-destructive/60"},
			"gsxui-recipe-button-variant-outline":     {"border-border", "bg-background", "dark:bg-input/30"},
			"gsxui-recipe-button-variant-secondary":   {"bg-secondary", "hover:bg-secondary/80"},
			"gsxui-recipe-button-variant-ghost":       {"hover:bg-accent", "dark:hover:bg-accent/50"},
			"gsxui-recipe-button-variant-link":        {"text-primary", "hover:underline"},
			"gsxui-recipe-button-size-default":        {"h-8", "px-2.5", "has-[>svg]:px-2"},
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
			recipes, err := ParseRecipes(filename, src)
			if err != nil {
				t.Fatalf("ParseRecipes() error = %v", err)
			}
			for role, wantUtilities := range roles {
				recipe, ok := recipes.Lookup(role)
				if !ok {
					t.Errorf("recipe %q is missing", role)
					continue
				}
				for _, utility := range wantUtilities {
					if !containsString(recipe.Utilities, utility) {
						t.Errorf("recipe %q utilities = %q, want recognizable %q", role, recipe.Utilities, utility)
					}
				}
			}
		})
	}
}

func TestGenerateButtonArtifactsMatchCommittedResolvedSources(t *testing.T) {
	t.Parallel()

	root := repositoryRoot(t)
	if err := GenerateButton(root, true); err != nil {
		t.Fatalf("GenerateButton(check) error = %v", err)
	}

	canonicalPath := filepath.Join(root, "ui", "button.gsx")
	canonical, err := os.ReadFile(canonicalPath)
	if err != nil {
		t.Fatal(err)
	}
	generated := make(map[string][]byte)
	for _, style := range []string{"nova", "maia"} {
		recipePath := filepath.Join(root, "registry", "styles", style, "button.css")
		recipeSource, err := os.ReadFile(recipePath)
		if err != nil {
			t.Fatal(err)
		}
		recipes, err := ParseRecipes(recipePath, recipeSource)
		if err != nil {
			t.Fatalf("ParseRecipes(%s) error = %v", style, err)
		}
		want, report, err := Resolve(canonicalPath, canonical, recipes)
		if err != nil {
			t.Fatalf("Resolve(%s) error = %v", style, err)
		}
		if !reflect.DeepEqual(report.UsedTokens, buttonRecipeTokens) {
			t.Errorf("Resolve(%s) used tokens = %q, want %q", style, report.UsedTokens, buttonRecipeTokens)
		}

		generatedPath := filepath.Join(root, "registry", "generated", style, "button.gsx")
		got, err := os.ReadFile(generatedPath)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s committed artifact differs from resolved canonical source\n--- got ---\n%s\n--- want ---\n%s", style, got, want)
		}
		assertGeneratedButtonSource(t, generatedPath, got)
		generated[style] = got
	}

	assertButtonStylesDifferVisibly(t, generated["nova"], generated["maia"])
}

func TestGenerateButtonWritesDeterministicallyAndCheckNeverWrites(t *testing.T) {
	t.Parallel()

	root := buttonGeneratorFixture(t)
	if err := GenerateButton(root, false); err != nil {
		t.Fatalf("GenerateButton(write) error = %v", err)
	}
	first := readGeneratedButtons(t, root)

	if err := GenerateButton(root, false); err != nil {
		t.Fatalf("GenerateButton(second write) error = %v", err)
	}
	second := readGeneratedButtons(t, root)
	if !reflect.DeepEqual(second, first) {
		t.Error("repeated generation changed Button artifacts")
	}
	if err := GenerateButton(root, true); err != nil {
		t.Fatalf("GenerateButton(clean check) error = %v", err)
	}

	maiaPath := filepath.Join(root, "registry", "generated", "maia", "button.gsx")
	if err := os.WriteFile(maiaPath, []byte("locally modified\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeCheck := readGeneratedButtons(t, root)
	err := GenerateButton(root, true)
	if err == nil {
		t.Fatal("GenerateButton(drift check) error = nil, want drift")
	}
	if !strings.Contains(err.Error(), filepath.ToSlash("registry/generated/maia/button.gsx")) {
		t.Errorf("GenerateButton(drift check) error %q does not identify Maia artifact", err)
	}
	afterCheck := readGeneratedButtons(t, root)
	if !reflect.DeepEqual(afterCheck, beforeCheck) {
		t.Error("check mode mutated generated artifacts")
	}
}

func TestGenerateButtonResolvesEveryStyleBeforeWriting(t *testing.T) {
	t.Parallel()

	root := buttonGeneratorFixture(t)
	for _, style := range []string{"nova", "maia"} {
		path := filepath.Join(root, "registry", "generated", style, "button.gsx")
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(style+" sentinel\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	before := readGeneratedButtons(t, root)

	maiaRecipe := filepath.Join(root, "registry", "styles", "maia", "button.css")
	if err := os.WriteFile(maiaRecipe, []byte("@layer components {"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := GenerateButton(root, false); err == nil {
		t.Fatal("GenerateButton() error = nil, want Maia recipe failure")
	}
	after := readGeneratedButtons(t, root)
	if !reflect.DeepEqual(after, before) {
		t.Error("generation mutated an artifact before all styles resolved")
	}
}

func canonicalButtonSource(t *testing.T) (string, []byte) {
	t.Helper()

	filename := filepath.Join("..", "..", "ui", "button.gsx")
	src, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return filename, src
}

func repositoryRoot(t *testing.T) string {
	t.Helper()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return root
}

func buttonGeneratorFixture(t *testing.T) string {
	t.Helper()

	sourceRoot := repositoryRoot(t)
	root := t.TempDir()
	for _, relative := range []string{
		filepath.Join("ui", "button.gsx"),
		filepath.Join("registry", "styles", "nova", "button.css"),
		filepath.Join("registry", "styles", "maia", "button.css"),
	} {
		src, err := os.ReadFile(filepath.Join(sourceRoot, relative))
		if err != nil {
			t.Fatal(err)
		}
		dst := filepath.Join(root, relative)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(dst, src, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
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

	if bytes.Contains(src, []byte(RecipePrefix)) {
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

	classes := buttonClassAttrs(t, filename, src)
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

func buttonClassAttrs(t *testing.T, filename string, src []byte) map[string]*gsxast.ClassAttr {
	t.Helper()

	fset := token.NewFileSet()
	file, err := gsxparser.ParseFile(fset, filename, src, 0)
	if err != nil {
		t.Fatalf("parser.ParseFile() error = %v", err)
	}
	classes := make(map[string]*gsxast.ClassAttr)
	gsxast.Inspect(file, func(node gsxast.Node) bool {
		element, ok := node.(*gsxast.Element)
		if !ok || (element.Tag != "a" && element.Tag != "button") {
			return true
		}
		for _, attr := range element.Attrs {
			class, ok := attr.(*gsxast.ClassAttr)
			if ok && class.Name == "class" {
				classes[element.Tag] = class
			}
		}
		return true
	})
	return classes
}

func assertButtonRoleSwitch(t *testing.T, name string, cf *gsxast.ValueCF, tag string, want []buttonRoleCase) {
	t.Helper()

	if cf == nil || cf.Switch == nil {
		t.Errorf("%s is not an inline value switch", name)
		return
	}
	switchNode := cf.Switch
	if switchNode.Tag != tag {
		t.Errorf("%s tag = %q, want %q", name, switchNode.Tag, tag)
	}
	if len(switchNode.Cases) != len(want)+1 {
		t.Errorf("%s cases = %d, want %d role cases plus fallback", name, len(switchNode.Cases), len(want)+1)
		return
	}
	for i, role := range want {
		caseNode := switchNode.Cases[i]
		if caseNode.Default {
			t.Errorf("%s case %d is default, want %s", name, i, role.list)
			continue
		}
		if caseNode.List != role.list {
			t.Errorf("%s case %d list = %q, want %q", name, i, caseNode.List, role.list)
		}
		assertButtonRoleLiteral(t, name+" case "+role.list, caseNode.Value.Expr, role.token)
	}
	fallback := switchNode.Cases[len(switchNode.Cases)-1]
	if !fallback.Default {
		t.Errorf("%s final case is not default", name)
		return
	}
	assertButtonRoleLiteral(t, name+" fallback", fallback.Value.Expr, "")
}

func assertButtonRoleLiteral(t *testing.T, name, expr, want string) {
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
	if got != want {
		t.Errorf("%s = %q, want %q", name, got, want)
	}
	if got != "" && !strings.HasPrefix(got, RecipePrefix) {
		t.Errorf("%s contains concrete presentation utility %q", name, got)
	}
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
