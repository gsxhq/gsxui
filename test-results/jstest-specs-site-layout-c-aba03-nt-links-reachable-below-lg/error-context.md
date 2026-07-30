# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/site-layout.spec.ts >> compact docs navigation keeps docs and component links reachable below lg
- Location: jstest/specs/site-layout.spec.ts:104:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/docs/theming", waiting until "load"

```

# Test source

```ts
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
  106 | }) => {
  107 |   await page.setViewportSize({ width: 900, height: 900 });
> 108 |   await page.goto("/docs/theming");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  109 | 
  110 |   const headerPadding = await page
  111 |     .locator("header > div")
  112 |     .evaluate((element) => getComputedStyle(element).paddingLeft);
  113 |   const footerPadding = await page
  114 |     .locator("[data-site-footer] > div")
  115 |     .evaluate((element) => getComputedStyle(element).paddingLeft);
  116 |   expect(headerPadding).toBe("24px");
  117 |   expect(footerPadding).toBe(headerPadding);
  118 | 
  119 |   await expect(page.locator("[data-site-docs-sidebar]")).toBeHidden();
  120 |   await expect(page.locator("[data-site-docs-toc]")).toBeHidden();
  121 |   const compactNav = page.locator("[data-site-docs-mobile-nav]");
  122 |   await expect(compactNav).toBeVisible();
  123 | 
  124 |   const trigger = compactNav.getByRole("button", {
  125 |     name: "Open documentation navigation",
  126 |   });
  127 |   await trigger.focus();
  128 |   await trigger.press("Enter");
  129 |   await expect(trigger).toHaveAttribute("aria-expanded", "true");
  130 | 
  131 |   const gettingStarted = compactNav.locator(
  132 |     'a[href="/docs/getting-started"]',
  133 |   );
  134 |   const button = compactNav.locator('a[href="/components/button"]');
  135 |   await expect(gettingStarted).toBeVisible();
  136 |   await expect(gettingStarted).toBeFocused();
  137 |   await expect(button).toBeVisible();
  138 | 
  139 |   await gettingStarted.click();
  140 |   await expect(page).toHaveURL(/\/docs\/getting-started$/);
  141 | 
  142 |   const nextCompactNav = page.locator("[data-site-docs-mobile-nav]");
  143 |   await nextCompactNav
  144 |     .getByRole("button", { name: "Open documentation navigation" })
  145 |     .press("Enter");
  146 |   await nextCompactNav.locator('a[href="/components/button"]').click();
  147 |   await expect(page).toHaveURL(/\/components\/button$/);
  148 | });
  149 | 
  150 | test("theme workspace contains desktop overflow and preserves narrow document order", async ({
  151 |   page,
  152 | }) => {
  153 |   await page.setViewportSize({ width: 1280, height: 900 });
  154 |   await page.goto("/site/theme");
  155 | 
  156 |   const stylePanel = page.locator("[data-theme-style-panel]");
  157 |   const previewPanel = page.locator("[data-theme-preview-panel]");
  158 |   const controlsPanel = page.locator("[data-theme-controls-panel]");
  159 | 
  160 |   await expect(stylePanel).toBeVisible();
  161 |   await expect(previewPanel).toBeVisible();
  162 |   await expect(controlsPanel).toBeVisible();
  163 |   await expect(page.locator("[data-theme-preview-frame]")).toHaveCount(1);
  164 | 
  165 |   const desktopGeometry = await page.evaluate(() => {
  166 |     const controls = document.querySelector<HTMLElement>(
  167 |       "[data-theme-controls-panel]",
  168 |     );
  169 |     const preview = document.querySelector<HTMLElement>(
  170 |       "[data-theme-preview-panel]",
  171 |     );
  172 |     if (!controls || !preview) {
  173 |       throw new Error("theme workspace panels are missing");
  174 |     }
  175 |     return {
  176 |       documentOverflow:
  177 |         document.documentElement.scrollHeight >
  178 |         document.documentElement.clientHeight,
  179 |       controlsWidth: controls.getBoundingClientRect().width,
  180 |       previewWidth: preview.getBoundingClientRect().width,
  181 |       controlsOverflow: getComputedStyle(controls).overflowY,
  182 |       controlsScrollHeight: controls.scrollHeight,
  183 |       controlsClientHeight: controls.clientHeight,
  184 |     };
  185 |   });
  186 |   expect(desktopGeometry.documentOverflow).toBe(false);
  187 |   expect(desktopGeometry.controlsWidth).toBeLessThan(
  188 |     desktopGeometry.previewWidth,
  189 |   );
  190 |   expect(desktopGeometry.controlsOverflow).toBe("auto");
  191 |   expect(desktopGeometry.controlsScrollHeight).toBeGreaterThan(
  192 |     desktopGeometry.controlsClientHeight,
  193 |   );
  194 | 
  195 |   await page.setViewportSize({ width: 900, height: 900 });
  196 | 
  197 |   const narrowGeometry = await page.evaluate(() => {
  198 |     const top = (selector: string) => {
  199 |       const element = document.querySelector<HTMLElement>(selector);
  200 |       if (!element) {
  201 |         throw new Error(`missing ${selector}`);
  202 |       }
  203 |       return element.getBoundingClientRect().top;
  204 |     };
  205 |     return {
  206 |       documentOverflow:
  207 |         document.documentElement.scrollHeight >
  208 |         document.documentElement.clientHeight,
```