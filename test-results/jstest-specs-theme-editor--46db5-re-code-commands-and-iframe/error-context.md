# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> click commits palette state to JSON, share code, commands, and iframe
- Location: jstest/specs/theme-editor.spec.ts:355:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme", waiting until "load"

```

# Test source

```ts
  258 |   page,
  259 | }) => {
  260 |   await page.goto("/theme");
  261 |   const catalog = await schema(page);
  262 |   const committedPrimary = await iframeVariable(page, "primary");
  263 |   const committedCommands = await commands(page);
  264 |   const committedJSON = (await downloadText(page, "json")).text;
  265 | 
  266 |   const theme = picker(page, "theme");
  267 |   await theme.locator("[data-theme-picker-trigger]").click();
  268 |   const rose = theme.getByRole("radio", { name: "Rose", exact: true }).locator("..");
  269 |   await rose.locator("[data-theme-choice-swatch]").hover();
  270 | 
  271 |   await expect
  272 |     .poll(() => iframeVariable(page, "primary"))
  273 |     .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  274 |   const roseBox = await rose.boundingBox();
  275 |   if (!roseBox) throw new Error("Rose choice has no layout box");
  276 |   await page.mouse.move(
  277 |     roseBox.x + roseBox.width - 2,
  278 |     roseBox.y + roseBox.height / 2,
  279 |   );
  280 |   await expect
  281 |     .poll(() => iframeVariable(page, "primary"))
  282 |     .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  283 |   await expect(
  284 |     theme.getByRole("radio", { name: "Neutral", exact: true }),
  285 |   ).toBeFocused();
  286 |   expect(await commands(page)).toEqual(committedCommands);
  287 |   expect((await downloadText(page, "json", false)).text).toBe(committedJSON);
  288 | 
  289 |   await page.keyboard.press("Escape");
  290 |   await expect
  291 |     .poll(() =>
  292 |       theme
  293 |         .locator("[data-gsxui-popover-content]")
  294 |         .evaluate((element) => element.matches(":popover-open")),
  295 |     )
  296 |     .toBe(false);
  297 |   await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
  298 | 
  299 |   await theme.locator("[data-theme-picker-trigger]").click();
  300 |   await rose.hover();
  301 |   await expect
  302 |     .poll(() => iframeVariable(page, "primary"))
  303 |     .toBe(catalog.palette.resolved.neutral.rose.light.primary);
  304 |   await page.mouse.move(1, 1);
  305 |   await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
  306 |   expect(await commands(page)).toEqual(committedCommands);
  307 | });
  308 | 
  309 | test("touch pointerenter does not preview on a fine hover-capable device", async ({ page }) => {
  310 |   await page.goto("/theme");
  311 |   const committedPrimary = await iframeVariable(page, "primary");
  312 |   const theme = picker(page, "theme");
  313 |   await theme.locator("[data-theme-picker-trigger]").click();
  314 |   const rose = theme.getByRole("radio", { name: "Rose", exact: true }).locator("..");
  315 | 
  316 |   expect(
  317 |     await page.evaluate(() =>
  318 |       matchMedia("(hover: hover) and (pointer: fine)").matches,
  319 |     ),
  320 |   ).toBe(true);
  321 |   const pointerType = await rose.evaluate((element) => {
  322 |     const event = new PointerEvent("pointerenter", {
  323 |       bubbles: false,
  324 |       composed: true,
  325 |       isPrimary: true,
  326 |       pointerId: 17,
  327 |       pointerType: "touch",
  328 |     });
  329 |     element.dispatchEvent(event);
  330 |     return event.pointerType;
  331 |   });
  332 |   expect(pointerType).toBe("touch");
  333 | 
  334 |   await expect.poll(() => iframeVariable(page, "primary")).toBe(committedPrimary);
  335 | });
  336 | 
  337 | test("radius hover previews only the iframe and restores on pointer leave", async ({ page }) => {
  338 |   await page.goto("/theme");
  339 |   const committedRadius = await iframeVariable(page, "radius");
  340 |   const committedCommands = await commands(page);
  341 |   const committedJSON = (await downloadText(page, "json")).text;
  342 |   const radius = picker(page, "radius");
  343 |   await radius.locator("[data-theme-picker-trigger]").click();
  344 | 
  345 |   await radius.getByRole("radio", { name: "Large", exact: true }).locator("..").hover();
  346 |   await expect.poll(() => iframeVariable(page, "radius")).toBe("0.875rem");
  347 |   expect(await commands(page)).toEqual(committedCommands);
  348 |   expect((await downloadText(page, "json", false)).text).toBe(committedJSON);
  349 | 
  350 |   await page.mouse.move(1, 1);
  351 |   await expect.poll(() => iframeVariable(page, "radius")).toBe(committedRadius);
  352 |   expect(await commands(page)).toEqual(committedCommands);
  353 | });
  354 | 
  355 | test("click commits palette state to JSON, share code, commands, and iframe", async ({
  356 |   page,
  357 | }) => {
> 358 |   await page.goto("/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  359 |   const catalog = await schema(page);
  360 |   const beforeJSON = (await downloadText(page, "json")).text;
  361 |   const beforeCommands = await commands(page);
  362 |   const beforeShare = shareFromInit(beforeCommands.init);
  363 | 
  364 |   await choose(page, "theme", "Blue");
  365 | 
  366 |   const afterJSON = (await downloadText(page, "json")).text;
  367 |   const afterCommands = await commands(page);
  368 |   expect(afterJSON).not.toBe(beforeJSON);
  369 |   expect(JSON.parse(afterJSON).theme.light.primary).toBe(
  370 |     catalog.palette.resolved.neutral.blue.light.primary,
  371 |   );
  372 |   expect(shareFromInit(afterCommands.init)).not.toBe(beforeShare);
  373 |   expect(afterCommands.init).not.toBe(beforeCommands.init);
  374 |   expect(afterCommands.apply).not.toBe(beforeCommands.apply);
  375 |   await expect
  376 |     .poll(() => iframeVariable(page, "primary"))
  377 |     .toBe(catalog.palette.resolved.neutral.blue.light.primary);
  378 | });
  379 | 
  380 | test("share commands use compact built-ins and full custom imports", async ({ page }) => {
  381 |   await page.goto("/theme");
  382 | 
  383 |   const initialShare = shareFromInit((await commands(page)).init);
  384 |   expect(initialShare).toMatch(/^gsxui:p1:/);
  385 |   expect(initialShare.length).toBeLessThanOrEqual(12);
  386 | 
  387 |   await choose(page, "theme", "Blue");
  388 |   const selectedShare = shareFromInit((await commands(page)).init);
  389 |   expect(selectedShare).toMatch(/^gsxui:p1:/);
  390 |   expect(selectedShare.length).toBeLessThanOrEqual(12);
  391 |   expect(selectedShare).not.toBe(initialShare);
  392 | 
  393 |   const imported = JSON.parse((await downloadText(page, "json")).text);
  394 |   imported.theme.light.primary = "rgb(1 2 3)";
  395 |   await page.locator('[data-theme-import="json"]').fill(`${JSON.stringify(imported)}\n`);
  396 |   await page.locator('[data-theme-import-apply="json"]').click();
  397 | 
  398 |   const customShare = shareFromInit((await commands(page)).init);
  399 |   expect(customShare).toMatch(/^gsxui:v1:/);
  400 |   expect(JSON.parse((await downloadText(page, "json")).text)).toEqual(imported);
  401 | });
  402 | 
  403 | test("historical full built-in share URLs still load", async ({ page }) => {
  404 |   const source = readFileSync(
  405 |     `${repoRoot}/internal/preset/testdata/default-nova.json`,
  406 |     "utf8",
  407 |   );
  408 |   const fullCode = `gsxui:v1:${Buffer.from(source).toString("base64url")}`;
  409 |   await page.goto(`/theme?preset=${encodeURIComponent(fullCode)}`);
  410 | 
  411 |   await expect(page.locator("[data-theme-status]")).toHaveText("Loaded shared preset.");
  412 |   expect(JSON.parse((await downloadText(page, "json")).text)).toEqual(JSON.parse(source));
  413 |   expect(shareFromInit((await commands(page)).init)).toMatch(/^gsxui:p1:/);
  414 | });
  415 | 
  416 | test("custom JSON import is lossless and a built-in base replaces it atomically", async ({
  417 |   page,
  418 | }) => {
  419 |   await page.goto("/theme");
  420 |   const catalog = await schema(page);
  421 |   const imported = JSON.parse((await downloadText(page, "json")).text);
  422 |   imported.radius = "1rem";
  423 |   imported.theme.light.primary = "rgb(1 2 3)";
  424 |   imported.theme.dark.primary = "rgb(4 5 6)";
  425 | 
  426 |   await page.locator('[data-theme-import="json"]').fill(`${JSON.stringify(imported)}\n`);
  427 |   await page.locator('[data-theme-import-apply="json"]').click();
  428 | 
  429 |   await expect(selectionValue(page, "baseColor")).resolves.toBe("Custom");
  430 |   await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  431 |   await expect(selectionValue(page, "radius")).resolves.toBe("Custom");
  432 |   await expect.poll(() => iframeVariable(page, "primary")).toBe("rgb(1 2 3)");
  433 |   expect(JSON.parse((await downloadText(page, "json")).text)).toEqual(imported);
  434 | 
  435 |   await choose(page, "baseColor", "Stone");
  436 |   await expect(selectionValue(page, "baseColor")).resolves.toBe("Stone");
  437 |   await expect(selectionValue(page, "theme")).resolves.toBe("Stone");
  438 |   await expect(selectionValue(page, "radius")).resolves.toBe("Custom");
  439 |   await expect
  440 |     .poll(() => iframeVariable(page, "primary"))
  441 |     .toBe(catalog.palette.resolved.stone.stone.light.primary);
  442 |   const replaced = JSON.parse((await downloadText(page, "json")).text);
  443 |   expect(replaced.theme.light).toEqual(catalog.palette.resolved.stone.stone.light);
  444 |   expect(replaced.theme.dark).toEqual(catalog.palette.resolved.stone.stone.dark);
  445 |   expect(replaced.radius).toBe("1rem");
  446 | });
  447 | 
  448 | test("custom CSS import is lossless while a named radius preserves its colors", async ({
  449 |   page,
  450 | }) => {
  451 |   await page.goto("/theme");
  452 |   await page
  453 |     .locator('[data-theme-import="css"]')
  454 |     .fill(
  455 |       ":root { --radius: 1rem; --primary: rgb(10 20 30); }\n" +
  456 |         ".dark { --primary: rgb(40 50 60); }",
  457 |     );
  458 |   await page.locator('[data-theme-import-apply="css"]').click();
```