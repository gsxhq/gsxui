# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/layer-precedence.spec.ts >> sweep all fixtures (light)
- Location: jstest/specs/layer-precedence.spec.ts:28:3

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/x/accordion/basic", waiting until "load"

```

# Test source

```ts
  1  | import type { Page } from "@playwright/test";
  2  | 
  3  | export const SWEPT_PROPERTIES = [
  4  |   "display", "position", "borderRadius", "borderWidth", "borderColor",
  5  |   "width", "height", "paddingLeft", "paddingRight", "paddingTop",
  6  |   "paddingBottom", "gap", "fontSize", "fontWeight", "lineHeight",
  7  |   "color", "backgroundColor", "opacity", "boxShadow", "textDecorationLine",
  8  |   "alignItems", "justifyContent",
  9  | ] as const;
  10 | 
  11 | /**
  12 |  * sweepComputedStyles records the resting computed style of every element
  13 |  * carrying a `data-gsxui-slot-*` attribute on the given page. It is the only
  14 |  * thing that reliably catches CSS-layer precedence regressions: they change
  15 |  * rendering without failing any assertion. Repeats of the same marker on a
  16 |  * page are disambiguated with a `#n` suffix in DOM order.
  17 |  */
  18 | export async function sweepComputedStyles(
  19 |   page: Page,
  20 |   url: string,
  21 | ): Promise<Record<string, Record<string, string>>> {
> 22 |   await page.goto(url);
     |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  23 |   return page.evaluate((props: readonly string[]) => {
  24 |     const out: Record<string, Record<string, string>> = {};
  25 |     const seen = new Map<string, number>();
  26 |     for (const el of document.querySelectorAll("*")) {
  27 |       const marker = [...el.attributes]
  28 |         .map((a) => a.name)
  29 |         .find((n) => n.startsWith("data-gsxui-slot-"));
  30 |       if (!marker) continue;
  31 |       const n = (seen.get(marker) ?? 0) + 1;
  32 |       seen.set(marker, n);
  33 |       const computed = getComputedStyle(el);
  34 |       const record: Record<string, string> = {};
  35 |       for (const p of props) record[p] = computed[p as never];
  36 |       out[`${marker}#${n}`] = record;
  37 |     }
  38 |     return out;
  39 |   }, SWEPT_PROPERTIES);
  40 | }
  41 | 
```