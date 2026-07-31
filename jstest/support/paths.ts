import path from "node:path";
import { fileURLToPath } from "node:url";

const here = path.dirname(fileURLToPath(import.meta.url));

/** Repo root — jstest/support is two levels down. */
export const repoRoot = path.resolve(here, "..", "..");
export const jstestDir = path.join(repoRoot, "jstest");

/** Generated, gitignored: the example manifest and the compiled stylesheet. */
export const tmpDir = path.join(jstestDir, ".tmp");
export const manifestPath = path.join(tmpDir, "examples.json");
export const cssPath = path.join(tmpDir, "site.css");
export const foundationCSSPath = path.join(tmpDir, "foundation.css");

/**
 * 7799 is deliberately clear of the dev loop's 7777 and Vite's 5173, so a
 * running `make site-dev` never collides with a test run.
 *
 * This is the one source of truth for the harness's TCP port: the Go server
 * (jstest/harness) never reads GSXUI_HARNESS_PORT itself — playwright.config.ts
 * passes this value to it via `-addr 127.0.0.1:<port>` on the command line it
 * builds for `webServer`. Override with GSXUI_HARNESS_PORT to let multiple
 * concurrent worktrees each run Playwright / `make sweep-baseline` /
 * `make sweep-compare` without colliding on the same port.
 */
export const harnessPort = Number(process.env.GSXUI_HARNESS_PORT) || 7799;
export const baseURL = `http://127.0.0.1:${harnessPort}`;
