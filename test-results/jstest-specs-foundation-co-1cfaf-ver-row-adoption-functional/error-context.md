# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/foundation-contract.spec.ts >> foundation keeps Sonner fallback creation and server-row adoption functional
- Location: jstest/specs/foundation-contract.spec.ts:135:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sonner/types?css=foundation", waiting until "load"

```

# Test source

```ts
  38  |   await expect.poll(() => popover.evaluate((el) => getComputedStyle(el).display)).toBe("none");
  39  |   await page.locator("[data-gsxui-slot-popover-trigger]").first().click();
  40  |   await expect.poll(() => popover.evaluate((el) => el.matches(":popover-open"))).toBe(true);
  41  |   expect(
  42  |     await popover.evaluate((el) => {
  43  |       const rect = el.getBoundingClientRect();
  44  |       return rect.width > 0 && rect.height > 0;
  45  |     }),
  46  |   ).toBe(true);
  47  |   await page.keyboard.press("Escape");
  48  |   await expect.poll(() => popover.evaluate((el) => getComputedStyle(el).display)).toBe("none");
  49  | });
  50  | 
  51  | test("foundation keeps Carousel navigation geometrically functional", async ({ page }) => {
  52  |   await page.goto(foundation("/x/carousel/basic"));
  53  |   const root = page.locator("[data-gsxui-carousel]");
  54  |   const viewport = page.locator("[data-gsxui-carousel-content]");
  55  |   const items = page.locator("[data-gsxui-carousel-item]");
  56  |   await expect(items).toHaveCount(5);
  57  |   const geometry = await viewport.evaluate((el) => {
  58  |     const first = el.querySelector<HTMLElement>("[data-gsxui-carousel-item]")!;
  59  |     return {
  60  |       viewportWidth: el.getBoundingClientRect().width,
  61  |       itemWidth: first.getBoundingClientRect().width,
  62  |       scrollWidth: el.scrollWidth,
  63  |     };
  64  |   });
  65  |   expect(geometry.viewportWidth).toBeGreaterThan(0);
  66  |   expect(geometry.itemWidth).toBeCloseTo(geometry.viewportWidth, 0);
  67  |   expect(geometry.scrollWidth).toBeGreaterThan(geometry.viewportWidth * 4);
  68  | 
  69  |   await expect(root).toHaveAttribute("data-current-index", "0");
  70  |   await page.locator("[data-gsxui-carousel-next]").click();
  71  |   await expect.poll(() => viewport.evaluate((el) => el.scrollLeft)).toBeGreaterThan(0);
  72  |   await expect(root).toHaveAttribute("data-current-index", "1");
  73  | });
  74  | 
  75  | for (const width of [640, 900]) {
  76  |   test(`foundation keeps Resizable keyboard and pointer geometry at ${width}px`, async ({
  77  |     page,
  78  |   }) => {
  79  |     await page.setViewportSize({ width, height: 700 });
  80  |     await page.goto(foundation("/x/resizable/handle"));
  81  |     const group = page.locator("[data-gsxui-resizable]").first();
  82  |     const panels = group.locator(":scope > [data-gsxui-resizable-panel]");
  83  |     const handle = group.locator(":scope > [data-gsxui-resizable-handle]");
  84  |     await expect(panels).toHaveCount(2);
  85  | 
  86  |     const before = await panels.first().evaluate((el) => el.getBoundingClientRect().width);
  87  |     expect(before).toBeGreaterThan(0);
  88  |     expect(
  89  |       await handle.evaluate((el) => {
  90  |         const rect = el.getBoundingClientRect();
  91  |         const hit = getComputedStyle(el, "::after");
  92  |         return rect.height > 0 && Number.parseFloat(hit.width) >= 4;
  93  |       }),
  94  |     ).toBe(true);
  95  | 
  96  |     await handle.focus();
  97  |     await page.keyboard.press("ArrowRight");
  98  |     await expect(handle).toHaveAttribute("aria-valuenow", "35");
  99  |     const keyboardWidth = await panels.first().evaluate(
  100 |       (el) => el.getBoundingClientRect().width,
  101 |     );
  102 |     expect(keyboardWidth).toBeGreaterThan(before);
  103 | 
  104 |     const box = await handle.boundingBox();
  105 |     if (!box) throw new Error("Resizable handle has no pointer geometry");
  106 |     await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2);
  107 |     await page.mouse.down();
  108 |     await page.mouse.move(box.x + box.width / 2 + 30, box.y + box.height / 2);
  109 |     await page.mouse.up();
  110 |     await expect
  111 |       .poll(() => panels.first().evaluate((el) => el.getBoundingClientRect().width))
  112 |       .toBeGreaterThan(keyboardWidth);
  113 |   });
  114 | }
  115 | 
  116 | test("foundation keeps InputOTP entry and client caret functional", async ({ page }) => {
  117 |   await page.goto(foundation("/x/input-otp/basic"));
  118 |   const input = page.locator("[data-gsxui-input-otp-input]");
  119 |   const slots = page.locator("[data-gsxui-input-otp-slot]");
  120 |   await input.fill("123");
  121 |   await expect(slots.nth(0)).toHaveText("1");
  122 |   await expect(slots.nth(2)).toHaveText("3");
  123 |   await expect(slots.nth(3)).toHaveAttribute("data-active", "true");
  124 |   await expect(slots.nth(3).locator("[data-gsxui-input-otp-caret-overlay]")).toHaveCount(1);
  125 |   await expect(slots.nth(3).locator("[data-gsxui-input-otp-caret]")).toHaveCount(1);
  126 |   expect(
  127 |     await input.evaluate((el) => {
  128 |       const inputRect = el.getBoundingClientRect();
  129 |       const rootRect = el.parentElement!.getBoundingClientRect();
  130 |       return inputRect.width === rootRect.width && inputRect.height === rootRect.height;
  131 |     }),
  132 |   ).toBe(true);
  133 | });
  134 | 
  135 | test("foundation keeps Sonner fallback creation and server-row adoption functional", async ({
  136 |   page,
  137 | }) => {
> 138 |   await page.goto(foundation("/x/sonner/types"));
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  139 |   await page.evaluate(() => {
  140 |     const section = document.querySelector('[aria-label="Notifications"]')!;
  141 |     for (const template of [...section.querySelectorAll("template")]) {
  142 |       document.body.append(template);
  143 |     }
  144 |     section.remove();
  145 |     window.gsxui.toast.success("Fallback", { duration: 60_000 });
  146 |     const template = document.querySelector<HTMLTemplateElement>(
  147 |       'template[data-gsxui-toast-template="info"]',
  148 |     )!;
  149 |     const row = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
  150 |     row.dataset.foundationServer = "true";
  151 |     row.dataset.duration = "60000";
  152 |     document.querySelector("[data-gsxui-toaster]")!.append(row);
  153 |   });
  154 | 
  155 |   const region = page.locator("[data-gsxui-toaster]");
  156 |   await expect(region).toHaveCount(1);
  157 |   await expect(region.locator("li[data-gsxui-toast]")).toHaveCount(2);
  158 |   await expect(region.locator('[data-foundation-server="true"]')).toHaveAttribute(
  159 |     "data-state",
  160 |     "open",
  161 |   );
  162 |   expect(
  163 |     await region.evaluate((el) => {
  164 |       const css = getComputedStyle(el);
  165 |       return {
  166 |         position: css.position,
  167 |         pointerEvents: css.pointerEvents,
  168 |         inViewport:
  169 |           el.getBoundingClientRect().right <= innerWidth &&
  170 |           el.getBoundingClientRect().bottom <= innerHeight,
  171 |       };
  172 |     }),
  173 |   ).toEqual({ position: "fixed", pointerEvents: "none", inViewport: true });
  174 | });
  175 | 
  176 | test("full style leaves caller utilities later in the cascade", async ({ page }) => {
  177 |   await page.goto("/f/style-contract");
  178 |   expect(
  179 |     await page.evaluate(() =>
  180 |       getComputedStyle(document.documentElement)
  181 |         .getPropertyValue("--gsxui-default-style")
  182 |         .trim(),
  183 |     ),
  184 |   ).toBe("default");
  185 |   const button = page.getByRole("button", { name: "Caller override" });
  186 |   expect(
  187 |     await button.evaluate((el) => {
  188 |       const css = getComputedStyle(el);
  189 |       return { height: css.height, radius: css.borderRadius };
  190 |     }),
  191 |   ).toEqual({ height: "48px", radius: "0px" });
  192 | });
  193 | 
```