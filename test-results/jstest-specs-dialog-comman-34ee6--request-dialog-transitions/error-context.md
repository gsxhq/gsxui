# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/dialog.spec.ts >> command palette shortcuts and navigable selections request dialog transitions
- Location: jstest/specs/dialog.spec.ts:62:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/dialog/basic", waiting until "load"

```

# Test source

```ts
  1   | import type { Page } from "@playwright/test";
  2   | import { expect, test } from "../support/fixtures";
  3   | 
  4   | const BASIC = "/x/dialog/basic";
  5   | const DIALOG = "dialog[data-gsxui-dialog-content]";
  6   | const COMMAND_DIALOG = "dialog[data-gsxui-command-dialog]";
  7   | 
  8   | async function dispatch(page: Page, selector: string, type: string, detail = {}) {
  9   |   return page.locator(selector).evaluate(
  10  |     (element, { type, detail }) =>
  11  |       element.dispatchEvent(
  12  |         new CustomEvent(type, { bubbles: true, cancelable: true, detail }),
  13  |       ),
  14  |     { type, detail },
  15  |   );
  16  | }
  17  | 
  18  | async function open(page: Page, selector = DIALOG) {
  19  |   await dispatch(page, selector, "gsxui:request-open");
  20  |   await expect(page.locator(selector)).toHaveJSProperty("open", true);
  21  | }
  22  | 
  23  | async function close(page: Page, selector = DIALOG) {
  24  |   await dispatch(page, selector, "gsxui:request-close");
  25  |   await expect(page.locator(selector)).toHaveJSProperty("open", false);
  26  | }
  27  | 
  28  | test("the proximity trigger requests one targeted open with its stable reason", async ({ page }) => {
  29  |   await page.goto(BASIC);
  30  |   await page.evaluate(() => {
  31  |     (window as any).__dialogRequests = [];
  32  |     document.addEventListener("gsxui:request-open", (event) => {
  33  |       const request = event as CustomEvent;
  34  |       (window as any).__dialogRequests.push({
  35  |         detail: request.detail,
  36  |         target: (request.target as Element).matches(
  37  |           "dialog[data-gsxui-dialog-content]",
  38  |         ),
  39  |       });
  40  |     });
  41  |   });
  42  | 
  43  |   await page.getByRole("button", { name: "Delete account", exact: true }).click();
  44  | 
  45  |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  46  |   expect(await page.evaluate(() => (window as any).__dialogRequests)).toEqual([
  47  |     { detail: { reason: "trigger" }, target: true },
  48  |   ]);
  49  | });
  50  | 
  51  | test("descendant request events drive the native dialog state", async ({ page }) => {
  52  |   await page.goto(BASIC);
  53  |   const source = `${DIALOG} [data-gsxui-dialog-title]`;
  54  | 
  55  |   await dispatch(page, source, "gsxui:request-open", { reason: "application" });
  56  |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", true);
  57  | 
  58  |   await dispatch(page, source, "gsxui:request-close", { reason: "application" });
  59  |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
  60  | });
  61  | 
  62  | test("command palette shortcuts and navigable selections request dialog transitions", async ({ page }) => {
> 63  |   await page.goto(BASIC);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  64  |   await page.evaluate(() => {
  65  |     document.body.insertAdjacentHTML(
  66  |       "beforeend",
  67  |       `
  68  |         <div data-gsxui-dialog>
  69  |           <dialog data-gsxui-dialog-content data-gsxui-command-dialog data-state="closed">
  70  |             <div data-gsxui-command>
  71  |               <input data-gsxui-command-input>
  72  |               <div data-gsxui-command-list>
  73  |                 <button data-gsxui-command-item data-href="/command-selected">Open selected page</button>
  74  |               </div>
  75  |             </div>
  76  |           </dialog>
  77  |         </div>
  78  |       `,
  79  |     );
  80  |     (window as any).__commandDialogEvents = [];
  81  |     document.addEventListener("gsxui:request-open", (event) => {
  82  |       if ((event.target as Element).matches("dialog[data-gsxui-command-dialog]")) {
  83  |         (window as any).__commandDialogEvents.push({
  84  |           type: event.type,
  85  |           detail: (event as CustomEvent).detail,
  86  |           target: (event.target as Element).matches("dialog[data-gsxui-command-dialog]"),
  87  |         });
  88  |       }
  89  |     });
  90  |     document.addEventListener("gsxui:request-close", (event) => {
  91  |       if ((event.target as Element).matches("dialog[data-gsxui-command-dialog]")) {
  92  |         (window as any).__commandDialogEvents.push({
  93  |           type: event.type,
  94  |           detail: (event as CustomEvent).detail,
  95  |           target: (event.target as Element).matches("dialog[data-gsxui-command-dialog]"),
  96  |         });
  97  |       }
  98  |     });
  99  |   });
  100 | 
  101 |   const dialog = page.locator(COMMAND_DIALOG);
  102 |   await page.keyboard.press("Control+k");
  103 |   await expect(dialog).toHaveJSProperty("open", true);
  104 |   expect(await page.evaluate(() => (window as any).__commandDialogEvents)).toEqual([
  105 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  106 |   ]);
  107 | 
  108 |   await dialog.evaluate((element) => {
  109 |     element.animate([{ opacity: 1 }, { opacity: 0 }], { duration: 120 });
  110 |   });
  111 |   await page.keyboard.press("Control+k");
  112 |   await expect(dialog).toHaveAttribute("data-state", "closed");
  113 |   await expect(dialog).toHaveJSProperty("open", true);
  114 |   await expect(dialog).toHaveJSProperty("open", false);
  115 |   expect(await page.evaluate(() => (window as any).__commandDialogEvents)).toEqual([
  116 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  117 |     { type: "gsxui:request-close", detail: { reason: "shortcut" }, target: true },
  118 |   ]);
  119 | 
  120 |   await page.keyboard.press("Control+k");
  121 |   await expect(dialog).toHaveJSProperty("open", true);
  122 |   await page.locator(`${COMMAND_DIALOG} [data-gsxui-command-item][data-href]`).first().evaluate((item) => {
  123 |     item.setAttribute("data-href", "#command-selected");
  124 |     (window as any).__commandNavigationOrder = [];
  125 |     document.addEventListener(
  126 |       "gsxui:request-close",
  127 |       (event) => {
  128 |         if ((event.target as Element).matches("dialog[data-gsxui-command-dialog]"))
  129 |           (window as any).__commandNavigationOrder.push("request-close");
  130 |       },
  131 |       { once: true },
  132 |     );
  133 |     window.addEventListener(
  134 |       "hashchange",
  135 |       () => (window as any).__commandNavigationOrder.push("navigation"),
  136 |       { once: true },
  137 |     );
  138 |   });
  139 |   await page.locator(`${COMMAND_DIALOG} [data-gsxui-command-item][data-href]`).first().click();
  140 |   await expect.poll(() => page.evaluate(() => (window as any).__commandNavigationOrder)).toEqual([
  141 |     "request-close",
  142 |     "navigation",
  143 |   ]);
  144 |   expect(await page.evaluate(() => (window as any).__commandDialogEvents)).toEqual([
  145 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  146 |     { type: "gsxui:request-close", detail: { reason: "shortcut" }, target: true },
  147 |     { type: "gsxui:request-open", detail: { reason: "shortcut" }, target: true },
  148 |     { type: "gsxui:request-close", detail: { reason: "select" }, target: true },
  149 |   ]);
  150 | });
  151 | 
  152 | test("a late document listener can cancel both deferred request defaults", async ({ page }) => {
  153 |   await page.goto(BASIC);
  154 |   await page.evaluate(() => {
  155 |     (window as any).__cancelledDialogRequests = [];
  156 |     document.addEventListener("gsxui:request-open", (event) => {
  157 |       event.preventDefault();
  158 |       (window as any).__cancelledDialogRequests.push(event.type);
  159 |     });
  160 |   });
  161 | 
  162 |   await page.getByRole("button", { name: "Delete account", exact: true }).click();
  163 |   await expect(page.locator(DIALOG)).toHaveJSProperty("open", false);
```