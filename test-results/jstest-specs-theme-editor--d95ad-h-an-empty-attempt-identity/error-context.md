# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/theme-editor.spec.ts >> preview rejects state messages with an empty attempt identity
- Location: jstest/specs/theme-editor.spec.ts:679:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/theme", waiting until "load"

```

# Test source

```ts
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
  678 | 
  679 | test("preview rejects state messages with an empty attempt identity", async ({
  680 |   page,
  681 | }) => {
> 682 |   await page.goto("/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  683 |   const frame = page.locator("[data-theme-preview-frame]");
  684 |   const preview = page.frameLocator("[data-theme-preview-frame]");
  685 |   const status = page.locator("[data-theme-preview-status]");
  686 |   const retry = page.locator("[data-theme-preview-retry]");
  687 |   await expect(status).toHaveText("Live");
  688 | 
  689 |   await page.evaluate(() => {
  690 |     addEventListener("message", (event) => {
  691 |       if (event.data?.attempt === "") {
  692 |         (
  693 |           globalThis as typeof globalThis & {
  694 |             invalidAttemptResponse?: unknown;
  695 |           }
  696 |         ).invalidAttemptResponse = event.data;
  697 |       }
  698 |     });
  699 |   });
  700 |   await preview.locator("html").evaluate(() => {
  701 |     addEventListener("message", (event) => {
  702 |       if (event.data?.type === "gsxui:theme-preview:v1") {
  703 |         (
  704 |           globalThis as typeof globalThis & {
  705 |             capturedThemePreviewMessage?: unknown;
  706 |           }
  707 |         ).capturedThemePreviewMessage = structuredClone(event.data);
  708 |       }
  709 |     });
  710 |   });
  711 |   await page.locator('[data-theme-mode-tab="dark"]').click();
  712 |   await expect
  713 |     .poll(() =>
  714 |       preview.locator("html").evaluate(
  715 |         () =>
  716 |           (
  717 |             globalThis as typeof globalThis & {
  718 |               capturedThemePreviewMessage?: unknown;
  719 |             }
  720 |           ).capturedThemePreviewMessage,
  721 |       ),
  722 |     )
  723 |     .not.toBeUndefined();
  724 |   const payload = await preview.locator("html").evaluate(
  725 |     () =>
  726 |       (
  727 |         globalThis as typeof globalThis & {
  728 |           capturedThemePreviewMessage: unknown;
  729 |         }
  730 |       ).capturedThemePreviewMessage,
  731 |   );
  732 | 
  733 |   await frame.evaluate(
  734 |     (
  735 |       element: HTMLIFrameElement,
  736 |       message: Record<string, unknown>,
  737 |     ) => {
  738 |       element.contentWindow?.postMessage(
  739 |         { ...message, attempt: "", mode: "light" },
  740 |         location.origin,
  741 |       );
  742 |     },
  743 |     payload,
  744 |   );
  745 | 
  746 |   await expect
  747 |     .poll(() =>
  748 |       page.evaluate(
  749 |         () =>
  750 |           (
  751 |             globalThis as typeof globalThis & {
  752 |               invalidAttemptResponse?: {
  753 |                 type?: string;
  754 |                 message?: string;
  755 |               };
  756 |             }
  757 |           ).invalidAttemptResponse,
  758 |       ),
  759 |     )
  760 |     .toMatchObject({
  761 |       type: "gsxui:theme-preview-error:v1",
  762 |       message: "attempt must be a non-empty string",
  763 |     });
  764 |   await expect
  765 |     .poll(() =>
  766 |       preview
  767 |         .locator("html")
  768 |         .evaluate((element) => element.classList.contains("dark")),
  769 |     )
  770 |     .toBe(true);
  771 |   await expect(status).toHaveText("Live");
  772 |   await expect(retry).toBeHidden();
  773 | });
  774 | 
```