import { expect, test } from "../support/fixtures";

const route = "/x/chart/basic";

/**
 * ui/chart.js's stub lazily imports ui/chart.render.js on the first
 * [data-gsxui-slot-chart] it finds and hands off to renderChart(), which
 * reads the sibling <script data-gsxui-chart-model> JSON and draws the SVG.
 * This is the one check Task 5 owns: the client actually drew something,
 * themed by the server's --color-<key> variables (ChartConfig.styleBlock).
 *
 * chart.render.js is a byte-faithful port of templui's own chart.js (see
 * that file's header): bars paint as <path class="recharts-rectangle">
 * (Recharts' own Bar shape, rounded-corner capable, even at radius 0), not
 * <rect> — the selector below matches the real port, not a stand-in.
 */
test("BarChart renders themed SVG bars", async ({ page }) => {
  const response = await page.goto(route);
  expect(response?.status()).toBe(200);

  const chart = page.locator("[data-gsxui-slot-chart]").first();
  // site/examples/chart/basic.gsx: 3 months (Jan/Feb/Mar) x 2 series
  // (desktop, mobile) = 6 bars, one <path class="recharts-rectangle"> each.
  const bars = chart.locator("svg path.recharts-rectangle");
  await expect(bars).toHaveCount(6);

  const fill = await bars.first().evaluate((n) => getComputedStyle(n).fill);
  expect(fill).not.toBe("");

  // Declaration order is paint order (ChartBar("desktop") comes first in
  // basic.gsx), so the first painted bar is desktop's, themed through
  // styleBlock's --color-desktop: var(--chart-1) — resolve both sides to
  // rgb() and compare, the same probe-element technique toaster.spec.ts
  // uses for its own CSS-variable assertions.
  const chart1Rgb = await page.evaluate(() => {
    const probe = document.createElement("span");
    probe.style.color = "var(--chart-1)";
    document.body.append(probe);
    const rgb = getComputedStyle(probe).color;
    probe.remove();
    return rgb;
  });
  expect(fill).toBe(chart1Rgb);
});
