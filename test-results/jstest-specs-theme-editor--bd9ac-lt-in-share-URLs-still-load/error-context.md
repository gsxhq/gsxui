# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> historical full built-in share URLs still load
- Location: jstest/specs/theme-editor.spec.ts:403:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme?preset=gsxui%3Av1%3AewogICIkc2NoZW1hIjogImh0dHBzOi8vdWkuZ3N4aHEuZGV2L3NjaGVtYXMvcHJlc2V0LXYxLmpzb24iLAogICJzY2hlbWFWZXJzaW9uIjogMSwKICAic3R5bGUiOiAibm92YSIsCiAgInJhZGl1cyI6ICIwLjYyNXJlbSIsCiAgInRoZW1lIjogewogICAgImxpZ2h0IjogewogICAgICAiYmFja2dyb3VuZCI6ICJva2xjaCgxIDAgMCkiLAogICAgICAiZm9yZWdyb3VuZCI6ICJva2xjaCgwLjE0NSAwIDApIiwKICAgICAgImNhcmQiOiAib2tsY2goMSAwIDApIiwKICAgICAgImNhcmQtZm9yZWdyb3VuZCI6ICJva2xjaCgwLjE0NSAwIDApIiwKICAgICAgInBvcG92ZXIiOiAib2tsY2goMSAwIDApIiwKICAgICAgInBvcG92ZXItZm9yZWdyb3VuZCI6ICJva2xjaCgwLjE0NSAwIDApIiwKICAgICAgInByaW1hcnkiOiAib2tsY2goMC4yMDUgMCAwKSIsCiAgICAgICJwcmltYXJ5LWZvcmVncm91bmQiOiAib2tsY2goMC45ODUgMCAwKSIsCiAgICAgICJzZWNvbmRhcnkiOiAib2tsY2goMC45NyAwIDApIiwKICAgICAgInNlY29uZGFyeS1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuMjA1IDAgMCkiLAogICAgICAibXV0ZWQiOiAib2tsY2goMC45NyAwIDApIiwKICAgICAgIm11dGVkLWZvcmVncm91bmQiOiAib2tsY2goMC41NTYgMCAwKSIsCiAgICAgICJhY2NlbnQiOiAib2tsY2goMC45NyAwIDApIiwKICAgICAgImFjY2VudC1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuMjA1IDAgMCkiLAogICAgICAiZGVzdHJ1Y3RpdmUiOiAib2tsY2goMC41NzcgMC4yNDUgMjcuMzI1KSIsCiAgICAgICJkZXN0cnVjdGl2ZS1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuOTcgMC4wMSAxNykiLAogICAgICAiYm9yZGVyIjogIm9rbGNoKDAuOTIyIDAgMCkiLAogICAgICAiaW5wdXQiOiAib2tsY2goMC45MjIgMCAwKSIsCiAgICAgICJyaW5nIjogIm9rbGNoKDAuNzA4IDAgMCkiLAogICAgICAic2lkZWJhciI6ICJva2xjaCgwLjk4NSAwIDApIiwKICAgICAgInNpZGViYXItZm9yZWdyb3VuZCI6ICJva2xjaCgwLjE0NSAwIDApIiwKICAgICAgInNpZGViYXItcHJpbWFyeSI6ICJva2xjaCgwLjIwNSAwIDApIiwKICAgICAgInNpZGViYXItcHJpbWFyeS1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuOTg1IDAgMCkiLAogICAgICAic2lkZWJhci1hY2NlbnQiOiAib2tsY2goMC45NyAwIDApIiwKICAgICAgInNpZGViYXItYWNjZW50LWZvcmVncm91bmQiOiAib2tsY2goMC4yMDUgMCAwKSIsCiAgICAgICJzaWRlYmFyLWJvcmRlciI6ICJva2xjaCgwLjkyMiAwIDApIiwKICAgICAgInNpZGViYXItcmluZyI6ICJva2xjaCgwLjcwOCAwIDApIiwKICAgICAgInN1Y2Nlc3MiOiAib2tsY2goNjkuNiUgMC4xNyAxNjIuNDgpIiwKICAgICAgImluZm8iOiAib2tsY2goNjguNSUgMC4xNjkgMjM3LjMyMykiLAogICAgICAid2FybmluZyI6ICJva2xjaCg3Ni45JSAwLjE4OCA3MC4wOCkiLAogICAgICAib3ZlcmxheSI6ICJva2xjaCgwJSAwIDAgLyAxMCUpIiwKICAgICAgImNvbnRyYXN0IjogIm9rbGNoKDEwMCUgMCAwKSIKICAgIH0sCiAgICAiZGFyayI6IHsKICAgICAgImJhY2tncm91bmQiOiAib2tsY2goMC4xNDUgMCAwKSIsCiAgICAgICJmb3JlZ3JvdW5kIjogIm9rbGNoKDAuOTg1IDAgMCkiLAogICAgICAiY2FyZCI6ICJva2xjaCgwLjIwNSAwIDApIiwKICAgICAgImNhcmQtZm9yZWdyb3VuZCI6ICJva2xjaCgwLjk4NSAwIDApIiwKICAgICAgInBvcG92ZXIiOiAib2tsY2goMC4yMDUgMCAwKSIsCiAgICAgICJwb3BvdmVyLWZvcmVncm91bmQiOiAib2tsY2goMC45ODUgMCAwKSIsCiAgICAgICJwcmltYXJ5IjogIm9rbGNoKDAuOTIyIDAgMCkiLAogICAgICAicHJpbWFyeS1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuMjA1IDAgMCkiLAogICAgICAic2Vjb25kYXJ5IjogIm9rbGNoKDAuMjY5IDAgMCkiLAogICAgICAic2Vjb25kYXJ5LWZvcmVncm91bmQiOiAib2tsY2goMC45ODUgMCAwKSIsCiAgICAgICJtdXRlZCI6ICJva2xjaCgwLjI2OSAwIDApIiwKICAgICAgIm11dGVkLWZvcmVncm91bmQiOiAib2tsY2goMC43MDggMCAwKSIsCiAgICAgICJhY2NlbnQiOiAib2tsY2goMC4yNjkgMCAwKSIsCiAgICAgICJhY2NlbnQtZm9yZWdyb3VuZCI6ICJva2xjaCgwLjk4NSAwIDApIiwKICAgICAgImRlc3RydWN0aXZlIjogIm9rbGNoKDAuNzA0IDAuMTkxIDIyLjIxNikiLAogICAgICAiZGVzdHJ1Y3RpdmUtZm9yZWdyb3VuZCI6ICJva2xjaCgwLjU4IDAuMjIgMjcpIiwKICAgICAgImJvcmRlciI6ICJva2xjaCgxIDAgMCAvIDEwJSkiLAogICAgICAiaW5wdXQiOiAib2tsY2goMSAwIDAgLyAxNSUpIiwKICAgICAgInJpbmciOiAib2tsY2goMC41NTYgMCAwKSIsCiAgICAgICJzaWRlYmFyIjogIm9rbGNoKDAuMjA1IDAgMCkiLAogICAgICAic2lkZWJhci1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuOTg1IDAgMCkiLAogICAgICAic2lkZWJhci1wcmltYXJ5IjogIm9rbGNoKDAuNDg4IDAuMjQzIDI2NC4zNzYpIiwKICAgICAgInNpZGViYXItcHJpbWFyeS1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuOTg1IDAgMCkiLAogICAgICAic2lkZWJhci1hY2NlbnQiOiAib2tsY2goMC4yNjkgMCAwKSIsCiAgICAgICJzaWRlYmFyLWFjY2VudC1mb3JlZ3JvdW5kIjogIm9rbGNoKDAuOTg1IDAgMCkiLAogICAgICAic2lkZWJhci1ib3JkZXIiOiAib2tsY2goMSAwIDAgLyAxMCUpIiwKICAgICAgInNpZGViYXItcmluZyI6ICJva2xjaCgwLjU1NiAwIDApIiwKICAgICAgInN1Y2Nlc3MiOiAib2tsY2goNjkuNiUgMC4xNyAxNjIuNDgpIiwKICAgICAgImluZm8iOiAib2tsY2goNjguNSUgMC4xNjkgMjM3LjMyMykiLAogICAgICAid2FybmluZyI6ICJva2xjaCg3Ni45JSAwLjE4OCA3MC4wOCkiLAogICAgICAib3ZlcmxheSI6ICJva2xjaCgwJSAwIDAgLyAxMCUpIiwKICAgICAgImNvbnRyYXN0IjogIm9rbGNoKDEwMCUgMCAwKSIKICAgIH0KICB9Cn0K", waiting until "load"

```

# Test source

```ts
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
  358 |   await page.goto("/theme");
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
> 409 |   await page.goto(`/theme?preset=${encodeURIComponent(fullCode)}`);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  459 | 
  460 |   await expect(selectionValue(page, "baseColor")).resolves.toBe("Custom");
  461 |   await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  462 |   await expect(selectionValue(page, "radius")).resolves.toBe("Custom");
  463 |   await expect.poll(() => iframeVariable(page, "primary")).toBe("rgb(10 20 30)");
  464 | 
  465 |   await choose(page, "radius", "Large");
  466 |   await expect(selectionValue(page, "baseColor")).resolves.toBe("Custom");
  467 |   await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  468 |   await expect(selectionValue(page, "radius")).resolves.toBe("Large");
  469 |   const exported = JSON.parse((await downloadText(page, "json")).text);
  470 |   expect(exported.theme.light.primary).toBe("rgb(10 20 30)");
  471 |   expect(exported.theme.dark.primary).toBe("rgb(40 50 60)");
  472 |   expect(exported.radius).toBe("0.875rem");
  473 | });
  474 | 
  475 | for (const rejection of [
  476 |   {
  477 |     name: "duplicate recognized declarations",
  478 |     css: ":root { --primary: red; --primary: blue; }",
  479 |     message: "duplicated",
  480 |   },
  481 |   {
  482 |     name: "malformed unrelated syntax",
  483 |     css: "body { color red; } :root { --primary: green; }",
  484 |     message: "malformed",
  485 |   },
  486 |   {
  487 |     name: "selector identifiers split across comments",
  488 |     css: ":r/**/oot { --primary: red; }",
  489 |     message: "must belong to :root or .dark",
  490 |   },
  491 |   {
  492 |     name: "important recognized declarations",
  493 |     css: ":root { --primary: red !important; }",
  494 |     message: "theme.light.primary",
  495 |   },
  496 | ]) {
  497 |   test(`CSS import rejects ${rejection.name} atomically`, async ({ page }) => {
  498 |     await page.goto("/theme");
  499 |     const before = await commands(page);
  500 |     await page.locator('[data-theme-import="css"]').fill(rejection.css);
  501 |     await page.locator('[data-theme-import-apply="css"]').click();
  502 | 
  503 |     await expect(page.locator("[data-theme-status]")).toContainText(rejection.message);
  504 |     expect(await commands(page)).toEqual(before);
  505 |   });
  506 | }
  507 | 
  508 | test("CSS import ignores valid unrelated CSS without mutating the preset", async ({
  509 |   page,
```