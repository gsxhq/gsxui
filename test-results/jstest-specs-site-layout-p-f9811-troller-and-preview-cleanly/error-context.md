# Instructions

- Following Playwright test failed.
- Explain why, be concise, respect Playwright best practices.
- Provide a snippet of code with the fix, if possible.

# Test info

- Name: jstest/specs/site-layout.spec.ts >> production-composed theme workspace loads its controller and preview cleanly
- Location: jstest/specs/site-layout.spec.ts:298:1

# Error details

```
Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
Call log:
  - navigating to "/site/theme", waiting until "load"

```

# Test source

```ts
  221 | 
  222 | test("short desktop theme workspaces keep every control region reachable", async ({
  223 |   page,
  224 | }) => {
  225 |   await page.setViewportSize({ width: 1280, height: 360 });
  226 |   await page.goto("/site/theme");
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
> 321 |   await page.goto("/site/theme");
      |              ^ Error: page.goto: Protocol error (Page.navigate): Cannot navigate to invalid URL
  322 |   await expect(page.locator("[data-theme-preview-status]")).toHaveText("Live");
  323 |   await expect(page.locator("[data-theme-preview-retry]")).toBeHidden();
  324 | 
  325 |   const theme = page.locator('[data-theme-picker="theme"]');
  326 |   await theme.locator("[data-theme-picker-trigger]").click();
  327 |   await theme.getByRole("radio", { name: "Blue", exact: true }).click();
  328 |   await expect(
  329 |     theme.locator("[data-theme-selection-value]"),
  330 |   ).toHaveText("Blue");
  331 | 
  332 |   const catalog = JSON.parse(
  333 |     (await page.locator("[data-theme-schema]").textContent()) ?? "{}",
  334 |   );
  335 |   await expect
  336 |     .poll(() =>
  337 |       page
  338 |         .frameLocator("[data-theme-preview-frame]")
  339 |         .locator("html")
  340 |         .evaluate((element) =>
  341 |           element.style.getPropertyValue("--primary").trim(),
  342 |         ),
  343 |     )
  344 |     .toBe(catalog.palette.resolved.neutral.blue.light.primary);
  345 | 
  346 |   expect(consoleErrors).toEqual([]);
  347 |   expect(pageErrors).toEqual([]);
  348 |   expect(failedRequests).toEqual([]);
  349 |   expect(failedResponses).toEqual([]);
  350 | });
  351 | 
```