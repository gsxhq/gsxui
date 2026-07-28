package stylegen

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type buttonStyleSource struct {
	name       string
	recipePath string
	outputPath string
}

var buttonStyleSources = []buttonStyleSource{
	{
		name:       "nova",
		recipePath: filepath.Join("registry", "styles", "nova", "button.css"),
		outputPath: filepath.Join("registry", "generated", "nova", "button.gsx"),
	},
	{
		name:       "maia",
		recipePath: filepath.Join("registry", "styles", "maia", "button.css"),
		outputPath: filepath.Join("registry", "generated", "maia", "button.gsx"),
	},
}

type generatedButtonSource struct {
	relativePath string
	content      []byte
}

// GenerateButton resolves the canonical Button against every authored style
// recipe before it checks or mutates any generated artifact.
func GenerateButton(root string, check bool) error {
	outputs, err := resolveButtonSources(root)
	if err != nil {
		return err
	}
	if check {
		return checkButtonSources(root, outputs)
	}
	for _, output := range outputs {
		if err := writeGeneratedButton(filepath.Join(root, output.relativePath), output.content); err != nil {
			return fmt.Errorf("write %s: %w", filepath.ToSlash(output.relativePath), err)
		}
	}
	return nil
}

func resolveButtonSources(root string) ([]generatedButtonSource, error) {
	canonicalPath := filepath.Join(root, "ui", "button.gsx")
	canonical, err := os.ReadFile(canonicalPath)
	if err != nil {
		return nil, fmt.Errorf("read canonical Button %s: %w", canonicalPath, err)
	}

	outputs := make([]generatedButtonSource, 0, len(buttonStyleSources))
	for _, style := range buttonStyleSources {
		recipePath := filepath.Join(root, style.recipePath)
		recipeSource, err := os.ReadFile(recipePath)
		if err != nil {
			return nil, fmt.Errorf("read %s Button recipe %s: %w", style.name, recipePath, err)
		}
		recipes, err := ParseRecipes(recipePath, recipeSource)
		if err != nil {
			return nil, fmt.Errorf("parse %s Button recipe: %w", style.name, err)
		}
		resolved, _, err := Resolve(canonicalPath, canonical, recipes)
		if err != nil {
			return nil, fmt.Errorf("resolve %s Button: %w", style.name, err)
		}
		outputs = append(outputs, generatedButtonSource{
			relativePath: style.outputPath,
			content:      resolved,
		})
	}
	return outputs, nil
}

func checkButtonSources(root string, outputs []generatedButtonSource) error {
	var drift []string
	for _, output := range outputs {
		path := filepath.Join(root, output.relativePath)
		got, err := os.ReadFile(path)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("read generated Button %s: %w", filepath.ToSlash(output.relativePath), err)
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
			"generated Button sources are out of date: %s; run go run ./cmd/stylegen",
			strings.Join(drift, ", "),
		)
	}
	return nil
}

func writeGeneratedButton(path string, content []byte) error {
	if got, err := os.ReadFile(path); err == nil && bytes.Equal(got, content) {
		return nil
	} else if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".button.gsx-*")
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
