# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/menus-style-contract.spec.ts >> navigation viewport chrome and trigger rotation follow reflected state
- Location: jstest/specs/menus-style-contract.spec.ts:103:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/f/style-contract", waiting until "load"

```

# Test source

```ts
  1   | import { expect, test } from "../support/fixtures";
  2   | 
  3   | async function openStyleContract(page: import("@playwright/test").Page) {
> 4   |   const response = await page.goto("/f/style-contract");
      |                               ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  5   |   expect(response?.status(), "style contract fixture response").toBe(200);
  6   | }
  7   | 
  8   | test("menu destructive focus state wins without important specificity", async ({
  9   |   page,
  10  | }) => {
  11  |   await openStyleContract(page);
  12  | 
  13  |   const reference = page.locator(
  14  |     '[data-style-contract-reference="menu-destructive-focus"]',
  15  |   );
  16  |   const expected = await reference.evaluate((element) => {
  17  |     const css = getComputedStyle(element);
  18  |     return { backgroundColor: css.backgroundColor, color: css.color };
  19  |   });
  20  | 
  21  |   for (const family of ["dropdown", "context", "menubar"]) {
  22  |     const item = page.locator(`[data-style-contract="${family}-destructive"]`);
  23  |     await item.focus();
  24  |     expect(
  25  |       await item.evaluate((element) => {
  26  |         const css = getComputedStyle(element);
  27  |         return { backgroundColor: css.backgroundColor, color: css.color };
  28  |       }),
  29  |     ).toEqual(expected);
  30  |   }
  31  | });
  32  | 
  33  | test("checked menu and combobox indicators follow semantic owner state", async ({
  34  |   page,
  35  | }) => {
  36  |   await openStyleContract(page);
  37  | 
  38  |   const cases = [
  39  |     ["dropdown", "[data-gsxui-dropdown-checkbox-indicator]"],
  40  |     ["context", "[data-gsxui-contextmenu-radio-indicator]"],
  41  |     ["menubar", "[data-gsxui-menubar-checkbox-indicator]"],
  42  |     ["combobox", "[data-gsxui-combobox-item-indicator]"],
  43  |   ] as const;
  44  | 
  45  |   for (const [family, indicator] of cases) {
  46  |     const display = async (state: "checked" | "unchecked") =>
  47  |       page
  48  |         .locator(`[data-style-contract="${family}-${state}"] ${indicator}`)
  49  |         .evaluate((element) => getComputedStyle(element).display);
  50  |     expect(await display("checked")).toBe("flex");
  51  |     expect(await display("unchecked")).toBe("none");
  52  |   }
  53  | });
  54  | 
  55  | test("submenu uses the emitted right-side directional enter offset", async ({
  56  |   page,
  57  | }) => {
  58  |   await openStyleContract(page);
  59  | 
  60  |   const content = page.locator('[data-style-contract="submenu-right"]');
  61  |   const keyframes = await content.evaluate((element) => {
  62  |     (element as HTMLElement).showPopover();
  63  |     return element.getAnimations().flatMap((animation) => {
  64  |       const effect = animation.effect as KeyframeEffect | null;
  65  |       return effect?.getKeyframes().map((frame) => frame.translate) ?? [];
  66  |     });
  67  |   });
  68  |   expect(keyframes).toContain("-8px");
  69  |   await content.evaluate((element) => (element as HTMLElement).hidePopover());
  70  | });
  71  | 
  72  | test("CommandDialog sizing is supplied by CSS ancestry", async ({ page }) => {
  73  |   await openStyleContract(page);
  74  |   await page.keyboard.press(process.platform === "darwin" ? "Meta+K" : "Control+K");
  75  | 
  76  |   const dialog = page.locator('[data-style-contract="command-dialog"]');
  77  |   await expect(dialog).toBeVisible();
  78  |   expect(
  79  |     await dialog.evaluate((element) => {
  80  |       const command = element.querySelector("[data-gsxui-command]");
  81  |       const wrapper = element.querySelector("[data-gsxui-command-input-wrapper]");
  82  |       const input = element.querySelector("[data-gsxui-command-input]");
  83  |       const css = getComputedStyle(element);
  84  |       return {
  85  |         maxWidth: css.maxWidth,
  86  |         overflow: css.overflow,
  87  |         padding: css.padding,
  88  |         commandPadding: command ? getComputedStyle(command).padding : null,
  89  |         wrapperHeight: wrapper ? getComputedStyle(wrapper).height : null,
  90  |         inputHeight: input ? getComputedStyle(input).height : null,
  91  |       };
  92  |     }),
  93  |   ).toEqual({
  94  |     maxWidth: "512px",
  95  |     overflow: "hidden",
  96  |     padding: "0px",
  97  |     commandPadding: "4px",
  98  |     wrapperHeight: "48px",
  99  |     inputHeight: "48px",
  100 |   });
  101 | });
  102 | 
  103 | test("navigation viewport chrome and trigger rotation follow reflected state", async ({
  104 |   page,
```