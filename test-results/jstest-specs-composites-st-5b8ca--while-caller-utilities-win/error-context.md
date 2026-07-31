# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/composites-style-contract.spec.ts >> ScrollArea keeps native overflow while caller utilities win
- Location: jstest/specs/composites-style-contract.spec.ts:254:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/scroll-area/basic", waiting until "load"

```

# Test source

```ts
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
> 255 |   const response = await page.goto("/x/scroll-area/basic");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  256 |   expect(response?.status()).toBe(200);
  257 | 
  258 |   const area = page.locator('[data-gsxui-slot-scroll-area]');
  259 |   expect(
  260 |     await area.evaluate((el) => {
  261 |       const css = getComputedStyle(el);
  262 |       return {
  263 |         overflowX: css.overflowX,
  264 |         overflowY: css.overflowY,
  265 |         radius: css.borderRadius,
  266 |       };
  267 |     }),
  268 |   ).toEqual({
  269 |     overflowX: "auto",
  270 |     overflowY: "auto",
  271 |     radius: "8px",
  272 |   });
  273 | });
  274 | 
```