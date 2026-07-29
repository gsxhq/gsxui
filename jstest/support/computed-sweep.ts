import type { Page } from "@playwright/test";

export const SWEPT_PROPERTIES = [
  "display", "position", "borderRadius", "borderWidth", "borderColor",
  "width", "height", "paddingLeft", "paddingRight", "paddingTop",
  "paddingBottom", "gap", "fontSize", "fontWeight", "lineHeight",
  "color", "backgroundColor", "opacity", "boxShadow", "textDecorationLine",
  "alignItems", "justifyContent",
] as const;

/**
 * sweepComputedStyles records the resting computed style of every element
 * carrying a `data-gsxui-slot-*` attribute on the given page. It is the only
 * thing that reliably catches CSS-layer precedence regressions: they change
 * rendering without failing any assertion. Repeats of the same marker on a
 * page are disambiguated with a `#n` suffix in DOM order.
 */
export async function sweepComputedStyles(
  page: Page,
  url: string,
): Promise<Record<string, Record<string, string>>> {
  await page.goto(url);
  return page.evaluate((props: readonly string[]) => {
    const out: Record<string, Record<string, string>> = {};
    const seen = new Map<string, number>();
    for (const el of document.querySelectorAll("*")) {
      const marker = [...el.attributes]
        .map((a) => a.name)
        .find((n) => n.startsWith("data-gsxui-slot-"));
      if (!marker) continue;
      const n = (seen.get(marker) ?? 0) + 1;
      seen.set(marker, n);
      const computed = getComputedStyle(el);
      const record: Record<string, string> = {};
      for (const p of props) record[p] = computed[p as never];
      out[`${marker}#${n}`] = record;
    }
    return out;
  }, SWEPT_PROPERTIES);
}
