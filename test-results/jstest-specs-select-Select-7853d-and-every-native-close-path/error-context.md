# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/select.spec.ts >> Select trigger state follows open and every native close path
- Location: jstest/specs/select.spec.ts:3:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/select/basic", waiting until "load"

```

# Test source

```ts
  1  | import { expect, test } from "../support/fixtures";
  2  | 
  3  | test("Select trigger state follows open and every native close path", async ({
  4  |   page,
  5  | }) => {
> 6  |   const response = await page.goto("/x/select/basic");
     |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  7  |   expect(response?.status(), "select example response").toBe(200);
  8  | 
  9  |   const trigger = page.locator("[data-gsxui-slot-select-trigger]");
  10 |   const content = page.locator("[data-gsxui-slot-select-content]");
  11 | 
  12 |   await expect(trigger).toHaveAttribute("data-state", "closed");
  13 |   await expect(trigger).toHaveAttribute("aria-expanded", "false");
  14 | 
  15 |   const open = async () => {
  16 |     await trigger.click();
  17 |     await expect(content).toHaveAttribute("data-state", "open");
  18 |     await expect(trigger).toHaveAttribute("data-state", "open");
  19 |     await expect(trigger).toHaveAttribute("aria-expanded", "true");
  20 |   };
  21 | 
  22 |   const expectClosed = async () => {
  23 |     await expect(content).toHaveAttribute("data-state", "closed");
  24 |     await expect(trigger).toHaveAttribute("data-state", "closed");
  25 |     await expect(trigger).toHaveAttribute("aria-expanded", "false");
  26 |   };
  27 | 
  28 |   await open();
  29 |   await page.locator("[data-gsxui-slot-select-item]").first().click();
  30 |   await expectClosed();
  31 | 
  32 |   await open();
  33 |   await page.keyboard.press("Escape");
  34 |   await expectClosed();
  35 | 
  36 |   await open();
  37 |   await page.locator("[data-harness-root]").click({ position: { x: 1, y: 1 } });
  38 |   await expectClosed();
  39 | 
  40 |   await open();
  41 |   await content.evaluate((element: HTMLElement) => element.hidePopover());
  42 |   await expectClosed();
  43 | });
  44 | 
```