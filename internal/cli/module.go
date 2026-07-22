package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// modulePath reads the module directive from dir's go.mod.
func modulePath(dir string) (string, error) {
	f, err := os.Open(filepath.Join(dir, "go.mod"))
	if err != nil {
		return "", fmt.Errorf("run gsxui from your module root (go.mod not found)")
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.Trim(strings.TrimSpace(rest), `"`), nil
		}
	}
	return "", fmt.Errorf("no module directive in go.mod")
}
