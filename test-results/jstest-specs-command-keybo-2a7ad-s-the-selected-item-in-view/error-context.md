# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/command.spec.ts >> keyboard selection keeps the selected item in view
- Location: jstest/specs/command.spec.ts:26:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/command/basic", waiting until "load"

```

# Test source

```ts
  1  | import type { Page } from "@playwright/test";
  2  | import { expect, test } from "../support/fixtures";
  3  | 
  4  | async function recordCommandScrollRequests(page: Page) {
  5  |   await page.addInitScript(() => {
  6  |     (window as any).__commandScrollRequests = 0;
  7  |     const scrollIntoView = Element.prototype.scrollIntoView;
  8  |     Element.prototype.scrollIntoView = function (...args) {
  9  |       if ((this as Element).matches("[data-gsxui-command-item]")) {
  10 |         (window as any).__commandScrollRequests++;
  11 |       }
  12 |       return scrollIntoView.apply(this, args as [ScrollIntoViewOptions]);
  13 |     };
  14 |   });
  15 | }
  16 | 
  17 | test("initial selection does not request a document scroll", async ({ page }) => {
  18 |   await recordCommandScrollRequests(page);
  19 |   await page.goto("/x/command/basic");
  20 | 
  21 |   expect(
  22 |     await page.evaluate(() => (window as any).__commandScrollRequests),
  23 |   ).toBe(0);
  24 | });
  25 | 
  26 | test("keyboard selection keeps the selected item in view", async ({ page }) => {
  27 |   await recordCommandScrollRequests(page);
> 28 |   await page.goto("/x/command/basic");
     |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  29 |   await page.evaluate(() => ((window as any).__commandScrollRequests = 0));
  30 | 
  31 |   await page.locator("[data-gsxui-command-input]").first().press("ArrowDown");
  32 | 
  33 |   expect(
  34 |     await page.evaluate(() => (window as any).__commandScrollRequests),
  35 |   ).toBe(1);
  36 | });
  37 | 
```