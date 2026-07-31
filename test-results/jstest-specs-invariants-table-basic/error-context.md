# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/invariants.spec.ts >> table/basic
- Location: jstest/specs/invariants.spec.ts:14:3

# Error details

```
Error: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7799/registrations
Call log:
  - navigating to "http://127.0.0.1:7799/registrations", waiting until "load"

```

# Test source

```ts
  1   | import { readdirSync } from "node:fs";
  2   | import path from "node:path";
  3   | import { test as base } from "@playwright/test";
  4   | import type { Registration } from "./globals";
  5   | import { baseURL, repoRoot } from "./paths";
  6   | 
  7   | /**
  8   |  * The behavior modules the recorded registry must account for: every ui/*.js
  9   |  * except the delegation core (gsxui.js, which registers nothing itself) and
  10  |  * the barrel (index.js, which only imports).
  11  |  *
  12  |  * Derived from the directory rather than hardcoded, and deliberately NOT
  13  |  * derived from ui/index.js's import list — a module missing from the barrel
  14  |  * is exactly one of the failures this is meant to catch, and reading the
  15  |  * barrel to build the expectation would make that failure invisible.
  16  |  *
  17  |  * The standing assumption is that every behavior module registers at least
  18  |  * one delegated handler; all 21 do today. A future ui/*.js that binds
  19  |  * nothing through on() would fail this, and belongs here as a named
  20  |  * exception with a reason rather than silently tolerated.
  21  |  */
  22  | export function behaviorModules(): string[] {
  23  |   const uiDir = path.join(repoRoot, "ui");
  24  |   const names = readdirSync(uiDir)
  25  |     .filter((name) => name.endsWith(".js"))
  26  |     .filter((name) => name !== "gsxui.js" && name !== "index.js")
  27  |     .sort();
  28  |   if (names.length === 0) {
  29  |     throw new Error(`no behavior modules found under ${uiDir}`);
  30  |   }
  31  |   return names;
  32  | }
  33  | 
  34  | /**
  35  |  * `registrations` is the full delegation registry, recorded once per worker
  36  |  * by loading /registrations — a blank page whose only script is
  37  |  * /shim/index.js, where on() records instead of binding.
  38  |  *
  39  |  * Worker-scoped rather than collected in globalSetup because globalSetup's
  40  |  * ordering against webServer startup is not something to depend on: a
  41  |  * fixture runs after the server is definitely up.
  42  |  *
  43  |  * The completeness guard is a module-SET comparison, not a non-empty check.
  44  |  * Non-empty is nearly vacuous: adding ui/newthing.js and forgetting its
  45  |  * ui/index.js import line leaves the module out of the registry AND out of
  46  |  * production, silently; and a throw while evaluating the last import
  47  |  * (tooltip.js) drops its 4 registrations and leaves 122 of 126 — still
  48  |  * "non-empty". The pageerror listener covers that second case from the other
  49  |  * side, so an evaluation throw fails loudly here instead of quietly
  50  |  * shortening the list.
  51  |  */
  52  | export const test = base.extend<{}, { registrations: Registration[] }>({
  53  |   registrations: [
  54  |     async ({ browser }, use) => {
  55  |       const page = await browser.newPage();
  56  |       const pageErrors: string[] = [];
  57  |       page.on("pageerror", (err) => pageErrors.push(String(err)));
  58  |       // Explicit absolute URL. A page from browser.newPage() only resolves a
  59  |       // relative path because Playwright back-fills use.baseURL onto
  60  |       // manually created contexts — an internal, not a contract.
> 61  |       await page.goto(`${baseURL}/registrations`);
      |                  ^ Error: page.goto: net::ERR_CONNECTION_REFUSED at http://127.0.0.1:7799/registrations
  62  |       const recorded = await page.evaluate(() => window.__gsxuiRegistrations);
  63  |       await page.close();
  64  | 
  65  |       if (pageErrors.length > 0) {
  66  |         throw new Error(
  67  |           `/registrations threw while evaluating the shimmed module graph, so the ` +
  68  |             `recorded registry is truncated: ${pageErrors.join("; ")}`,
  69  |         );
  70  |       }
  71  |       if (!recorded || recorded.length === 0) {
  72  |         throw new Error("/registrations recorded nothing — the shim did not run");
  73  |       }
  74  | 
  75  |       const expected = behaviorModules();
  76  |       const seen = new Set(recorded.map((r) => r.module));
  77  |       const missing = expected.filter((name) => !seen.has(name));
  78  |       const unrecognised = [...seen].filter((name) => !expected.includes(name)).sort();
  79  |       if (missing.length > 0 || unrecognised.length > 0) {
  80  |         throw new Error(
  81  |           `/registrations recorded ${seen.size} of ${expected.length} behavior modules ` +
  82  |             `(${recorded.length} registrations).` +
  83  |             (missing.length > 0
  84  |               ? ` Missing: ${missing.join(", ")} — either the module is absent from ` +
  85  |                 `ui/index.js (so it is absent from production too), or it registered nothing.`
  86  |               : "") +
  87  |             (unrecognised.length > 0
  88  |               ? ` Unrecognised: ${unrecognised.join(", ")} — not a ui/*.js behavior module, ` +
  89  |                 `so shim.js's stack walk mis-attributed them.`
  90  |               : ""),
  91  |         );
  92  |       }
  93  | 
  94  |       await use(recorded);
  95  |     },
  96  |     { scope: "worker" },
  97  |   ],
  98  | });
  99  | 
  100 | export { expect } from "@playwright/test";
  101 | 
```