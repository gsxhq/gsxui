# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/runtime-style-contract.spec.ts >> real interactions cover the exact runtime-owned style contract
- Location: jstest/specs/runtime-style-contract.spec.ts:29:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/alert-dialog/basic", waiting until "load"

```

# Test source

```ts
  1   | import { readFileSync } from "node:fs";
  2   | 
  3   | import type { Locator, Page } from "@playwright/test";
  4   | 
  5   | import { expect, test } from "../support/fixtures";
  6   | import { jstestDir } from "../support/paths";
  7   | 
  8   | type Entry = {
  9   |   component: string;
  10  |   slot: string;
  11  |   attribute: string;
  12  |   value: string;
  13  |   scenario: string;
  14  | };
  15  | 
  16  | const manifest = JSON.parse(
  17  |   readFileSync(`${jstestDir}/runtime-style-contract.json`, "utf8"),
  18  | ) as Entry[];
  19  | 
  20  | const key = (entry: Entry) =>
  21  |   [
  22  |     entry.component,
  23  |     entry.slot,
  24  |     entry.attribute,
  25  |     entry.value,
  26  |     entry.scenario,
  27  |   ].join("\u0000");
  28  | 
  29  | test("real interactions cover the exact runtime-owned style contract", async ({
  30  |   page,
  31  | }) => {
  32  |   test.setTimeout(120_000);
  33  |   const observed = new Set<string>();
  34  | 
  35  |   async function observe(
  36  |     scenario: string,
  37  |     component: string,
  38  |     slot: string,
  39  |     locator: Locator,
  40  |     attribute: string,
  41  |     value: string,
  42  |   ) {
  43  |     const entry = { component, slot, attribute, value, scenario };
  44  |     expect(manifest.map(key)).toContain(key(entry));
  45  |     if (value === "") {
  46  |       await expect(locator).toHaveAttribute(attribute, "");
  47  |     } else {
  48  |       await expect(locator).toHaveAttribute(attribute, value);
  49  |     }
  50  |     observed.add(key(entry));
  51  |   }
  52  | 
  53  |   async function dialogLifecycle(
  54  |     scenario: string,
  55  |     component: string,
  56  |     route: string,
  57  |     triggerName: string,
  58  |     triggerSlot: string,
  59  |     contentSlot: string,
  60  |   ) {
> 61  |     await page.goto(route);
      |                ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  62  |     const trigger = page.locator(`[data-gsxui-slot-${triggerSlot}]`).filter({
  63  |       hasText: triggerName,
  64  |     });
  65  |     const content = page.locator(`[data-gsxui-slot-${contentSlot}]`).first();
  66  |     await trigger.click();
  67  |     await observe(
  68  |       scenario,
  69  |       component,
  70  |       triggerSlot,
  71  |       trigger,
  72  |       "aria-expanded",
  73  |       "true",
  74  |     );
  75  |     await observe(
  76  |       scenario,
  77  |       component,
  78  |       contentSlot,
  79  |       content,
  80  |       "data-state",
  81  |       "open",
  82  |     );
  83  |   }
  84  | 
  85  |   await dialogLifecycle(
  86  |     "alert-dialog-lifecycle",
  87  |     "alert-dialog",
  88  |     "/x/alert-dialog/basic",
  89  |     "Show dialog",
  90  |     "alert-dialog-trigger",
  91  |     "alert-dialog-content",
  92  |   );
  93  |   await dialogLifecycle(
  94  |     "dialog-lifecycle",
  95  |     "dialog",
  96  |     "/x/dialog/basic",
  97  |     "Delete account",
  98  |     "dialog-trigger",
  99  |     "dialog-content",
  100 |   );
  101 |   await dialogLifecycle(
  102 |     "drawer-lifecycle",
  103 |     "drawer",
  104 |     "/x/drawer/basic",
  105 |     "Open Drawer",
  106 |     "drawer-trigger",
  107 |     "drawer-content",
  108 |   );
  109 |   await dialogLifecycle(
  110 |     "sheet-lifecycle",
  111 |     "sheet",
  112 |     "/x/sheet/basic",
  113 |     "Edit Profile",
  114 |     "sheet-trigger",
  115 |     "sheet-content",
  116 |   );
  117 | 
  118 |   await page.goto("/x/hover-card/basic");
  119 |   const hoverTrigger = page.locator("[data-gsxui-hovercard-trigger]");
  120 |   const hoverContent = page.locator("[data-gsxui-slot-hover-card-content]");
  121 |   await hoverTrigger.hover();
  122 |   await observe(
  123 |     "hover-card-lifecycle",
  124 |     "hover-card",
  125 |     "hover-card-content",
  126 |     hoverContent,
  127 |     "data-state",
  128 |     "open",
  129 |   );
  130 | 
  131 |   await page.goto("/x/popover/basic");
  132 |   const popoverTrigger = page
  133 |     .locator("[data-gsxui-slot-popover-trigger]")
  134 |     .first();
  135 |   const popoverContent = page.locator("[data-gsxui-slot-popover-content]");
  136 |   await popoverTrigger.click();
  137 |   await observe(
  138 |     "popover-lifecycle",
  139 |     "popover",
  140 |     "popover-trigger",
  141 |     popoverTrigger,
  142 |     "aria-expanded",
  143 |     "true",
  144 |   );
  145 |   await observe(
  146 |     "popover-lifecycle",
  147 |     "popover",
  148 |     "popover-content",
  149 |     popoverContent,
  150 |     "data-state",
  151 |     "open",
  152 |   );
  153 | 
  154 |   await page.goto("/x/tooltip/basic");
  155 |   const tooltipTrigger = page
  156 |     .locator("[data-gsxui-slot-tooltip-trigger]")
  157 |     .first();
  158 |   const tooltipContent = page.locator("[data-gsxui-slot-tooltip-content]");
  159 |   await tooltipTrigger.hover();
  160 |   await observe(
  161 |     "tooltip-lifecycle",
```