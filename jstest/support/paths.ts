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
 */
export const harnessPort = 7799;
export const baseURL = `http://127.0.0.1:${harnessPort}`;
