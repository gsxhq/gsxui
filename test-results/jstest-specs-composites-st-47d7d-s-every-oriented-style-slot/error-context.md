# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/composites-style-contract.spec.ts >> registered vertical Carousel covers every oriented style slot
- Location: jstest/specs/composites-style-contract.spec.ts:158:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/carousel/vertical", waiting until "load"

```

# Test source

```ts
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
  108 |   const response = await page.goto("/x/carousel/sizes");
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
> 159 |   const response = await page.goto("/x/carousel/vertical");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  209 |     await group.evaluate((el) => getComputedStyle(el).display),
  210 |   ).toBe("flex");
  211 |   expect(
  212 |     await panels.first().evaluate((el) => {
  213 |       const css = getComputedStyle(el);
  214 |       return { basis: css.flexBasis, grow: css.flexGrow, overflow: css.overflow };
  215 |     }),
  216 |   ).toEqual({ basis: "0px", grow: "25", overflow: "hidden" });
  217 | 
  218 |   await expect.poll(() => handle.evaluate((el) => el.getBoundingClientRect().width)).toBe(1);
  219 |   await page.addStyleTag({
  220 |     content: `
  221 |       /* No @layer wrapper, not :where(): Resizable is migrated to the slot
  222 |          axis, so its width now comes from a literal Tailwind utility class
  223 |          (aria-[orientation=vertical]:w-px) compiled into Tailwind's own
  224 |          internal @layer utilities — which always outranks @layer components
  225 |          regardless of specificity or source order, by CSS cascade-layer
  226 |          order, not just specificity. Before migration this override lived
  227 |          in the same @layer components as default.css's rule and won on a
  228 |          same-layer specificity/order tie; that tie no longer exists post-
  229 |          migration, so the override must be unlayered (author styles beat
  230 |          every layered rule unconditionally) to still take effect. This
  231 |          test's intent — JS reads live computed width, not whichever rule
  232 |          authored it — is unaffected by the change. */
  233 |       [data-gsxui-slot-resizable-handle][aria-orientation="vertical"] {
  234 |         width: 6px;
  235 |       }
  236 |     `,
  237 |   });
  238 |   await expect.poll(() => handle.evaluate((el) => el.getBoundingClientRect().width)).toBe(6);
  239 |   expect(
  240 |     await handle.evaluate((el) => {
  241 |       const hitTarget = getComputedStyle(el, "::after");
  242 |       return { flex: getComputedStyle(el).flex, hitTargetWidth: hitTarget.width };
  243 |     }),
  244 |   ).toEqual({ flex: "0 0 auto", hitTargetWidth: "4px" });
  245 | 
  246 |   await handle.focus();
  247 |   await page.keyboard.press("ArrowRight");
  248 |   await expect(handle).toHaveAttribute("aria-valuenow", "35");
  249 |   await expect
  250 |     .poll(() => panels.first().evaluate((el: HTMLElement) => el.style.flexGrow))
  251 |     .toBe("35");
  252 | });
  253 | 
  254 | test("ScrollArea keeps native overflow while caller utilities win", async ({ page }) => {
  255 |   const response = await page.goto("/x/scroll-area/basic");
  256 |   expect(response?.status()).toBe(200);
  257 | 
  258 |   const area = page.locator('[data-gsxui-slot-scroll-area]');
  259 |   expect(
```