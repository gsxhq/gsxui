# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/dropdown.spec.ts >> one checkbox click changes state and emits exactly once
- Location: jstest/specs/dropdown.spec.ts:5:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/dropdown-menu/checkboxes", waiting until "load"

```

# Test source

```ts
  1  | import { expect, test } from "../support/fixtures";
  2  | 
  3  | const CHECKBOXES = "/x/dropdown-menu/checkboxes";
  4  | 
  5  | test("one checkbox click changes state and emits exactly once", async ({ page }) => {
> 6  |   await page.goto(CHECKBOXES);
     |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  7  |   await page.evaluate(() => {
  8  |     (window as any).__dropdownChanges = [];
  9  |     document.addEventListener("gsxui:change", (event) => {
  10 |       (window as any).__dropdownChanges.push((event as CustomEvent).detail);
  11 |     });
  12 |   });
  13 | 
  14 |   const trigger = page.locator("[data-gsxui-dropdown-trigger]");
  15 |   const content = page.locator("[data-gsxui-dropdown-content]");
  16 |   const item = page.locator(
  17 |     '[data-gsxui-dropdown-checkbox-item][data-value="statusbar"]',
  18 |   );
  19 | 
  20 |   await trigger.click();
  21 |   await expect(trigger).toHaveAttribute("aria-expanded", "true");
  22 |   expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  23 |   await expect(item).toHaveAttribute("aria-checked", "false");
  24 | 
  25 |   await item.click();
  26 | 
  27 |   await expect(item).toHaveAttribute("aria-checked", "true");
  28 |   expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  29 |   expect(await page.evaluate(() => (window as any).__dropdownChanges)).toEqual([
  30 |     { checked: true, value: "statusbar" },
  31 |   ]);
  32 | });
  33 | 
```