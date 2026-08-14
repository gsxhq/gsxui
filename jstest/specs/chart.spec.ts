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

/**
 * Task 7: the network-level proof behind ui/chart.js's own header comment
 * ("pages without charts never parse it") — the BarChart test above only
 * shows chart.render.js's body ran on a chart page, never that it was
 * skipped anywhere else. jstest/harness/modules.go serves ui/*.js
 * byte-for-byte with no bundler/hashing, so the request URL is the literal
 * /ui/chart.render.js — matched here on the substring "chart.render" so the
 * assertion survives if that path ever grows a prefix or transform.
 */
test("chart.render.js loads only on pages with a chart", async ({ page }) => {
  const requests: string[] = [];
  page.on("request", (r) => requests.push(r.url()));

  const noChart = await page.goto("/x/button/variants");
  expect(noChart?.status()).toBe(200);
  expect(requests.filter((u) => u.includes("chart.render"))).toEqual([]);

  requests.length = 0;
  const withChart = await page.goto(route);
  expect(withChart?.status()).toBe(200);
  const bars = page.locator("[data-gsxui-slot-chart] svg path.recharts-rectangle");
  await expect(bars).toHaveCount(6);
  expect(requests.filter((u) => u.includes("chart.render"))).toHaveLength(1);

  // A second chart-page visit in the same browsing context: a navigation
  // tears down the prior document (and with it chart.js's own module-level
  // bodyPromise guard), so this is a fresh module instantiation, not proof
  // the first fetch was reused — what it does prove is that nothing about
  // the stub only works once. It still fetches chart.render.js exactly
  // once and still renders every bar.
  requests.length = 0;
  const secondVisit = await page.goto(route);
  expect(secondVisit?.status()).toBe(200);
  await expect(bars).toHaveCount(6);
  expect(requests.filter((u) => u.includes("chart.render"))).toHaveLength(1);
});

/**
 * Task 6's own fixture (jstest/harness/chart_contract.gsx, served at
 * /f/chart-tooltip) — NOT the public docs example above: it registers a
 * ChartTooltip (basic.gsx doesn't) and gives its "desktop" series a
 * per-scheme ChartSeriesTheme so the dark-mode-flip spec below has a bar
 * whose fill genuinely differs between light and dark (gsxui's shared
 * placeholder theme keeps every --chart-N token identical across schemes
 * on purpose, so a plain Color: "var(--chart-N)" series never could).
 */
const contractRoute = "/f/chart-tooltip";

/**
 * The owner ruling superseding this task's brief: the tooltip's chrome
 * (ui.ChartTooltipTemplate, registry/canonical/chart.gsx) is a real,
 * server-rendered <template> whose shell div carries the tooltip slot's
 * compiled recipe class — like every other gsxui recipe accessor, stylegen
 * desugars class={ chart.Tooltip() } into the resolved style's expanded
 * utility string at generation time (registry/generated/<style>/chart.gsx,
 * ui/chart.gsx), not a runtime reference to the named
 * .gsxui-recipe-chart-tooltip rule — that name exists only as the pack
 * CSS's own authoring/conformance-check vocabulary (registry/styles/
 * <style>/chart.css), never emitted into served HTML. chart.render.js
 * clones the shell's own class attribute back onto the live tooltip on
 * first hover (chartTooltipClasses(), see that file's own doc comment)
 * rather than typing the utilities as a JS literal — asserting the live
 * tooltip's class against the template's own class (not a hardcoded
 * string) stays correct regardless of which style the harness compiles by
 * default. This is the RED this task watched fail before the
 * template/CSS/JS landed: with no template to read from,
 * chartTooltipClasses(container) returned null and tc.shell threw inside
 * tooltipHTML, so hovering never produced a visible tooltip at all.
 */
test("tooltip appears on hover with the per-style recipe class", async ({ page }) => {
  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);

  const expectedClass = await page
    .locator("template[data-gsxui-chart-tooltip-template]")
    .evaluate((tpl) => {
      const shell = (tpl as HTMLTemplateElement).content.querySelector("[data-gsxui-slot-chart-tooltip]");
      return shell?.getAttribute("class") ?? "";
    });
  expect(expectedClass).not.toBe("");

  const chart = page.locator("[data-gsxui-slot-chart]").first();
  const bar = chart.locator("svg path.recharts-rectangle").first();
  await bar.hover();

  const tooltip = page.locator("[data-gsxui-slot-chart-tooltip]");
  await expect(tooltip).toBeVisible();
  await expect(tooltip).toHaveClass(expectedClass);

  // The class match alone doesn't prove the utilities actually compiled —
  // confirm the recipe's own padding utility (px-2.5 py-1.5, every style's
  // chart.css) resolved to a real computed value, not 0px.
  const paddingTop = await tooltip.evaluate((n) => getComputedStyle(n).paddingTop);
  expect(paddingTop).not.toBe("0px");
});

/**
 * Legend renders server-side: block chart.render.js's own network request
 * entirely and confirm the legend text is still present in the HTML
 * Playwright receives — proving ChartLegendContent (registry/canonical/
 * chart.gsx) is real server markup, not something the client draws.
 */
test("legend renders server-side, before chart.render.js loads", async ({ page }) => {
  await page.route("**/ui/chart.render.js", (route) => route.abort());

  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);

  const legend = page.locator("[data-gsxui-slot-chart-legend]");
  await expect(legend).toContainText("Desktop");
  await expect(legend).toContainText("Mobile");
});

/**
 * Dark-mode flip: toggling .dark on <html> must change a bar's computed
 * fill through the CSS cascade alone — Chart.styleBlock's own
 * [data-chart=ID] / .dark [data-chart=ID] rule pair, not a client
 * re-render. evaluateHandle pins down the exact same DOM node across the
 * flip: if chart.render.js re-rendered (rebuilt the SVG), a fresh query
 * for the same selector would return a DIFFERENT node and this equality
 * check would fail.
 */
test("dark-mode flip changes a bar's fill with no re-render", async ({ page }) => {
  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);

  const barSelector = "[data-gsxui-slot-chart] svg path.recharts-rectangle";
  const barHandle = await page.evaluateHandle(
    (sel) => document.querySelector(sel) as Element,
    barSelector,
  );
  const fillBefore = await page.evaluate((el) => getComputedStyle(el as Element).fill, barHandle);

  await page.evaluate(() => document.documentElement.classList.add("dark"));

  const fillAfter = await page.evaluate((el) => getComputedStyle(el as Element).fill, barHandle);
  expect(fillAfter).not.toBe(fillBefore);

  const sameNode = await page.evaluate(
    ([el, sel]) => document.querySelector(sel as string) === el,
    [barHandle, barSelector] as const,
  );
  expect(sameNode).toBe(true);
});

/**
 * htmx-style swap: an out-of-band swap replaces a subtree with a FRESH
 * server render (never the same-document DOM re-captured — that would
 * carry chart.render.js's own data-gsxui-chart-init="true" marker forward
 * on the <script> tag and mask a real re-init, since it's a real DOM
 * attribute chart.render.js sets via dataset). ui/gsxui.js's init()
 * MutationObserver contract re-inits any freshly added
 * [data-gsxui-slot-chart] with no re-scan needed, so the chart's SVG
 * reappears once for the fresh subtree.
 */
test("innerHTML swap re-inits and the chart reappears", async ({ page, request }) => {
  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);
  await expect(page.locator("[data-gsxui-slot-chart] svg")).toHaveCount(1);

  const fresh = await request.get(contractRoute);
  const freshHTML = await fresh.text();

  await page.evaluate((html) => {
    const doc = new DOMParser().parseFromString(html, "text/html");
    const freshChart = doc.querySelector("[data-gsxui-slot-chart]");
    const liveChart = document.querySelector("[data-gsxui-slot-chart]");
    if (!freshChart || !liveChart?.parentElement) {
      throw new Error("fixture missing [data-gsxui-slot-chart]");
    }
    liveChart.parentElement.innerHTML = freshChart.outerHTML;
  }, freshHTML);

  await expect(page.locator("[data-gsxui-slot-chart] svg")).toHaveCount(1);
});
