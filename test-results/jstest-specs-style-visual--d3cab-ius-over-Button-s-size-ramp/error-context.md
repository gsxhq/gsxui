# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/style-visual.spec.ts >> InputGroupButton keeps its own type scale and radius over Button's size ramp
- Location: jstest/specs/style-visual.spec.ts:731:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/input-group/basic", waiting until "load"

```

# Test source

```ts
  634 | }) => {
  635 |   const response = await page.goto("/f/style-contract");
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
> 734 |   const response = await page.goto("/x/input-group/basic");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  735 |   expect(response?.status(), "input-group/basic fixture response").toBe(200);
  736 | 
  737 |   const button = page.locator("[data-gsxui-slot-input-group-button]").first();
  738 |   await expect(button).toHaveAttribute("data-size", "xs");
  739 |   expect(
  740 |     await button.evaluate((element) => {
  741 |       const css = getComputedStyle(element);
  742 |       return {
  743 |         borderRadius: css.borderTopLeftRadius,
  744 |         fontSize: css.fontSize,
  745 |         paddingLeft: css.paddingLeft,
  746 |       };
  747 |     }),
  748 |   ).toEqual({
  749 |     // InputGroup's own xs ramp: rounded-[calc(var(--radius)-3px)] and text-sm,
  750 |     // NOT Button's xs ramp (rounded-[min(var(--radius-md),10px)] = 8px and
  751 |     // text-xs = 12px).
  752 |     borderRadius: "7px",
  753 |     fontSize: "14px",
  754 |     paddingLeft: "6px",
  755 |   });
  756 | });
  757 | 
```