# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/navigation-menu.spec.ts >> keyboard Tab focus opens the trigger panel
- Location: jstest/specs/navigation-menu.spec.ts:28:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/navigation-menu/basic", waiting until "load"

```

# Test source

```ts
  1  | import { expect, test } from "../support/fixtures";
  2  | 
  3  | const BASIC = "/x/navigation-menu/basic";
  4  | 
  5  | async function mouseClick(
  6  |   page: import("@playwright/test").Page,
  7  |   locator: import("@playwright/test").Locator,
  8  | ) {
  9  |   const box = await locator.boundingBox();
  10 |   if (!box) throw new Error("navigation-menu trigger has no bounding box");
  11 |   await page.mouse.click(box.x + box.width / 2, box.y + box.height / 2);
  12 | }
  13 | 
  14 | test("a mouse click opens and the second click closes", async ({ page }) => {
  15 |   await page.goto(BASIC);
  16 |   const trigger = page.locator("[data-gsxui-navigation-menu-trigger]");
  17 |   const content = page.locator("[data-gsxui-navigation-menu-content]");
  18 | 
  19 |   await mouseClick(page, trigger);
  20 |   await expect(trigger).toHaveAttribute("aria-expanded", "true");
  21 |   expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  22 | 
  23 |   await mouseClick(page, trigger);
  24 |   await expect(trigger).toHaveAttribute("aria-expanded", "false");
  25 |   expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(false);
  26 | });
  27 | 
  28 | test("keyboard Tab focus opens the trigger panel", async ({ page }) => {
> 29 |   await page.goto(BASIC);
     |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  30 |   const home = page.locator('[data-gsxui-navigation-menu-link]').first();
  31 |   const trigger = page.locator("[data-gsxui-navigation-menu-trigger]");
  32 |   const content = page.locator("[data-gsxui-navigation-menu-content]");
  33 | 
  34 |   await home.focus();
  35 |   await page.keyboard.press("Tab");
  36 | 
  37 |   await expect(trigger).toBeFocused();
  38 |   await expect(trigger).toHaveAttribute("aria-expanded", "true");
  39 |   expect(await content.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  40 | });
  41 | 
```