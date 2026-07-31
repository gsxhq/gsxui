import type { Page } from "@playwright/test";

export const SWEPT_PROPERTIES = [
  "display",
  "position",
  "borderRadius",
  "borderWidth",
  "borderColor",
  "width",
  "height",
  "paddingLeft",
  "paddingRight",
  "paddingTop",
  "paddingBottom",
  "gap",
  "fontSize",
  "fontWeight",
  "lineHeight",
  "color",
  "backgroundColor",
  "opacity",
  "boxShadow",
  "textDecorationLine",
  "alignItems",
  "justifyContent",
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
    // Sample animations at a fixed point instead of whenever the sweep happens
    // to run. Sidebar's skeletons pulse, so an unpaused sweep caught opacity
    // mid-cycle (0.999919 rather than 1) about one run in five — and because
    // this sweep is the migration's acceptance gate, that flake read as a
    // regression and cost an investigation each time. Pausing at currentTime 0
    // makes "resting" mean the animation's start value, deterministically.
    // Pausing is not the same as disabling: the animation still applies, so a
    // rule that changes its start value still shows up here. Timing itself —
    // animation-name, duration — is outside SWEPT_PROPERTIES entirely and is
    // covered by pins instead (Drawer's 200ms duration, Spinner's spin).
    for (const animation of document.getAnimations()) {
      animation.pause();
      animation.currentTime = 0;
    }
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
