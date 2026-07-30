# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/composites-style-contract.spec.ts >> Carousel keeps native-scroll mechanics and caller spacing overrides
- Location: jstest/specs/composites-style-contract.spec.ts:107:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/carousel/sizes", waiting until "load"

```

# Test source

```ts
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
  107 | test("Carousel keeps native-scroll mechanics and caller spacing overrides", async ({ page }) => {
> 108 |   const response = await page.goto("/x/carousel/sizes");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  109 |   expect(response?.status()).toBe(200);
  110 | 
  111 |   const root = page.locator("[data-gsxui-carousel]");
  112 |   const viewport = page.locator("[data-gsxui-carousel-content]");
  113 |   const track = page.locator('[data-gsxui-slot-carousel-track]');
  114 |   const item = page.locator("[data-gsxui-carousel-item]").first();
  115 |   const previous = page.locator("[data-gsxui-carousel-prev]");
  116 |   const next = page.locator("[data-gsxui-carousel-next]");
  117 | 
  118 |   expect(
  119 |     await viewport.evaluate((el) => {
  120 |       const css = getComputedStyle(el);
  121 |       return { overflowX: css.overflowX, snap: css.scrollSnapType };
  122 |     }),
  123 |   ).toEqual({ overflowX: "auto", snap: "x mandatory" });
  124 |   expect(
  125 |     await track.evaluate((el) => {
  126 |       const css = getComputedStyle(el);
  127 |       return { display: css.display, marginLeft: css.marginLeft };
  128 |     }),
  129 |   ).toEqual({ display: "flex", marginLeft: "-4px" });
  130 |   expect(
  131 |     await item.evaluate((el) => {
  132 |       const css = getComputedStyle(el);
  133 |       return {
  134 |         flexBasis: css.flexBasis,
  135 |         flexGrow: css.flexGrow,
  136 |         flexShrink: css.flexShrink,
  137 |         paddingLeft: css.paddingLeft,
  138 |       };
  139 |     }),
  140 |   ).toEqual({
  141 |     flexBasis: "33.3333%",
  142 |     flexGrow: "0",
  143 |     flexShrink: "0",
  144 |     paddingLeft: "4px",
  145 |   });
  146 | 
  147 |   await expect(root).toHaveAttribute("data-current-index", "0");
  148 |   await expect(previous).toBeDisabled();
  149 |   await expect(next).toBeEnabled();
  150 |   await root.evaluate(
  151 |     (el: HTMLElement & { gsxuiCarousel: { scrollTo(index: number): void } }) =>
  152 |       el.gsxuiCarousel.scrollTo(4),
  153 |   );
  154 |   await expect(previous).toBeEnabled();
  155 |   await expect(next).toBeDisabled();
  156 | });
  157 | 
  158 | test("registered vertical Carousel covers every oriented style slot", async ({ page }) => {
  159 |   const response = await page.goto("/x/carousel/vertical");
  160 |   expect(response?.status()).toBe(200);
  161 | 
  162 |   const orientedSlots = [
  163 |     "carousel",
  164 |     "carousel-content",
  165 |     "carousel-track",
  166 |     "carousel-item",
  167 |     "carousel-previous",
  168 |     "carousel-next",
  169 |   ];
  170 |   for (const slot of orientedSlots) {
  171 |     await expect(
  172 |       page.locator(`[data-gsxui-slot-${slot}]`).first(),
  173 |       `${slot} must reflect the production example's vertical axis`,
  174 |     ).toHaveAttribute("data-orientation", "vertical");
  175 |   }
  176 | 
  177 |   const viewport = page.locator("[data-gsxui-carousel-content]");
  178 |   const track = page.locator('[data-gsxui-slot-carousel-track]');
  179 |   expect(
  180 |     await viewport.evaluate((el) => {
  181 |       const css = getComputedStyle(el);
  182 |       return { overflowY: css.overflowY, snap: css.scrollSnapType };
  183 |     }),
  184 |   ).toEqual({ overflowY: "auto", snap: "y mandatory" });
  185 |   await expect
  186 |     .poll(() => track.evaluate((el) => getComputedStyle(el).flexDirection))
  187 |     .toBe("column");
  188 | 
  189 |   const itemGeometry = await page
  190 |     .locator("[data-gsxui-carousel-item]")
  191 |     .evaluateAll((items) =>
  192 |       items.slice(0, 2).map((item) => item.getBoundingClientRect().height),
  193 |     );
  194 |   expect(itemGeometry).toHaveLength(2);
  195 |   expect(itemGeometry[0]).toBeCloseTo(100, 0);
  196 |   expect(itemGeometry[1]).toBeCloseTo(100, 0);
  197 | });
  198 | 
  199 | test("Resizable consumes dynamic flex values and remains keyboard operable", async ({ page }) => {
  200 |   const response = await page.goto("/x/resizable/handle");
  201 |   expect(response?.status()).toBe(200);
  202 | 
  203 |   const group = page.locator("[data-gsxui-resizable]").first();
  204 |   const panels = group.locator(":scope > [data-gsxui-resizable-panel]");
  205 |   const handle = group.locator(":scope > [data-gsxui-resizable-handle]");
  206 |   await expect(panels).toHaveCount(2);
  207 | 
  208 |   expect(
```