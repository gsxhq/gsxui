# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/style-visual.spec.ts >> Tooltip Kbd relationship and Accordion native open mechanics compute
- Location: jstest/specs/style-visual.spec.ts:632:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/style-contract", waiting until "load"

```

# Test source

```ts
  535 |         name: `Open ${overlay.family} ${overlay.fixture}`,
  536 |       })
  537 |       .click();
  538 |     await expect(dialog).toBeVisible();
  539 |     await expect(dialog).toHaveAttribute("open", "");
  540 |     await expect(dialog).toHaveAttribute("data-state", "open");
  541 |     await dialog.evaluate((element) => {
  542 |       for (const animation of element.getAnimations()) animation.finish();
  543 |     });
  544 | 
  545 |     const geometry = await dialog.evaluate(
  546 |       (element, enterProperty) => {
  547 |         const rect = element.getBoundingClientRect();
  548 |         const css = getComputedStyle(element);
  549 |         return {
  550 |           side: element.getAttribute("data-side"),
  551 |           display: css.display,
  552 |           animationName: css.animationName,
  553 |           enterValue: css.getPropertyValue(enterProperty).trim(),
  554 |           left: rect.left,
  555 |           right: rect.right,
  556 |           top: rect.top,
  557 |           bottom: rect.bottom,
  558 |           width: rect.width,
  559 |           height: rect.height,
  560 |           viewportWidth: window.innerWidth,
  561 |           viewportHeight: window.innerHeight,
  562 |         };
  563 |       },
  564 |       overlay.enterProperty,
  565 |     );
  566 |     expect(geometry.side).toBe(overlay.side);
  567 |     expect(geometry.display).toBe("flex");
  568 |     expect(geometry.animationName).toContain("enter");
  569 |     expect(geometry.enterValue).toBe(overlay.enterValue);
  570 | 
  571 |     const subpixelEdgeTolerance = 0.5;
  572 |     const expectAtEdge = (actual: number, expected: number) =>
  573 |       expect(Math.abs(actual - expected)).toBeLessThanOrEqual(
  574 |         subpixelEdgeTolerance,
  575 |       );
  576 |     if (overlay.side === "top" || overlay.side === "bottom") {
  577 |       expectAtEdge(geometry.left, 0);
  578 |       expectAtEdge(geometry.right, geometry.viewportWidth);
  579 |       expectAtEdge(geometry.width, geometry.viewportWidth);
  580 |       expect(geometry.height).toBeGreaterThan(0);
  581 |       expect(geometry.height).toBeLessThan(geometry.viewportHeight);
  582 |       if (overlay.side === "top") {
  583 |         expectAtEdge(geometry.top, 0);
  584 |       } else {
  585 |         expectAtEdge(geometry.bottom, geometry.viewportHeight);
  586 |       }
  587 |     } else {
  588 |       expectAtEdge(geometry.top, 0);
  589 |       expectAtEdge(geometry.bottom, geometry.viewportHeight);
  590 |       expectAtEdge(geometry.height, geometry.viewportHeight);
  591 |       expectAtEdge(geometry.width, 384);
  592 |       if (overlay.side === "left") {
  593 |         expectAtEdge(geometry.left, 0);
  594 |       } else {
  595 |         expectAtEdge(geometry.right, geometry.viewportWidth);
  596 |       }
  597 |     }
  598 | 
  599 |     await page
  600 |       .getByRole("button", {
  601 |         name: `Close ${overlay.family} ${overlay.fixture}`,
  602 |       })
  603 |       .click();
  604 |     await expect(dialog).toHaveAttribute("data-state", "closed");
  605 |     const closing = await dialog.evaluate(
  606 |       (element: HTMLDialogElement, exitProperty) => {
  607 |         const css = getComputedStyle(element);
  608 |         return {
  609 |           open: element.open,
  610 |           animationName: css.animationName,
  611 |           exitValue: css.getPropertyValue(exitProperty).trim(),
  612 |         };
  613 |       },
  614 |       overlay.enterProperty.replace("--tw-enter-", "--tw-exit-"),
  615 |     );
  616 |     expect(closing.open).toBe(true);
  617 |     expect(closing.animationName).toContain("exit");
  618 |     expect(closing.exitValue).toBe(overlay.enterValue);
  619 |     await expect(dialog).toBeHidden();
  620 |     expect(
  621 |       await dialog.evaluate((element: HTMLDialogElement) => ({
  622 |         open: element.open,
  623 |         display: getComputedStyle(element).display,
  624 |       })),
  625 |     ).toEqual({
  626 |       open: false,
  627 |       display: "none",
  628 |     });
  629 |   });
  630 | }
  631 | 
  632 | test("Tooltip Kbd relationship and Accordion native open mechanics compute", async ({
  633 |   page,
  634 | }) => {
> 635 |   const response = await page.goto("/f/style-contract");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  636 |   expect(response?.status(), "style contract fixture response").toBe(200);
  637 | 
  638 |   await page.getByRole("button", { name: "Show contract tooltip" }).focus();
  639 |   const tooltip = page.locator('[data-style-contract="tooltip-kbd"]');
  640 |   await expect(tooltip).toBeVisible();
  641 |   expect(
  642 |     await tooltip.evaluate((element) => getComputedStyle(element).paddingRight),
  643 |   ).toBe("6px");
  644 | 
  645 |   const trigger = page.getByText("Contract accordion", { exact: true });
  646 |   await trigger.click();
  647 |   const details = trigger.locator("..");
  648 |   await expect(details).toHaveAttribute("open", "");
  649 |   await page.evaluate(() => {
  650 |     for (const animation of document.getAnimations()) animation.finish();
  651 |   });
  652 |   expect(
  653 |     await details.evaluate((element) => {
  654 |       const icon = element.querySelector(
  655 |         '[data-gsxui-slot-accordion-trigger-icon]',
  656 |       );
  657 |       return {
  658 |         detailsContentDisplay: getComputedStyle(element, "::details-content").display,
  659 |         iconRotate: icon ? getComputedStyle(icon).rotate : null,
  660 |       };
  661 |     }),
  662 |   ).toEqual({
  663 |     detailsContentDisplay: "grid",
  664 |     iconRotate: "180deg",
  665 |   });
  666 | });
  667 | 
  668 | for (const route of desktopRoutes) {
  669 |   test(`${route} keeps its desktop light and dark visual contract`, async ({ page }) => {
  670 |     for (const theme of ["light", "dark"] as const) {
  671 |       await prepareVisualRoute(page, route, theme);
  672 |       await expect(page.locator("body")).toHaveScreenshot(
  673 |         `desktop-${theme}-${snapshotSlug(route)}.png`,
  674 |         screenshotOptions,
  675 |       );
  676 |     }
  677 |   });
  678 | }
  679 | 
  680 | test.describe("mobile visual contracts", () => {
  681 |   test.use({ viewport: { width: 390, height: 844 } });
  682 | 
  683 |   for (const route of mobileRoutes) {
  684 |     test(`${route} keeps its mobile visual contract`, async ({ page }) => {
  685 |       await prepareVisualRoute(page, route, "light");
  686 |       await expect(page.locator("body")).toHaveScreenshot(
  687 |         `mobile-light-${snapshotSlug(route)}.png`,
  688 |         screenshotOptions,
  689 |       );
  690 |     });
  691 |   }
  692 | });
  693 | 
  694 | // The three tests below pin presentation that a component's own CSS supplies
  695 | // ON TOP OF the Button it renders through. ui.Button ships concrete Tailwind
  696 | // utilities compiled from its style recipe, and those land in @layer
  697 | // utilities; the cascade compares whole layers before specificity, so a rule
  698 | // in @layer components can never override them. Every property asserted here
  699 | // silently reverted to Button's own value when the Button switched from a
  700 | // recipe token to concrete utilities, and the 313-test suite did not notice.
  701 | // If one of these fails, check which layer the rule lives in before touching
  702 | // the expectation.
  703 | 
  704 | test("Carousel arrows stay circular over Button's own radius", async ({
  705 |   page,
  706 | }) => {
  707 |   const response = await page.goto("/x/carousel/basic");
  708 |   expect(response?.status(), "carousel/basic fixture response").toBe(200);
  709 | 
  710 |   for (const marker of [
  711 |     "[data-gsxui-slot-carousel-previous]",
  712 |     "[data-gsxui-slot-carousel-next]",
  713 |   ]) {
  714 |     const arrow = page.locator(marker).first();
  715 |     const { radius, height } = await arrow.evaluate((element) => {
  716 |       const css = getComputedStyle(element);
  717 |       return {
  718 |         radius: Number.parseFloat(css.borderTopLeftRadius),
  719 |         height: Number.parseFloat(css.height),
  720 |       };
  721 |     });
  722 |     // rounded-full, not Button's rounded-lg (10px): a pill radius is clamped
  723 |     // to half the box, so anything >= height/2 renders as a circle here.
  724 |     expect(
  725 |       radius,
  726 |       `${marker} must be circular, not Button's rounded-lg`,
  727 |     ).toBeGreaterThanOrEqual(height / 2);
  728 |   }
  729 | });
  730 | 
  731 | test("InputGroupButton keeps its own type scale and radius over Button's size ramp", async ({
  732 |   page,
  733 | }) => {
  734 |   const response = await page.goto("/x/input-group/basic");
  735 |   expect(response?.status(), "input-group/basic fixture response").toBe(200);
```