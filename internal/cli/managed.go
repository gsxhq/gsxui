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

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

type artifact struct {
	RelativePath string
	Content      []byte
	Managed      bool
}

var portableCaseFolder = cases.Fold()

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
	return safeProjectPath(root, relative, false)
}

func artifactDirectoryPath(root, relative string) (string, error) {
	return safeProjectPath(root, relative, true)
}

func safeProjectPath(root, relative string, directory bool) (string, error) {
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

	target := filepath.Join(rootResolved, filepath.FromSlash(relative))
	current := rootResolved
	segments := strings.Split(relative, "/")
	for index, segment := range segments {
		current = filepath.Join(current, segment)
		info, err := os.Lstat(current)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("inspect artifact path %q: %w", relative, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", fmt.Errorf("artifact path %q contains unsafe symlink %s", relative, current)
		}
		if index < len(segments)-1 && !info.IsDir() {
			return "", fmt.Errorf("artifact path %q has non-directory ancestor %s", relative, current)
		}
		if index == len(segments)-1 {
			switch {
			case directory && !info.IsDir():
				return "", fmt.Errorf("artifact directory %q is not a directory", relative)
			case !directory && !info.Mode().IsRegular():
				return "", fmt.Errorf("artifact destination %q is not a regular file", relative)
			}
		}
	}
	return target, nil
}

func validateArtifactPlan(root string, cfg Config, artifacts []artifact, overwrite bool) error {
	type destination struct {
		relative string
		path     string
		portable []string
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
		portable := portableArtifactIdentity(planned.RelativePath)
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
			switch {
			case portableArtifactIdentitiesEqual(existing.portable, portable):
				return fmt.Errorf(
					"artifact paths %q and %q are portable aliases",
					existing.relative,
					planned.RelativePath,
				)
			case portableArtifactIdentityIsAncestor(existing.portable, portable):
				return fmt.Errorf(
					"portable planned file %q is an ancestor of destination %q",
					existing.relative,
					planned.RelativePath,
				)
			case portableArtifactIdentityIsAncestor(portable, existing.portable):
				return fmt.Errorf(
					"portable planned file %q is an ancestor of destination %q",
					planned.RelativePath,
					existing.relative,
				)
			}
		}
		destinations = append(destinations, destination{
			relative: planned.RelativePath,
			path:     target,
			portable: portable,
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

func portableArtifactIdentity(relative string) []string {
	segments := strings.Split(relative, "/")
	for index, segment := range segments {
		canonical := norm.NFC.String(segment)
		segments[index] = norm.NFC.String(portableCaseFolder.String(canonical))
	}
	return segments
}

func portableArtifactIdentitiesEqual(left, right []string) bool {
	return len(left) == len(right) && portableArtifactIdentityIsPrefix(left, right)
}

func portableArtifactIdentityIsAncestor(ancestor, descendant []string) bool {
	return len(ancestor) < len(descendant) && portableArtifactIdentityIsPrefix(ancestor, descendant)
}

func portableArtifactIdentityIsPrefix(prefix, full []string) bool {
	if len(prefix) > len(full) {
		return false
	}
	for index, segment := range prefix {
		if segment != full[index] {
			return false
		}
	}
	return true
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
		if err := writeFileAtomic(target, planned.Content, 0o644); err != nil {
			return fmt.Errorf("write artifact %s: %w", target, err)
		}
	}
	return nil
}

func writeFileAtomic(target string, content []byte, defaultMode fs.FileMode) error {
	mode := defaultMode
	info, err := os.Lstat(target)
	switch {
	case err == nil:
		if !info.Mode().IsRegular() {
			return fmt.Errorf("destination is not a regular file")
		}
		mode = info.Mode().Perm()
	case os.IsNotExist(err):
	default:
		return fmt.Errorf("inspect destination: %w", err)
	}

	temp, err := os.CreateTemp(filepath.Dir(target), "."+filepath.Base(target)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temporary file: %w", err)
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath)
	defer temp.Close()

	if err := temp.Chmod(mode); err != nil {
		return fmt.Errorf("set temporary file mode: %w", err)
	}
	written, err := temp.Write(content)
	if err != nil {
		return fmt.Errorf("write temporary file: %w", err)
	}
	if written != len(content) {
		return fmt.Errorf("write temporary file: wrote %d of %d bytes", written, len(content))
	}
	if err := temp.Sync(); err != nil {
		return fmt.Errorf("sync temporary file: %w", err)
	}
	if err := temp.Close(); err != nil {
		return fmt.Errorf("close temporary file: %w", err)
	}
	if err := os.Rename(tempPath, target); err != nil {
		return fmt.Errorf("replace destination: %w", err)
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
