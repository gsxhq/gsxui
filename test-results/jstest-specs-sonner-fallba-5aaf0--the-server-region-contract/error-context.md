# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/sonner.spec.ts >> fallback toaster matches the server region contract
- Location: jstest/specs/sonner.spec.ts:88:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/sonner/types", waiting until "load"

```

# Test source

```ts
  1   | import type { Page } from "@playwright/test";
  2   | 
  3   | import { expect, test } from "../support/fixtures";
  4   | 
  5   | const route = "/x/sonner/types";
  6   | 
  7   | async function expectSameTurnToastLifecycle(page: Page, label: string) {
  8   |   const card = page.locator(`[data-same-turn-toast="${label}"]`);
  9   |   await expect(card).toHaveAttribute("data-state", "open");
  10  |   await expect(card).toHaveAttribute("data-visible", "true");
  11  |   expect(
  12  |     await card.evaluate((element) => {
  13  |       const probe = (window as any).__sonnerSameTurnProbe;
  14  |       const row = element as HTMLElement;
  15  |       return {
  16  |         expectedParent: row.parentElement === probe.expectedRegion,
  17  |         opacity: row.style.opacity,
  18  |         transform: row.style.transform,
  19  |         pointerEvents: row.style.pointerEvents,
  20  |       };
  21  |     }),
  22  |   ).toEqual({
  23  |     expectedParent: true,
  24  |     opacity: "1",
  25  |     transform: "translateY(0px) scale(1)",
  26  |     pointerEvents: "auto",
  27  |   });
  28  |   await expect
  29  |     .poll(() =>
  30  |       page.evaluate(
  31  |         () => (window as any).__sonnerSameTurnProbe.activeTimers.size,
  32  |       ),
  33  |     )
  34  |     .toBe(1);
  35  |   expect(
  36  |     await page.evaluate(
  37  |       () => (window as any).__sonnerSameTurnProbe.createdTimers,
  38  |     ),
  39  |   ).toBe(1);
  40  | 
  41  |   await page.evaluate(() => {
  42  |     const probe = (window as any).__sonnerSameTurnProbe;
  43  |     probe.row.addEventListener("gsxui:toast-action", () => {
  44  |       probe.directActionEvents++;
  45  |     });
  46  |     probe.actionButton.click();
  47  |   });
  48  |   await expect
  49  |     .poll(() =>
  50  |       page.evaluate(() => (window as any).__sonnerSameTurnProbe.actionCalls),
  51  |     )
  52  |     .toBe(1);
  53  |   await expect
  54  |     .poll(() =>
  55  |       page.evaluate(() => (window as any).__sonnerSameTurnProbe.actionEvents),
  56  |     )
  57  |     .toBe(1);
  58  |   await expect(card).toHaveCount(0);
  59  |   await expect
  60  |     .poll(() =>
  61  |       page.evaluate(
  62  |         () => (window as any).__sonnerSameTurnProbe.activeTimers.size,
  63  |       ),
  64  |     )
  65  |     .toBe(0);
  66  | 
  67  |   await page.evaluate(() => {
  68  |     (window as any).__sonnerSameTurnProbe.actionButton.click();
  69  |   });
  70  |   expect(
  71  |     await page.evaluate(() => {
  72  |       const probe = (window as any).__sonnerSameTurnProbe;
  73  |       return {
  74  |         actionCalls: probe.actionCalls,
  75  |         directActionEvents: probe.directActionEvents,
  76  |         connected: probe.row.isConnected,
  77  |         parent: probe.row.parentNode,
  78  |       };
  79  |     }),
  80  |   ).toEqual({
  81  |     actionCalls: 1,
  82  |     directActionEvents: 1,
  83  |     connected: false,
  84  |     parent: null,
  85  |   });
  86  | }
  87  | 
  88  | test("fallback toaster matches the server region contract", async ({ page }) => {
> 89  |   await page.goto(route);
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  90  |   await page.evaluate(() => {
  91  |     const section = document.querySelector('[aria-label="Notifications"]')!;
  92  |     for (const template of [...section.querySelectorAll("template")]) {
  93  |       document.body.append(template);
  94  |     }
  95  |     section.remove();
  96  |     window.gsxui.toast.success("Fallback toast", { duration: 60_000 });
  97  |     const template = document.querySelector<HTMLTemplateElement>(
  98  |       'template[data-gsxui-toast-template="info"]',
  99  |     )!;
  100 |     const serverRow = template.content.firstElementChild!.cloneNode(true) as HTMLElement;
  101 |     serverRow.dataset.fallbackServer = "true";
  102 |     serverRow.dataset.duration = "60000";
  103 |     document.querySelector("#gsxui-toaster")!.append(serverRow);
  104 |   });
  105 | 
  106 |   const region = page.locator("#gsxui-toaster");
  107 |   await expect(region).toHaveAttribute("data-gsxui-toaster", "");
  108 |   await expect(region).toHaveAttribute("data-gsxui-slot-toaster", "");
  109 |   await expect(region).not.toHaveAttribute("class", /.+/);
  110 |   await expect(region.locator("li[data-gsxui-toast]")).toHaveCount(2);
  111 |   await expect(
  112 |     region.locator('li[data-gsxui-toast][data-fallback-server="true"]'),
  113 |   ).toHaveAttribute("data-state", "open");
  114 |   expect(
  115 |     await region.evaluate((element) => {
  116 |       const section = element.parentElement!;
  117 |       const css = getComputedStyle(element);
  118 |       return {
  119 |         sectionLabel: section.getAttribute("aria-label"),
  120 |         sectionTabIndex: section.getAttribute("tabindex"),
  121 |         position: css.position,
  122 |         right: css.right,
  123 |         bottom: css.bottom,
  124 |         pointerEvents: css.pointerEvents,
  125 |       };
  126 |     }),
  127 |   ).toEqual({
  128 |     sectionLabel: "Notifications",
  129 |     sectionTabIndex: "-1",
  130 |     position: "fixed",
  131 |     right: "0px",
  132 |     bottom: "0px",
  133 |     pointerEvents: "none",
  134 |   });
  135 | });
  136 | 
  137 | test("every toast type uses semantic icon color and dedicated hooks", async ({ page }) => {
  138 |   await page.goto(route);
  139 |   await page.evaluate(() => {
  140 |     const api = window.gsxui.toast;
  141 |     api("Default", { duration: 60_000 });
  142 |     api.success("Success", { duration: 60_000 });
  143 |     api.info("Info", { duration: 60_000 });
  144 |     api.warning("Warning", { duration: 60_000 });
  145 |     api.error("Error", { duration: 60_000 });
  146 |     api.loading("Loading", { duration: Infinity });
  147 |   });
  148 | 
  149 |   const cards = page.locator("li[data-gsxui-toast]");
  150 |   await expect(cards).toHaveCount(6);
  151 |   const result = await cards.evaluateAll((elements) =>
  152 |     elements.map((element) => {
  153 |       const type = (element as HTMLElement).dataset.type!;
  154 |       const icon = element.querySelector("[data-gsxui-toast-icon]");
  155 |       let semanticColor: string | null = null;
  156 |       if (!["default", "loading"].includes(type)) {
  157 |         const probe = document.createElement("span");
  158 |         probe.style.color = `var(--${type === "error" ? "destructive" : type})`;
  159 |         document.body.append(probe);
  160 |         semanticColor = getComputedStyle(probe).color;
  161 |         probe.remove();
  162 |       }
  163 |       return {
  164 |         type,
  165 |         toastSlot: element.hasAttribute("data-gsxui-slot-toast"),
  166 |         className: element.getAttribute("class"),
  167 |         toastIconSlot: icon?.hasAttribute("data-gsxui-slot-toast-icon") ?? false,
  168 |         primitiveIconSlot: icon?.hasAttribute("data-gsxui-slot-icon") ?? false,
  169 |         iconColor: icon ? getComputedStyle(icon).color : null,
  170 |         semanticColor,
  171 |       };
  172 |     }),
  173 |   );
  174 |   expect(result.map(({ type }) => type)).toEqual([
  175 |     "default",
  176 |     "success",
  177 |     "info",
  178 |     "warning",
  179 |     "error",
  180 |     "loading",
  181 |   ]);
  182 |   for (const card of result) {
  183 |     expect(card.toastSlot).toBe(true);
  184 |     expect(card.className).toBeNull();
  185 |     if (card.semanticColor) {
  186 |       expect(card.toastIconSlot).toBe(true);
  187 |       expect(card.primitiveIconSlot).toBe(true);
  188 |       expect(card.iconColor).toBe(card.semanticColor);
  189 |     }
```