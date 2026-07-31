# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sheet.spec.ts >> top direction demo renders a complete, padded sheet
- Location: jstest/specs/sheet.spec.ts:4:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sheet/directions", waiting until "load"

```

# Test source

```ts
  1  | import { expect, test } from "../support/fixtures";
  2  | 
  3  | for (const side of ["top", "bottom"]) {
  4  |   test(`${side} direction demo renders a complete, padded sheet`, async ({
  5  |     page,
  6  |   }) => {
> 7  |     const response = await page.goto("/x/sheet/directions");
     |                                 ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  8  |     expect(response?.status(), "sheet directions fixture response").toBe(200);
  9  | 
  10 |     await page.getByRole("button", { name: side, exact: true }).click();
  11 |     const content = page.locator(
  12 |       `[data-gsxui-slot-sheet-content][data-side="${side}"]`,
  13 |     );
  14 |     await expect(content).toBeVisible();
  15 |     await expect(
  16 |       content.locator('[data-gsxui-slot-sheet-header]'),
  17 |     ).toHaveCount(1);
  18 |     await expect(
  19 |       content.locator('[data-gsxui-slot-sheet-footer]'),
  20 |     ).toHaveCount(1);
  21 |     await expect(content.locator("input")).toHaveCount(2);
  22 | 
  23 |     const geometry = await content.evaluate((element) => {
  24 |       const panel = element.getBoundingClientRect();
  25 |       const title = element
  26 |         .querySelector('[data-gsxui-slot-sheet-title]')
  27 |         ?.getBoundingClientRect();
  28 |       return {
  29 |         height: panel.height,
  30 |         titleInsidePanel:
  31 |           !!title && title.top >= panel.top && title.bottom <= panel.bottom,
  32 |       };
  33 |     });
  34 |     expect(geometry.height).toBeGreaterThan(160);
  35 |     expect(geometry.titleInsidePanel).toBe(true);
  36 |   });
  37 | }
  38 | 
```