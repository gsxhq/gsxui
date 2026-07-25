import { readdirSync } from "node:fs";
import path from "node:path";
import { test as base } from "@playwright/test";
import type { Registration } from "./globals";
import { baseURL, repoRoot } from "./paths";

/**
 * The behavior modules the recorded registry must account for: every ui/*.js
 * except the delegation core (gsxui.js, which registers nothing itself) and
 * the barrel (index.js, which only imports).
 *
 * Derived from the directory rather than hardcoded, and deliberately NOT
 * derived from ui/index.js's import list — a module missing from the barrel
 * is exactly one of the failures this is meant to catch, and reading the
 * barrel to build the expectation would make that failure invisible.
 *
 * The standing assumption is that every behavior module registers at least
 * one delegated handler; all 21 do today. A future ui/*.js that binds
 * nothing through on() would fail this, and belongs here as a named
 * exception with a reason rather than silently tolerated.
 */
export function behaviorModules(): string[] {
  const uiDir = path.join(repoRoot, "ui");
  const names = readdirSync(uiDir)
    .filter((name) => name.endsWith(".js"))
    .filter((name) => name !== "gsxui.js" && name !== "index.js")
    .sort();
  if (names.length === 0) {
    throw new Error(`no behavior modules found under ${uiDir}`);
  }
  return names;
}

/**
 * `registrations` is the full delegation registry, recorded once per worker
 * by loading /registrations — a blank page whose only script is
 * /shim/index.js, where on() records instead of binding.
 *
 * Worker-scoped rather than collected in globalSetup because globalSetup's
 * ordering against webServer startup is not something to depend on: a
 * fixture runs after the server is definitely up.
 *
 * The completeness guard is a module-SET comparison, not a non-empty check.
 * Non-empty is nearly vacuous: adding ui/newthing.js and forgetting its
 * ui/index.js import line leaves the module out of the registry AND out of
 * production, silently; and a throw while evaluating the last import
 * (tooltip.js) drops its 4 registrations and leaves 122 of 126 — still
 * "non-empty". The pageerror listener covers that second case from the other
 * side, so an evaluation throw fails loudly here instead of quietly
 * shortening the list.
 */
export const test = base.extend<{}, { registrations: Registration[] }>({
  registrations: [
    async ({ browser }, use) => {
      const page = await browser.newPage();
      const pageErrors: string[] = [];
      page.on("pageerror", (err) => pageErrors.push(String(err)));
      // Explicit absolute URL. A page from browser.newPage() only resolves a
      // relative path because Playwright back-fills use.baseURL onto
      // manually created contexts — an internal, not a contract.
      await page.goto(`${baseURL}/registrations`);
      const recorded = await page.evaluate(() => window.__gsxuiRegistrations);
      await page.close();

      if (pageErrors.length > 0) {
        throw new Error(
          `/registrations threw while evaluating the shimmed module graph, so the ` +
            `recorded registry is truncated: ${pageErrors.join("; ")}`,
        );
      }
      if (!recorded || recorded.length === 0) {
        throw new Error("/registrations recorded nothing — the shim did not run");
      }

      const expected = behaviorModules();
      const seen = new Set(recorded.map((r) => r.module));
      const missing = expected.filter((name) => !seen.has(name));
      const unrecognised = [...seen].filter((name) => !expected.includes(name)).sort();
      if (missing.length > 0 || unrecognised.length > 0) {
        throw new Error(
          `/registrations recorded ${seen.size} of ${expected.length} behavior modules ` +
            `(${recorded.length} registrations).` +
            (missing.length > 0
              ? ` Missing: ${missing.join(", ")} — either the module is absent from ` +
                `ui/index.js (so it is absent from production too), or it registered nothing.`
              : "") +
            (unrecognised.length > 0
              ? ` Unrecognised: ${unrecognised.join(", ")} — not a ui/*.js behavior module, ` +
                `so shim.js's stack walk mis-attributed them.`
              : ""),
        );
      }

      await use(recorded);
    },
    { scope: "worker" },
  ],
});

export { expect } from "@playwright/test";
