# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/style-visual.spec.ts >> Dialog keeps dedicated a11y hooks, semantic backdrop, and caller cascade
- Location: jstest/specs/style-visual.spec.ts:339:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/style-contract", waiting until "load"

```

# Test source

```ts
  242 |   );
  243 | 
  244 |   await page.evaluate(() => document.documentElement.classList.add("dark"));
  245 |   await page.evaluate(() => {
  246 |     for (const animation of document.getAnimations()) {
  247 |       animation.finish();
  248 |     }
  249 |   });
  250 |   expect(await colors(slot)).toEqual(
  251 |     await colors('[data-style-contract-reference="otp-invalid-dark"]'),
  252 |   );
  253 | });
  254 | 
  255 | test("joined ToggleGroup items override composed Toggle sizing and borders", async ({
  256 |   page,
  257 | }) => {
  258 |   const response = await page.goto("/f/style-contract");
  259 |   expect(response?.status(), "style contract fixture response").toBe(200);
  260 | 
  261 |   const metrics = async (selector: string) =>
  262 |     page.locator(selector).evaluate((element) => {
  263 |       const css = getComputedStyle(element);
  264 |       const box = element.getBoundingClientRect();
  265 |       return {
  266 |         x: box.x,
  267 |         width: box.width,
  268 |         height: css.height,
  269 |         paddingLeft: css.paddingLeft,
  270 |         paddingRight: css.paddingRight,
  271 |         borderLeftWidth: css.borderLeftWidth,
  272 |         borderTopLeftRadius: css.borderTopLeftRadius,
  273 |         borderTopRightRadius: css.borderTopRightRadius,
  274 |       };
  275 |     });
  276 | 
  277 |   const first = await metrics('[data-style-contract="toggle-group-sm-first"]');
  278 |   const iconItem = await metrics('[data-style-contract="toggle-group-sm-icon"]');
  279 |   const last = await metrics('[data-style-contract="toggle-group-sm-last"]');
  280 |   expect(first).toMatchObject({
  281 |     height: "28px",
  282 |     paddingLeft: "12px",
  283 |     paddingRight: "12px",
  284 |     borderLeftWidth: "1px",
  285 |     borderTopRightRadius: "0px",
  286 |   });
  287 |   expect(iconItem).toMatchObject({
  288 |     height: "28px",
  289 |     paddingLeft: "6px",
  290 |     paddingRight: "6px",
  291 |     borderLeftWidth: "0px",
  292 |     borderTopLeftRadius: "0px",
  293 |     borderTopRightRadius: "0px",
  294 |   });
  295 |   expect(last).toMatchObject({
  296 |     height: "28px",
  297 |     paddingLeft: "12px",
  298 |     paddingRight: "12px",
  299 |     borderLeftWidth: "0px",
  300 |     borderTopLeftRadius: "0px",
  301 |   });
  302 |   expect(iconItem.x).toBeCloseTo(first.x + first.width, 5);
  303 |   expect(last.x).toBeCloseTo(iconItem.x + iconItem.width, 5);
  304 | 
  305 |   expect(
  306 |     await metrics('[data-style-contract="toggle-group-default"]'),
  307 |   ).toMatchObject({
  308 |     height: "32px",
  309 |     paddingLeft: "12px",
  310 |     paddingRight: "12px",
  311 |   });
  312 |   expect(
  313 |     await metrics('[data-style-contract="toggle-group-large"]'),
  314 |   ).toMatchObject({
  315 |     height: "36px",
  316 |     paddingLeft: "12px",
  317 |     paddingRight: "12px",
  318 |   });
  319 |   expect(
  320 |     await metrics('[data-style-contract="toggle-group-caller"]'),
  321 |   ).toMatchObject({
  322 |     height: "28px",
  323 |     paddingLeft: "32px",
  324 |     paddingRight: "32px",
  325 |   });
  326 | });
  327 | 
  328 | test("an outline FieldGroup applies the separator overlap", async ({ page }) => {
  329 |   const response = await page.goto("/f/style-contract");
  330 |   expect(response?.status(), "style contract fixture response").toBe(200);
  331 | 
  332 |   expect(
  333 |     await page
  334 |       .locator('[data-style-contract="field-outline-separator"]')
  335 |       .evaluate((element) => getComputedStyle(element).marginBottom),
  336 |   ).toBe("-8px");
  337 | });
  338 | 
  339 | test("Dialog keeps dedicated a11y hooks, semantic backdrop, and caller cascade", async ({
  340 |   page,
  341 | }) => {
> 342 |   const response = await page.goto("/f/style-contract");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  343 |   expect(response?.status(), "style contract fixture response").toBe(200);
  344 | 
  345 |   await page.getByRole("button", { name: "Open contract dialog" }).click();
  346 |   const dialog = page.locator('[data-style-contract="dialog-caller"]');
  347 |   await expect(dialog).toBeVisible();
  348 |   await expect(dialog).toHaveAttribute("data-state", "open");
  349 | 
  350 |   expect(
  351 |     await dialog.evaluate((element) => {
  352 |       const title = element.querySelector("[data-gsxui-dialog-title]");
  353 |       const description = element.querySelector("[data-gsxui-dialog-description]");
  354 |       const css = getComputedStyle(element);
  355 |       return {
  356 |         labelledBy: element.getAttribute("aria-labelledby"),
  357 |         titleID: title?.id,
  358 |         describedBy: element.getAttribute("aria-describedby"),
  359 |         descriptionID: description?.id,
  360 |         borderRadius: css.borderRadius,
  361 |         display: css.display,
  362 |         backdrop: getComputedStyle(element, "::backdrop").backgroundColor,
  363 |       };
  364 |     }),
  365 |   ).toEqual({
  366 |     labelledBy: expect.any(String),
  367 |     titleID: expect.any(String),
  368 |     describedBy: expect.any(String),
  369 |     descriptionID: expect.any(String),
  370 |     borderRadius: "0px",
  371 |     display: "grid",
  372 |     backdrop: await page
  373 |       .locator('[data-style-contract-reference="overlay"]')
  374 |       .evaluate((element) => getComputedStyle(element).backgroundColor),
  375 |   });
  376 |   const relationships = await dialog.evaluate((element) => ({
  377 |     labelledBy: element.getAttribute("aria-labelledby"),
  378 |     titleID: element.querySelector("[data-gsxui-dialog-title]")?.id,
  379 |     describedBy: element.getAttribute("aria-describedby"),
  380 |     descriptionID: element.querySelector("[data-gsxui-dialog-description]")?.id,
  381 |   }));
  382 |   expect(relationships.labelledBy).toBe(relationships.titleID);
  383 |   expect(relationships.describedBy).toBe(relationships.descriptionID);
  384 | });
  385 | 
  386 | test("bottom Drawer uses Dialog mechanics and content-side header alignment", async ({
  387 |   page,
  388 | }) => {
  389 |   const response = await page.goto("/f/style-contract");
  390 |   expect(response?.status(), "style contract fixture response").toBe(200);
  391 | 
  392 |   await page.getByRole("button", { name: "Open contract drawer" }).click();
  393 |   const drawer = page.locator('[data-style-contract="drawer-bottom"]');
  394 |   await expect(drawer).toBeVisible();
  395 |   await expect(drawer).toHaveAttribute("data-state", "open");
  396 |   await page.evaluate(() => {
  397 |     for (const animation of document.getAnimations()) animation.finish();
  398 |   });
  399 | 
  400 |   expect(
  401 |     await drawer.evaluate((element) => {
  402 |       const rect = element.getBoundingClientRect();
  403 |       const header = element.querySelector('[data-style-contract="drawer-header"]');
  404 |       return {
  405 |         display: getComputedStyle(element).display,
  406 |         side: element.getAttribute("data-side"),
  407 |         left: rect.left,
  408 |         right: rect.right,
  409 |         bottom: rect.bottom,
  410 |         viewportWidth: window.innerWidth,
  411 |         viewportHeight: window.innerHeight,
  412 |         headerAlign: header ? getComputedStyle(header).textAlign : null,
  413 |       };
  414 |     }),
  415 |   ).toEqual({
  416 |     display: "flex",
  417 |     side: "bottom",
  418 |     left: 0,
  419 |     right: 1280,
  420 |     bottom: 900,
  421 |     viewportWidth: 1280,
  422 |     viewportHeight: 900,
  423 |     headerAlign: "center",
  424 |   });
  425 | });
  426 | 
  427 | test("Accordion caller padding overrides the inner content default", async ({
  428 |   page,
  429 | }) => {
  430 |   const response = await page.goto("/f/style-contract");
  431 |   expect(response?.status(), "style contract fixture response").toBe(200);
  432 | 
  433 |   const content = page.locator(
  434 |     '[data-style-contract="accordion-caller-content"]',
  435 |   );
  436 |   const inner = content.locator(
  437 |     ':scope > [data-gsxui-slot-accordion-content-inner]',
  438 |   );
  439 | 
  440 |   await expect(content).toHaveAttribute("id", "accordion-caller-content");
  441 |   await expect(content).not.toHaveClass(/\bpb-8\b/);
  442 |   await expect(inner).toHaveClass(/\bpb-8\b/);
```