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
 * ui/chart.js's bodyPromise cache used to hold the REJECTED promise
 * forever once a dynamic import() of chart.render.js failed (a transient
 * network blip, an ad blocker, a flaky path in front of /ui/): every chart
 * on the page, including ones ui/gsxui.js's own init() re-init contract
 * would otherwise recover (a later swap/morph re-invoking initRoot for the
 * same or a different match), silently drew nothing for the rest of the
 * page's life. The fix clears the cache in a .catch so the NEXT initRoot
 * call gets a fresh import() attempt. Simulated here by aborting the
 * request once, then letting a later one through and forcing a re-init via
 * a harmless attribute mutation on the chart root — the same
 * self-healing re-init every other ui/*.js module relies on.
 */
test("a failed chart.render.js load does not wedge later charts on the same page", async ({ page }) => {
  let blocked = true;
  await page.route("**/ui/chart.render.js", (r) => (blocked ? r.abort() : r.continue()));

  const response = await page.goto(route);
  expect(response?.status()).toBe(200);
  await expect(page.locator("[data-gsxui-slot-chart] svg")).toHaveCount(0);

  blocked = false;
  await page.evaluate(() => {
    const el = document.querySelector("[data-gsxui-slot-chart]")!;
    el.setAttribute("data-gsxui-retry-probe", "1");
    el.removeAttribute("data-gsxui-retry-probe");
  });

  await expect(page.locator("[data-gsxui-slot-chart] svg.recharts-surface")).toHaveCount(1);
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
 * jstest/harness/chart_contract.gsx gives desktop a ChartSeries.Icon (raw
 * SVG, no width/height of its own — sized entirely by the tooltip row's
 * own SVG-child sizing utilities, chart.render.js's tooltipHTML). Those
 * utilities used to live only as a JS template literal inside
 * ui/chart.render.js — a string web/site.css's Tailwind scan (which reads
 * .gsx files, never ui/*.js) never saw, so no CSS rule for them ever
 * compiled and the icon rendered at the SVG replaced-element default size
 * instead of 10px. Moving the literal into ui.ChartTooltipTemplate's
 * server-rendered <template> (a real class="..." the scan does see) is
 * what makes this pass: chart.render.js now only ever echoes the class
 * back, never types it. This comment deliberately doesn't spell out the
 * utility literal itself, to avoid tainting any stylesheet build that ever
 * widens its Tailwind @source glob to include spec files.
 */
test("a series Icon in the tooltip is sized by the compiled row utilities, not left at its SVG default", async ({
  page,
}) => {
  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);

  const chart = page.locator("[data-gsxui-slot-chart]").first();
  const bar = chart.locator("svg path.recharts-rectangle").first();
  await bar.hover();

  const icon = page.locator("[data-gsxui-slot-chart-tooltip] svg");
  await expect(icon).toBeVisible();
  const size = await icon.evaluate((n) => {
    const cs = getComputedStyle(n);
    return { width: cs.width, height: cs.height };
  });
  expect(size).toEqual({ width: "10px", height: "10px" });
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
 * templui's ChartContainer carries a wall of [&_.recharts-*] arbitrary-
 * variant selectors (/tmp/templui-ref/chart.templ:126) that re-color
 * Recharts' own hardcoded SVG sentinels (fill="#666" on axis-tick text,
 * stroke="#ccc" on grid lines — see chart.render.js's renderCartesian,
 * ported byte-faithful from upstream's chart.js) into theme tokens. Without
 * it, axis ticks and gridlines paint upstream's literal light-grey no
 * matter the resolved theme or color scheme. The wall lives in
 * registry/styles/<style>/chart.css's .gsxui-recipe-chart rule (chart.Root(),
 * stamped on ui.Chart's own div by registry/canonical/chart.gsx) so every
 * style pack inherits it.
 */
test("axis ticks and gridlines are themed, not upstream's hardcoded grey", async ({ page }) => {
  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);

  const tick = page.locator("[data-gsxui-slot-chart] .recharts-cartesian-axis-tick text").first();
  await expect(tick).toHaveCount(1);
  const gridLine = page.locator("[data-gsxui-slot-chart] .recharts-cartesian-grid line").first();
  await expect(gridLine).toHaveCount(1);

  const mutedForegroundRgb = async () =>
    page.evaluate(() => {
      const probe = document.createElement("span");
      probe.style.color = "var(--muted-foreground)";
      document.body.append(probe);
      const rgb = getComputedStyle(probe).color;
      probe.remove();
      return rgb;
    });

  const tickFillLight = await tick.evaluate((n) => getComputedStyle(n).fill);
  expect(tickFillLight).toBe(await mutedForegroundRgb());

  const gridStroke = await gridLine.evaluate((n) => getComputedStyle(n).stroke);
  expect(gridStroke).not.toBe("rgb(204, 204, 204)");

  await page.evaluate(() => document.documentElement.classList.add("dark"));

  const tickFillDark = await tick.evaluate((n) => getComputedStyle(n).fill);
  expect(tickFillDark).toBe(await mutedForegroundRgb());
  expect(tickFillDark).not.toBe(tickFillLight);
});

/**
 * htmx-style swap: an out-of-band swap replaces a subtree with a FRESH
 * server render (a brand new script/panel/container node triple, never the
 * same-document DOM re-captured). ui/chart.render.js's own once() guard
 * (keyed on the model script node's identity, ui/gsxui.js's own once()
 * helper) has never seen this fresh script node, so it inits normally.
 * ui/gsxui.js's init() MutationObserver contract re-inits any freshly
 * added [data-gsxui-slot-chart] with no re-scan needed, so the chart's SVG
 * reappears once for the fresh subtree. Scoped to svg.recharts-surface
 * (the chart's own drawn SVG), not a bare svg selector: the fixture's
 * desktop series carries a ChartSeries.Icon that ChartLegendContent also
 * renders server-side into a plain <svg> beside the recharts one.
 */
test("innerHTML swap re-inits and the chart reappears", async ({ page, request }) => {
  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);
  await expect(page.locator("[data-gsxui-slot-chart] svg.recharts-surface")).toHaveCount(1);

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

  await expect(page.locator("[data-gsxui-slot-chart] svg.recharts-surface")).toHaveCount(1);
});

/**
 * The morph-back case a real htmx/idiomorph-style reconciliation produces:
 * the SAME script/panel/container node objects, patched in place to match
 * fresh server markup rather than replaced — unlike the swap test above.
 * Before the final-review fix wave, ui/chart.render.js's own re-init guard
 * read script.dataset.gsxuiChartInit, an attribute the server never emits;
 * a morph that reconciles this container back to server markup strips it
 * (nothing in the fresh HTML sets it), so initPanel() ran a second full
 * bind — a second ResizeObserver, a second document keydown listener, a
 * second set of panel listeners, all stacked on the SAME still-live nodes.
 * The fix (ui/chart.render.js's initPanelOnce) keys re-init on the script
 * node's own identity via ui/gsxui.js's once() (a WeakSet, immune to
 * attribute stripping), so this same-node case is now blocked outright,
 * verified here by counting real ResizeObserver construction calls — the
 * most direct proxy available for "did initPanel's resource-bind path run
 * again," since the tooltip/cursor HTML side effects downstream of it are
 * idempotent per call and wouldn't visibly differ if doubled. Counted as a
 * DELTA across the mutation, not an absolute value: ui/carousel.js
 * instantiates its own module-level ResizeObserver at load time (the
 * shared-instance pattern that file's own init() comment documents), and
 * /ui/index.js — the bundle this harness page loads — pulls every ui/*.js
 * module in, carousel.js included, so the window-global count this test
 * observes is never "just the chart's."
 */
test("a same-node attribute mutation (morph-back) does not double-bind ResizeObserver", async ({ page }) => {
  await page.addInitScript(() => {
    (window as any).__gsxuiChartROCount = 0;
    const NativeRO = window.ResizeObserver;
    window.ResizeObserver = class extends NativeRO {
      constructor(...args: ConstructorParameters<typeof ResizeObserver>) {
        (window as any).__gsxuiChartROCount++;
        super(...args);
      }
    };
  });

  const response = await page.goto(contractRoute);
  expect(response?.status()).toBe(200);
  await expect(page.locator("[data-gsxui-slot-chart] svg.recharts-surface")).toHaveCount(1);

  // Settle past the entrance animation (400ms, same duration the narrow-
  // flex spec above waits out) so no in-flight animation-driven mutation
  // is still scheduling an init pass when the baseline is captured.
  await page.waitForTimeout(600);
  const before = await page.evaluate(() => (window as any).__gsxuiChartROCount);

  // Mutate an attribute on the SAME script node a morph-back would patch —
  // no node is replaced, so a genuine swap's "fresh script node" escape
  // hatch does not apply here.
  await page.evaluate(() => {
    const script = document.querySelector("[data-gsxui-slot-chart] script[data-gsxui-chart-model]")!;
    script.setAttribute("data-gsxui-chart-init", "true"); // the pre-fix guard's own flag: proves it's no longer read
    script.removeAttribute("data-gsxui-chart-init");
  });

  // ui/gsxui.js's init() flushes on the next microtask; give it a real
  // task boundary before asserting the count never moved past its baseline.
  await page.waitForTimeout(50);
  expect(await page.evaluate(() => (window as any).__gsxuiChartROCount)).toBe(before);
});

/**
 * jstest/harness/chart_contract.gsx's ChartNarrowFlexFixture (served at
 * /f/chart-narrow-flex): ui.Chart mounted with NO caller class override,
 * inside a 340px-wide flex column — the same shape site/stylepreview/
 * gallery.gsx.src's galleryChartCard hits. Without container-level
 * protection, ui/chart.render.js's renderCartesian re-measures
 * panel.clientHeight on every entrance-animation frame and draws the SVG at
 * that height; Chart's own div is a flex item of the narrow column, and a
 * flex item's default min-height:auto lets its own drawn content override
 * the aspect-ratio-derived height, so each frame's slightly-taller draw
 * feeds the next frame's measurement — a real runaway (190px climbing past
 * 900px) for the whole ~400ms bar entrance. min-h-0 on the recipe
 * (registry/styles/<style>/chart.css's .gsxui-recipe-chart) is the fix:
 * sampling clientHeight once early and once after the entrance animation
 * must finish (bar duration 400ms) has to land on the SAME, aspect-derived
 * value both times.
 */
test("chart height stays at its aspect-derived value in a narrow flex column, not growing during the entrance animation", async ({
  page,
}) => {
  const response = await page.goto("/f/chart-narrow-flex");
  expect(response?.status()).toBe(200);

  const panel = page.locator("[data-chart-narrow-flex] [data-gsxui-slot-chart] > div").first();
  await panel.waitFor();

  const early = await panel.evaluate((el) => el.clientHeight);
  // The bar entrance animation runs 400ms (ui/chart.render.js); wait past
  // it so a runaway (still climbing every frame) has nowhere left to hide.
  await page.waitForTimeout(600);
  const settled = await panel.evaluate((el) => el.clientHeight);

  expect(settled).toBe(early);
  // aspect-video at the fixture's 340px width: 340 * 9 / 16 ≈ 191. A
  // generous ceiling well under the runaway's own >900px failure mode.
  expect(settled).toBeLessThan(250);
});
