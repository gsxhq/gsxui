# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> CSS import rejects selector identifiers split across comments atomically
- Location: jstest/specs/theme-editor.spec.ts:497:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme", waiting until "load"

```

# Test source

```ts
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
> 498 |     await page.goto("/theme");
      |                ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  510 | }) => {
  511 |   await page.goto("/theme");
  512 |   const before = await commands(page);
  513 |   await page
  514 |     .locator('[data-theme-import="css"]')
  515 |     .fill("body { color: red; --unrelated: blue; }");
  516 |   await page.locator('[data-theme-import-apply="css"]').click();
  517 | 
  518 |   expect(await commands(page)).toEqual(before);
  519 | });
  520 | 
  521 | test("CSS import accepts comments between supported selector tokens", async ({
  522 |   page,
  523 | }) => {
  524 |   await page.goto("/theme");
  525 |   await page
  526 |     .locator('[data-theme-import="css"]')
  527 |     .fill(":/* comment */root { --primary: red; }");
  528 |   await page.locator('[data-theme-import-apply="css"]').click();
  529 | 
  530 |   await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  531 |   await expect.poll(() => iframeVariable(page, "primary")).toBe("red");
  532 | });
  533 | 
  534 | test("CSS import ignores valid unrelated implicit nesting", async ({ page }) => {
  535 |   await page.goto("/theme");
  536 |   await page
  537 |     .locator('[data-theme-import="css"]')
  538 |     .fill(
  539 |       "body { .nested { color: red; --unrelated: blue; } } :root { --primary: green; }",
  540 |     );
  541 |   await page.locator('[data-theme-import-apply="css"]').click();
  542 | 
  543 |   await expect(selectionValue(page, "theme")).resolves.toBe("Custom");
  544 |   await expect.poll(() => iframeVariable(page, "primary")).toBe("green");
  545 | });
  546 | 
  547 | test("theme editor exposes Retry when the preview never handshakes", async ({
  548 |   page,
  549 | }) => {
  550 |   let responsive = false;
  551 |   await page.route("**/theme/preview/button", async (route) => {
  552 |     if (responsive) {
  553 |       await route.fallback();
  554 |       return;
  555 |     }
  556 |     await route.fulfill({
  557 |       contentType: "text/html",
  558 |       body: "<!doctype html><title>Unresponsive preview</title>",
  559 |     });
  560 |   });
  561 |   await page.goto("/theme");
  562 | 
  563 |   await expect(page.locator("[data-theme-preview-status]")).toHaveText(
  564 |     "Preview did not respond.",
  565 |   );
  566 |   await expect(page.locator("[data-theme-preview-retry]")).toBeVisible();
  567 | 
  568 |   responsive = true;
  569 |   await page.locator("[data-theme-preview-retry]").click();
  570 |   await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  571 |   await expect(page.locator("[data-theme-preview-retry]")).toBeHidden();
  572 | });
  573 | 
  574 | test("preview acknowledgement survives a late parent-observed iframe load", async ({
  575 |   page,
  576 | }) => {
  577 |   await page.goto("/theme");
  578 |   const status = page.locator("[data-theme-preview-status]");
  579 |   const retry = page.locator("[data-theme-preview-retry]");
  580 |   await expect(status).toHaveText("Live");
  581 |   await expect(retry).toBeHidden();
  582 | 
  583 |   await page
  584 |     .locator("[data-theme-preview-frame]")
  585 |     .evaluate((frame) => frame.dispatchEvent(new Event("load")));
  586 | 
  587 |   await page.waitForTimeout(2_100);
  588 |   await expect(status).toHaveText("Live");
  589 |   await expect(retry).toBeHidden();
  590 | });
  591 | 
  592 | test("stale previous-document responses cannot complete a fresh preview attempt", async ({
  593 |   page,
  594 | }) => {
  595 |   let responsive = true;
  596 |   await page.route("**/theme/preview/button", async (route) => {
  597 |     if (responsive) {
  598 |       await route.fallback();
```