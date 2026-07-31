# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/dialog.spec.ts >> close controls, Escape, and the backdrop request their stable reasons
- Location: jstest/specs/dialog.spec.ts:188:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/dialog/basic", waiting until "load"

```

# Test source

```ts
  89  |     });
  90  |     document.addEventListener("gsxui:request-close", (event) => {
  91  |       if ((event.target as Element).matches("dialog[data-gsxui-command-dialog]")) {
  92  |         (window as any).__commandDialogEvents.push({
  93  |           type: event.type,
  94  |           detail: (event as CustomEvent).detail,
  95  |           target: (event.target as Element).matches("dialog[data-gsxui-command-dialog]"),
  96  |         });
  97  |       }
  98  |     });
  99  |   });
  100 | 
  101 |   const dialog = page.locator(COMMAND_DIALOG);
  102 |   await page.keyboard.press("Control+k");
  103 |   await expect(dialog).toHaveJSProperty("open", true);
  104 |   expect(await page.evaluate(() => (window as any).__commandDialogEvents)).toEqual([
  105 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  106 |   ]);
  107 | 
  108 |   await dialog.evaluate((element) => {
  109 |     element.animate([{ opacity: 1 }, { opacity: 0 }], { duration: 120 });
  110 |   });
  111 |   await page.keyboard.press("Control+k");
  112 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  113 |   await expect(dialog).toHaveJSProperty("open", true);
  114 |   await expect(dialog).toHaveJSProperty("open", false);
  115 |   expect(await page.evaluate(() => (window as any).__commandDialogEvents)).toEqual([
  116 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  117 |     { type: "gsxui:request-close", detail: { reason: "shortcut" }, target: true },
  118 |   ]);
  119 | 
  120 |   await page.keyboard.press("Control+k");
  121 |   await expect(dialog).toHaveJSProperty("open", true);
  122 |   await page.locator(`${COMMAND_DIALOG} [data-gsxui-command-item][data-href]`).first().evaluate((item) => {
  123 |     item.setAttribute("data-href", "#command-selected");
  124 |     (window as any).__commandNavigationOrder = [];
  125 |     document.addEventListener(
  126 |       "gsxui:request-close",
  127 |       (event) => {
  128 |         if ((event.target as Element).matches("dialog[data-gsxui-command-dialog]"))
  129 |           (window as any).__commandNavigationOrder.push("request-close");
  130 |       },
  131 |       { once: true },
  132 |     );
  133 |     window.addEventListener(
  134 |       "hashchange",
  135 |       () => (window as any).__commandNavigationOrder.push("navigation"),
  136 |       { once: true },
  137 |     );
  138 |   });
  139 |   await page.locator(`${COMMAND_DIALOG} [data-gsxui-command-item][data-href]`).first().click();
  140 |   await expect.poll(() => page.evaluate(() => (window as any).__commandNavigationOrder)).toEqual([
  141 |     "request-close",
  142 |     "navigation",
  143 |   ]);
  144 |   expect(await page.evaluate(() => (window as any).__commandDialogEvents)).toEqual([
  145 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  146 |     { type: "gsxui:request-close", detail: { reason: "shortcut" }, target: true },
  147 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  148 |     { type: "gsxui:request-close", detail: { reason: "select" }, target: true },
  149 |   ]);
  150 | });
  151 | 
  152 | test("a late document listener can cancel both deferred request defaults", async ({ page }) => {
  153 |   await page.goto(BASIC);
  154 |   await page.evaluate(() => {
  155 |     (window as any).__cancelledDialogRequests = [];
  156 |     document.addEventListener("gsxui:request-open", (event) => {
  157 |       event.preventDefault();
  158 |       (window as any).__cancelledDialogRequests.push(event.type);
  159 |     });
  160 |   });
  161 | 
  162 |   await page.getByRole("button", { name: "Delete account", exact: true }).click();
  163 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  164 | 
  165 |   await page.evaluate(() => {
  166 |     document.addEventListener("gsxui:request-open", (event) => event.stopImmediatePropagation(), {
  167 |       capture: true,
  168 |       once: true,
  169 |     });
  170 |   });
  171 |   await page.locator(DIALOG).evaluate((dialog: HTMLDialogElement) => dialog.showModal());
  172 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  173 | 
  174 |   await page.evaluate(() => {
  175 |     document.addEventListener("gsxui:request-close", (event) => {
  176 |       event.preventDefault();
  177 |       (window as any).__cancelledDialogRequests.push(event.type);
  178 |     });
  179 |   });
  180 |   await page.locator("[data-gsxui-dialog-close]").first().click();
  181 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  182 |   expect(await page.evaluate(() => (window as any).__cancelledDialogRequests)).toEqual([
  183 |     "gsxui:request-open",
  184 |     "gsxui:request-close",
  185 |   ]);
  186 | });
  187 | 
  188 | test("close controls, Escape, and the backdrop request their stable reasons", async ({ page }) => {
> 189 |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  190 |   await page.evaluate(() => {
  191 |     (window as any).__dialogCloseReasons = [];
  192 |     document.addEventListener("gsxui:request-close", (event) => {
  193 |       (window as any).__dialogCloseReasons.push((event as CustomEvent).detail.reason);
  194 |     });
  195 |   });
  196 | 
  197 |   await open(page);
  198 |   await page.locator("[data-gsxui-dialog-close]").first().click();
  199 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  200 | 
  201 |   await open(page);
  202 |   await page.keyboard.press("Escape");
  203 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  204 | 
  205 |   await open(page);
  206 |   await page.mouse.click(1, 1);
  207 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  208 | 
  209 |   expect(await page.evaluate(() => (window as any).__dialogCloseReasons)).toEqual([
  210 |     "close-button",
  211 |     "cancel",
  212 |     "backdrop",
  213 |   ]);
  214 | });
  215 | 
  216 | test("the finite exit stays open while closed and emits one notification after native close", async ({
  217 |   page,
  218 | }) => {
  219 |   await page.goto(BASIC);
  220 |   await page.evaluate(() => {
  221 |     (window as any).__dialogCloses = [];
  222 |     document.addEventListener("gsxui:close", (event) => {
  223 |       const dialog = event.target as HTMLDialogElement;
  224 |       (window as any).__dialogCloses.push({
  225 |         open: dialog.open,
  226 |         state: dialog.dataset.state,
  227 |       });
  228 |     });
  229 |   });
  230 | 
  231 |   await open(page);
  232 |   await dispatch(page, DIALOG, "gsxui:request-close");
  233 | 
  234 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  235 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  236 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  237 |   expect(await page.evaluate(() => (window as any).__dialogCloses)).toEqual([
  238 |     { open: false, state: "closed" },
  239 |   ]);
  240 | });
  241 | 
  242 | test("an open request during exit aborts the pending close", async ({ page }) => {
  243 |   await page.goto(BASIC);
  244 | 
  245 |   await open(page);
  246 |   await dispatch(page, DIALOG, "gsxui:request-close");
  247 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "closed");
  248 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  249 | 
  250 |   await dispatch(page, DIALOG, "gsxui:request-open");
  251 |   await expect(page.locator(DIALOG)).toHaveAttribute("data-state", "open");
  252 |   await page.waitForTimeout(350);
  253 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  254 | });
  255 | 
  256 | test("overlapping close generations wait for the latest finite dialog animation", async ({
  257 |   page,
  258 | }) => {
  259 |   await page.goto(BASIC);
  260 |   const dialog = page.locator(DIALOG);
  261 | 
  262 |   await open(page);
  263 |   await dialog.evaluate((element) => {
  264 |     const animation = element.animate([{ opacity: 1 }, { opacity: 0 }], {
  265 |       duration: 60_000,
  266 |     });
  267 |     animation.pause();
  268 |     animation.currentTime = 0;
  269 |     (window as any).__dialogAnimationA = animation;
  270 |   });
  271 |   await dispatch(page, DIALOG, "gsxui:request-close");
  272 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  273 |   await expect(dialog).toHaveJSProperty("open", true);
  274 | 
  275 |   await dispatch(page, DIALOG, "gsxui:request-open");
  276 |   await expect(dialog).toHaveAttribute("data-state", "open");
  277 |   await dialog.evaluate((element) => {
  278 |     (window as any).__dialogAnimationB = element.animate(
  279 |       [{ opacity: 1 }, { opacity: 0 }],
  280 |       { duration: 60_000 },
  281 |     );
  282 |   });
  283 |   await dispatch(page, DIALOG, "gsxui:request-close");
  284 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  285 |   await expect(dialog).toHaveJSProperty("open", true);
  286 | 
  287 |   await page.evaluate(async () => {
  288 |     const animation = (window as any).__dialogAnimationA as Animation;
  289 |     animation.finish();
```