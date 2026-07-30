# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/style-visual.spec.ts >> sidebar/variants?_preview=none keeps its desktop light and dark visual contract
- Location: jstest/specs/style-visual.spec.ts:669:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sidebar/variants?_preview=none", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | const desktopRoutes = [
  4   |   "accordion/basic",
  5   |   "alert/variants",
  6   |   "button/variants",
  7   |   "card/compound",
  8   |   "calendar/loadedrange",
  9   |   "checkbox/states",
  10  |   "combobox/basic",
  11  |   "dialog/basic",
  12  |   "dropdown-menu/basic",
  13  |   "field/invalid",
  14  |   "navigation-menu/mega",
  15  |   "sidebar/variants",
  16  |   "sidebar/variants?_preview=floating",
  17  |   "sidebar/variants?_preview=inset",
  18  |   "sidebar/variants?_preview=offcanvas",
  19  |   "sidebar/variants?_preview=icon",
  20  |   "sidebar/variants?_preview=none",
  21  |   "sidebar/variants?_preview=right-collapsed",
  22  |   "sidebar/variants?_preview=icon-collapsed",
  23  |   "sonner/types",
  24  |   "tabs/basic",
  25  | ] as const;
  26  | 
  27  | const mobileRoutes = [
  28  |   "button/variants",
  29  |   "dialog/basic",
  30  |   "sidebar/variants",
  31  |   "calendar/basic",
  32  | ] as const;
  33  | 
  34  | const screenshotOptions = {
  35  |   animations: "disabled" as const,
  36  |   caret: "hide" as const,
  37  |   maxDiffPixelRatio: 0.01,
  38  | };
  39  | 
  40  | type VisualRoute = (typeof desktopRoutes)[number] | (typeof mobileRoutes)[number];
  41  | 
  42  | function snapshotSlug(route: VisualRoute) {
  43  |   return route.replace("?_preview=", "-").replace("/", "-");
  44  | }
  45  | 
  46  | async function prepareVisualRoute(
  47  |   page: import("@playwright/test").Page,
  48  |   route: VisualRoute,
  49  |   theme: "light" | "dark",
  50  | ) {
> 51  |   const response = await page.goto(`/x/${route}`);
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  52  |   expect(response?.status(), `${route} fixture response`).toBe(200);
  53  | 
  54  |   await page.evaluate(async (isDark) => {
  55  |     document.documentElement.classList.toggle("dark", isDark);
  56  |     await document.fonts.ready;
  57  |     await new Promise<void>((resolve) => requestAnimationFrame(() => resolve()));
  58  |   }, theme === "dark");
  59  | 
  60  |   if (route === "dialog/basic") {
  61  |     await page.getByRole("button", { name: "Delete account" }).click();
  62  |     await expect(page.getByRole("dialog")).toBeVisible();
  63  |   }
  64  | 
  65  |   if (route === "dropdown-menu/basic") {
  66  |     await page.getByRole("button", { name: "Options" }).click();
  67  |     await expect(page.locator("[data-gsxui-dropdown-content]")).toBeVisible();
  68  |   }
  69  | 
  70  |   if (route === "sonner/types") {
  71  |     const triggers = page.locator("button[data-gsxui-toast]");
  72  |     await expect(triggers).toHaveCount(5);
  73  |     for (let i = 0; i < 5; i++) {
  74  |       await triggers.nth(i).click();
  75  |     }
  76  |     const toasts = page.locator("#gsxui-toaster > [data-gsxui-toast]");
  77  |     await expect(toasts).toHaveCount(5);
  78  |     await toasts.last().hover();
  79  |     await expect(toasts.nth(2)).toBeVisible();
  80  |   }
  81  | 
  82  |   if (route === "sidebar/variants" && (page.viewportSize()?.width ?? 0) < 768) {
  83  |     await page.getByRole("button", { name: "Toggle Sidebar" }).first().click();
  84  |     await expect(page.getByRole("dialog")).toBeVisible();
  85  |   }
  86  | }
  87  | 
  88  | test("caller utilities override the Button defaults", async ({ page }) => {
  89  |   const response = await page.goto("/f/style-contract");
  90  |   expect(response?.status(), "style contract fixture response").toBe(200);
  91  | 
  92  |   const override = page.getByRole("button", { name: "Caller override" });
  93  |   expect(
  94  |     await override.evaluate((el) => {
  95  |       const css = getComputedStyle(el);
  96  |       return {
  97  |         height: css.height,
  98  |         borderRadius: css.borderRadius,
  99  |         display: css.display,
  100 |       };
  101 |     }),
  102 |   ).toEqual({
  103 |     height: "48px",
  104 |     borderRadius: "0px",
  105 |     display: "inline-flex",
  106 |   });
  107 | });
  108 | 
  109 | test("dark primitive states use their dark semantic colors", async ({ page }) => {
  110 |   const response = await page.goto("/f/style-contract");
  111 |   expect(response?.status(), "style contract fixture response").toBe(200);
  112 |   await page.evaluate(() => document.documentElement.classList.add("dark"));
  113 |   const finishTransitions = () =>
  114 |     page.evaluate(() => {
  115 |       for (const animation of document.getAnimations()) {
  116 |         animation.finish();
  117 |       }
  118 |     });
  119 |   await finishTransitions();
  120 | 
  121 |   const computed = async (
  122 |     selector: string,
  123 |     property: "backgroundColor" | "boxShadow",
  124 |   ) =>
  125 |     page.locator(selector).evaluate(
  126 |       (element, name) => getComputedStyle(element)[name],
  127 |       property,
  128 |     );
  129 | 
  130 |   const darkDestructive = await computed(
  131 |     '[data-style-contract-reference="dark-destructive"]',
  132 |     "backgroundColor",
  133 |   );
  134 |   expect(
  135 |     await computed(
  136 |       '[data-style-contract="dark-button-destructive"]',
  137 |       "backgroundColor",
  138 |     ),
  139 |   ).toBe(darkDestructive);
  140 |   expect(
  141 |     await computed(
  142 |       '[data-style-contract="dark-badge-destructive"]',
  143 |       "backgroundColor",
  144 |     ),
  145 |   ).toBe(darkDestructive);
  146 |   const destructive = page.locator(
  147 |     '[data-style-contract="dark-button-destructive"]',
  148 |   );
  149 |   await destructive.hover();
  150 |   await finishTransitions();
  151 |   expect(
```