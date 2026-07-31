# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> preview acknowledgement survives a late parent-observed iframe load
- Location: jstest/specs/theme-editor.spec.ts:574:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme", waiting until "load"

```

# Test source

```ts
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
> 577 |   await page.goto("/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
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
  599 |       return;
  600 |     }
  601 |     await route.fulfill({
  602 |       contentType: "text/html",
  603 |       body: "<!doctype html><body data-unresponsive-preview>Unresponsive preview</body>",
  604 |     });
  605 |   });
  606 | 
  607 |   await page.goto("/theme");
  608 |   const frame = page.locator("[data-theme-preview-frame]");
  609 |   const preview = page.frameLocator("[data-theme-preview-frame]");
  610 |   const status = page.locator("[data-theme-preview-status]");
  611 |   const retry = page.locator("[data-theme-preview-retry]");
  612 |   await expect(status).toHaveText("Live");
  613 | 
  614 |   await preview.locator("html").evaluate(() => {
  615 |     addEventListener("message", (event) => {
  616 |       if (event.data?.type === "gsxui:theme-preview:v1") {
  617 |         (
  618 |           globalThis as typeof globalThis & {
  619 |             capturedThemePreviewAttempt?: string;
  620 |           }
  621 |         ).capturedThemePreviewAttempt = event.data.attempt;
  622 |       }
  623 |     });
  624 |   });
  625 |   await page.locator('[data-theme-mode-tab="dark"]').click();
  626 |   await expect
  627 |     .poll(() =>
  628 |       preview.locator("html").evaluate(
  629 |         () =>
  630 |           (
  631 |             globalThis as typeof globalThis & {
  632 |               capturedThemePreviewAttempt?: string;
  633 |             }
  634 |           ).capturedThemePreviewAttempt,
  635 |       ),
  636 |     )
  637 |     .toBeTruthy();
  638 |   const staleAttempt = await preview.locator("html").evaluate(
  639 |     () =>
  640 |       (
  641 |         globalThis as typeof globalThis & {
  642 |           capturedThemePreviewAttempt: string;
  643 |         }
  644 |       ).capturedThemePreviewAttempt,
  645 |   );
  646 |   await expect(status).toHaveText("Live");
  647 |   await expect(retry).toBeHidden();
  648 | 
  649 |   responsive = false;
  650 |   await frame.evaluate((element: HTMLIFrameElement) => {
  651 |     element.src = element.src;
  652 |   });
  653 |   await expect(preview.locator("[data-unresponsive-preview]")).toBeVisible();
  654 | 
  655 |   await preview.locator("html").evaluate((_, attempt) => {
  656 |     parent.postMessage(
  657 |       { type: "gsxui:theme-preview-applied:v1", attempt },
  658 |       location.origin,
  659 |     );
  660 |     parent.postMessage(
  661 |       {
  662 |         type: "gsxui:theme-preview-error:v1",
  663 |         attempt,
  664 |         message: "stale preview error",
  665 |       },
  666 |       location.origin,
  667 |     );
  668 |   }, staleAttempt);
  669 | 
  670 |   await expect(status).toHaveText("Preview did not respond.");
  671 |   await expect(retry).toBeVisible();
  672 | 
  673 |   responsive = true;
  674 |   await retry.click();
  675 |   await expect(status).toHaveText("Live");
  676 |   await expect(retry).toBeHidden();
  677 | });
```