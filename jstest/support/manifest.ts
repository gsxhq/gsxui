import { readFileSync } from "node:fs";
import { manifestPath } from "./paths";

export type ExampleEntry = {
  component: string;
  example: string;
  url: string;
};

/**
 * Reads the manifest Playwright's globalSetup generated. Synchronous on
 * purpose: spec files call this at module scope to generate one test per
 * example, and Playwright's test declaration is not async.
 */
export function examples(): ExampleEntry[] {
  let raw: string;
  try {
    raw = readFileSync(manifestPath, "utf8");
  } catch (err) {
    throw new Error(
      `example manifest missing at ${manifestPath}. It is written by ` +
        `jstest/global-setup.ts, which only runs under Playwright — use ` +
        `\`make test-js\` or \`npx playwright test --config jstest/playwright.config.ts\`. ` +
        `(${err})`,
    );
  }
  const entries = JSON.parse(raw) as ExampleEntry[];
  if (!Array.isArray(entries) || entries.length === 0) {
    throw new Error(`example manifest at ${manifestPath} is empty`);
  }
  return entries;
}
