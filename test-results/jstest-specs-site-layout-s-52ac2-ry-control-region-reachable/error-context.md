# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/site-layout.spec.ts >> short desktop theme workspaces keep every control region reachable
- Location: jstest/specs/site-layout.spec.ts:222:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/site/theme", waiting until "load"

```

# Test source

```ts
  126 |   });
  127 |   await trigger.focus();
  128 |   await trigger.press("Enter");
  129 |   await expect(trigger).toHaveAttribute("aria-expanded", "true");
  130 | 
  131 |   const gettingStarted = compactNav.locator(
  132 |     'a[href="/docs/getting-started"]',
  133 |   );
  134 |   const button = compactNav.locator('a[href="/components/button"]');
  135 |   await expect(gettingStarted).toBeVisible();
  136 |   await expect(gettingStarted).toBeFocused();
  137 |   await expect(button).toBeVisible();
  138 | 
  139 |   await gettingStarted.click();
  140 |   await expect(page).toHaveURL(/\/docs\/getting-started$/);
  141 | 
  142 |   const nextCompactNav = page.locator("[data-site-docs-mobile-nav]");
  143 |   await nextCompactNav
  144 |     .getByRole("button", { name: "Open documentation navigation" })
  145 |     .press("Enter");
  146 |   await nextCompactNav.locator('a[href="/components/button"]').click();
  147 |   await expect(page).toHaveURL(/\/components\/button$/);
  148 | });
  149 | 
  150 | test("theme workspace contains desktop overflow and preserves narrow document order", async ({
  151 |   page,
  152 | }) => {
  153 |   await page.setViewportSize({ width: 1280, height: 900 });
  154 |   await page.goto("/site/theme");
  155 | 
  156 |   const stylePanel = page.locator("[data-theme-style-panel]");
  157 |   const previewPanel = page.locator("[data-theme-preview-panel]");
  158 |   const controlsPanel = page.locator("[data-theme-controls-panel]");
  159 | 
  160 |   await expect(stylePanel).toBeVisible();
  161 |   await expect(previewPanel).toBeVisible();
  162 |   await expect(controlsPanel).toBeVisible();
  163 |   await expect(page.locator("[data-theme-preview-frame]")).toHaveCount(1);
  164 | 
  165 |   const desktopGeometry = await page.evaluate(() => {
  166 |     const controls = document.querySelector<HTMLElement>(
  167 |       "[data-theme-controls-panel]",
  168 |     );
  169 |     const preview = document.querySelector<HTMLElement>(
  170 |       "[data-theme-preview-panel]",
  171 |     );
  172 |     if (!controls || !preview) {
  173 |       throw new Error("theme workspace panels are missing");
  174 |     }
  175 |     return {
  176 |       documentOverflow:
  177 |         document.documentElement.scrollHeight >
  178 |         document.documentElement.clientHeight,
  179 |       controlsWidth: controls.getBoundingClientRect().width,
  180 |       previewWidth: preview.getBoundingClientRect().width,
  181 |       controlsOverflow: getComputedStyle(controls).overflowY,
  182 |       controlsScrollHeight: controls.scrollHeight,
  183 |       controlsClientHeight: controls.clientHeight,
  184 |     };
  185 |   });
  186 |   expect(desktopGeometry.documentOverflow).toBe(false);
  187 |   expect(desktopGeometry.controlsWidth).toBeLessThan(
  188 |     desktopGeometry.previewWidth,
  189 |   );
  190 |   expect(desktopGeometry.controlsOverflow).toBe("auto");
  191 |   expect(desktopGeometry.controlsScrollHeight).toBeGreaterThan(
  192 |     desktopGeometry.controlsClientHeight,
  193 |   );
  194 | 
  195 |   await page.setViewportSize({ width: 900, height: 900 });
  196 | 
  197 |   const narrowGeometry = await page.evaluate(() => {
  198 |     const top = (selector: string) => {
  199 |       const element = document.querySelector<HTMLElement>(selector);
  200 |       if (!element) {
  201 |         throw new Error(`missing ${selector}`);
  202 |       }
  203 |       return element.getBoundingClientRect().top;
  204 |     };
  205 |     return {
  206 |       documentOverflow:
  207 |         document.documentElement.scrollHeight >
  208 |         document.documentElement.clientHeight,
  209 |       styleTop: top("[data-theme-style-panel]"),
  210 |       previewTop: top("[data-theme-preview-panel]"),
  211 |       controlsTop: top("[data-theme-controls-panel]"),
  212 |       iframeCount: document.querySelectorAll("[data-theme-preview-frame]")
  213 |         .length,
  214 |     };
  215 |   });
  216 |   expect(narrowGeometry.documentOverflow).toBe(true);
  217 |   expect(narrowGeometry.styleTop).toBeLessThan(narrowGeometry.previewTop);
  218 |   expect(narrowGeometry.previewTop).toBeLessThan(narrowGeometry.controlsTop);
  219 |   expect(narrowGeometry.iframeCount).toBe(1);
  220 | });
  221 | 
  222 | test("short desktop theme workspaces keep every control region reachable", async ({
  223 |   page,
  224 | }) => {
  225 |   await page.setViewportSize({ width: 1280, height: 360 });
> 226 |   await page.goto("/site/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  227 | 
  228 |   const workspaceMain = page.locator(
  229 |     'body[data-site-layout="workspace"] > div > main',
  230 |   );
  231 |   const controls = page.locator("[data-theme-controls-panel]");
  232 |   const geometry = await page.evaluate(() => {
  233 |     const main = document.querySelector<HTMLElement>(
  234 |       'body[data-site-layout="workspace"] > div > main',
  235 |     );
  236 |     const controlsPanel = document.querySelector<HTMLElement>(
  237 |       "[data-theme-controls-panel]",
  238 |     );
  239 |     if (!main || !controlsPanel) {
  240 |       throw new Error("theme workspace scrolling regions are missing");
  241 |     }
  242 |     return {
  243 |       documentWidth: document.documentElement.scrollWidth,
  244 |       viewportWidth: document.documentElement.clientWidth,
  245 |       mainOverflowY: getComputedStyle(main).overflowY,
  246 |       mainScrollHeight: main.scrollHeight,
  247 |       mainClientHeight: main.clientHeight,
  248 |       controlsClientHeight: controlsPanel.clientHeight,
  249 |       controlsScrollHeight: controlsPanel.scrollHeight,
  250 |       controlsOverflowY: getComputedStyle(controlsPanel).overflowY,
  251 |       controlsScrollWidth: controlsPanel.scrollWidth,
  252 |       controlsClientWidth: controlsPanel.clientWidth,
  253 |     };
  254 |   });
  255 | 
  256 |   expect(geometry.documentWidth).toBeLessThanOrEqual(geometry.viewportWidth);
  257 |   expect(geometry.mainOverflowY).toBe("auto");
  258 |   expect(geometry.mainScrollHeight).toBeGreaterThan(geometry.mainClientHeight);
  259 |   expect(geometry.controlsClientHeight).toBeGreaterThan(0);
  260 |   expect(geometry.controlsOverflowY).toBe("auto");
  261 |   expect(geometry.controlsScrollHeight).toBeGreaterThan(
  262 |     geometry.controlsClientHeight,
  263 |   );
  264 |   expect(geometry.controlsScrollWidth).toBeLessThanOrEqual(
  265 |     geometry.controlsClientWidth,
  266 |   );
  267 | 
  268 |   for (const name of [
  269 |     "Mode and palette",
  270 |     "Preset JSON",
  271 |     "Theme CSS",
  272 |     "Share and install",
  273 |   ]) {
  274 |     const heading = controls.getByRole("heading", { name, exact: true });
  275 |     await heading.evaluate((element) =>
  276 |       element.scrollIntoView({ block: "nearest" }),
  277 |     );
  278 |     expect(
  279 |       await heading.evaluate((element) => {
  280 |         const rect = element.getBoundingClientRect();
  281 |         const scroller = element
  282 |           .closest("[data-theme-controls-panel]")
  283 |           ?.getBoundingClientRect();
  284 |         return (
  285 |           !!scroller &&
  286 |           rect.bottom > Math.max(0, scroller.top) &&
  287 |           rect.top < Math.min(innerHeight, scroller.bottom)
  288 |         );
  289 |       }),
  290 |       `${name} should be reachable in the short workspace`,
  291 |     ).toBe(true);
  292 |   }
  293 | 
  294 |   await expect(workspaceMain).toBeVisible();
  295 |   await expect(page.locator("[data-theme-preview-frame]")).toHaveCount(1);
  296 | });
  297 | 
  298 | test("production-composed theme workspace loads its controller and preview cleanly", async ({
  299 |   page,
  300 | }) => {
  301 |   const consoleErrors: string[] = [];
  302 |   const pageErrors: string[] = [];
  303 |   const failedRequests: string[] = [];
  304 |   const failedResponses: string[] = [];
  305 |   page.on("console", (message) => {
  306 |     if (message.type() === "error") consoleErrors.push(message.text());
  307 |   });
  308 |   page.on("pageerror", (error) => pageErrors.push(String(error)));
  309 |   page.on("requestfailed", (request) => {
  310 |     failedRequests.push(
  311 |       `${request.method()} ${request.url()} ${request.failure()?.errorText ?? ""}`,
  312 |     );
  313 |   });
  314 |   page.on("response", (response) => {
  315 |     if (response.status() >= 400) {
  316 |       failedResponses.push(`${response.status()} ${response.url()}`);
  317 |     }
  318 |   });
  319 | 
  320 |   await page.setViewportSize({ width: 1280, height: 900 });
  321 |   await page.goto("/site/theme");
  322 |   await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  323 |   await expect(page.locator("[data-theme-preview-retry]")).toBeHidden();
  324 | 
  325 |   const theme = page.locator('[data-theme-picker="theme"]');
  326 |   await theme.locator("[data-theme-picker-trigger]").click();
```