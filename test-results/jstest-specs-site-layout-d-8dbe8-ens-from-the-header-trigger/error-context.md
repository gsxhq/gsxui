# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/site-layout.spec.ts >> documentation search opens from the header trigger
- Location: jstest/specs/site-layout.spec.ts:3:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/components", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | test("documentation search opens from the header trigger", async ({ page }) => {
  4   |   await page.setViewportSize({ width: 1280, height: 900 });
> 5   |   await page.goto("/components");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  6   | 
  7   |   const trigger = page.getByRole("button", { name: /Search docs/ });
  8   |   const dialog = page.locator("dialog[data-gsxui-command-dialog]");
  9   |   const input = dialog.getByPlaceholder("Search documentation...");
  10  | 
  11  |   await expect(dialog).not.toBeVisible();
  12  |   await trigger.click();
  13  |   await expect(dialog).toBeVisible();
  14  |   await expect(trigger).toHaveAttribute("aria-expanded", "true");
  15  |   await expect(input).toBeFocused();
  16  | 
  17  |   await page.keyboard.press("Escape");
  18  |   await expect(dialog).not.toBeVisible();
  19  |   await expect(trigger).toHaveAttribute("aria-expanded", "false");
  20  | });
  21  | 
  22  | test("documentation rails respond around the fixed 640px article", async ({
  23  |   page,
  24  | }) => {
  25  |   await page.setViewportSize({ width: 1280, height: 900 });
  26  |   await page.goto("/docs/theming");
  27  | 
  28  |   await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
  29  |   await expect(page.locator("[data-site-docs-toc]")).toBeVisible();
  30  |   expect(
  31  |     await page
  32  |       .locator("[data-site-docs-article]")
  33  |       .evaluate((element) => element.getBoundingClientRect().width),
  34  |   ).toBe(640);
  35  | 
  36  |   await page.setViewportSize({ width: 1100, height: 900 });
  37  |   await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
  38  |   await expect(page.locator("[data-site-docs-toc]")).toBeHidden();
  39  | 
  40  |   await page.setViewportSize({ width: 900, height: 900 });
  41  |   await expect(page.locator("[data-site-docs-sidebar]")).toBeHidden();
  42  |   await expect(page.locator("[data-site-docs-toc]")).toBeHidden();
  43  | 
  44  |   await page.setViewportSize({ width: 1280, height: 900 });
  45  |   await page.goto("/components");
  46  |   await expect(page.locator("[data-site-docs-sidebar]")).toBeVisible();
  47  |   await expect(page.locator("[data-site-docs-toc]")).toHaveCount(0);
  48  | 
  49  |   const sidebarNav = page.locator("[data-site-docs-sidebar] nav");
  50  |   await sidebarNav.evaluate((element) => {
  51  |     element.scrollTop = element.scrollHeight;
  52  |   });
  53  |   const sidebarEnd = await page.evaluate(() => {
  54  |     const lastLink = document.querySelector<HTMLElement>(
  55  |       '[data-site-docs-sidebar] a[href="/components/tooltip"]',
  56  |     );
  57  |     if (!lastLink) {
  58  |       throw new Error("final documentation sidebar link is missing");
  59  |     }
  60  |     return {
  61  |       documentScrollY: window.scrollY,
  62  |       linkBottom: lastLink.getBoundingClientRect().bottom,
  63  |       viewportBottom: document.documentElement.clientHeight,
  64  |     };
  65  |   });
  66  |   expect(sidebarEnd.documentScrollY).toBe(0);
  67  |   expect(sidebarEnd.linkBottom).toBeLessThanOrEqual(sidebarEnd.viewportBottom);
  68  | });
  69  | 
  70  | test("component table of contents follows existing hashes and observed headings", async ({
  71  |   page,
  72  | }) => {
  73  |   await page.setViewportSize({ width: 1280, height: 900 });
  74  |   await page.goto("/components/button#example-variants");
  75  | 
  76  |   const basicLink = page.locator(
  77  |     '[data-site-toc-link][href="#example-basic"]',
  78  |   );
  79  |   const variantsLink = page.locator(
  80  |     '[data-site-toc-link][href="#example-variants"]',
  81  |   );
  82  |   await expect(page.locator("#example-variants")).toHaveText("Variants");
  83  |   await expect(variantsLink).toHaveAttribute("data-active", "");
  84  |   await expect(variantsLink).toHaveAttribute("aria-current", "location");
  85  |   await expect(basicLink).not.toHaveAttribute("data-active", "");
  86  |   await expect(basicLink).not.toHaveAttribute("aria-current", "location");
  87  | 
  88  |   await basicLink.click();
  89  |   await expect(page).toHaveURL(/#example-basic$/);
  90  |   await expect(basicLink).toHaveAttribute("data-active", "");
  91  |   await expect(basicLink).toHaveAttribute("aria-current", "location");
  92  |   await expect(variantsLink).not.toHaveAttribute("data-active", "");
  93  | 
  94  |   await page.locator("#example-variants").evaluate((heading) => {
  95  |     window.scrollTo({
  96  |       top: heading.getBoundingClientRect().top + window.scrollY - 100,
  97  |     });
  98  |   });
  99  |   await expect(variantsLink).toHaveAttribute("data-active", "");
  100 |   await expect(variantsLink).toHaveAttribute("aria-current", "location");
  101 |   await expect(basicLink).not.toHaveAttribute("data-active", "");
  102 | });
  103 | 
  104 | test("compact docs navigation keeps docs and component links reachable below lg", async ({
  105 |   page,
```