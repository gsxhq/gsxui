# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/smoke.spec.ts >> no failed subresource requests
- Location: jstest/specs/smoke.spec.ts:55:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/toggle/basic", waiting until "load"

```

# Test source

```ts
  1  | import { expect, test } from "../support/fixtures";
  2  | import { examples } from "../support/manifest";
  3  | 
  4  | test("the manifest reaches the specs", () => {
  5  |   const all = examples();
  6  |   expect(all.length).toBeGreaterThan(100);
  7  |   expect(all).toContainEqual({
  8  |     component: "dropdown-menu",
  9  |     example: "checkboxes",
  10 |     url: "/x/dropdown-menu/checkboxes",
  11 |   });
  12 | });
  13 | 
  14 | test("markup, stylesheet and behavior JS all arrive", async ({ page }) => {
  15 |   await page.goto("/x/toggle/basic");
  16 | 
  17 |   // Markup: rendered through the real gsx component.
  18 |   const toggle = page.locator("[data-gsxui-toggle]").first();
  19 |   await expect(toggle).toBeVisible();
  20 | 
  21 |   // CSS: bg-background resolves to an opaque colour. Without the stylesheet
  22 |   // the body background is rgba(0, 0, 0, 0).
  23 |   const background = await page.evaluate(
  24 |     () => getComputedStyle(document.body).backgroundColor,
  25 |   );
  26 |   expect(background).not.toBe("rgba(0, 0, 0, 0)");
  27 | 
  28 |   // JS: ui/toggle.js is bound and flips both attributes on click.
  29 |   await expect(toggle).toHaveAttribute("aria-pressed", "false");
  30 |   await toggle.click();
  31 |   await expect(toggle).toHaveAttribute("aria-pressed", "true");
  32 |   await expect(toggle).toHaveAttribute("data-state", "on");
  33 | });
  34 | 
  35 | test("the shim records what the real modules would have registered", async ({ page }) => {
  36 |   await page.goto("/registrations");
  37 | 
  38 |   const registrations = await page.evaluate(() => window.__gsxuiRegistrations);
  39 |   expect(registrations.length).toBeGreaterThan(50);
  40 | 
  41 |   // Every entry is attributed to a real module, not "unknown" — the stack
  42 |   // walk in shim.js is the only thing that can break this.
  43 |   const modules = new Set(registrations.map((r) => r.module));
  44 |   expect(modules).not.toContain("unknown");
  45 |   expect(modules).toContain("dropdown-menu.js");
  46 |   expect(modules).toContain("toggle.js");
  47 | });
  48 | 
  49 | // Confirms the /static/ root choice in Task 1's harness: Tailwind bundles
  50 | // @fontsource-variable/geist, whose CSS carries url() references, and every
  51 | // page's shell links the compiled stylesheet. If those font requests (or any
  52 | // other subresource) 404, every page load would log a failed request and
  53 | // Task 4's clean-load invariant would be noisy — this keeps that invariant
  54 | // meaningful by catching real breakage here instead of filtering it out.
  55 | test("no failed subresource requests", async ({ page }) => {
  56 |   const failed: string[] = [];
  57 |   page.on("requestfailed", (r) => failed.push(r.url()));
  58 |   page.on("response", (r) => {
  59 |     if (r.status() >= 400) failed.push(`${r.status()} ${r.url()}`);
  60 |   });
> 61 |   await page.goto("/x/toggle/basic");
     |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  62 |   expect(failed).toEqual([]);
  63 | });
  64 | 
```