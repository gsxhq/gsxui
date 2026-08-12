import path from "node:path";
import { defineConfig, devices } from "@playwright/test";
import { baseURL, harnessPort, jstestDir, repoRoot } from "./support/paths";

export default defineConfig({
  testDir: "./specs",
  // Browser pixels are characterized in Linux CI but developers replay them
  // on other hosts. Keep the project discriminator while deliberately
  // omitting Playwright's default platform suffix from snapshot paths.
  snapshotPathTemplate: "{testDir}/{testFilePath}-snapshots/{arg}-{projectName}{ext}",
  // Playwright's built-in default resolves outputDir against the nearest
  // package.json directory (the repo root), not the config file's
  // directory — that would leak a top-level /test-results/ next to
  // everything else. Force it under jstest/ so .gitignore's
  // `jstest/test-results/` entry actually covers it.
  outputDir: path.join(jstestDir, "test-results"),
  globalSetup: "./global-setup.ts",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  // CI gets "github" for inline annotations plus "html" — the html reporter
  // is what actually copies trace attachments into its output folder, so
  // the workflow's on-failure artifact upload has real content to find.
  // Local runs stay on plain "list"; forcing an html report there would
  // leave an untracked jstest/playwright-report/ directory to open manually.
  reporter: process.env.CI
    ? [
        ["github"],
        ["html", { outputFolder: path.join(jstestDir, "playwright-report"), open: "never" }],
      ]
    : "list",
  use: {
    baseURL,
    trace: "on-first-retry",
  },
  projects: [
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        viewport: { width: 1280, height: 900 },
        // Cross-platform screenshot determinism: FreeType (Linux CI) hints
        // and stem-darkens variable fonts where CoreText (macOS, where the
        // committed baselines are generated) ignores hinting entirely, so
        // text-dense cards rendered 2-5% bolder ON CI ONLY — past the 1%
        // visual gate — from the per-style gate's very first Linux run.
        // --font-render-hinting=none converges the rasterizers on the
        // dominant divergence (hinted stem darkening). Deliberately NOT
        // paired with --disable-lcd-text/--disable-font-subpixel-positioning:
        // measured on CI, that trio fixed the bold-text cards but thinned
        // small/light text enough to push a different set (feedback,
        // lyra/mira calendar) past the same 1% gate.
        launchOptions: {
          args: ["--font-render-hinting=none"],
        },
      },
    },
  ],
  webServer: {
    command: `go run ./jstest/harness -addr 127.0.0.1:${harnessPort} -root .`,
    cwd: repoRoot,
    url: `${baseURL}/healthz`,
    reuseExistingServer: !process.env.CI,
    timeout: 120_000,
    stdout: "pipe",
    stderr: "pipe",
  },
});
