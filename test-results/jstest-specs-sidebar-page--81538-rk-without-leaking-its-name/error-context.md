# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sidebar-page.spec.ts >> sidebar icon collapse keeps a compact brand mark without leaking its name
- Location: jstest/specs/sidebar-page.spec.ts:68:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/components/sidebar", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | test("sidebar basic keeps separators and trailing controls inside their rows", async ({
  4   |   page,
  5   | }) => {
  6   |   await page.goto("/components/sidebar");
  7   | 
  8   |   const basic = page.frameLocator('iframe[src="/examples/sidebar/basic"]');
  9   |   const desktop = basic.locator("[data-gsxui-slot-sidebar-desktop]");
  10  |   const inner = desktop.locator("[data-gsxui-slot-sidebar-inner]");
  11  |   const separator = desktop.locator("[data-gsxui-slot-sidebar-separator]");
  12  | 
  13  |   const separatorGeometry = await Promise.all([
  14  |     inner.evaluate((element) => {
  15  |       const rect = element.getBoundingClientRect();
  16  |       return { left: rect.left, right: rect.right };
  17  |     }),
  18  |     separator.evaluate((element) => {
  19  |       const rect = element.getBoundingClientRect();
  20  |       return { left: rect.left, right: rect.right };
  21  |     }),
  22  |   ]);
  23  |   expect(separatorGeometry[1].left - separatorGeometry[0].left).toBeCloseTo(8);
  24  |   expect(separatorGeometry[0].right - separatorGeometry[1].right).toBeCloseTo(8);
  25  | 
  26  |   const decorationGeometry = await desktop
  27  |     .locator("[data-gsxui-slot-sidebar-menu-item]")
  28  |     .evaluateAll((items) =>
  29  |       items.flatMap((item) => {
  30  |         const itemRect = item.getBoundingClientRect();
  31  |         if (itemRect.width === 0 || itemRect.height === 0) return [];
  32  | 
  33  |         const button = item.querySelector("[data-gsxui-slot-sidebar-menu-button]");
  34  |         if (!(button instanceof HTMLElement)) return [];
  35  |         const buttonRect = button.getBoundingClientRect();
  36  |         const decorations = [
  37  |           ...item.querySelectorAll(
  38  |             ":scope > [data-gsxui-slot-sidebar-menu-action], " +
  39  |               ":scope > [data-gsxui-slot-sidebar-menu-badge]",
  40  |           ),
  41  |         ].filter((element) => element.getBoundingClientRect().width > 0);
  42  | 
  43  |         return decorations.map((decoration, index) => {
  44  |           const rect = decoration.getBoundingClientRect();
  45  |           const overlapsAnother = decorations.some((other, otherIndex) => {
  46  |             if (index === otherIndex) return false;
  47  |             const otherRect = other.getBoundingClientRect();
  48  |             return (
  49  |               rect.left < otherRect.right &&
  50  |               rect.right > otherRect.left &&
  51  |               rect.top < otherRect.bottom &&
  52  |               rect.bottom > otherRect.top
  53  |             );
  54  |           });
  55  |           return {
  56  |             withinButtonRow:
  57  |               rect.top >= buttonRect.top && rect.bottom <= buttonRect.bottom,
  58  |             overlapsAnother,
  59  |           };
  60  |         });
  61  |       }),
  62  |     );
  63  |   expect(decorationGeometry.length).toBeGreaterThan(0);
  64  |   expect(decorationGeometry.every((entry) => entry.withinButtonRow)).toBe(true);
  65  |   expect(decorationGeometry.every((entry) => !entry.overlapsAnother)).toBe(true);
  66  | });
  67  | 
  68  | test("sidebar icon collapse keeps a compact brand mark without leaking its name", async ({
  69  |   page,
  70  | }) => {
> 71  |   await page.goto("/components/sidebar");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  72  | 
  73  |   const icon = page.frameLocator(
  74  |     'iframe[src="/examples/sidebar/variants?_preview=icon-collapsed"]',
  75  |   );
  76  |   await page
  77  |     .locator('iframe[src="/examples/sidebar/variants?_preview=icon-collapsed"]')
  78  |     .scrollIntoViewIfNeeded();
  79  |   const brand = icon
  80  |     .locator("[data-gsxui-slot-sidebar-desktop]")
  81  |     .locator("[data-sidebar-example-brand]");
  82  |   const mark = brand.locator("[data-sidebar-example-brand-mark]");
  83  |   const name = brand.locator("[data-sidebar-example-brand-name]");
  84  | 
  85  |   await expect(mark).toHaveText("A");
  86  |   await expect(name).toHaveText("Acme Inc");
  87  | 
  88  |   const geometry = await brand.evaluate((element) => {
  89  |     const mark = element.querySelector("[data-sidebar-example-brand-mark]");
  90  |     const name = element.querySelector("[data-sidebar-example-brand-name]");
  91  |     if (!(mark instanceof HTMLElement) || !(name instanceof HTMLElement)) {
  92  |       throw new Error("sidebar example brand is missing its mark or name");
  93  |     }
  94  |     const rect = element.getBoundingClientRect();
  95  |     const markRect = mark.getBoundingClientRect();
  96  |     const nameRect = name.getBoundingClientRect();
  97  |     return {
  98  |       width: rect.width,
  99  |       overflowX: getComputedStyle(element).overflowX,
  100 |       markInside: markRect.left >= rect.left && markRect.right <= rect.right,
  101 |       nameStartsOutside: nameRect.left >= rect.right,
  102 |     };
  103 |   });
  104 |   expect(geometry.width).toBe(32);
  105 |   expect(geometry.overflowX).toBe("hidden");
  106 |   expect(geometry.markInside).toBe(true);
  107 |   expect(geometry.nameStartsOutside).toBe(true);
  108 | });
  109 | 
  110 | test("sidebar documentation contains each app shell inside its own viewport", async ({ page }) => {
  111 |   await page.goto("/components/sidebar");
  112 | 
  113 |   const previews = page.locator("iframe[data-site-isolated-preview]");
  114 |   await expect(previews).toHaveCount(10);
  115 |   await expect(page.locator("[data-gsxui-slot-sidebar-container]")).toHaveCount(0);
  116 |   await expect(
  117 |     page.locator('[data-site-docs-sidebar] a[href="/components/button"]'),
  118 |   ).toBeVisible();
  119 | 
  120 |   const basicFrame = page.locator('iframe[src="/examples/sidebar/basic"]');
  121 |   const basic = page.frameLocator('iframe[src="/examples/sidebar/basic"]');
  122 |   const wrapper = basic.locator("[data-gsxui-sidebar-wrapper]");
  123 |   const container = basic.locator("[data-gsxui-slot-sidebar-container]");
  124 |   await expect(basicFrame).toHaveAttribute("width", "1024");
  125 |   await expect
  126 |     .poll(() =>
  127 |       basic
  128 |         .locator("html")
  129 |         .evaluate(() => document.documentElement.clientWidth),
  130 |     )
  131 |     .toBe(1024);
  132 |   const previewGeometry = await basicFrame
  133 |     .locator("..")
  134 |     .evaluate((surface) => {
  135 |       const frame = surface.querySelector("iframe");
  136 |       if (!(frame instanceof HTMLIFrameElement)) {
  137 |         throw new Error("isolated preview surface is missing its iframe");
  138 |       }
  139 |       return {
  140 |         surfaceWidth: surface.getBoundingClientRect().width,
  141 |         frameWidth: frame.getBoundingClientRect().width,
  142 |         overflowX: getComputedStyle(surface).overflowX,
  143 |       };
  144 |     });
  145 |   expect(previewGeometry.surfaceWidth).toBe(640);
  146 |   expect(previewGeometry.frameWidth).toBe(1024);
  147 |   expect(previewGeometry.overflowX).toBe("hidden");
  148 |   const outerDocumentWidth = await page.locator("html").evaluate(() => ({
  149 |     client: document.documentElement.clientWidth,
  150 |     scroll: document.documentElement.scrollWidth,
  151 |   }));
  152 |   expect(outerDocumentWidth.scroll).toBe(outerDocumentWidth.client);
  153 |   await expect(wrapper).toHaveAttribute("data-state", "expanded");
  154 |   await expect(container).toBeVisible();
  155 | 
  156 |   const geometry = await container.evaluate((element) => {
  157 |     const rect = element.getBoundingClientRect();
  158 |     return {
  159 |       position: getComputedStyle(element).position,
  160 |       left: rect.left,
  161 |       right: rect.right,
  162 |       viewportWidth: document.documentElement.clientWidth,
  163 |     };
  164 |   });
  165 |   expect(geometry.position).toBe("fixed");
  166 |   expect(geometry.left).toBeGreaterThanOrEqual(0);
  167 |   expect(geometry.right).toBeLessThanOrEqual(geometry.viewportWidth);
  168 | 
  169 |   await basic.getByRole("button", { name: "Toggle Sidebar" }).first().click();
  170 |   await expect(wrapper).toHaveAttribute("data-state", "collapsed");
  171 | 
```