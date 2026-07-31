# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/composites-style-contract.spec.ts >> ButtonGroup restores only the visible tail before an aria-hidden select
- Location: jstest/specs/composites-style-contract.spec.ts:3:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/style-contract", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | test("ButtonGroup restores only the visible tail before an aria-hidden select", async ({
  4   |   page,
  5   | }) => {
> 6   |   const response = await page.goto("/f/style-contract");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  7   |   expect(response?.status()).toBe(200);
  8   | 
  9   |   const group = page.locator('[data-style-contract="button-group-tail"]');
  10  |   const radii = await group.evaluate((element) => {
  11  |     const corners = (selector: string) => {
  12  |       const css = getComputedStyle(element.querySelector(selector)!);
  13  |       return {
  14  |         topRight: css.borderTopRightRadius,
  15  |         bottomRight: css.borderBottomRightRadius,
  16  |       };
  17  |     };
  18  |     return {
  19  |       earlier: corners('[data-style-contract="button-group-earlier"]'),
  20  |       visibleTail: corners('[data-style-contract="button-group-visible-tail"]'),
  21  |     };
  22  |   });
  23  | 
  24  |   expect(radii).toEqual({
  25  |     earlier: { topRight: "0px", bottomRight: "0px" },
  26  |     visibleTail: { topRight: "10px", bottomRight: "10px" },
  27  |   });
  28  | });
  29  | 
  30  | test("horizontal ButtonGroup separators stay on the cross-axis", async ({ page }) => {
  31  |   const response = await page.goto("/x/button-group/basic");
  32  |   expect(response?.status()).toBe(200);
  33  | 
  34  |   const groups = page.locator('[data-gsxui-slot-button-group]');
  35  |   await expect(groups).toHaveCount(4);
  36  |   const clipboard = groups.nth(1);
  37  |   const separator = clipboard.locator(
  38  |     ':scope > [data-gsxui-slot-button-group-separator]',
  39  |   );
  40  | 
  41  |   await expect(separator).toHaveAttribute("data-orientation", "vertical");
  42  |   expect(
  43  |     await clipboard.evaluate((group) => {
  44  |       const groupRect = group.getBoundingClientRect();
  45  |       const separatorRect = group
  46  |         .querySelector('[data-gsxui-slot-button-group-separator]')!
  47  |         .getBoundingClientRect();
  48  |       return {
  49  |         separatorWidth: separatorRect.width,
  50  |         separatorHeight: separatorRect.height,
  51  |         groupHeight: groupRect.height,
  52  |         childrenInside: [...group.children].every((child) => {
  53  |           const rect = child.getBoundingClientRect();
  54  |           return (
  55  |             rect.left >= groupRect.left &&
  56  |             rect.right <= groupRect.right &&
  57  |             rect.top >= groupRect.top &&
  58  |             rect.bottom <= groupRect.bottom
  59  |           );
  60  |         }),
  61  |       };
  62  |     }),
  63  |   ).toEqual({
  64  |     separatorWidth: 1,
  65  |     separatorHeight: 32,
  66  |     groupHeight: 32,
  67  |     childrenInside: true,
  68  |   });
  69  | });
  70  | 
  71  | test("vertical ButtonGroup separators stay on the cross-axis", async ({ page }) => {
  72  |   const response = await page.goto("/x/button-group/basic");
  73  |   expect(response?.status()).toBe(200);
  74  | 
  75  |   const group = page.getByRole("group", { name: "Media controls" });
  76  |   const children = group.locator(":scope > *");
  77  |   await expect(children).toHaveCount(3);
  78  |   const separator = group.locator(
  79  |     ':scope > [data-gsxui-slot-button-group-separator]',
  80  |   );
  81  |   await expect(separator).toHaveAttribute("data-orientation", "horizontal");
  82  |   expect(
  83  |     await group.evaluate((element) => {
  84  |       const groupRect = element.getBoundingClientRect();
  85  |       const separatorRect = element
  86  |         .querySelector('[data-gsxui-slot-button-group-separator]')!
  87  |         .getBoundingClientRect();
  88  |       const button = element.lastElementChild!;
  89  |       const css = getComputedStyle(button);
  90  |       return {
  91  |         separatorWidth: separatorRect.width,
  92  |         separatorHeight: separatorRect.height,
  93  |         groupWidth: groupRect.width,
  94  |         bottomLeft: css.borderBottomLeftRadius,
  95  |         bottomRight: css.borderBottomRightRadius,
  96  |       };
  97  |     }),
  98  |   ).toEqual({
  99  |     separatorWidth: 32,
  100 |     separatorHeight: 1,
  101 |     groupWidth: 32,
  102 |     bottomLeft: "10px",
  103 |     bottomRight: "10px",
  104 |   });
  105 | });
  106 | 
```