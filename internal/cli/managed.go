package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type artifact struct {
	RelativePath string
	Content      []byte
	Managed      bool
}

func contentHash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func validateManaged(managed map[string]string) error {
	for relative, hash := range managed {
		if err := validateManagedPath(relative); err != nil {
			return fmt.Errorf("managed path %q: %w", relative, err)
		}
		if relative == "gsxui.json" {
			return fmt.Errorf("managed path %q: configuration cannot hash itself", relative)
		}
		if len(hash) != sha256.Size*2 {
			return fmt.Errorf("managed path %q: hash must be %d lowercase hexadecimal characters", relative, sha256.Size*2)
		}
		decoded, err := hex.DecodeString(hash)
		if err != nil || hex.EncodeToString(decoded) != hash {
			return fmt.Errorf("managed path %q: hash must be %d lowercase hexadecimal characters", relative, sha256.Size*2)
		}
	}
	return nil
}

func validateManagedPath(relative string) error {
	if relative == "" {
		return fmt.Errorf("must not be empty")
	}
	if strings.Contains(relative, `\`) {
		return fmt.Errorf("must use slash-normalized separators")
	}
	if strings.IndexByte(relative, 0) >= 0 {
		return fmt.Errorf("must not contain NUL")
	}
	if filepath.IsAbs(relative) || filepath.VolumeName(relative) != "" {
		return fmt.Errorf("must be relative to the module root")
	}
	cleaned := path.Clean(relative)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		return fmt.Errorf("escapes module root")
	}
	if cleaned != relative || !fs.ValidPath(relative) {
		return fmt.Errorf("must be a clean slash-normalized path")
	}
	return nil
}

func artifactPath(root, relative string) (string, error) {
	if err := validateManagedPath(relative); err != nil {
		return "", fmt.Errorf("artifact path %q: %w", relative, err)
	}

	rootAbsolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve module root: %w", err)
	}
	rootResolved, err := filepath.EvalSymlinks(rootAbsolute)
	if err != nil {
		return "", fmt.Errorf("resolve module root: %w", err)
	}

	target := filepath.Join(rootAbsolute, filepath.FromSlash(relative))
	existing := target
	var missing []string
	for {
		_, err := os.Lstat(existing)
		if err == nil {
			break
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("inspect artifact path %q: %w", relative, err)
		}
		parent := filepath.Dir(existing)
		if parent == existing {
			return "", fmt.Errorf("inspect artifact path %q: no existing ancestor", relative)
		}
		missing = append(missing, filepath.Base(existing))
		existing = parent
	}

	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		return "", fmt.Errorf("resolve artifact path %q: %w", relative, err)
	}
	for index := len(missing) - 1; index >= 0; index-- {
		resolved = filepath.Join(resolved, missing[index])
	}
	within, err := filepath.Rel(rootResolved, resolved)
	if err != nil {
		return "", fmt.Errorf("compare artifact path %q with module root: %w", relative, err)
	}
	if within == ".." || strings.HasPrefix(within, ".."+string(filepath.Separator)) || filepath.IsAbs(within) {
		return "", fmt.Errorf("artifact path %q escapes module root through a symlink", relative)
	}
	return resolved, nil
}

func validateArtifactPlan(root string, cfg Config, artifacts []artifact, overwrite bool) error {
	type destination struct {
		relative string
		path     string
		info     fs.FileInfo
	}
	destinations := make([]destination, 0, len(artifacts))
	for _, planned := range artifacts {
		target, err := artifactPath(root, planned.RelativePath)
		if err != nil {
			return err
		}
		info, err := os.Stat(target)
		if err != nil && !errorsIsNotExist(err) {
			return fmt.Errorf("inspect artifact %s: %w", target, err)
		}
		for _, existing := range destinations {
			samePath := existing.path == target
			sameFile := info != nil && existing.info != nil && os.SameFile(info, existing.info)
			if samePath || sameFile {
				return fmt.Errorf(
					"artifact paths %q and %q resolve to the same destination",
					existing.relative,
					planned.RelativePath,
				)
			}
		}
		destinations = append(destinations, destination{
			relative: planned.RelativePath,
			path:     target,
			info:     info,
		})

		existing, err := os.ReadFile(target)
		if errorsIsNotExist(err) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read artifact %s: %w", target, err)
		}
		if bytes.Equal(existing, planned.Content) {
			continue
		}
		if overwrite {
			continue
		}
		if !planned.Managed {
			continue
		}
		if planned.RelativePath == "gsxui.preset.json" {
			return fmt.Errorf("%s differs from the selected preset — pass --overwrite to replace it", target)
		}
		recorded, managed := cfg.Managed[planned.RelativePath]
		if managed && contentHash(existing) == recorded {
			continue
		}
		return fmt.Errorf("%s differs from the gsxui version — pass --overwrite to replace it", target)
	}
	return nil
}

func writeArtifactPlan(root string, artifacts []artifact) error {
	for _, planned := range artifacts {
		target, err := artifactPath(root, planned.RelativePath)
		if err != nil {
			return err
		}
		existing, err := os.ReadFile(target)
		if err == nil && bytes.Equal(existing, planned.Content) {
			continue
		}
		if err != nil && !errorsIsNotExist(err) {
			return fmt.Errorf("read artifact %s: %w", target, err)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return fmt.Errorf("create artifact directory for %s: %w", target, err)
		}
		if err := os.WriteFile(target, planned.Content, 0o644); err != nil {
			return fmt.Errorf("write artifact %s: %w", target, err)
		}
	}
	return nil
}

func managedAfter(cfg Config, artifacts []artifact) map[string]string {
	managed := make(map[string]string, len(cfg.Managed)+len(artifacts))
	maps.Copy(managed, cfg.Managed)
	for _, planned := range artifacts {
		if planned.Managed {
			managed[planned.RelativePath] = contentHash(planned.Content)
		}
	}
	return managed
}

func artifactPlanWithConfig(cfg Config, artifacts []artifact) (Config, []artifact, error) {
	next := cfg
	next.Managed = managedAfter(cfg, artifacts)
	configContent, err := next.canonicalJSON()
	if err != nil {
		return Config{}, nil, fmt.Errorf("plan gsxui.json: %w", err)
	}
	complete := make([]artifact, len(artifacts), len(artifacts)+1)
	copy(complete, artifacts)
	complete = append(complete, artifact{
		RelativePath: "gsxui.json",
		Content:      configContent,
	})
	return next, complete, nil
}

func errorsIsNotExist(err error) bool {
	return err != nil && (os.IsNotExist(err) || err == fs.ErrNotExist)
}
