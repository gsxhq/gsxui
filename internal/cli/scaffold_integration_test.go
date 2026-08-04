package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testGSXMainPristine = `// gsx dev panel: Cmd-D / Ctrl-D
import "virtual:gsx-devpanel";
import "./style.css";
import { setupCounter } from "./counter.js";

setupCounter(document.querySelector("#counter"));
`

const testGSXMainIntegrated = `// gsx dev panel: Cmd-D / Ctrl-D
import "virtual:gsx-devpanel";
import "./gsxui/index.js";
import "./gsxui/index.css";
import "./style.css";
import { setupCounter } from "./counter.js";

setupCounter(document.querySelector("#counter"));
`

const testGSXVitePristine = `import { defineConfig, loadEnv, createLogger } from "vite";
import { gsx, devFallback } from "@gsxhq/vite-plugin-gsx";

export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const goPort = env.GO_PORT || "7777";
  const vitePort = parseInt(env.VITE_PORT || "5173", 10);
  // Single-computed upstream: ` + "`gsx dev`" + ` injects GSX_DEV_UPSTREAM into the real
  // process env, injected by gsx dev at spawn; everything else
  // (standalone ` + "`vite`" + `, ` + "`vite build`" + `) falls back to GO_PORT. Passed explicitly
  // to devFallback AND the proxy below so evaluation order can't matter — a
  // no-args devFallback call throws when neither is set, which would break
  // ` + "`vite build`" + ` (devFallback is inert there, but the zero-arg constructor
  // still throws) and standalone ` + "`vite serve`" + `.
  const upstream = process.env.GSX_DEV_UPSTREAM ?? ` + "`http://localhost:${goPort}`" + `;
  // Serve a self-recovering interstitial while the Go server is down/restarting
  // (instead of a raw proxy error).
  const fallback = devFallback({ target: upstream });

  // While the Go server is down/restarting, the dev-fallback interstitial already
  // shows it — so drop Vite's redundant "http proxy error … ECONNREFUSED" spam.
  const logger = createLogger();
  const baseError = logger.error;
  logger.error = (msg, opts) => {
    if (typeof msg === "string" && msg.includes("http proxy error")) return;
    baseError(msg, opts);
  };

  return {
    clearScreen: false,
    // Dev serves all Vite assets under /__vite/ (matches gsxhq/vite DevBase);
    // prod uses the default base ("/") since hashed assets are served from
    // /static via gsxhq/vite.
    base: command === "serve" ? "/__vite/" : "/",
    publicDir: false,
    customLogger: logger,
    plugins: [
      gsx(),
      fallback.plugin,
      {
        // Log the app front door ("/"); Vite's own line shows the empty /__vite/ base.
        name: "gsx-dev-url",
        configureServer(server) {
          server.printUrls = () => {
            server.config.logger.info(
              ` + "`\\n  \\x1b[32m➜\\x1b[0m  Open \\x1b[36mhttp://localhost:${vitePort}/\\x1b[0m to view your app\\n`" + `,
            );
          };
        },
      },
    ],
    server: {
      port: vitePort,
      // ` + "`gsx dev`" + ` chooses an available VITE_PORT and matching VITE_DEV_URL before
      // starting Vite. Keep Vite strict so proxy, overlay, and printed URL agree.
      strictPort: true,
      proxy: {
        // Everything except /__vite/ (Vite's own dev-asset namespace) and
        // /__dev/ (fallback plugin status) goes to the Go server.
        // Single exclusion — no source-dir denylist, no app-route collisions.
        // No ` + "`ws: true`" + ` — the Go server has no WebSocket; proxying ws would
        // capture Vite's HMR socket.
        "^(?!/__vite/|/__dev).*": {
          target: upstream,
          changeOrigin: true,
          configure: fallback.configureProxy,
        },
      },
    },
    build: {
      manifest: true,
      outDir: "dist",
      rollupOptions: { input: "web/main.js" },
    },
  };
});
`

const testGSXViteIntegrated = `import { defineConfig, loadEnv, createLogger } from "vite";
import tailwindcss from "@tailwindcss/vite";
import { gsx, devFallback } from "@gsxhq/vite-plugin-gsx";

export default defineConfig(({ command, mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const goPort = env.GO_PORT || "7777";
  const vitePort = parseInt(env.VITE_PORT || "5173", 10);
  // Single-computed upstream: ` + "`gsx dev`" + ` injects GSX_DEV_UPSTREAM into the real
  // process env, injected by gsx dev at spawn; everything else
  // (standalone ` + "`vite`" + `, ` + "`vite build`" + `) falls back to GO_PORT. Passed explicitly
  // to devFallback AND the proxy below so evaluation order can't matter — a
  // no-args devFallback call throws when neither is set, which would break
  // ` + "`vite build`" + ` (devFallback is inert there, but the zero-arg constructor
  // still throws) and standalone ` + "`vite serve`" + `.
  const upstream = process.env.GSX_DEV_UPSTREAM ?? ` + "`http://localhost:${goPort}`" + `;
  // Serve a self-recovering interstitial while the Go server is down/restarting
  // (instead of a raw proxy error).
  const fallback = devFallback({ target: upstream });

  // While the Go server is down/restarting, the dev-fallback interstitial already
  // shows it — so drop Vite's redundant "http proxy error … ECONNREFUSED" spam.
  const logger = createLogger();
  const baseError = logger.error;
  logger.error = (msg, opts) => {
    if (typeof msg === "string" && msg.includes("http proxy error")) return;
    baseError(msg, opts);
  };

  return {
    clearScreen: false,
    // Dev serves all Vite assets under /__vite/ (matches gsxhq/vite DevBase);
    // prod uses the default base ("/") since hashed assets are served from
    // /static via gsxhq/vite.
    base: command === "serve" ? "/__vite/" : "/",
    publicDir: false,
    customLogger: logger,
    plugins: [
      tailwindcss(),
      gsx(),
      fallback.plugin,
      {
        // Log the app front door ("/"); Vite's own line shows the empty /__vite/ base.
        name: "gsx-dev-url",
        configureServer(server) {
          server.printUrls = () => {
            server.config.logger.info(
              ` + "`\\n  \\x1b[32m➜\\x1b[0m  Open \\x1b[36mhttp://localhost:${vitePort}/\\x1b[0m to view your app\\n`" + `,
            );
          };
        },
      },
    ],
    server: {
      port: vitePort,
      // ` + "`gsx dev`" + ` chooses an available VITE_PORT and matching VITE_DEV_URL before
      // starting Vite. Keep Vite strict so proxy, overlay, and printed URL agree.
      strictPort: true,
      proxy: {
        // Everything except /__vite/ (Vite's own dev-asset namespace) and
        // /__dev/ (fallback plugin status) goes to the Go server.
        // Single exclusion — no source-dir denylist, no app-route collisions.
        // No ` + "`ws: true`" + ` — the Go server has no WebSocket; proxying ws would
        // capture Vite's HMR socket.
        "^(?!/__vite/|/__dev).*": {
          target: upstream,
          changeOrigin: true,
          configure: fallback.configureProxy,
        },
      },
    },
    build: {
      manifest: true,
      outDir: "dist",
      rollupOptions: { input: "web/main.js" },
    },
  };
});
`

func TestPlanScaffoldIntegrationRecognizesExactDocumentsIndependently(t *testing.T) {
	tests := []struct {
		name           string
		vite           string
		main           string
		wantArtifacts  int
		wantViteChange bool
		wantMainChange bool
	}{
		{
			name:           "pristine",
			vite:           testGSXVitePristine,
			main:           testGSXMainPristine,
			wantArtifacts:  2,
			wantViteChange: true,
			wantMainChange: true,
		},
		{
			name:          "integrated",
			vite:          testGSXViteIntegrated,
			main:          testGSXMainIntegrated,
			wantArtifacts: 0,
		},
		{
			name:           "vite only pristine",
			vite:           testGSXVitePristine,
			main:           testGSXMainIntegrated,
			wantArtifacts:  1,
			wantViteChange: true,
		},
		{
			name:           "main only pristine",
			vite:           testGSXViteIntegrated,
			main:           testGSXMainPristine,
			wantArtifacts:  1,
			wantMainChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := writeScaffoldFixture(t, tt.vite, tt.main)
			artifacts, err := planScaffoldIntegration(dir, DefaultConfig())
			if err != nil {
				t.Fatal(err)
			}
			if len(artifacts) != tt.wantArtifacts {
				t.Fatalf("artifacts = %#v, want %d", artifacts, tt.wantArtifacts)
			}
			byPath := make(map[string]artifact, len(artifacts))
			for _, planned := range artifacts {
				byPath[planned.RelativePath] = planned
				if planned.Managed {
					t.Fatalf("integration artifact %s is managed", planned.RelativePath)
				}
			}
			if tt.wantViteChange && !bytes.Equal(byPath["vite.config.ts"].Content, []byte(testGSXViteIntegrated)) {
				t.Fatal("vite replacement is not the exact integrated document")
			}
			if tt.wantMainChange && !bytes.Equal(byPath["web/main.js"].Content, []byte(testGSXMainIntegrated)) {
				t.Fatal("main replacement is not the exact integrated document")
			}
		})
	}
}

func TestSupportedScaffoldDocumentFingerprintsArePinned(t *testing.T) {
	want := map[string]struct {
		pristine   string
		integrated string
	}{
		"vite.config.ts": {
			pristine:   "e7bccf275de4ab9210c168503022de4868f3186977d245a9c1b7e06bf65241d0",
			integrated: "ea4df2b92c4f9327b6365353d36b6ab7cca439f9bef207a61cf9ba35608529ca",
		},
		"web/main.js": {
			pristine:   "441ab2eca4fef1f817cab40d7db56a85ab6efc66ebfc4793b968455abb507b10",
			integrated: "cc572e73a839016d5b52943d1f0bb79048f28ab7a38e55a52dabcad5924d4e1f",
		},
	}
	for _, document := range supportedScaffoldDocuments() {
		fingerprints, ok := want[document.target]
		if !ok {
			t.Fatalf("unpinned scaffold document %q", document.target)
		}
		if got := contentHash(document.pristine); got != fingerprints.pristine {
			t.Errorf("%s pristine fingerprint = %s, want %s", document.target, got, fingerprints.pristine)
		}
		if got := contentHash(document.integrated); got != fingerprints.integrated {
			t.Errorf("%s integrated fingerprint = %s, want %s", document.target, got, fingerprints.integrated)
		}
		delete(want, document.target)
	}
	if len(want) != 0 {
		t.Fatalf("missing supported scaffold documents: %v", want)
	}
}

func TestPlanScaffoldIntegrationRefusesUnknownConsumerFiles(t *testing.T) {
	for _, tt := range []struct {
		name string
		vite string
		main string
		want string
	}{
		{name: "custom vite", vite: testGSXVitePristine + "\n// custom\n", main: testGSXMainPristine, want: "vite.config.ts"},
		{name: "custom main", vite: testGSXVitePristine, main: testGSXMainPristine + "\n// custom\n", want: "web/main.js"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := writeScaffoldFixture(t, tt.vite, tt.main)
			_, err := planScaffoldIntegration(dir, DefaultConfig())
			if err == nil || !strings.Contains(err.Error(), tt.want) ||
				!strings.Contains(err.Error(), "npm install --save-dev") ||
				!strings.Contains(err.Error(), `import tailwindcss from "@tailwindcss/vite"`) ||
				!strings.Contains(err.Error(), `import "./gsxui/index.css"`) {
				t.Fatalf("planScaffoldIntegration error = %v", err)
			}
		})
	}
}

func TestPlanScaffoldIntegrationValidatesPackageContract(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*testing.T, string)
		want   string
	}{
		{
			name: "missing package json",
			mutate: func(t *testing.T, dir string) {
				t.Helper()
				if err := os.Remove(filepath.Join(dir, "package.json")); err != nil {
					t.Fatal(err)
				}
			},
			want: "package.json",
		},
		{
			name: "wrong dev script",
			mutate: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "package.json", strings.ReplaceAll(testGSXPackageJSON, "go tool gsx dev", "vite"))
			},
			want: "scripts.dev",
		},
		{
			name: "wrong build script",
			mutate: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "package.json", strings.ReplaceAll(testGSXPackageJSON, "vite build", "other build"))
			},
			want: "scripts.build",
		},
		{
			name: "missing vite plugin",
			mutate: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "package.json", strings.ReplaceAll(testGSXPackageJSON, `"@gsxhq/vite-plugin-gsx": "^0.10.0",`, `"other": "1",`))
			},
			want: "@gsxhq/vite-plugin-gsx",
		},
		{
			name: "missing package lock",
			mutate: func(t *testing.T, dir string) {
				t.Helper()
				if err := os.Remove(filepath.Join(dir, "package-lock.json")); err != nil {
					t.Fatal(err)
				}
			},
			want: "package-lock.json",
		},
		{
			name: "alternate lockfile",
			mutate: func(t *testing.T, dir string) {
				t.Helper()
				writeFile(t, dir, "pnpm-lock.yaml", "lockfileVersion: 9\n")
			},
			want: "pnpm-lock.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := writeScaffoldFixture(t, testGSXVitePristine, testGSXMainPristine)
			tt.mutate(t, dir)
			if _, err := planScaffoldIntegration(dir, DefaultConfig()); err == nil ||
				!strings.Contains(err.Error(), tt.want) {
				t.Fatalf("planScaffoldIntegration error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestPlanScaffoldIntegrationRejectsCustomEntryPaths(t *testing.T) {
	dir := writeScaffoldFixture(t, testGSXVitePristine, testGSXMainPristine)
	cfg := DefaultConfig()
	cfg.JS = "web/behavior"
	cfg.CSS = "web/styles/app.css"
	if _, err := planScaffoldIntegration(dir, cfg); err == nil ||
		!strings.Contains(err.Error(), "default gsxui JS and CSS paths") {
		t.Fatalf("planScaffoldIntegration error = %v", err)
	}
}

func TestManualScaffoldInstructionsMatchDocumentation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "..", "site", "snippets", "manual-integration.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(content)) != manualScaffoldIntegrationInstructions {
		t.Fatalf(
			"manual integration documentation drifted from CLI refusal:\n%s",
			content,
		)
	}
}

func TestNonViteInitSummaryMatchesDocumentation(t *testing.T) {
	content, err := os.ReadFile(filepath.Join("..", "..", "site", "snippets", "nonvite-init.output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	cfg := DefaultConfig()
	got := fmt.Sprintf(nonViteSummaryFormat, cfg.CSS, cfg.JS)
	if strings.TrimSpace(string(content)) != strings.TrimSpace(got) {
		t.Fatalf(
			"non-Vite init summary drifted from documentation:\n%s",
			content,
		)
	}
}

func TestInitRefusesUnknownScaffoldBeforeCommandsOrWrites(t *testing.T) {
	for _, overwrite := range []bool{false, true} {
		name := "plain"
		args := []string{"init"}
		if overwrite {
			name = "overwrite"
			args = append(args, "--overwrite")
		}
		t.Run(name, func(t *testing.T) {
			dir, commands := initTestModule(t)
			writeFile(t, dir, "web/main.js", testGSXMainPristine+"\n// consumer edit\n")
			before := snapshotFiles(t, dir)

			err := Run(args)
			if err == nil || !strings.Contains(err.Error(), "web/main.js") ||
				!strings.Contains(err.Error(), manualScaffoldIntegrationInstructions) {
				t.Fatalf("Run(%v) error = %v", args, err)
			}
			if len(*commands) != 0 {
				t.Fatalf("preflight refusal ran commands: %v", *commands)
			}
			if after := snapshotFiles(t, dir); !equalFileSnapshots(before, after) {
				t.Fatal("preflight refusal changed the project")
			}
		})
	}
}

func TestInitRestoresPackageMetadataWhenNPMFails(t *testing.T) {
	dir, _ := initTestModule(t)
	before := snapshotFiles(t, dir)
	previous := runCommand
	runCommand = func(dir, name string, args ...string) error {
		if name != "npm" {
			return nil
		}
		writeFile(t, dir, "package.json", `{"changed":true}`+"\n")
		writeFile(t, dir, "package-lock.json", `{"changed":true}`+"\n")
		return errors.New("injected npm failure")
	}
	t.Cleanup(func() { runCommand = previous })

	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "injected npm failure") {
		t.Fatalf("Run(init) error = %v", err)
	}
	if after := snapshotFiles(t, dir); !equalFileSnapshots(before, after) {
		t.Fatal("npm failure did not restore exact authored project state")
	}
}

func TestInitRestoresPackageMetadataWhenGoDependencyFails(t *testing.T) {
	dir, _ := initTestModule(t)
	before := snapshotFiles(t, dir)
	previous := runCommand
	var npmCalls int
	runCommand = func(dir, name string, args ...string) error {
		if name == "npm" {
			npmCalls++
			writeFile(t, dir, "package.json", `{"npmChanged":true}`+"\n")
			writeFile(t, dir, "package-lock.json", `{"npmChanged":true}`+"\n")
			return nil
		}
		return errors.New("injected Go dependency failure")
	}
	t.Cleanup(func() { runCommand = previous })

	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "injected Go dependency failure") {
		t.Fatalf("Run(init) error = %v", err)
	}
	if after := snapshotFiles(t, dir); !equalFileSnapshots(before, after) {
		t.Fatal("Go dependency failure did not restore npm metadata and integration state")
	}
	if npmCalls != 1 {
		t.Fatalf("npm calls = %d, want 1 before the Go dependency failure", npmCalls)
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestInitRestoresPackageMetadataAndArtifactsWhenGenerationFails(t *testing.T) {
	dir, _ := initTestModule(t)
	before := snapshotFiles(t, dir)
	previous := runCommand
	var generateCalls int
	runCommand = func(dir, name string, args ...string) error {
		if name == "npm" {
			writeFile(t, dir, "package.json", `{"npmChanged":true}`+"\n")
			writeFile(t, dir, "package-lock.json", `{"npmChanged":true}`+"\n")
			return nil
		}
		if name == "go" && len(args) >= 3 && args[0] == "tool" && args[1] == "gsx" &&
			args[2] == "generate" {
			generateCalls++
			if generateCalls == 1 {
				return errors.New("injected generation failure after npm")
			}
		}
		return nil
	}
	t.Cleanup(func() { runCommand = previous })

	err := Run([]string{"init"})
	if err == nil || !strings.Contains(err.Error(), "injected generation failure after npm") {
		t.Fatalf("Run(init) error = %v", err)
	}
	if generateCalls != 2 {
		t.Fatalf("generation calls = %d, want failed generation plus restored-source generation", generateCalls)
	}
	if after := snapshotFiles(t, dir); !equalFileSnapshots(before, after) {
		t.Fatal("generation failure did not restore package metadata and integration artifacts")
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestInterruptedNPMIsRecoveredByArtifactJournal(t *testing.T) {
	dir, _ := initTestModule(t)
	before := snapshotFiles(t, dir)
	previous := runCommand
	runCommand = func(dir, name string, args ...string) error {
		if name == "npm" {
			writeFile(t, dir, "package.json", `{"interrupted":true}`+"\n")
			writeFile(t, dir, "package-lock.json", `{"interrupted":true}`+"\n")
			panic("simulated process interruption during npm")
		}
		return nil
	}
	t.Cleanup(func() { runCommand = previous })

	func() {
		defer func() {
			if recovered := recover(); recovered != "simulated process interruption during npm" {
				t.Fatalf("panic = %#v", recovered)
			}
		}()
		_ = Run([]string{"init"})
	}()
	if err := recoverArtifactTransaction(dir); err != nil {
		t.Fatal(err)
	}
	if after := snapshotFiles(t, dir); !equalFileSnapshots(before, after) {
		t.Fatal("journal recovery did not restore interrupted npm and integration state")
	}
	assertNoTransactionArtifacts(t, dir)
}

func TestInitRerunKeepsIntegrationDocumentsByteStable(t *testing.T) {
	dir, commands := initTestModule(t)
	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	vite := readFile(t, dir, "vite.config.ts")
	main := readFile(t, dir, "web/main.js")

	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, dir, "vite.config.ts"); got != vite {
		t.Fatal("rerun changed integrated vite.config.ts")
	}
	if got := readFile(t, dir, "web/main.js"); got != main {
		t.Fatal("rerun changed integrated web/main.js")
	}
	if strings.Count(vite, `import tailwindcss from "@tailwindcss/vite"`) != 1 ||
		strings.Count(vite, "tailwindcss()") != 1 {
		t.Fatalf("vite integration is duplicated:\n%s", vite)
	}
	if strings.Count(main, `import "./gsxui/index.js"`) != 1 ||
		strings.Count(main, `import "./gsxui/index.css"`) != 1 {
		t.Fatalf("entry integration is duplicated:\n%s", main)
	}
	var npmCalls int
	for _, command := range *commands {
		if len(command) != 0 && command[0] == "npm" {
			npmCalls++
		}
	}
	if npmCalls != 2 {
		t.Fatalf("npm calls = %d, want one per init", npmCalls)
	}
}

func TestInitPrintsCompletedIntegrationSummary(t *testing.T) {
	_, _ = initTestModule(t)
	var output bytes.Buffer
	previous := commandStdout
	commandStdout = &output
	t.Cleanup(func() { commandStdout = previous })

	if err := Run([]string{"init"}); err != nil {
		t.Fatal(err)
	}
	const want = `gsxui initialized.
  css:  web/gsxui/index.css
  js:   web/gsxui/index.js
  vite: Tailwind CSS configured
  next: gsxui add button
`
	if output.String() != want {
		t.Fatalf("init output:\n%s\nwant:\n%s", output.String(), want)
	}
}

const testGSXPackageJSON = `{
  "name": "example",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "go tool gsx dev",
    "build": "vite build"
  },
  "devDependencies": {
    "@gsxhq/vite-plugin-gsx": "^0.10.0",
    "vite": "^6.0.0"
  }
}
`

func writeScaffoldFixture(t *testing.T, vite, main string) string {
	t.Helper()
	dir := t.TempDir()
	writeFile(t, dir, "package.json", testGSXPackageJSON)
	writeFile(t, dir, "package-lock.json", "{}\n")
	writeFile(t, dir, "vite.config.ts", vite)
	writeFile(t, dir, "web/main.js", main)
	return dir
}

func writeFile(t *testing.T, dir, relative, content string) {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, dir, relative string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(relative)))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func TestScaffoldAbsent(t *testing.T) {
	dir := t.TempDir()
	absent, err := scaffoldAbsent(dir)
	if err != nil || !absent {
		t.Fatalf("empty dir: absent=%v err=%v, want true nil", absent, err)
	}

	withVite := t.TempDir()
	writeFile(t, withVite, "vite.config.ts", "export default {}\n")
	absent, err = scaffoldAbsent(withVite)
	if err != nil || absent {
		t.Fatalf("vite.config.ts present: absent=%v err=%v, want false nil", absent, err)
	}

	withMain := t.TempDir()
	writeFile(t, withMain, "web/main.js", "console.log(1)\n")
	absent, err = scaffoldAbsent(withMain)
	if err != nil || absent {
		t.Fatalf("web/main.js present: absent=%v err=%v, want false nil", absent, err)
	}
}

func TestRefusalIncludesNonViteGuidance(t *testing.T) {
	dir := writeScaffoldFixture(t, "tampered", testGSXMainPristine)
	_, err := planScaffoldIntegration(dir, DefaultConfig())
	if err == nil || !strings.Contains(err.Error(), "without Vite") {
		t.Fatalf("refusal must mention the non-Vite path, got: %v", err)
	}
}
