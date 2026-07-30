# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sidebar-style-contract.spec.ts >> desktop offcanvas state drives gap, fixed edge, rail cursor, event, and cookie ownership
- Location: jstest/specs/sidebar-style-contract.spec.ts:15:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/sidebar-contract?case=offcanvas-left", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | async function openFixture(
  4   |   page: import("@playwright/test").Page,
  5   |   name: string,
  6   |   viewport = { width: 1024, height: 768 },
  7   | ) {
  8   |   await page.setViewportSize(viewport);
> 9   |   const response = await page.goto(`/f/sidebar-contract?case=${name}`);
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  10  |   expect(response?.status(), `${name} fixture response`).toBe(200);
  11  | }
  12  | 
  13  | const slot = (name: string) => `[data-gsxui-slot-${name}]`;
  14  | 
  15  | test("desktop offcanvas state drives gap, fixed edge, rail cursor, event, and cookie ownership", async ({
  16  |   page,
  17  | }) => {
  18  |   await openFixture(page, "offcanvas-left");
  19  |   const wrapper = page.locator(slot("sidebar-wrapper"));
  20  |   const desktop = page.locator(slot("sidebar-desktop"));
  21  |   const gap = page.locator(slot("sidebar-gap"));
  22  |   const container = page.locator(slot("sidebar-container"));
  23  |   const rail = desktop.locator(slot("sidebar-rail"));
  24  | 
  25  |   expect(
  26  |     await wrapper.evaluate((element) => {
  27  |       const css = getComputedStyle(element);
  28  |       return {
  29  |         state: element.getAttribute("data-state"),
  30  |         width: css.getPropertyValue("--sidebar-width").trim(),
  31  |         iconWidth: css.getPropertyValue("--sidebar-width-icon").trim(),
  32  |       };
  33  |     }),
  34  |   ).toEqual({ state: "expanded", width: "16rem", iconWidth: "3rem" });
  35  |   expect(await gap.evaluate((element) => getComputedStyle(element).width)).toBe("256px");
  36  |   expect(
  37  |     await container.evaluate((element) => {
  38  |       const css = getComputedStyle(element);
  39  |       return { left: css.left, width: css.width, position: css.position };
  40  |     }),
  41  |   ).toEqual({ left: "0px", width: "256px", position: "fixed" });
  42  |   expect(await rail.evaluate((element) => getComputedStyle(element).cursor)).toBe("w-resize");
  43  | 
  44  |   await page.evaluate(() => {
  45  |     document.cookie = "sidebar_owner=caller; path=/";
  46  |     const wrapper = document.querySelector('[data-gsxui-slot-sidebar-wrapper]')!;
  47  |     (window as typeof window & { sidebarChanges?: boolean[] }).sidebarChanges = [];
  48  |     wrapper.addEventListener("gsxui:change", (event) => {
  49  |       (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges.push(
  50  |         Boolean((event as CustomEvent).detail.open),
  51  |       );
  52  |     });
  53  |   });
  54  |   await page.getByRole("button", { name: "Toggle Sidebar" }).first().click();
  55  |   await page.evaluate(() => {
  56  |     for (const animation of document.getAnimations()) animation.finish();
  57  |   });
  58  | 
  59  |   await expect(wrapper).toHaveAttribute("data-state", "collapsed");
  60  |   await expect(desktop).toHaveAttribute("data-state", "collapsed");
  61  |   await expect(desktop).toHaveAttribute("data-collapsible", "offcanvas");
  62  |   expect(await gap.evaluate((element) => getComputedStyle(element).width)).toBe("0px");
  63  |   expect(await container.evaluate((element) => getComputedStyle(element).left)).toBe("-256px");
  64  |   expect(await rail.evaluate((element) => getComputedStyle(element).cursor)).toBe("e-resize");
  65  |   expect(
  66  |     await page.evaluate(
  67  |       () => (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges,
  68  |     ),
  69  |   ).toEqual([false]);
  70  |   expect(await page.evaluate(() => document.cookie)).toContain("sidebar_owner=caller");
  71  | 
  72  |   await rail.click({ force: true });
  73  |   await expect(wrapper).toHaveAttribute("data-state", "expanded");
  74  |   await expect(desktop).toHaveAttribute("data-collapsible", "");
  75  |   expect(
  76  |     await page.evaluate(
  77  |       () => (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges,
  78  |     ),
  79  |   ).toEqual([false, true]);
  80  | 
  81  |   await page.getByRole("button", { name: "Custom Sidebar Trigger" }).click();
  82  |   await expect(wrapper).toHaveAttribute("data-state", "collapsed");
  83  |   expect(
  84  |     await page.evaluate(
  85  |       () => (window as typeof window & { sidebarChanges: boolean[] }).sidebarChanges,
  86  |     ),
  87  |   ).toEqual([false, true, false]);
  88  | });
  89  | 
  90  | test("right offcanvas and icon/floating/inset geometry follow explicit axes", async ({
  91  |   page,
  92  | }) => {
  93  |   await openFixture(page, "offcanvas-right");
  94  |   const right = page.locator(slot("sidebar-container"));
  95  |   expect(
  96  |     await right.evaluate((element) => {
  97  |       const css = getComputedStyle(element);
  98  |       return { right: css.right, borderLeft: css.borderLeftWidth };
  99  |     }),
  100 |   ).toEqual({ right: "-256px", borderLeft: "1px" });
  101 |   expect(
  102 |     await page
  103 |       .locator(slot("sidebar-desktop"))
  104 |       .locator(slot("sidebar-rail"))
  105 |       .evaluate((element) => getComputedStyle(element).cursor),
  106 |   ).toBe("w-resize");
  107 | 
  108 |   await openFixture(page, "icon-sidebar");
  109 |   expect(await page.locator(slot("sidebar-gap")).evaluate((el) => getComputedStyle(el).width)).toBe(
```