package stylegen

import (
	"bytes"
	"errors"
	"fmt"
	goparser "go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	gsxparser "github.com/gsxhq/gsx/parser"

	"github.com/gsxhq/gsxui/internal/recipe"
	"github.com/gsxhq/gsxui/merge"
	"github.com/gsxhq/gsxui/registry/canonical/shapes"
)

// DefaultStyle names the style whose generated output ships as package ui. It
// is the single place the shipped style is chosen.
const DefaultStyle = "nova"

// consumerPackage is the package clause every generated component carries.
const consumerPackage = "ui"

type generatedSource struct {
	relativePath string
	content      []byte
}

// GenerateAll resolves every canonical component against every authored style
// before it checks or mutates any artifact. A failing run writes nothing.
func GenerateAll(root string, check bool) error {
	outputs, err := resolveAll(root)
	if err != nil {
		return err
	}
	if check {
		return checkOutputs(root, outputs)
	}
	for _, output := range outputs {
		if err := writeGenerated(filepath.Join(root, output.relativePath), output.content); err != nil {
			return fmt.Errorf("write %s: %w", filepath.ToSlash(output.relativePath), err)
		}
	}
	return nil
}

// resolveAll discovers every canonical component and every authored style and
// resolves the full cross product in memory. Every parse, conformance check,
// conflict check and resolve completes here, so GenerateAll never writes a
// partial result.
func resolveAll(root string) ([]generatedSource, error) {
	components, err := canonicalComponents(root)
	if err != nil {
		return nil, err
	}
	styles, err := styleNames(root)
	if err != nil {
		return nil, err
	}
	declared := shapes.All()

	shapesByComponent := make(map[string]recipe.Shape, len(components))
	resolvedByStyle := make(map[string]map[string]recipe.Resolved, len(styles))
	for _, style := range styles {
		resolvedByStyle[style] = make(map[string]recipe.Resolved, len(components))
	}

	var outputs []generatedSource
	for _, component := range components {
		canonicalPath := filepath.Join(root, "registry", "canonical", component+".gsx")
		src, err := os.ReadFile(canonicalPath)
		if err != nil {
			return nil, fmt.Errorf("read canonical %s: %w", component, err)
		}
		shape, ok := declared[component]
		if !ok {
			return nil, fmt.Errorf("canonical %s declares no shape in registry/canonical", component)
		}
		if err := checkHelperCalls(canonicalPath, src, shape); err != nil {
			return nil, err
		}
		shapesByComponent[component] = shape

		// The typed accessor set is a generated artifact like any other, so it
		// flows through the same validation-before-mutation and --check drift
		// machinery as the resolved sources.
		accessors, err := GenerateAccessors(shape)
		if err != nil {
			return nil, fmt.Errorf("generate %s accessors: %w", component, err)
		}
		outputs = append(outputs, generatedSource{
			relativePath: filepath.Join("registry", "canonical", component+"_recipe.gen.go"),
			content:      accessors,
		})

		for _, style := range styles {
			stylePath := filepath.Join(root, "registry", "styles", style, component+".css")
			styleSource, err := os.ReadFile(stylePath)
			if err != nil {
				return nil, fmt.Errorf("read %s %s recipe: %w", style, component, err)
			}
			parsed, err := recipe.ParseStyle(stylePath, styleSource)
			if err != nil {
				return nil, fmt.Errorf("parse %s %s recipe: %w", style, component, err)
			}
			resolved, err := recipe.Conform(stylePath, shape, parsed)
			if err != nil {
				return nil, fmt.Errorf("conform %s %s recipe: %w", style, component, err)
			}
			if err := recipe.CheckConflicts(stylePath, resolved, merge.Merge); err != nil {
				return nil, fmt.Errorf("check %s %s recipe: %w", style, component, err)
			}
			resolvedByStyle[style][component] = resolved
			desugared, err := Resolve(canonicalPath, src, resolved)
			if err != nil {
				return nil, fmt.Errorf("resolve %s %s: %w", style, component, err)
			}
			// The canonical is authored in package canonical, which is never
			// shipped; consumers receive package ui and the previews receive
			// one package per style.
			generated, err := rewriteGSXPackage(canonicalPath, desugared, consumerPackage)
			if err != nil {
				return nil, fmt.Errorf("derive %s %s consumer source: %w", style, component, err)
			}
			preview, err := rewriteGSXPackage(canonicalPath, desugared, style)
			if err != nil {
				return nil, fmt.Errorf("derive %s %s preview: %w", style, component, err)
			}
			outputs = append(outputs, generatedSource{
				relativePath: filepath.Join("registry", "generated", style, component+".gsx"),
				content:      generated,
			})
			outputs = append(outputs, generatedSource{
				relativePath: filepath.Join("site", "stylepreview", style, component+".gsx"),
				content:      preview,
			})
			if style == DefaultStyle {
				// ui/ is the default style's output, not a fourth artifact kind.
				outputs = append(outputs, generatedSource{
					relativePath: filepath.Join("ui", component+".gsx"),
					content:      generated,
				})
			}
		}
	}

	contract := recipe.BuildContract(shapesByComponent, resolvedByStyle)
	contractJSON, err := contract.MarshalIndent()
	if err != nil {
		return nil, fmt.Errorf("marshal recipe contract: %w", err)
	}
	outputs = append(outputs, generatedSource{
		relativePath: filepath.Join("registry", "generated", "recipes.json"),
		content:      contractJSON,
	})

	return outputs, nil
}

// checkHelperCalls cross-checks a canonical's accessor calls against its
// declared shape, in both directions. It runs before any style is touched, so
// either mismatch fails generation rather than producing an artifact:
//
//   - every accessor call resolves to a declared slot and dimension, and
//   - every slot, base rule and dimension the shape declares is actually
//     reached by the component's GSX.
//
// The second direction is what keeps a shape honest as slot counts grow. A
// declared-but-unrendered slot costs every style a set of rules conformance
// forces it to author, all of them dead, and nothing else in the pipeline
// notices: the style conforms, the contract lists the slot, and no element
// ever carries the utilities.
func checkHelperCalls(filename string, src []byte, shape recipe.Shape) error {
	calls, err := HelperCalls(filename, src, shape)
	if err != nil {
		return fmt.Errorf("scan %s accessor calls: %w", filepath.Base(filename), err)
	}
	for _, call := range calls {
		if call.Component != shape.Component {
			return fmt.Errorf(
				"%s: recipe accessor call on component %q, but this source declares component %q",
				filename, call.Component, shape.Component,
			)
		}
		slot, ok := shape.Slot(call.Slot)
		if !ok {
			return fmt.Errorf(
				"%s: component %q declares no slot %q",
				filename, shape.Component, call.Slot,
			)
		}
		if call.Dimension == "" {
			if !slot.Base {
				return fmt.Errorf(
					"%s: component %q slot %q declares no base rule",
					filename, shape.Component, call.Slot,
				)
			}
			continue
		}
		if _, ok := slot.Dimension(call.Dimension); !ok {
			return fmt.Errorf(
				"%s: component %q slot %q declares no dimension %q",
				filename, shape.Component, call.Slot, call.Dimension,
			)
		}
	}
	return checkShapeCoverage(filename, shape, calls)
}

// checkShapeCoverage reports every part of the shape the component's GSX never
// reaches. It reports all of them at once: converging a 38-slot component one
// failed run at a time is not a workflow.
func checkShapeCoverage(filename string, shape recipe.Shape, calls []Call) error {
	usedSlots := make(map[string]struct{}, len(calls))
	usedBases := make(map[string]struct{}, len(calls))
	usedDimensions := make(map[[2]string]struct{}, len(calls))
	for _, call := range calls {
		usedSlots[call.Slot] = struct{}{}
		if call.Dimension == "" {
			usedBases[call.Slot] = struct{}{}
			continue
		}
		usedDimensions[[2]string{call.Slot, call.Dimension}] = struct{}{}
	}

	var gaps []string
	for _, slot := range shape.Slots {
		if _, ok := usedSlots[slot.Name]; !ok {
			gaps = append(gaps, fmt.Sprintf(
				"%s: shape declares slot %q but the component never renders it",
				filename, slot.Name))
			continue
		}
		if slot.Base {
			if _, ok := usedBases[slot.Name]; !ok {
				gaps = append(gaps, fmt.Sprintf(
					"%s: shape declares a base rule on slot %q but the component never applies it",
					filename, slot.Name))
			}
		}
		for _, dimension := range slot.Dimensions {
			if _, ok := usedDimensions[[2]string{slot.Name, dimension.Name}]; ok {
				continue
			}
			gaps = append(gaps, fmt.Sprintf(
				"%s: shape declares dimension %q on slot %q but the component never applies it",
				filename, dimension.Name, slot.Name))
		}
	}
	if len(gaps) == 0 {
		return nil
	}
	return errors.New(strings.Join(gaps, "\n"))
}

// canonicalComponents lists every component name authored under
// registry/canonical, sorted, so generation order is deterministic.
func canonicalComponents(root string) ([]string, error) {
	dir := filepath.Join(root, "registry", "canonical")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read canonical sources: %w", err)
	}
	var components []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".gsx" {
			continue
		}
		components = append(components, strings.TrimSuffix(entry.Name(), ".gsx"))
	}
	if len(components) == 0 {
		return nil, fmt.Errorf("no canonical components found in %s", filepath.ToSlash(filepath.Join("registry", "canonical")))
	}
	sort.Strings(components)
	return components, nil
}

// styleNames lists every authored style, sorted. The default style must exist:
// without it there is nothing to ship as package ui.
func styleNames(root string) ([]string, error) {
	dir := filepath.Join(root, "registry", "styles")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read styles: %w", err)
	}
	var styles []string
	for _, entry := range entries {
		if entry.IsDir() {
			styles = append(styles, entry.Name())
		}
	}
	sort.Strings(styles)
	if !slices.Contains(styles, DefaultStyle) {
		return nil, fmt.Errorf("default style %q is not authored under registry/styles", DefaultStyle)
	}
	return styles, nil
}

func rewriteGSXPackage(filename string, src []byte, packageName string) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := goparser.ParseFile(fset, filename, src, goparser.PackageClauseOnly)
	if err != nil {
		return nil, fmt.Errorf("parse package declaration: %w", err)
	}
	start := fset.Position(file.Name.Pos()).Offset
	end := fset.Position(file.Name.End()).Offset
	if start < 0 || end < start || end > len(src) {
		return nil, fmt.Errorf("package declaration range [%d:%d] is outside %d-byte source", start, end, len(src))
	}

	result := make([]byte, 0, len(src)-end+start+len(packageName))
	result = append(result, src[:start]...)
	result = append(result, packageName...)
	result = append(result, src[end:]...)

	parsed, err := gsxparser.ParseFile(token.NewFileSet(), filename, result, 0)
	if err != nil {
		return nil, fmt.Errorf("validate rewritten GSX: %w", err)
	}
	if parsed.Package != packageName {
		return nil, fmt.Errorf("rewritten package = %q, want %q", parsed.Package, packageName)
	}
	return result, nil
}

func checkOutputs(root string, outputs []generatedSource) error {
	var drift []string
	for _, output := range outputs {
		path := filepath.Join(root, output.relativePath)
		got, err := os.ReadFile(path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("read generated source %s: %w", filepath.ToSlash(output.relativePath), err)
			}
			drift = append(drift, filepath.ToSlash(output.relativePath)+" (missing)")
			continue
		}
		if !bytes.Equal(got, output.content) {
			drift = append(drift, filepath.ToSlash(output.relativePath))
		}
	}
	if len(drift) != 0 {
		return fmt.Errorf(
			"generated style sources are out of date: %s; run go run ./cmd/stylegen",
			strings.Join(drift, ", "),
		)
	}
	return nil
}

func writeGenerated(path string, content []byte) error {
	if got, err := os.ReadFile(path); err == nil && bytes.Equal(got, content) {
		return nil
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".gsx-*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)

	if err := temp.Chmod(0o644); err != nil {
		temp.Close()
		return err
	}
	if _, err := temp.Write(content); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Sync(); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tempPath, path); err != nil {
		return err
	}
	return nil
}
