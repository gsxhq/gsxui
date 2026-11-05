package cli

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// scaffoldAbsent reports whether the project has no gsx npm/Vite scaffold at
// all — neither vite.config.ts nor web/main.js exists. That state selects
// init's non-Vite mode; a partially present or modified scaffold still goes
// through planScaffoldIntegration's refusal so tampering never silently
// downgrades an existing Vite project.
func scaffoldAbsent(dir string) (bool, error) {
	for _, relative := range []string{"vite.config.ts", "web/main.js"} {
		path, err := artifactPath(dir, relative)
		if err != nil {
			return false, err
		}
		if _, err := os.Stat(path); err == nil {
			return false, nil
		} else if !os.IsNotExist(err) {
			return false, err
		}
	}
	return true, nil
}

const manualScaffoldIntegrationInstructions = `Automatic integration currently supports an unmodified gsx init npm/Vite scaffold with the default gsxui JS and CSS paths.

Integrate this project manually:
  npm install --save-dev tailwindcss@^4.3.3 @tailwindcss/vite@^4.3.3 tw-animate-css@^1.4.0
  vite.config.ts: import tailwindcss from "@tailwindcss/vite" and add tailwindcss() to plugins
  web/main.js: import "./gsxui/index.js" and import "./gsxui/index.css"

Or, to run without Vite entirely, delete vite.config.ts and web/main.js and
re-run gsxui init: it will vendor everything self-contained (no npm), and you
serve the gsxui JS directory statically, load index.js with one
<script type="module"> tag, and build index.css with any Tailwind v4 tool
(e.g. npx @tailwindcss/cli -i web/gsxui/index.css -o dist.css).`

// scaffoldDocuments holds one directory per gsx init scaffold revision. A
// revision directory carries only the documents that changed in it.
//
//go:embed scaffold
var scaffoldDocuments embed.FS

// scaffoldDocument is one consumer file gsxui integrates into, with every gsx
// scaffold revision of it that init recognizes.
type scaffoldDocument struct {
	target    string
	revisions []scaffoldRevision
}

// scaffoldRevision is the exact before/after pair for one gsx scaffold
// revision of a document.
type scaffoldRevision struct {
	pristine         []byte
	integrated       []byte
	pristineDigest   [sha256.Size]byte
	integratedDigest [sha256.Size]byte
}

type scaffoldPackage struct {
	Scripts         map[string]string `json:"scripts"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func planScaffoldIntegration(dir string, cfg Config) ([]artifact, error) {
	var issues []string
	if cfg.JS != "web/gsxui" || cfg.CSS != "web/gsxui/index.css" {
		issues = append(issues, "automatic integration requires the default gsxui JS and CSS paths")
	}
	if err := validateScaffoldPackage(dir); err != nil {
		issues = append(issues, err.Error())
	}

	artifacts := make([]artifact, 0, 2)
	for _, document := range supportedScaffoldDocuments() {
		planned, err := planScaffoldDocument(dir, document)
		if err != nil {
			issues = append(issues, err.Error())
			continue
		}
		if planned != nil {
			artifacts = append(artifacts, *planned)
		}
	}
	if len(issues) != 0 {
		return nil, fmt.Errorf(
			"automatic gsx scaffold integration refused:\n  - %s\n\n%s",
			strings.Join(issues, "\n  - "),
			manualScaffoldIntegrationInstructions,
		)
	}
	return artifacts, nil
}

// supportedScaffoldDocuments lists, per document, every gsx scaffold revision
// init integrates: gsx-v1 is gsx init through v0.4.4; gsx-v2 (v0.4.5+) changed
// only vite.config.ts.
func supportedScaffoldDocuments() []scaffoldDocument {
	return []scaffoldDocument{
		newScaffoldDocument("vite.config.ts",
			newScaffoldRevision("gsx-v1/vite.config.ts", "gsx-v1/vite.gsxui.config.ts"),
			newScaffoldRevision("gsx-v2/vite.config.ts", "gsx-v2/vite.gsxui.config.ts"),
		),
		newScaffoldDocument("web/main.js",
			newScaffoldRevision("gsx-v1/main.js", "gsx-v1/main.gsxui.js"),
		),
	}
}

func newScaffoldDocument(target string, revisions ...scaffoldRevision) scaffoldDocument {
	return scaffoldDocument{target: target, revisions: revisions}
}

func newScaffoldRevision(pristineName, integratedName string) scaffoldRevision {
	pristine := mustScaffoldDocument(pristineName)
	integrated := mustScaffoldDocument(integratedName)
	return scaffoldRevision{
		pristine:         pristine,
		integrated:       integrated,
		pristineDigest:   sha256.Sum256(pristine),
		integratedDigest: sha256.Sum256(integrated),
	}
}

func mustScaffoldDocument(name string) []byte {
	content, err := fs.ReadFile(scaffoldDocuments, "scaffold/"+name)
	if err != nil {
		panic(err)
	}
	return content
}

func planScaffoldDocument(dir string, document scaffoldDocument) (*artifact, error) {
	target, err := artifactPath(dir, document.target)
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(target)
	if err != nil {
		return nil, fmt.Errorf("%s is not a recognized gsx scaffold document: %w", document.target, err)
	}
	digest := sha256.Sum256(content)
	for _, revision := range document.revisions {
		switch digest {
		case revision.integratedDigest:
			if !bytes.Equal(content, revision.integrated) {
				return nil, fmt.Errorf("%s did not match its integrated fingerprint", document.target)
			}
			return nil, nil
		case revision.pristineDigest:
			if !bytes.Equal(content, revision.pristine) {
				return nil, fmt.Errorf("%s did not match its pristine fingerprint", document.target)
			}
			return &artifact{
				RelativePath: document.target,
				Content:      revision.integrated,
			}, nil
		}
	}
	return nil, fmt.Errorf("%s is not an exact supported gsx scaffold document", document.target)
}

func validateScaffoldPackage(dir string) error {
	packagePath, err := artifactPath(dir, "package.json")
	if err != nil {
		return err
	}
	content, err := os.ReadFile(packagePath)
	if err != nil {
		return fmt.Errorf("package.json is required: %w", err)
	}
	var manifest scaffoldPackage
	if err := json.Unmarshal(content, &manifest); err != nil {
		return fmt.Errorf("package.json is not valid JSON: %w", err)
	}
	if manifest.Scripts["dev"] != "go tool gsx dev" {
		return errors.New(`package.json scripts.dev must be "go tool gsx dev"`)
	}
	if manifest.Scripts["build"] != "vite build" {
		return errors.New(`package.json scripts.build must be "vite build"`)
	}
	if manifest.DevDependencies["@gsxhq/vite-plugin-gsx"] == "" {
		return errors.New("package.json devDependencies must include @gsxhq/vite-plugin-gsx")
	}

	lockPath, err := artifactPath(dir, "package-lock.json")
	if err != nil {
		return err
	}
	if info, err := os.Stat(lockPath); err != nil {
		return fmt.Errorf("package-lock.json from gsx init --yes is required: %w", err)
	} else if !info.Mode().IsRegular() {
		return errors.New("package-lock.json must be a regular file")
	}
	for _, relative := range []string{"pnpm-lock.yaml", "yarn.lock", "bun.lock", "bun.lockb"} {
		path, err := artifactPath(dir, relative)
		if err != nil {
			return err
		}
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s is not supported; automatic integration requires npm", relative)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("inspect %s: %w", relative, err)
		}
	}
	return nil
}
