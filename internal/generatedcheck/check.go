package generatedcheck

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type inventory map[string][sha256.Size]byte

// Check verifies the repository's exact source-derived generated set and
// content stability around one generator invocation.
func Check(root string, generate func() error) error {
	expectedBefore, err := expectedOutputs(root)
	if err != nil {
		return err
	}
	actualBefore, err := generatedOutputs(root)
	if err != nil {
		return err
	}

	var problems []error
	problems = append(problems, compareSet("before generation", expectedBefore, actualBefore)...)

	if err := generate(); err != nil {
		problems = append(problems, fmt.Errorf("generator failed: %w", err))
	}

	expectedAfter, err := expectedOutputs(root)
	if err != nil {
		problems = append(problems, err)
	} else if !sameSet(expectedBefore, expectedAfter) {
		problems = append(problems, errors.New("GSX source-derived output set changed during generation"))
	}
	actualAfter, err := generatedOutputs(root)
	if err != nil {
		problems = append(problems, err)
	} else {
		problems = append(problems, compareSet("after generation", expectedAfter, actualAfter)...)
		if !sameInventory(actualBefore, actualAfter) {
			problems = append(problems, errors.New("generated output inventory changed during generation"))
		}
	}

	return errors.Join(problems...)
}

func expectedOutputs(root string) (map[string]struct{}, error) {
	expected := make(map[string]struct{})
	err := walkFiles(root, func(path string) error {
		if !strings.HasSuffix(path, ".gsx") {
			return nil
		}
		rel, err := filepath.Rel(root, strings.TrimSuffix(path, ".gsx")+".x.go")
		if err != nil {
			return err
		}
		expected[filepath.ToSlash(rel)] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("discovering GSX sources: %w", err)
	}
	return expected, nil
}

func generatedOutputs(root string) (inventory, error) {
	outputs := make(inventory)
	err := walkFiles(root, func(path string) error {
		if !strings.HasSuffix(path, ".x.go") {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		outputs[filepath.ToSlash(rel)] = sha256.Sum256(content)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("inventorying generated outputs: %w", err)
	}
	return outputs, nil
}

func walkFiles(root string, visit func(string) error) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && entry.Name() == ".git" && path != root {
			return filepath.SkipDir
		}
		if entry.Type().IsRegular() {
			return visit(path)
		}
		return nil
	})
}

func compareSet(when string, expected map[string]struct{}, actual inventory) []error {
	var missing, unexpected []string
	for path := range expected {
		if _, ok := actual[path]; !ok {
			missing = append(missing, path)
		}
	}
	for path := range actual {
		if _, ok := expected[path]; !ok {
			unexpected = append(unexpected, path)
		}
	}
	slices.Sort(missing)
	slices.Sort(unexpected)

	problems := make([]error, 0, len(missing)+len(unexpected))
	for _, path := range missing {
		problems = append(problems, fmt.Errorf("%s: missing generated output %q", when, path))
	}
	for _, path := range unexpected {
		problems = append(problems, fmt.Errorf("%s: unexpected generated output %q", when, path))
	}
	return problems
}

func sameSet(a, b map[string]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for path := range a {
		if _, ok := b[path]; !ok {
			return false
		}
	}
	return true
}

func sameInventory(a, b inventory) bool {
	if len(a) != len(b) {
		return false
	}
	for path, digest := range a {
		if b[path] != digest {
			return false
		}
	}
	return true
}
